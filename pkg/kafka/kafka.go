package kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers      []string
	Topic        string
	GroupID      string
	BatchSize    int
	BatchTimeout time.Duration
	Async        bool
}

func NewWriter(cfg Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		Async:        cfg.Async,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
	}
}

func NewReader(cfg Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  200 * time.Millisecond,
	})
}

func EnsureTopic(broker string, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	// More robust topic creation (handling controllers)
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1, // Set to 1 for dev/minikube, 3 for prod if needed
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}
