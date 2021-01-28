package httpserver

import (
	"context"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/logger/zaplog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func RunHTTPServer(ctx context.Context, zaplog zaplog.Logger, envconfig *Configuration) {
	// configure HTTP Routes
	routesHandler := newRouterHandler(envconfig, zaplog)

	// create HTTP Server with handler from router
	httpServer := &http.Server{
		Addr:    envconfig.Server.HTTP.ListenAddr,
		Handler: routesHandler,
	}

	// signal sent from terminal (Ctrl+C = os.Interrupt) or Kubernetes Signal (SIGTERM)
	interruptServerSignal := make(chan os.Signal, 1)
	signal.Notify(interruptServerSignal, os.Interrupt, syscall.SIGTERM)

	// create channel to do graceful shutdown
	done := make(chan bool, 1)
	go gracefulShutdown(httpServer, zaplog, interruptServerSignal, done)

	zaplog.Info("HTTP Server is ready to handle request",
		zap.String("listenAddr", httpServer.Addr),
		zap.String("environment", envconfig.EnvironmentName))
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zaplog.Fatal("Couldn't listen on selected port",
			zap.String("listenAddr", envconfig.Server.HTTP.ListenAddr), zap.Error(err),
		)
	}

	<-done
	zaplog.Warn("HTTP Server stopped")
}

func gracefulShutdown(server *http.Server, zaplog zaplog.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	//logger.Info("HTTP Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		zaplog.Fatal("Could not gracefully shutdown the HTTP server", zap.Error(err))
	}

	close(done)
}
