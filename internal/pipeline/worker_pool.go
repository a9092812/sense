package pipeline

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/repository"
	"github.com/Kartik30R/sense/pkg/logger"
)

// WorkerPool consumes decoded packets from Kafka and persists them to the database.
// It is the "source of truth" writer — no live/Redis concern here.
// Live dashboard updates are handled upstream by IngestHub via WebSocket.
type WorkerPool struct {
	workerCount int
	jobChannel  *JobChannel
	aggregator  *BatchAggregator
	repos       *repository.Registry
}

func NewWorkerPool(
	workerCount int,
	jobChannel *JobChannel,
	aggregator *BatchAggregator,
	repos *repository.Registry,
) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobChannel:  jobChannel,
		aggregator:  aggregator,
		repos:       repos,
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	logger.Info("Starting worker pool", zap.Int("workers", w.workerCount))

	var wg sync.WaitGroup

	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go w.startWorker(ctx, &wg, i)
	}

	go func() {
		wg.Wait()
		logger.Info("Worker pool stopped")
	}()
}

func (w *WorkerPool) startWorker(ctx context.Context, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	logger.Debug("Worker started", zap.Int("workerID", workerID))

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-w.jobChannel.Channel():
			if !ok {
				logger.Debug("Worker job channel closed", zap.Int("workerID", workerID))
				return
			}

			packet, err := models.Decode(job)
			if err != nil {
				logger.Error("packet decode error", zap.Error(err))
				continue
			}

			packet.ParsedType = strings.TrimSpace(packet.ParsedType)
			packet.DeviceID = strings.TrimSpace(packet.DeviceID)

			// 🔍 Log processing
			logger.Info("Worker processing packet",
				zap.Int("workerID", workerID),
				zap.String("deviceId", packet.DeviceID),
				zap.String("type", packet.ParsedType),
			)

			// Ensure device exists (Worker side discovery)
			err = w.repos.Device.EnsureDeviceExists(
				packet.DeviceID,
				packet.MobileID,
				packet.DeviceAddress,
				packet.ParsedType,
			)
			if err != nil {
				logger.Warn("Failed to ensure device exists", zap.String("deviceId", packet.DeviceID), zap.Error(err))
			}

			w.aggregator.Add(packet)

			// Update latest device state snapshot in DB
			w.updateLatestState(packet)
		}
	}
}

func (w *WorkerPool) updateLatestState(packet *models.SensorPacketMessage) {
	state := &models.DeviceLatestState{
		DeviceID: packet.DeviceID,
		LastSeen: packet.GetTime(),
	}

	switch packet.ParsedType {
	case "SHT40":
		state.Temperature = packet.Float("temperature")
		state.Humidity = packet.Float("humidity")
	case "LuxSensor":
		state.Lux = packet.Float("lux")
	case "TempLogger":
		state.Temperature = packet.Float("temperature")
		state.Humidity = packet.Float("humidity")
	case "LIS2DH":
		state.X = packet.Float("x")
		state.Y = packet.Float("y")
		state.Z = packet.Float("z")
	case "SoilSensor":
		state.Nitrogen = packet.Float("nitrogen")
		state.Phosphorus = packet.Float("phosphorus")
		state.Potassium = packet.Float("potassium")
		state.Moisture = packet.Float("moisture")
		state.Temperature = packet.Float("temperature")
		state.EC = packet.Float("ec")
		state.PH = packet.Float("pH")
		state.Salinity = packet.Float("salinity")
	case "SpeedDistance":
		state.Speed = packet.Float("speed")
		state.Distance = packet.Float("distance")
	case "AmmoniaSensor":
		ammoniaStr := packet.String("ammonia")
		ammoniaStr = strings.TrimSpace(strings.ReplaceAll(ammoniaStr, "ppm", ""))
		if val, err := strconv.ParseFloat(ammoniaStr, 64); err == nil {
			state.Ammonia = val
		}
	case "DataLogger":
		state.LastPacketID = packet.Int("currentPacketId")
	}

	w.repos.DeviceLatestState.Upsert(state)
}
