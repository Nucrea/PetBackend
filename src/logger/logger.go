package logger

import (
	"backend/src/request_context"
	"context"

	"github.com/rs/zerolog"
)

type Logger interface {
	Log() Event
	Warning() Event
	Error() Event
	Fatal() Event

	Printf(format string, v ...any)
}

type logger struct {
	prefix     string
	requestCtx request_context.RequestContext
	zeroLogger *zerolog.Logger
}

func (l logger) Log() Event {
	return l.wrapEvent(l.zeroLogger.Log())
}

func (l logger) Warning() Event {
	return l.wrapEvent(l.zeroLogger.Warn())
}

func (l logger) Error() Event {
	return l.wrapEvent(l.zeroLogger.Error())
}

func (l logger) Fatal() Event {
	return l.wrapEvent(l.zeroLogger.Fatal())
}

func (l logger) Printf(format string, v ...any) {
	l.zeroLogger.Printf(format, v...)
}

func (l logger) wrapEvent(zerologEvent *zerolog.Event) Event {
	var e Event = event{zerologEvent}

	if l.requestCtx != nil {
		e = e.Str("requestId", l.requestCtx.RequestId())
		e = e.Str("userId", l.requestCtx.UserId())
		if l.prefix != "" {
			e = e.Str("prefix", l.prefix)
		}
	}

	return e
}

func (l logger) WithContext(ctx context.Context) Logger {
	if rctx, ok := ctx.(request_context.RequestContext); ok {
		return logger{
			prefix:     l.prefix,
			requestCtx: rctx,
			zeroLogger: l.zeroLogger,
		}
	}

	return l
}

func (l logger) WithPrefix(prefix string) Logger {
	return logger{
		prefix:     prefix,
		requestCtx: l.requestCtx,
		zeroLogger: l.zeroLogger,
	}
}
