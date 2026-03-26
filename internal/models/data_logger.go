package models

import (
	"time"

	"gorm.io/datatypes"
)

type DataLoggerData struct {
	Time            time.Time      `gorm:"index" json:"timestamp"`
	DeviceID        string         `gorm:"index" json:"deviceId"`
	CurrentPacketID int            `json:"currentPacketId"`
	LastPacketID    int            `json:"lastPacketId"`
	AccelData       datatypes.JSON `json:"payloadAccel"`
	RawData         string         `json:"rawData"`
}

func (DataLoggerData) TableName() string {
	return "data_logger_data"
}
