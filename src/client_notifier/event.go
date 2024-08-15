package client_notifier

type Event struct {
	Type EventType
	Data []byte
}

type EventType string

const (
	EventTypeEmailConfirmed EventType = "event_email_confirmed"
)
