package repository

import (
	"time"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// EnsureDeviceExists upserts the sensor row, always stamping mobile_id and last_seen.
// Uses raw SQL ON CONFLICT to reliably update the two columns even when the
// sensor_type ENUM is already set (GORM's clause helper can silently skip the
// UPDATE when the row already exists with a custom Postgres type).
func (r *DeviceRepository) EnsureDeviceExists(deviceID, mobileID, address, sensorType string) error {
	// 1. Ensure the device itself is tracked globally
	err := r.db.Exec(`
		INSERT INTO devices (device_id, mobile_id, device_address, sensor_type, last_seen, created_at)
		VALUES (?, ?, ?, ?::sensor_type, ?, NOW())
		ON CONFLICT (device_id) DO UPDATE
		  SET mobile_id  = EXCLUDED.mobile_id,
		      last_seen  = EXCLUDED.last_seen
	`, deviceID, mobileID, address, sensorType, time.Now()).Error
	if err != nil {
		return err
	}

	// 2. Update the mobile-to-sensor affinity (so multiple mobiles can "own" the same sensor)
	return r.db.Exec(`
		INSERT INTO mobile_sensor_affinity (mobile_id, device_id, last_seen)
		VALUES (?, ?, ?)
		ON CONFLICT (mobile_id, device_id) DO UPDATE
		  SET last_seen = EXCLUDED.last_seen
	`, mobileID, deviceID, time.Now()).Error
}

func (r *DeviceRepository) ListDevices() ([]DeviceRow, error) {
	var rows []DeviceRow
	err := r.db.Raw(`SELECT id, device_id, mobile_id, device_address, sensor_type, last_seen, created_at FROM devices ORDER BY id`).Scan(&rows).Error
	return rows, err
}

func (r *DeviceRepository) ListSensorsByMobile(mobileID string) ([]DeviceRow, error) {
	var rows []DeviceRow
	// We join with mobile_sensor_affinity to get the last_seen specifically for this mobile.
	// We consider a sensor "Live" if seen in the last 15 seconds.
	err := r.db.Raw(`
		SELECT 
			d.id, 
			d.device_id, 
			a.mobile_id, 
			d.device_address, 
			d.sensor_type, 
			a.last_seen as last_seen, 
			d.created_at,
			(a.last_seen >= NOW() - INTERVAL '15 seconds') as is_live
		FROM devices d
		JOIN mobile_sensor_affinity a ON d.device_id = a.device_id
		WHERE a.mobile_id = ?
		ORDER BY a.last_seen DESC
	`, mobileID).Scan(&rows).Error
	return rows, err
}

// DeviceRow is a plain struct used only for decoding query results —
// avoids the GORM model overhead and the enum type casting issue.
type DeviceRow struct {
	ID            uint      `json:"id"`
	DeviceID      string    `json:"deviceId"`
	MobileID      string    `json:"mobileId"`
	DeviceAddress string    `json:"deviceAddress"`
	SensorType    string    `json:"sensorType"`
	LastSeen      time.Time `json:"lastSeen"`
	CreatedAt     time.Time `json:"createdAt"`
	IsLive        bool      `json:"isLive"`
}