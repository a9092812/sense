package services

import (
	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/repository"
)

type MobileDeviceService struct {
	repo *repository.MobileDeviceRepository
}

func NewMobileDeviceService(repo *repository.MobileDeviceRepository) *MobileDeviceService {
	return &MobileDeviceService{repo: repo}
}

func (s *MobileDeviceService) EnsureMobile(mobileID string) error {
	return s.repo.EnsureMobileExists(mobileID)
}

func (s *MobileDeviceService) ListMobiles() ([]models.MobileDevice, error) {
	return s.repo.ListMobiles()
}
