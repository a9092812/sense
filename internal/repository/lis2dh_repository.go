package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type LIS2DHRepository struct {
	db *gorm.DB
}

func NewLIS2DHRepository(db *gorm.DB) *LIS2DHRepository {
	return &LIS2DHRepository{db: db}
}

func (r *LIS2DHRepository) CreateBatch(data []models.LIS2DHData) error {

	return r.db.CreateInBatches(data, 500).Error
}
