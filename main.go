package main

import (
	"context"
	"log"

	"github.com/ameykpatil/github-data-analyzer/cmd"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
