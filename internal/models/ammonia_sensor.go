package models

import "time"

type AmmoniaSensorData struct {
	Time     time.Time `gorm:"index" json:"timestamp"`
	DeviceID string    `gorm:"index" json:"deviceId"`
	Ammonia  float64   `json:"ammonia"`
	RawData  string    `json:"rawData"`
}

func (AmmoniaSensorData) TableName() string {
	return "ammonia_sensor_data"
}
