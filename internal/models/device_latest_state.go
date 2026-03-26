package models

import "time"

type DeviceLatestState struct {
	DeviceID string `gorm:"primaryKey"`

	Temperature float64
	Humidity    float64
	Lux         float64

	X float64
	Y float64
	Z float64

	Ammonia float64

	Nitrogen   float64
	Phosphorus float64
	Potassium  float64

	Moisture float64
	EC       float64
	PH       float64
	Salinity float64

	Speed    float64
	Distance float64

	LastPacketID int
	LastSeen     time.Time
}
