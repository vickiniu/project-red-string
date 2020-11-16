package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type config struct {
	dbURL          string
	airtableBaseID string
	airtableAPIKey string
}

func main() {
	ctx := context.Background()
	cfg := config{
		dbURL:          envString("DATABASE_URL", "postgres:///redstring?sslmode=disable"),
		airtableBaseID: envString("AIRTABLE_BASE_ID", ""),
		airtableAPIKey: envString("AIRTABLE_API_KEY", ""),
	}
	c := airtableClient{
		baseID: cfg.airtableBaseID,
		apiKey: cfg.airtableAPIKey,
	}

	// Open DB connection
	db, err := sql.Open("postgres", cfg.dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %v\n", err)
	}

	err = c.forEachRecord(func(r record) error {
		return insertAnnotation(ctx, db, r.Fields)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func insertAnnotation(ctx context.Context, db *sql.DB, a annotation) error {
	// Upsert individual
	individualID, err := a.upsertIndividual(ctx, db)
	if err != nil {
		return errors.Wrap(err, "upserting individual")
	}

	// Now, upsert associations
	associationIDs, err := a.upsertAssociations(ctx, db)
	if err != nil {
		return errors.Wrap(err, "upserting associaitons")
	}

	// Insert individual : association mappings
	for _, aid := range associationIDs {
		const insertQ = `
			INSERT INTO individual_associations (
				individual_id,
				association_id,
				updated_ts
			) VALUES (
				$1,
				$2,
				current_timestamp
			)
		`
		_, err := db.ExecContext(ctx, insertQ, individualID, aid)
		if err != nil {
			return errors.Wrap(err, "inserting individual association")
		}
	}
	return nil
}

func (a annotation) upsertIndividual(ctx context.Context, db *sql.DB) (string, error) {
	const individualQ = `SELECT id FROM individuals WHERE first_name = $1 AND last_name = $2`
	var individualID string
	err := db.QueryRowContext(ctx, individualQ, a.FirstName, a.LastName).Scan(&individualID)
	if err == sql.ErrNoRows {
		const insertQ = `
			INSERT INTO individuals (
				first_name,
				last_name,
				cfb_name,
				role,
				updated_ts
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				current_timestamp
			) RETURNING id
		`
		err := db.QueryRowContext(
			ctx, insertQ, a.FirstName, a.LastName, cfbName(a.FirstName, a.LastName), a.Role,
		).Scan(&individualID)
		if err != nil {
			return "", errors.Wrap(err, "inserting individual")
		}
	} else if err != nil {
		return "", errors.Wrap(err, "querying individual")
	}
	return individualID, nil
}

func (a annotation) upsertAssociations(ctx context.Context, db *sql.DB) ([]string, error) {
	var associationIDs []string
	for _, assoc := range a.Associations {
		const associationQ = `SELECT id FROM associations WHERE description = $1`
		var aid string
		err := db.QueryRowContext(ctx, associationQ, assoc).Scan(&aid)
		if err == sql.ErrNoRows {
			const insertQ = `
				INSERT INTO associations (
					description
				) VALUES (
					$1
				) RETURNING id
			`
			err = db.QueryRowContext(ctx, insertQ, assoc).Scan(&aid)
			if err != nil {
				return nil, errors.Wrap(err, "inserting association")
			}
		} else if err != nil {
			return nil, errors.Wrap(err, "querying association")
		}
		associationIDs = append(associationIDs, aid)
	}
	return associationIDs, nil
}

func cfbName(first, last string) string {
	return fmt.Sprintf("%s, %s", last, first)
}

// envString returns the value of the named environment variable.
// If name isn't in the environment os ir empty, it returns value.
func envString(name, value string) string {
	if s := os.Getenv(name); s != "" {
		value = s
	}
	return value
}
