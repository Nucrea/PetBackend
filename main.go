package main

import (
	"backend/src"
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

	app := &src.App{}
	app.Run(
		src.RunParams{
			Ctx:     context.Background(),
			OsArgs:  os.Args,
			EnvVars: env,
		},
	)
}
