package httpserver

import (
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/logger/zaplog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// newRouter creates the main HTTP router for this application and some middlewares
func newRouterHandler(envconfig *Configuration, logger zaplog.Logger) http.Handler {
	r := chi.NewRouter()

	// configure middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat(envconfig.HealthCheckEndpoint))
	r.Use(middleware.StripSlashes)

	// configure routes
	configureChiRoutes(r, logger)

	return r
}

func configureChiRoutes(handler *chi.Mux, logger zaplog.Logger) {
	//pageService := page.NewService(logger)
	//webhookSvc := webhook.NewService(logger)

	// routes
	//handler.Mount("/pages", page.NewHandler(pageService))
	//handler.Mount("/webhook", webhook.NewHandler(webhookSvc))
}
