package models

import "time"

type Device struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	DeviceID      string     `gorm:"uniqueIndex;not null" json:"deviceId"` // BLE Sensor MAC
	MobileID      string     `gorm:"index" json:"mobileId"`                // The associated mobile phone
	DeviceAddress string     `json:"deviceAddress"`
	SensorType    SensorType `gorm:"type:sensor_type" json:"sensorType"`
	LastSeen      time.Time  `json:"lastSeen"`                             // Updated on every HTTP packet received
	CreatedAt     time.Time  `json:"createdAt"`
}
