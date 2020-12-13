package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/vickiniu/project-red-string/api"

	_ "github.com/lib/pq"
)

func main() {
	// TODO(vicki): actually configure
	port := ":8080"
	dbURL := envString("DATABASE_URL", "postgres:///redstring?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %v\n", err)
	}
	s := api.NewServer(db)
	httpserver := http.Server{
		Addr:    port,
		Handler: s.API(),
	}
	log.Printf("Starting server on port %s", httpserver.Addr)
	httpserver.ListenAndServe()
}

// envString returns the value of the named environment variable.
// If name isn't in the environment os ir empty, it returns value.
func envString(name, value string) string {
	if s := os.Getenv(name); s != "" {
		value = s
	}
	return value
}
