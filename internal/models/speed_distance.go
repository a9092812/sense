package models

import "time"

type SpeedDistanceData struct {
	Time     time.Time `gorm:"index" json:"timestamp"`
	DeviceID string    `gorm:"index" json:"deviceId"`
	Speed    float64   `json:"speed"`
	Distance float64   `json:"distance"`
}

func (SpeedDistanceData) TableName() string {
	return "speed_distance_data"
}
