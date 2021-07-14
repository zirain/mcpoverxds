package main

import (
	"os"

	"github.com/zirain/mcpoverxds/cmd/app"
	"istio.io/pkg/log"
)

func main() {
	rootCmd := app.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
