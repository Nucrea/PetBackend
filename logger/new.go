package logger

import (
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

func New(opts NewLoggerOpts) (Logger, error) {
	writers := []io.Writer{}
	writers = append(writers, os.Stderr)

	if opts.OutputFile != "" {
		file, err := os.OpenFile(opts.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	level := zerolog.TraceLevel
	if opts.Debug {
		level = zerolog.DebugLevel
	}

	writer := io.MultiWriter(writers...)

	l := zerolog.New(writer).Level(level)
	return &logger{
		zeroLogger: &l,
	}, nil
}
