package repository

import (
	"github.com/Kartik30R/sense/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceLatestStateRepository struct {
	db *gorm.DB
}

func NewDeviceLatestStateRepository(db *gorm.DB) *DeviceLatestStateRepository {
	return &DeviceLatestStateRepository{db: db}
}

func (r *DeviceLatestStateRepository) Upsert(state *models.DeviceLatestState) error {

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		UpdateAll: true,
	}).Create(state).Error
}
