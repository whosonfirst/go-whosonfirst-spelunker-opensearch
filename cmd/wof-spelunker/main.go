package main

import (
	"context"
	"log"
	"log/slog"
	
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-opensearch"
	"github.com/whosonfirst/go-whosonfirst-spelunker/app/cli"
)

func main() {

	ctx := context.Background()
	err := cli.Run(ctx, slog.Default())

	if err != nil {
		log.Fatal(err)
	}
}
