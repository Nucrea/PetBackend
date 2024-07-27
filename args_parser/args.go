package args_parser

import (
	"github.com/akamensky/argparse"
)

type Args interface {
	GetConfigPath() string
}

func Parse(osArgs []string) (Args, error) {
	parser := argparse.NewParser("backend", "runs backend")

	s := parser.String("c", "config", &argparse.Options{Required: true, Help: "Path to a config file"})

	err := parser.Parse(osArgs)
	if err != nil {
		return nil, err
	}

	return &args{
		ConfigPath: *s,
	}, nil
}

type args struct {
	ConfigPath string
}

func (a *args) GetConfigPath() string {
	return a.ConfigPath
}
