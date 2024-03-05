package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/whosonfirst/go-whosonfirst-spelunker-opensearch"
	"github.com/whosonfirst/go-whosonfirst-spelunker/app/cli"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := cli.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run spelunker application", "error", err)
		os.Exit(1)
	}
}
