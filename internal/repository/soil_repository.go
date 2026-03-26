package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type SoilRepository struct {
	db *gorm.DB
}

func NewSoilRepository(db *gorm.DB) *SoilRepository {
	return &SoilRepository{db: db}
}

func (r *SoilRepository) CreateBatch(data []models.SoilSensorData) error {

	return r.db.CreateInBatches(data, 500).Error
}
