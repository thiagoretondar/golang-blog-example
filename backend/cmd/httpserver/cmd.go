package httpserver

import (
	"fmt"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/environment"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/logger"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/logger/zaplog"

	"github.com/spf13/cobra"
)

// Configuration contains the data structure for the environment configuration.
type Configuration struct {
	AppName string

	EnvironmentName string

	LogLevel string

	HealthCheckEndpoint string

	Server struct {
		HTTP struct {
			Network    string
			ListenAddr string
		}
	}
}

// HTTPServerCMD configures an HTTP Server with all dependencies necessary (connections, cache, ...)
var HTTPServerCMD = &cobra.Command{
	Use:   "httpserver",
	Short: "Starts application HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// get environment flag param - panic if any error
		envFlag, err := cmd.Flags().GetString("environment")
		if err != nil {
			panic(err)
		}

		// get environment configuration - panic if any error
		var envconfig = &Configuration{}
		err = environment.NewFromYAML("configs/environment", envFlag, envconfig)
		if err != nil {
			panic(fmt.Errorf("failed to load environment config: %s", err))
		}

		// forces EnvironmentName to be always equal to flag received
		envconfig.EnvironmentName = envFlag

		// configure logger (Zap Logger) - panic if any error
		zaplog, err := zaplog.NewCustomZap(logger.Config{
			LogLevel:   envconfig.LogLevel,
			AppName:    envconfig.AppName,
			Production: envFlag == "production",
		})
		if err != nil {
			panic(fmt.Errorf("failed to configure zaplog logger: %s", err))
		}

		// execute HTTP Server
		RunHTTPServer(ctx, zaplog, envconfig)
	},
}
