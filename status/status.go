package status

import (
	"context"
	"database/sql"

	"github.com/IdeaEvolver/cutter-pkg/cuterr"
)

type StatusStore struct {
	db *sql.DB

	PlatformHealthcheck    string
	FulfillmentHealthcheck string
	CrmHealthcheck         string
	StudyHealthcheck       string
}

type AllStatuses struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type Status struct {
	Status string `json:"status"`
}

func New(platform, fulfillment, crm, study string, db *sql.DB) *StatusStore {
	return &StatusStore{
		PlatformHealthcheck:    platform,
		FulfillmentHealthcheck: fulfillment,
		CrmHealthcheck:         crm,
		StudyHealthcheck:       study,
		db:                     db,
	}
}

func (s *StatusStore) UpdateStatus(ctx context.Context, service, status string) error {
	var query = `INSERT INTO statuses (service, status)
	VALUES ($1, $2) 
	ON CONFLICT DO UPDATE SET status = ($1)
	WHERE service = $2
	`
	_, err := s.db.ExecContext(ctx, query, service, status)
	if err != nil {
		return cuterr.FromDatabaseError("UpdateStatus", err)
	}

	return nil
}

func (s *StatusStore) GetAllStatuses(ctx context.Context) (*AllStatuses, error) {
	var query = `SELECT * FROM statuses`

	ret := &AllStatuses{}

	err := s.db.QueryRowContext(ctx, query).
		Scan(
			&ret.Service,
			&ret.Status,
		)

	if err != nil {
		return nil, cuterr.FromDatabaseError("GetAllStatuses", err)
	}

	return ret, nil
}

// might be useful to get individual service status

func (s *StatusStore) GetStatus(ctx context.Context, service string) (*Status, error) {
	var query = `SELECT status FROM statuses WHERE service = $1`

	ret := &Status{}
	err := s.db.QueryRowContext(ctx, query, service).
		Scan(
			&ret.Status,
		)

	if err != nil {
		return nil, cuterr.FromDatabaseError("GetStatus", err)
	}

	return ret, nil
}
