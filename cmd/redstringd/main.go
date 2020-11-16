package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type config struct {
	dbUrl      string
	apiAddr    string
	port       string
	workingDir string
}

func main() {
	cfg := config{
		dbUrl:      envString("DATABASE_URL", "postgres:///redstring?sslmode=disable"),
		apiAddr:    envString("API_ADDR", ""),
		port:       envString("PORT", "8080"),
		workingDir: envString("REDSTRING", ""),
	}
	// Open DB connection
	db, err := sql.Open("postgres", cfg.dbUrl)
	if err != nil {
		log.Fatalf("error opening database connection: %v\n", err)
	}

	// Initialize DB from schema
	schema, err := ioutil.ReadFile(cfg.workingDir + "/schema.sql")
	if err != nil {
		log.Fatalf("db schema not found: %v\n", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("error initializing db from schema: %v\n", err)
	}

	// // Check to see if DB exists — if not, create and set schema
	// var x interface{}
	// const q = `select true from pg_database where datname = $1`
	// err = db.QueryRow(q, dbName).Scan(&x)
	// if err == sql.ErrNoRows {
	// 	_, err = db.Exec("CREATE DATABASE " + dbName)
	// 	if err != nil {
	// 		log.Fatalf("error creating database: %v", err)
	// 	}
	// 	schema, err := ioutil.ReadFile(schemaFile)
	// 	if err != nil {
	// 		log.Fatalf("error reading schema: %v", err)
	// 	}
	// 	_, err = db.Exec(string(schema))
	// 	if err != nil {
	// 		log.Fatalf("error setting db schema: %v", err)
	// 	}
	// } else if err != nil {
	// 	log.Fatalf("error querying pg_database: %v", err)
	// }
}

// envString returns the value of the named environment variable.
// If name isn't in the environment os ir empty, it returns value.
func envString(name, value string) string {
	if s := os.Getenv(name); s != "" {
		value = s
	}
	return value
}
