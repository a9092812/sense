package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type AmmoniaRepository struct {
	db *gorm.DB
}

func NewAmmoniaRepository(db *gorm.DB) *AmmoniaRepository {
	return &AmmoniaRepository{db: db}
}

func (r *AmmoniaRepository) CreateBatch(data []models.AmmoniaSensorData) error {

	return r.db.CreateInBatches(data, 500).Error
}
