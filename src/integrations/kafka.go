package integrations

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	conn *kafka.Conn
}

func (k *Kafka) Connect(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "", "")
	if err != nil {
		return err
	}
	defer conn.Close()

	// w := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
	// 	Topic:    "topic-A",
	// 	Balancer: &kafka.LeastBytes{},
	// }

	return nil
}

func (k *Kafka) SendMessage() {
	k.conn.WriteMessages()

	msg := kafka.Message{
		Topic: "event",
		Key:   []byte("send_email"),
		Value: []byte("value"),
	}
}
