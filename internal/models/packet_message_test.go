package models

import (
	"encoding/json"
	"testing"
)

func TestSensorPacketMessage_Decoding(t *testing.T) {
	rawData := `{"deviceId":"DEV123","parsedType":"SHT40","parsedData":{"temperature":25.5,"humidity":60.2},"timestamp":1700000000}`
	
	msg, err := Decode([]byte(rawData))
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if msg.DeviceID != "DEV123" {
		t.Errorf("Expected DEV123, got %s", msg.DeviceID)
	}

	if msg.Float("temperature") != 25.5 {
		t.Errorf("Expected 25.5, got %f", msg.Float("temperature"))
	}

	if msg.Float("humidity") != 60.2 {
		t.Errorf("Expected 60.2, got %f", msg.Float("humidity"))
	}
}

func TestSensorPacketMessage_Helpers(t *testing.T) {
	msg := &SensorPacketMessage{
		ParsedData: json.RawMessage(`{"intKey":100, "strKey":"val", "floatKey":1.23}`),
	}

	if msg.Int("intKey") != 100 {
		t.Errorf("Expected 100, got %d", msg.Int("intKey"))
	}

	if msg.String("strKey") != "val" {
		t.Errorf("Expected val, got %s", msg.String("strKey"))
	}

	if msg.Float("floatKey") != 1.23 {
		t.Errorf("Expected 1.23, got %f", msg.Float("floatKey"))
	}
}
