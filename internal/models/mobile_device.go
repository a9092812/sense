package models

import "time"

// MobileDevice represents the mobile phone from which sensor data is coming
type MobileDevice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MobileID  string    `gorm:"uniqueIndex;not null" json:"mobileId"` // Unique ID of the mobile device
	Name      string    `json:"name"`                                 // Optional user-friendly name
	LastSeen  time.Time `json:"lastSeen"`                             // Updated on every HTTP batch
	CreatedAt time.Time `json:"createdAt"`
}

// MobileStatus represents the computed status for the UI
type MobileStatus struct {
	MobileDevice
	IsOnline bool `json:"isOnline"`
}
