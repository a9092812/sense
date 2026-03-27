package pipeline

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/repository"
	"github.com/Kartik30R/sense/pkg/logger"
	"go.uber.org/zap"
)

const batchSize = 500

type BatchAggregator struct {
	repos *repository.Registry

	sht40Batch      []models.SHT40Data
	luxBatch        []models.LuxSensorData
	lis2dhBatch     []models.LIS2DHData
	soilBatch       []models.SoilSensorData
	speedBatch      []models.SpeedDistanceData
	ammoniaBatch    []models.AmmoniaSensorData
	tempLoggerBatch []models.TempLoggerData
	dataLoggerBatch []models.DataLoggerData
	sen6xBatch      []models.Sen6xData

	mutex sync.Mutex
}

func NewBatchAggregator(repos *repository.Registry) *BatchAggregator {

	agg := &BatchAggregator{
		repos: repos,
	}

	go agg.flushLoop()

	return agg
}

func (a *BatchAggregator) Add(packet *models.SensorPacketMessage) {

	a.mutex.Lock()
	defer a.mutex.Unlock()

	switch packet.ParsedType {

	case "SHT40":
		data := models.SHT40Data{
			DeviceID:    packet.DeviceID,
			Time:        packet.GetTime(),
			Temperature: packet.Float("temperature"),
			Humidity:    packet.Float("humidity"),
		}
		a.sht40Batch = append(a.sht40Batch, data)
		if len(a.sht40Batch) >= batchSize {
			a.flushSHT40()
		}

	case "LuxSensor":
		data := models.LuxSensorData{
			DeviceID: packet.DeviceID,
			Time:     packet.GetTime(),
			Lux:      packet.Float("lux"),
			RawData:  packet.String("rawData"),
		}
		a.luxBatch = append(a.luxBatch, data)
		if len(a.luxBatch) >= batchSize {
			a.flushLux()
		}

	case "LIS2DH":
		data := models.LIS2DHData{
			DeviceID: packet.DeviceID,
			Time:     packet.GetTime(),
			X:        packet.Float("x"),
			Y:        packet.Float("y"),
			Z:        packet.Float("z"),
		}
		a.lis2dhBatch = append(a.lis2dhBatch, data)
		if len(a.lis2dhBatch) >= batchSize {
			a.flushLIS2DH()
		}

	case "SoilSensor":
		data := models.SoilSensorData{
			DeviceID:    packet.DeviceID,
			Time:        packet.GetTime(),
			Nitrogen:    packet.Float("nitrogen"),
			Phosphorus:  packet.Float("phosphorus"),
			Potassium:   packet.Float("potassium"),
			Moisture:    packet.Float("moisture"),
			Temperature: packet.Float("temperature"),
			EC:          packet.Float("ec"),
			PH:          packet.Float("pH"),
			Salinity:    packet.Float("salinity"),
		}
		a.soilBatch = append(a.soilBatch, data)
		if len(a.soilBatch) >= batchSize {
			a.flushSoil()
		}

	case "SpeedDistance":
		data := models.SpeedDistanceData{
			DeviceID: packet.DeviceID,
			Time:     packet.GetTime(),
			Speed:    packet.Float("speed"),
			Distance: packet.Float("distance"),
		}
		a.speedBatch = append(a.speedBatch, data)
		if len(a.speedBatch) >= batchSize {
			a.flushSpeed()
		}

	case "AmmoniaSensor":
		ammoniaStr := packet.String("ammonia")
		ammoniaStr = strings.TrimSpace(strings.ReplaceAll(ammoniaStr, "ppm", ""))
		ammoniaVal, _ := strconv.ParseFloat(ammoniaStr, 64)

		data := models.AmmoniaSensorData{
			DeviceID: packet.DeviceID,
			Time:     packet.GetTime(),
			Ammonia:  ammoniaVal,
			RawData:  packet.String("rawData"),
		}
		a.ammoniaBatch = append(a.ammoniaBatch, data)
		if len(a.ammoniaBatch) >= batchSize {
			a.flushAmmonia()
		}

	case "TempLogger":
		data := models.TempLoggerData{
			DeviceID:       packet.DeviceID,
			Time:           packet.GetTime(),
			Temperature:    packet.Float("temperature"),
			Humidity:       packet.Float("humidity"),
			RawTemperature: packet.Int("rawTemperature"),
			RawHumidity:    packet.Int("rawHumidity"),
			RawData:        packet.String("rawData"),
			DeviceAddress:  packet.String("deviceAddress"),
		}
		a.tempLoggerBatch = append(a.tempLoggerBatch, data)
		if len(a.tempLoggerBatch) >= batchSize {
			a.flushTempLogger()
		}

	case "DataLogger":
		data := models.DataLoggerData{
			DeviceID:        packet.DeviceID,
			Time:            packet.GetTime(),
			CurrentPacketID: packet.Int("currentPacketId"),
			LastPacketID:    packet.Int("lastPacketId"),
			AccelData:       packet.JSON("payloadAccel"),
			RawData:         packet.String("rawData"),
		}
		a.dataLoggerBatch = append(a.dataLoggerBatch, data)
		if len(a.dataLoggerBatch) >= batchSize {
			a.flushDataLogger()
		}
	case "SEN6x":
		data := models.Sen6xData{
			DeviceID:    packet.DeviceID,
			Time:        packet.GetTime(),
			PM1:         packet.Float("pm1"),
			PM25:        packet.Float("pm25"),
			PM40:        packet.Float("pm4"),
			PM100:       packet.Float("pm10"),
			Temperature: packet.Float("temperature"),
			Humidity:    packet.Float("humidity"),
			CO2:         packet.Float("co2"),
			VOC:         packet.Float("voc"),
			NOx:         packet.Float("nox"),
		}
		a.sen6xBatch = append(a.sen6xBatch, data)
		if len(a.sen6xBatch) >= batchSize {
			a.flushSen6x()
		}
	}
}

