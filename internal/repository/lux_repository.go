package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type LuxRepository struct {
	db *gorm.DB
}

func NewLuxRepository(db *gorm.DB) *LuxRepository {
	return &LuxRepository{db: db}
}

func (r *LuxRepository) CreateBatch(data []models.LuxSensorData) error {

	return r.db.CreateInBatches(data, 500).Error
}
