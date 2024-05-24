package main

import (
	"context"
	"log"

	_ "github.com/whosonfirst/go-whosonfirst-spelunker-opensearch"
	"github.com/whosonfirst/go-whosonfirst-spelunker/app/cli"
)

func main() {

	ctx := context.Background()
	err := cli.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
