package server

import (
	"context"
	"net/http"
	"time"

	"github.com/IdeaEvolver/cutter-pkg/clog"
)

type StatusRequest struct {
	Service string `json:"service"`
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	service := r.URL.Query().Get("service")

	return h.Statuses.GetStatus(r.Context(), service)
}

func (h *Handler) GetAllStatuses(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h.Statuses.GetAllStatuses(r.Context())
}

func (h *Handler) AllChecks(ctx context.Context) error {
	for {

		platformStatus, err := h.Healthchecks.PlatformStatus(ctx)
		if err != nil {
			clog.Fatalf("Error retrieving platform status", err)
		}

		if err := h.Statuses.UpdateStatus(ctx, "platform", platformStatus.Status); err != nil {
			clog.Fatalf("Error updating platform status", err)
		}

		time.Sleep(60 * time.Second)
	}
}
