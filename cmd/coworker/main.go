package main

import (
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/segmentio/kafka-go"
)

type SendEmailEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func main() {
	ctx := context.Background()

	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}

	emailer, err := NewEmailer(config.SMTP)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.ConsumerGroupId,
	})

	logger, err := logger.New(ctx, logger.NewLoggerOpts{
		Debug:      true,
		OutputFile: config.App.LogFile,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Printf("coworker service started\n")

	for {
		msg, err := r.FetchMessage(ctx)
		if err == io.EOF {
			log.Fatal("EOF")
			return
		}
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		log.Printf("offset: %d, partition: %d, key: %s, value: %s\n", msg.Offset, msg.Partition, string(msg.Key), string(msg.Value))

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Fatalf("failed to commit: %s\n", err.Error())
			continue
		}

		if err := handleEvent(config, emailer, msg); err != nil {
			log.Printf("failed to handle event: %s\n", err.Error())
			continue
		}
	}
}

func handleEvent(config Config, emailer *Emailer, msg kafka.Message) error {
	event := SendEmailEvent{}
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	switch string(msg.Key) {
	case "email_forgot_password":
		return emailer.SendRestorePassword(event.Email, event.Token)
	case "email_password_changed":
		return emailer.SendPasswordChanged(event.Email)
	case "email_verify_user":
		link := fmt.Sprintf("%s/verify-user?token=%s", config.App.ServiceUrl, event.Token)
		return emailer.SendVerifyUser(event.Email, link)
	}

	return fmt.Errorf("unknown event type")
}
