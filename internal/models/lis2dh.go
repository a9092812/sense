package models

import "time"

type LIS2DHData struct {
	Time     time.Time `gorm:"index" json:"timestamp"`
	DeviceID string    `gorm:"index" json:"deviceId"`
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	Z        float64   `json:"z"`
}

func (LIS2DHData) TableName() string {
	return "lis2dh_data"
}
