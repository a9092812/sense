package models

import "time"

type TempLoggerData struct {
	Time           time.Time `gorm:"index" json:"timestamp"`
	DeviceID       string    `gorm:"index" json:"deviceId"`
	Temperature    float64   `json:"temperature"`
	Humidity       float64   `json:"humidity"`
	RawTemperature int       `json:"rawTemperature"`
	RawHumidity    int       `json:"rawHumidity"`
	RawData        string    `json:"rawData"`
	DeviceAddress  string    `json:"deviceAddress"`
}

func (TempLoggerData) TableName() string {
	return "temp_logger_data"
}
