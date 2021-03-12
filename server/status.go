package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
