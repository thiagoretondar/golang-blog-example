package main

import (
	"context"
	"github.com/thiagoretondar/golang-blog-example/backend/cmd/httpserver"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "golang-blog-backend",
	Short: "golang-blog-backend control a backend for blog post",
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// flags for "httpserver" command
	httpserver.HTTPServerCMD.Flags().String("environment", "", "Define environment")
	httpserver.HTTPServerCMD.MarkFlagRequired("environment")
	rootCMD.AddCommand(httpserver.HTTPServerCMD)

	err := rootCMD.ExecuteContext(ctx)
	if err != nil {
		panic(err)
	}
}
