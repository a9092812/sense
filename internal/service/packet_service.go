package services

import (
	"context"

	"github.com/Kartik30R/sense/internal/dto"
	"github.com/Kartik30R/sense/internal/kafka"
	"github.com/Kartik30R/sense/internal/models"
)

type PacketService struct {
	producer            *kafka.Producer
	deviceService       *DeviceService
	mobileDeviceService *MobileDeviceService
}

func NewPacketService(p *kafka.Producer, d *DeviceService, m *MobileDeviceService) *PacketService {
	return &PacketService{
		producer:            p,
		deviceService:       d,
		mobileDeviceService: m,
	}
}

func (s *PacketService) Ingest(ctx context.Context, req dto.PacketUploadRequest) error {
	msg := models.FromDTO(req)
	_ = s.RegisterDiscovery(req.MobileID, req.DeviceID, req.DeviceAddress, string(msg.ParsedType))
	return s.producer.Publish(ctx, req.DeviceID, msg)
}

// IngestBatch handles a slice of packets
func (s *PacketService) IngestBatch(ctx context.Context, reqs []dto.PacketUploadRequest) error {
	seen := map[string]bool{}

	for _, req := range reqs {
		msg := models.FromDTO(req)

		if !seen[req.DeviceID] {
			seen[req.DeviceID] = true
			_ = s.RegisterDiscovery(req.MobileID, req.DeviceID, req.DeviceAddress, string(msg.ParsedType))
		}

		if err := s.producer.Publish(ctx, req.DeviceID, msg); err != nil {
			return err
		}
	}

	return nil
}

// RegisterDiscovery is a lightweight discovery update for the mobile-sensor affinity.
// Used by both HTTP and WebSocket pipelines to ensure sensors are mapped to all mobiles that see them.
func (s *PacketService) RegisterDiscovery(mobileID, deviceID, address, sensorType string) error {
	_ = s.mobileDeviceService.EnsureMobile(mobileID)
	return s.deviceService.EnsureDevice(
		deviceID,
		mobileID,
		address,
		models.SensorType(sensorType),
	)
}