package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Kartik30R/sense/internal/dto"
)

type SensorPacketMessage struct {
	MobileID         string          `json:"mobileId"`
	DeviceID         string          `json:"deviceId"`
	DeviceAddress    string          `json:"deviceAddress"`
	RSSI             int             `json:"rssi"`
	RawAdvertisement string          `json:"rawAdvertisement"`
	ParsedType       string          `json:"parsedType"`
	ParsedData       json.RawMessage `json:"parsedData"`
	Timestamp        int64           `json:"timestamp"`
	cache map[string]interface{} `json:"-"`
}

func (m *SensorPacketMessage) GetTime() time.Time {
	if m.Timestamp > 2500000000 { // Year 2049 in seconds, likely ms
		return time.UnixMilli(m.Timestamp)
	}
	return time.Unix(m.Timestamp, 0)
}

func FromDTO(req dto.PacketUploadRequest) SensorPacketMessage {

	msg := SensorPacketMessage{
		MobileID:         req.MobileID,
		DeviceID:         req.DeviceID,
		DeviceAddress:    req.DeviceAddress,
		RSSI:             req.RSSI,
		RawAdvertisement: req.RawAdvertisement,
		ParsedType:       req.ParsedType,
		ParsedData:       req.ParsedData,
		Timestamp:        req.Timestamp,
	}

	if msg.ParsedType == "" {
		msg.ParsedType = (&msg).String("type")
	}

	return msg
}

func Encode(msg SensorPacketMessage) ([]byte, error) {

	return json.Marshal(msg)
}

func Decode(data []byte) (*SensorPacketMessage, error) {

	var msg SensorPacketMessage

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	if msg.ParsedType == "" {
		msg.ParsedType = (&msg).String("type")
	}

	return &msg, nil
}

func (m *SensorPacketMessage) parse() {

	if m.cache != nil {
		return
	}

	var data map[string]interface{}

	_ = json.Unmarshal(m.ParsedData, &data)

	m.cache = data
}

func (m *SensorPacketMessage) String(key string) string {

	m.parse()

	v, ok := m.cache[key]

	if !ok {
		return ""
	}

	switch val := v.(type) {

	case string:
		return val

	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)

	default:
		return ""
	}
}

func (m *SensorPacketMessage) Float(key string) float64 {

	m.parse()

	v, ok := m.cache[key]

	if !ok {
		return 0
	}

	switch val := v.(type) {

	case float64:
		return val

	case string:

		f, _ := strconv.ParseFloat(val, 64)

		return f
	}

	return 0
}

func (m *SensorPacketMessage) Int(key string) int {

	m.parse()

	v, ok := m.cache[key]

	if !ok {
		return 0
	}

	switch val := v.(type) {

	case float64:
		return int(val)

	case string:

		i, _ := strconv.Atoi(val)

		return i
	}

	return 0
}

func (m *SensorPacketMessage) JSON(key string) []byte {

	m.parse()

	v, ok := m.cache[key]

	if !ok {
		return nil
	}

	data, _ := json.Marshal(v)

	return data
}
