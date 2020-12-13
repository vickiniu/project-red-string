package api

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type individual struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	ZIP       string    `json:"zip"`
	UpdatedTS time.Time `json:"updated_ts"`
	Role      string    `json:"role"`
	Title     string    `json:"title"`
	Twitter   string    `json:"twitter"`

	Associations []association `json:"associations"`
}

type association struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (s *Server) getIndividual(ctx context.Context, id string) (*individual, error) {
	i := &individual{}
	const q = `
		SELECT
			id,
			first_name,
			last_name,
			COALESCE(zip, ''),
			updated_ts,
			COALESCE(role, ''),
			COALESCE(title, ''),
			COALESCE(twitter, '')
		FROM individuals
		WHERE id = $1
	`
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&i.ID, &i.FirstName, &i.LastName, &i.ZIP, &i.UpdatedTS, &i.Role, &i.Title, &i.Twitter,
	)
	if err != nil {
		return nil, errors.Wrap(err, "querying individual from db")
	}
	const associationsQ = `
		SELECT
			id,
			description
		FROM associations
		JOIN individual_associations
		ON associations.id = individual_associations.association_id
		WHERE individual_associations.individual_id = $1
	`
	rows, err := s.db.QueryContext(ctx, associationsQ, id)
	if err != nil {
		return nil, errors.Wrap(err, "querying individual associations from db")
	}
	defer rows.Close()
	for rows.Next() {
		var a association
		err := rows.Scan(&a.ID, &a.Description)
		if err != nil {
			return nil, errors.Wrap(err, "scanning association row")
		}
		i.Associations = append(i.Associations, a)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading association rows")
	}
	return i, nil
}

// individualname contains only the ID and full name of an
// individual. Used to return list results from search queries
type individualname struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Server) searchIndividuals(ctx context.Context, query string) ([]individualname, error) {
	// Don't start showing suggestions until query is at least 3 chars
	if len(query) < 3 {
		return nil, nil
	}

	const searchQ = `
		SELECT id
		FROM individuals
		WHERE 
			cfb_name ILIKE '%' || $1 || '%' OR
			first_name ILIKE '%' || $1 || '%' OR
			last_name ILIKE '%' || $1 || '%' 
	`
	var individualIDs pq.StringArray
	rows, err := s.db.QueryContext(ctx, searchQ, query)
	if err != nil {
		return nil, errors.Wrap(err, "searching database for individuals by query")
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.Wrap(err, "scanning query row")
		}
		individualIDs = append(individualIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading search query rows")
	}

	// Get individual names matching query results
	const individualsQ = `
		SELECT id, first_name, last_name
		FROM individuals
		WHERE id = ANY($1::text[])
	`
	rows, err = s.db.QueryContext(ctx, individualsQ, individualIDs)
	if err != nil {
		return nil, errors.Wrap(err, "querying individual name results")
	}
	defer rows.Close()

	var resp []individualname
	for rows.Next() {
		var i individualname
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName)
		if err != nil {
			return nil, errors.Wrap(err, "scanning individual name row")
		}
		resp = append(resp, i)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading individual name rows")
	}
	return resp, nil
}
