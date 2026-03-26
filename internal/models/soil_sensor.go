package models

import "time"

type SoilSensorData struct {
	Time        time.Time `gorm:"index" json:"timestamp"`
	DeviceID    string    `gorm:"index" json:"deviceId"`
	Nitrogen    float64   `json:"nitrogen"`
	Phosphorus  float64   `json:"phosphorus"`
	Potassium   float64   `json:"potassium"`
	Moisture    float64   `json:"moisture"`
	Temperature float64   `json:"temperature"`
	EC          float64   `json:"ec"`
	PH          float64   `json:"pH"`
	Salinity    float64   `json:"salinity"`
}

func (SoilSensorData) TableName() string {
	return "soil_sensor_data"
}
