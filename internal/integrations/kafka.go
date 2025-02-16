package integrations

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka(addr, topic string) *Kafka {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		Balancer:               &kafka.RoundRobin{},
		AllowAutoTopicCreation: false,
		BatchSize:              100,
		BatchTimeout:           100 * time.Millisecond,
	}

	return &Kafka{
		writer: w,
	}
}

func (k *Kafka) SendMessage(ctx context.Context, key string, value []byte) error {
	return k.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
}
