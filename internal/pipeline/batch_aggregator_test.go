package pipeline

import (
	"testing"
	"time"

	"github.com/Kartik30R/sense/internal/models"
)

// Mock Registry or use memory? 
// For now, let's just test that the batch grows.

func TestBatchAggregator_Add(t *testing.T) {
	agg := &BatchAggregator{}
	
	msg := &models.SensorPacketMessage{
		DeviceID:   "DEV1",
		ParsedType: "SHT40",
		ParsedData: []byte(`{"temperature":20,"humidity":50}`),
		Timestamp:  time.Now().Unix(),
	}

	agg.Add(msg)

	if len(agg.sht40Batch) != 1 {
		t.Errorf("Expected 1 item in batch, got %d", len(agg.sht40Batch))
	}

	if agg.sht40Batch[0].Temperature != 20 {
		t.Errorf("Expected temperature 20, got %f", agg.sht40Batch[0].Temperature)
	}
}
