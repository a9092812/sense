package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type DataLoggerRepository struct {
	db *gorm.DB
}

func NewDataLoggerRepository(db *gorm.DB) *DataLoggerRepository {
	return &DataLoggerRepository{db: db}
}

func (r *DataLoggerRepository) CreateBatch(data []models.DataLoggerData) error {

	return r.db.CreateInBatches(data, 500).Error
}
