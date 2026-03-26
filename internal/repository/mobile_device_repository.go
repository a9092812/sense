package repository

import (
	"time"

	"github.com/Kartik30R/sense/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MobileDeviceRepository struct {
	db *gorm.DB
}

func NewMobileDeviceRepository(db *gorm.DB) *MobileDeviceRepository {
	return &MobileDeviceRepository{db: db}
}

func (r *MobileDeviceRepository) EnsureMobileExists(mobileID string) error {
	mobile := models.MobileDevice{
		MobileID: mobileID,
		LastSeen: time.Now(),
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mobile_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_seen"}),
	}).Create(&mobile).Error
}

func (r *MobileDeviceRepository) FindByMobileID(mobileID string) (*models.MobileDevice, error) {
	var mobile models.MobileDevice
	err := r.db.Where("mobile_id = ?", mobileID).First(&mobile).Error
	return &mobile, err
}

func (r *MobileDeviceRepository) ListMobiles() ([]models.MobileDevice, error) {
	var mobiles []models.MobileDevice
	err := r.db.Find(&mobiles).Error
	return mobiles, err
}
