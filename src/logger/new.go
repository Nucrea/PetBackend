package logger

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// type OutputStream int

// const (
// 	OutputStreamNone OutputStream = iota
// 	OutputStreamStderr
// 	OutputStreamStdout
// )

type NewLoggerOpts struct {
	Debug      bool
	OutputFile string
	// OutputStream OutputStream
}

func New(ctx context.Context, opts NewLoggerOpts) (Logger, error) {
	// TODO: pass output streams from opts
	writers := []io.Writer{}
	writers = append(writers, os.Stderr)

	if opts.OutputFile != "" {
		file, err := os.OpenFile(opts.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	// TODO: more log levels
	level := zerolog.TraceLevel
	if opts.Debug {
		level = zerolog.DebugLevel
	}

	// TODO: move to wrapper, determine optimal buffer size
	writer := io.MultiWriter(writers...)
	wrapper := newWrapper(writer)
	wrapper.FlushRoutine(ctx)

	l := zerolog.New(wrapper).Level(level).With().Timestamp().Logger()
	return &logger{
		zeroLogger: &l,
	}, nil
}