func (a *BatchAggregator) flushLoop() {

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {

		a.mutex.Lock()

		a.flushSHT40()
		a.flushLux()
		a.flushLIS2DH()
		a.flushSoil()
		a.flushSpeed()
		a.flushAmmonia()
		a.flushTempLogger()
		a.flushDataLogger()
		a.flushSen6x()

		a.mutex.Unlock()
	}
}

func (a *BatchAggregator) flushSHT40() {
	if len(a.sht40Batch) == 0 {
		return
	}
	if err := a.repos.SHT40.CreateBatch(a.sht40Batch); err != nil {
		logger.Error("Failed to flush SHT40 batch", zap.Error(err))
	}
	a.sht40Batch = nil
}

func (a *BatchAggregator) flushLux() {
	if len(a.luxBatch) == 0 {
		return
	}
	if err := a.repos.Lux.CreateBatch(a.luxBatch); err != nil {
		logger.Error("Failed to flush Lux batch", zap.Error(err))
	}
	a.luxBatch = nil
}

func (a *BatchAggregator) flushLIS2DH() {
	if len(a.lis2dhBatch) == 0 {
		return
	}
	if err := a.repos.LIS2DH.CreateBatch(a.lis2dhBatch); err != nil {
		logger.Error("Failed to flush LIS2DH batch", zap.Error(err))
	} else {
		logger.Debug("Flushed LIS2DH batch", zap.Int("count", len(a.lis2dhBatch)))
	}
	a.lis2dhBatch = nil
}

func (a *BatchAggregator) flushSoil() {
	if len(a.soilBatch) == 0 {
		return
	}
	if err := a.repos.Soil.CreateBatch(a.soilBatch); err != nil {
		logger.Error("Failed to flush Soil batch", zap.Error(err))
	}
	a.soilBatch = nil
}

func (a *BatchAggregator) flushSpeed() {
	if len(a.speedBatch) == 0 {
		return
	}
	if err := a.repos.Speed.CreateBatch(a.speedBatch); err != nil {
		logger.Error("Failed to flush Speed batch", zap.Error(err))
	}
	a.speedBatch = nil
}

func (a *BatchAggregator) flushAmmonia() {
	if len(a.ammoniaBatch) == 0 {
		return
	}
	if err := a.repos.Ammonia.CreateBatch(a.ammoniaBatch); err != nil {
		logger.Error("Failed to flush Ammonia batch", zap.Error(err))
	}
	a.ammoniaBatch = nil
}

func (a *BatchAggregator) flushTempLogger() {
	if len(a.tempLoggerBatch) == 0 {
		return
	}
	if err := a.repos.TempLogger.CreateBatch(a.tempLoggerBatch); err != nil {
		logger.Error("Failed to flush TempLogger batch", zap.Error(err))
	}
	a.tempLoggerBatch = nil
}

func (a *BatchAggregator) flushDataLogger() {
	if len(a.dataLoggerBatch) == 0 {
		return
	}
	if err := a.repos.DataLogger.CreateBatch(a.dataLoggerBatch); err != nil {
		logger.Error("Failed to flush DataLogger batch", zap.Error(err))
	}
	a.dataLoggerBatch = nil
}

func (a *BatchAggregator) flushSen6x() {
	if len(a.sen6xBatch) == 0 {
		return
	}
	if err := a.repos.Sen6x.CreateBatch(a.sen6xBatch); err != nil {
		logger.Error("Failed to flush SEN6x batch", zap.Error(err))
	}
	a.sen6xBatch = nil
}
