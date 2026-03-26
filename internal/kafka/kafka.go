package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Kartik30R/sense/config"
	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/pipeline"
	pkgkafka "github.com/Kartik30R/sense/pkg/kafka"
	"github.com/Kartik30R/sense/pkg/logger"
)

// Consts
const (
	SensorPacketsTopic = "sensor_packets"
)

// Producer wraps pkgkafka.Writer with application-specific logic
type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg config.KafkaConfig) *Producer {
	return &Producer{
		writer: pkgkafka.NewWriter(pkgkafka.Config{
			Brokers:      cfg.Brokers,
			Topic:        cfg.Topic,
			Async:        cfg.Async,
			BatchSize:    cfg.BatchSize,
			BatchTimeout: cfg.BatchTimeout,
		}),
	}
}

func (p *Producer) Publish(ctx context.Context, deviceID string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(deviceID),
		Value: data,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer wraps pkgkafka.Reader with application-specific processing logic
type Consumer struct {
	reader  *kafka.Reader
	limiter *pipeline.RateLimiter
}

func NewConsumer(cfg config.KafkaConfig, rateCfg config.RateLimitConfig) *Consumer {
	reader := pkgkafka.NewReader(pkgkafka.Config{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	return &Consumer{
		reader:  reader,
		limiter: pipeline.NewRateLimiter(rateCfg),
	}
}

func (c *Consumer) Start(ctx context.Context, jobChan chan<- []byte) {
	logger.Info("Kafka consumer started", zap.String("topic", c.reader.Config().Topic))

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Error("Kafka read error", zap.Error(err))
			continue
		}

		packet, err := models.Decode(msg.Value)
		if err != nil {
			logger.Warn("Failed to decode packet", zap.Error(err))
			continue
		}

		if !c.limiter.Allow(packet.DeviceID) {
			continue // Drop noisy packets
		}

		select {
		case <-ctx.Done():
			return
		case jobChan <- msg.Value:
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

// EnsureTopics is a convenience wrapper around pkgkafka.EnsureTopic
func EnsureTopics(broker string, topic string, partitions int) error {
	return pkgkafka.EnsureTopic(broker, topic, partitions)
}
