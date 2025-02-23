package main

import (
	"backend/internal/integrations"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
)

type SendEmailEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewEventHandler(
	config Config,
	logger logger.Logger,
	metrics *integrations.Metrics,
	emailer *Emailer,
) *EventHandler {
	eventsCounter := metrics.NewCounter("events_counter", "total events handled")
	return &EventHandler{
		config:        config,
		logger:        logger,
		emailer:       emailer,
		eventsCounter: eventsCounter,
	}
}

type EventHandler struct {
	config        Config
	logger        logger.Logger
	emailer       *Emailer
	eventsCounter integrations.Counter
}

func (e *EventHandler) eventLoop(ctx context.Context, kafkaReader *kafka.Reader) {
	for {
		msg, err := kafkaReader.FetchMessage(ctx)
		if err == io.EOF {
			e.logger.Fatal().Err(err)
		}
		if err != nil {
			e.logger.Fatal().Err(err)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		e.logger.Log().Msgf("event: %s", msg.Key)
		e.eventsCounter.Inc()

		if err := kafkaReader.CommitMessages(ctx, msg); err != nil {
			e.logger.Error().Err(err).Msg("failed to commit offset")
			continue
		}

		if err := e.handleEvent(ctx, msg); err != nil {
			e.logger.Error().Err(err).Msg("failed to handle event")
			continue
		}
	}
}

func (e *EventHandler) handleEvent(ctx context.Context, msg kafka.Message) error {
	event := SendEmailEvent{}
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	// TODO: add context somehow
	switch string(msg.Key) {
	case "email_forgot_password":
		return e.emailer.SendRestorePassword(event.Email, event.Token)
	case "email_password_changed":
		return e.emailer.SendPasswordChanged(event.Email)
	case "email_verify_user":
		link := fmt.Sprintf("%s/verify-user?token=%s", e.config.App.ServiceUrl, event.Token)
		return e.emailer.SendVerifyUser(event.Email, link)
	}

	return fmt.Errorf("unknown event type")
}
