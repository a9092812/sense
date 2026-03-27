package repository

import "gorm.io/gorm"

type Registry struct {
	Device            *DeviceRepository
	MobileDevice      *MobileDeviceRepository
	SHT40             *SHT40Repository
	Lux               *LuxRepository
	LIS2DH            *LIS2DHRepository
	Soil              *SoilRepository
	Speed             *SpeedRepository
	Ammonia           *AmmoniaRepository
	TempLogger        *TempLoggerRepository
	DataLogger        *DataLoggerRepository
	DeviceLatestState *DeviceLatestStateRepository
	Sen6x             *Sen6xRepository
}

func NewRegistry(db *gorm.DB) *Registry {

	return &Registry{
		Device:            NewDeviceRepository(db),
		MobileDevice:      NewMobileDeviceRepository(db),
		SHT40:             NewSHT40Repository(db),
		Lux:               NewLuxRepository(db),
		LIS2DH:            NewLIS2DHRepository(db),
		Soil:              NewSoilRepository(db),
		Speed:             NewSpeedRepository(db),
		Ammonia:           NewAmmoniaRepository(db),
		TempLogger:        NewTempLoggerRepository(db),
		DataLogger:        NewDataLoggerRepository(db),
		DeviceLatestState: NewDeviceLatestStateRepository(db),
		Sen6x:             NewSen6xRepository(db),
	}
}
