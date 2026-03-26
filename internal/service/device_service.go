package services

import (
	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/repository"
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService(repo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) EnsureDevice(deviceID string, mobileID string, address string, sensorType models.SensorType) error {
	return s.repo.EnsureDeviceExists(deviceID, mobileID, address, string(sensorType))
}

func (s *DeviceService) ListDevices() ([]repository.DeviceRow, error) {
	return s.repo.ListDevices()
}

func (s *DeviceService) ListSensorsByMobile(mobileID string) ([]repository.DeviceRow, error) {
	return s.repo.ListSensorsByMobile(mobileID)
}