package repos

import (
	"backend/internal/integrations"
	"context"
	"encoding/json"
)

func NewEventRepo(kafka *integrations.Kafka) *EventRepo {
	return &EventRepo{
		kafka: kafka,
	}
}

type EventRepo struct {
	kafka *integrations.Kafka
}

func (e *EventRepo) SendEmailForgotPassword(ctx context.Context, email, actionToken string) error {
	value := struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: email,
		Token: actionToken,
	}
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return e.kafka.SendMessage(ctx, "email_forgot_password", valueBytes)
}
