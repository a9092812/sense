package repository

import (
	"github.com/Kartik30R/sense/internal/models"
	"gorm.io/gorm"
)

type SHT40Repository struct {
	db *gorm.DB
}

func NewSHT40Repository(db *gorm.DB) *SHT40Repository {
	return &SHT40Repository{db: db}
}

func (r *SHT40Repository) CreateBatch(data []models.SHT40Data) error {

	return r.db.CreateInBatches(data, 500).Error
}
