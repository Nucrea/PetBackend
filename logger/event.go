package logger

import "github.com/rs/zerolog"

type Event interface {
	Err(err error) Event
	Msg(msg string)
	Msgf(format string, v ...interface{})
}

type event struct {
	zerologEvent *zerolog.Event
}

func (e event) Err(err error) Event {
	e.zerologEvent = e.zerologEvent.Err(err)
	return e
}

// Never call it twice: event gets disposed after first call
func (e event) Msg(msg string) {
	e.zerologEvent.Msg(msg)
}

// Never call it twice: event gets disposed after first call
func (e event) Msgf(format string, v ...interface{}) {
	e.zerologEvent.Msgf(format, v...)
}
