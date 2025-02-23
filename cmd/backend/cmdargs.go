package main

import (
	"github.com/akamensky/argparse"
)

type CmdArgs interface {
	GetProfilePath() string
	GetConfigPath() string
	GetLogPath() string
	GetSigningKeyPath() string
}

func CmdArgsParse(osArgs []string) (CmdArgs, error) {
	parser := argparse.NewParser("backend", "runs backend")

	s := parser.String("c", "config", &argparse.Options{Required: true, Help: "Path to a config file"})
	k := parser.String("k", "key", &argparse.Options{Required: false, Default: "", Help: "Path to a jwt signing key"})
	l := parser.String("o", "log", &argparse.Options{Required: false, Default: "", Help: "Path to a log file"})
	p := parser.String("p", "profile", &argparse.Options{Required: false, Default: "", Help: "Path to a cpu profile file"})

	err := parser.Parse(osArgs)
	if err != nil {
		return nil, err
	}

	return &args{
		ConfigPath:     *s,
		LogPath:        *l,
		ProfilePath:    *p,
		SigningKeyPath: *k,
	}, nil
}

type args struct {
	ProfilePath    string
	ConfigPath     string
	LogPath        string
	SigningKeyPath string
}

func (a *args) GetConfigPath() string {
	return a.ConfigPath
}

func (a *args) GetLogPath() string {
	return a.LogPath
}

func (a *args) GetProfilePath() string {
	return a.ProfilePath
}

func (a *args) GetSigningKeyPath() string {
	return a.SigningKeyPath
}
