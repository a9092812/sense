package dto

import "encoding/json"

type PacketUploadRequest struct {
	MobileID         string          `json:"mobileId" binding:"required"` // The phone sending the data
	DeviceID         string          `json:"deviceId" binding:"required"` // The sensor MAC
	DeviceAddress    string          `json:"deviceAddress" binding:"required"`
	RSSI             int             `json:"rssi"`
	RawAdvertisement string          `json:"rawAdvertisement"`
	ParsedType       string          `json:"parsedType"`
	ParsedData       json.RawMessage `json:"parsedData" binding:"required"`
	Timestamp        int64           `json:"timestamp" binding:"required"`
}
