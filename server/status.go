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

		fulfillmentStatus, err := h.Healthchecks.FulfillmentStatus(ctx)
		if err != nil {
			clog.Fatalf("Error retrieving fulfillment status", err)
		}

		if err := h.Statuses.UpdateStatus(ctx, "fulfillment", fulfillmentStatus.Status); err != nil {
			clog.Fatalf("Error updating fulfillment status", err)
		}

		crmStatus, err := h.Healthchecks.CrmStatus(ctx)
		if err != nil {
			clog.Fatalf("Error retrieving crm status", err)
		}

		if err := h.Statuses.UpdateStatus(ctx, "crm", crmStatus.Status); err != nil {
			clog.Fatalf("Error updating crm status", err)
		}

		studyStatus, err := h.Healthchecks.StudyStatus(ctx)
		if err != nil {
			clog.Fatalf("Error retrieving study status", err)
		}

		if err := h.Statuses.UpdateStatus(ctx, "study", studyStatus.Status); err != nil {
			clog.Fatalf("Error updating study status", err)
		}

		nodeMetrics, err := h.Metrics.GetNodeMetrics(ctx)
		if err != nil {
			clog.Fatalf("Error retrieving node metrics", err)
		}

		infra := "Ok"
		if !nodeMetrics.Healthy() {
			infra = "high utilization"
		}
		if err := h.Statuses.UpdateStatus(ctx, "infrastructure", infra); err != nil {
			clog.Fatalf("Error updating infra status", err)
		}

		time.Sleep(60 * time.Second)
	}
}
