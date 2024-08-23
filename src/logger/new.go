package logger

import (
	"bufio"
	"context"
	"io"
	"os"
	"sync"
	"time"

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

	writer := bufio.NewWriterSize(io.MultiWriter(writers...), 8*1024)
	wrapper := &bufioWrapper{writer, &sync.RWMutex{}}
	go func() {
		tmr := time.NewTicker(500 * time.Millisecond)
		defer tmr.Stop()

		for {
			wrapper.Flush()

			select {
			case <-context.Background().Done():
				wrapper.Flush()
				return
			case <-tmr.C:
			}
		}
	}()

	l := zerolog.New(wrapper).Level(level).With().Timestamp().Logger()
	return &logger{
		zeroLogger: &l,
	}, nil
}
