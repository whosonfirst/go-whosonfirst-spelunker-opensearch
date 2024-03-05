package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-opensearch"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}
