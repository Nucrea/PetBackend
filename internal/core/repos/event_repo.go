package repos

import (
	"backend/internal/integrations"
	"context"
	"encoding/json"
)

const (
	EventEmailPasswordChanged = "email_password_changed"
	EventEmailForgotPassword  = "email_forgot_password"
	EventEmailVerifyUser      = "email_verify_user"
)

func NewEventRepo(kafka *integrations.Kafka) *EventRepo {
	return &EventRepo{
		kafka: kafka,
	}
}

type EventRepo struct {
	kafka *integrations.Kafka
}

func (e *EventRepo) sendEmail(ctx context.Context, email, actionToken, eventType string) error {
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

	return e.kafka.SendMessage(ctx, eventType, valueBytes)
}

func (e *EventRepo) SendEmailPasswordChanged(ctx context.Context, email string) error {
	return e.sendEmail(ctx, email, "", EventEmailPasswordChanged)
}

func (e *EventRepo) SendEmailForgotPassword(ctx context.Context, email, actionToken string) error {
	return e.sendEmail(ctx, email, actionToken, EventEmailForgotPassword)
}

func (e *EventRepo) SendEmailVerifyUser(ctx context.Context, email, actionToken string) error {
	return e.sendEmail(ctx, email, actionToken, EventEmailVerifyUser)
}
