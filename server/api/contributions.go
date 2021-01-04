package api

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type contribution struct {
	ID              string    `json:"id"`
	Amount          int       `json:"amount"`
	Date            time.Time `json:"date"`
	ContributorName string    `json:"contributor_name"`
	ContributorID   string    `json:"contributor_id"`
	RecipientName   string    `json:"recipient_name"`
	RecipientID     string    `json:"recipient_id"`

	// TODO(vicki): optionally include other fields
}

func (s *Server) contributionsReceived(ctx context.Context, individualID string) ([]contribution, error) {
	const q = `
		SELECT
			id,
			amount,
			date,
			contributor_name,
			COALESCE(contributor_id, ''),
			recipient_name,
			recipient_id
		FROM contributions
		WHERE recipient_id = $1
	`
	res, err := s.getContributions(ctx, q, individualID)
	return res, errors.Wrap(err, "contributions received")
}

func (s *Server) contributionsGiven(ctx context.Context, individualID string) ([]contribution, error) {
	const q = `
		SELECT
			id,
			amount,
			date,
			contributor_name,
			contributor_id,
			recipient_name,
			COALESCE(recipient_id, '')
		FROM contributions
		WHERE contributor_id = $1
	`
	res, err := s.getContributions(ctx, q, individualID)
	return res, errors.Wrap(err, "contributions given")
}

func (s *Server) getContributions(ctx context.Context, query string, individualID string) ([]contribution, error) {
	rows, err := s.db.QueryContext(ctx, query, individualID)
	if err != nil {
		return nil, errors.Wrap(err, "querying contributions from db")
	}
	defer rows.Close()

	var res []contribution
	for rows.Next() {
		var c contribution
		err := rows.Scan(&c.ID, &c.Amount, &c.Date, &c.ContributorName, &c.ContributorID, &c.RecipientName, &c.RecipientID)
		if err != nil {
			return nil, errors.Wrap(err, "scanning contribution row")
		}
		res = append(res, c)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading contribution rows")
	}
	return res, nil
}
