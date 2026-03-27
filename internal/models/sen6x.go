package models

import "time"

type Sen6xData struct {
	Time        time.Time `gorm:"index" json:"timestamp"`
	DeviceID    string    `gorm:"index" json:"deviceId"`
	PM1         float64   `json:"pm1"`
	PM25        float64   `json:"pm25"`
	PM40        float64   `json:"pm4"`
	PM100       float64   `json:"pm10"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	CO2         float64   `json:"co2"`
	VOC         float64   `json:"voc"`
	NOx         float64   `json:"nox"`
}

func (Sen6xData) TableName() string {
	return "sen6x_data"
}
