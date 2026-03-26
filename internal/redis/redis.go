package redis

import (
	"context"
	"encoding/json"

	"github.com/Kartik30R/sense/config"
)

const (
	SensorDataChannel = "sensor_data"
)

type Publisher struct{}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(channel string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return config.RedisClient.Publish(context.Background(), channel, data).Err()
}

type Subscriber struct{}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

func (s *Subscriber) Subscribe(ctx context.Context, channel string) <-chan interface{} {
	pubsub := config.RedisClient.Subscribe(ctx, channel)
	ch := make(chan interface{})

	go func() {
		defer pubsub.Close()
		defer close(ch)

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				return
			}

			// Forward raw JSON bytes — avoid double-encoding for WebSocket
			select {
			case <-ctx.Done():
				return
			case ch <- []byte(msg.Payload):
			}
		}
	}()

	return ch
}
