package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Kartik30R/sense/internal/models"
)

type HistoryService struct {
	db *gorm.DB
}

func NewHistoryService(db *gorm.DB) *HistoryService {
	return &HistoryService{db: db}
}

func (s *HistoryService) GetDeviceHistory(deviceID string, sensor models.SensorType, startTime, endTime time.Time, limit int) (interface{}, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000 // capping max items
	}

	query := s.db.Where("LOWER(device_id) = LOWER(?)", deviceID)

	if !startTime.IsZero() {
		query = query.Where("time >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("time <= ?", endTime)
	}

	query = query.Order("time desc").Limit(limit)

	switch sensor {
	case models.SensorSHT40:
		var data []models.SHT40Data
		err := query.Find(&data).Error
		return data, err
	case models.SensorLux:
		var data []models.LuxSensorData
		err := query.Find(&data).Error
		return data, err
	case models.SensorLIS2DH:
		var data []models.LIS2DHData
		err := query.Find(&data).Error
		return data, err
	case models.SensorSoil:
		var data []models.SoilSensorData
		err := query.Find(&data).Error
		return data, err
	case models.SensorSpeedDistance:
		var data []models.SpeedDistanceData
		err := query.Find(&data).Error
		return data, err
	case models.SensorAmmonia:
		var data []models.AmmoniaSensorData
		err := query.Find(&data).Error
		return data, err
	case models.SensorTempLogger:
		var data []models.TempLoggerData
		err := query.Find(&data).Error
		return data, err
	case models.SensorDataLogger:
		var data []models.DataLoggerData
		err := query.Find(&data).Error
		return data, err
	default:
		return nil, errors.New("unsupported sensor type")
	}
}
