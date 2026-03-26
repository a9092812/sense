package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type TempLoggerRepository struct {
	db *gorm.DB
}

func NewTempLoggerRepository(db *gorm.DB) *TempLoggerRepository {
	return &TempLoggerRepository{db: db}
}

func (r *TempLoggerRepository) CreateBatch(data []models.TempLoggerData) error {

	return r.db.CreateInBatches(data, 500).Error
}
