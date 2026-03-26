package models

import "time"

type SHT40Data struct {
	Time        time.Time `gorm:"index" json:"timestamp"`
	DeviceID    string    `gorm:"index" json:"deviceId"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
}

func (SHT40Data) TableName() string {
	return "sht40_data"
}
