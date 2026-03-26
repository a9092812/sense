package models

import "time"

type LuxSensorData struct {
	Time     time.Time `gorm:"index" json:"timestamp"`
	DeviceID string    `gorm:"index" json:"deviceId"`
	Lux      float64   `json:"lux"`
	RawData  string    `json:"rawData"`
}

func (LuxSensorData) TableName() string {
	return "lux_sensor_data"
}
