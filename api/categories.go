package api

import (
	"context"

	"github.com/pkg/errors"
)

type category struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (s *Server) getCategories(ctx context.Context) ([]category, error) {
	const q = `
		SELECT id, description
		FROM categories
	`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "querying all categories from db")
	}
	defer rows.Close()

	var categories []category
	for rows.Next() {
		var c category
		err := rows.Scan(&c.ID, &c.Description)
		if err != nil {
			return nil, errors.Wrap(err, "scanning category row")
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading category rows")
	}
	return categories, nil
}

func (s *Server) getIndividualsByCategory(ctx context.Context, categoryID string) ([]individualname, error) {
	const q = `
		SELECT id, first_name, last_name
		FROM individuals
		JOIN individual_associations
		ON individuals.id = individual_associations.individual_id
		WHERE individual_associations.association_id IN (
			SELECT id 
			FROM associations
			WHERE category_id = $1
		)
	`
	rows, err := s.db.QueryContext(ctx, q, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "querying individuals for category from db")
	}
	defer rows.Close()

	var individuals []individualname
	for rows.Next() {
		var i individualname
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName)
		if err != nil {
			return nil, errors.Wrap(err, "scanning individual name row")
		}
		individuals = append(individuals, i)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading individual name rows")
	}
	return individuals, nil
}

func (s *Server) getAssociationsForCategory(ctx context.Context, categoryID string) ([]association, error) {
	const q = `
		SELECT id, description
		FROM associations
		WHERE category_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "querying associations for category from db")
	}
	defer rows.Close()

	var associations []association
	for rows.Next() {
		var a association
		err := rows.Scan(&a.ID, &a.Description)
		if err != nil {
			return nil, errors.Wrap(err, "scanning association row")
		}
		associations = append(associations, a)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading association rows")
	}
	return associations, nil
}

func (s *Server) getIndividualsByAssociation(ctx context.Context, associationID string) ([]individualname, error) {
	const q = `
		SELECT id, first_name, last_name
		FROM individuals
		JOIN individual_associations
		ON individuals.id = individual_associations.individual_id
		WHERE individual_associations.association_id = $1
	`
	rows, err := s.db.QueryContext(ctx, q, associationID)
	if err != nil {
		return nil, errors.Wrap(err, "querying individuals for association from db")
	}
	defer rows.Close()

	var individuals []individualname
	for rows.Next() {
		var i individualname
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName)
		if err != nil {
			return nil, errors.Wrap(err, "scanning individual name row")
		}
		individuals = append(individuals, i)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "reading individual name rows")
	}
	return individuals, nil
}
