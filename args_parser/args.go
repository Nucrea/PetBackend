package args_parser

import (
	"github.com/akamensky/argparse"
)

type Args interface {
	GetConfigPath() string
	GetLogPath() string
}

func Parse(osArgs []string) (Args, error) {
	parser := argparse.NewParser("backend", "runs backend")

	s := parser.String("c", "config", &argparse.Options{Required: true, Help: "Path to a config file"})
	l := parser.String("o", "log", &argparse.Options{Required: false, Default: "", Help: "Path to a log file"})

	err := parser.Parse(osArgs)
	if err != nil {
		return nil, err
	}

	return &args{
		ConfigPath: *s,
		LogPath:    *l,
	}, nil
}

type args struct {
	ConfigPath string
	LogPath    string
}

func (a *args) GetConfigPath() string {
	return a.ConfigPath
}

func (a *args) GetLogPath() string {
	return a.LogPath
}
