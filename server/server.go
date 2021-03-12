package server

import (
	"context"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/IdeaEvolver/cutter-pkg/service"
	"github.com/IdeaEvolver/cutter-status-dashboard/status"
	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
)

type StatusStore interface {
	UpdateStatus(ctx context.Context, service, status string) error
	GetAllStatuses(ctx context.Context) (*status.AllStatuses, error)
	GetStatus(ctx context.Context, service string) (*status.Status, error)
}

type Handler struct {
	Statuses StatusStore
}

func New(cfg *service.Config, handler *Handler) *service.Server {
	router := chi.NewRouter()

	router.Route("/api/v1", func(router chi.Router) {
		router.Route("/", func(router chi.Router) {
			router.Method("GET", "/get-all-statuses", service.JsonHandler(handler.GetAllStatuses))
			router.Method("GET", "/get-status", service.JsonHandler(handler.GetStatus))
		})
	})

	httpHandler := &ochttp.Handler{
		// Use the Google Cloud propagation format.
		Propagation: &propagation.HTTPFormat{},
		Handler:     router,
	}

	return service.GracefulServer(cfg, httpHandler)
}
