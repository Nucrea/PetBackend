package logger

import "github.com/rs/zerolog"

type Logger interface {
	Log() Event
	Warning() Event
	Error() Event
	Fatal() Event
}

type logger struct {
	zeroLogger *zerolog.Logger
}

func (l logger) Log() Event {
	return event{l.zeroLogger.Log()}
}

func (l *logger) Warning() Event {
	return event{l.zeroLogger.Warn()}
}

func (l *logger) Error() Event {
	return event{l.zeroLogger.Error()}
}

func (l *logger) Fatal() Event {
	return event{l.zeroLogger.Fatal()}
}
