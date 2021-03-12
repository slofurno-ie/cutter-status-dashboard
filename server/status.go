package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/IdeaEvolver/cutter-pkg/clog"
	"github.com/IdeaEvolver/cutter-pkg/cuterr"
)

type StatusRequest struct {
	Service string `json:"service"`
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req := &StatusRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, cuterr.New(cuterr.BadRequest, "", err)
	}

	return h.Statuses.GetStatus(r.Context(), req.Service)
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
