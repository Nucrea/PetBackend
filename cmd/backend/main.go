package main

import (
	"context"
	"os"
	"strings"
)

func main() {
	env := map[string]string{}
	for _, envEntry := range os.Environ() {
		kv := strings.SplitN(envEntry, "=", 2)
		env[kv[0]] = kv[1]
	}

	app := &App{}
	app.Run(
		RunParams{
			Ctx:     context.Background(),
			OsArgs:  os.Args,
			EnvVars: env,
		},
	)
}
