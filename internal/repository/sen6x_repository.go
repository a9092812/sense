package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type Sen6xRepository struct {
	db *gorm.DB
}

func NewSen6xRepository(db *gorm.DB) *Sen6xRepository {
	return &Sen6xRepository{db: db}
}

func (r *Sen6xRepository) CreateBatch(data []models.Sen6xData) error {
	return r.db.Create(&data).Error
}
