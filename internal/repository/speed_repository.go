package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type SpeedRepository struct {
	db *gorm.DB
}

func NewSpeedRepository(db *gorm.DB) *SpeedRepository {
	return &SpeedRepository{db: db}
}

func (r *SpeedRepository) CreateBatch(data []models.SpeedDistanceData) error {

	return r.db.CreateInBatches(data, 500).Error
}
