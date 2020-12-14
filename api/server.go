package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Server handles API requests and manages
// server state
type Server struct {
	db *sql.DB
}

// NewServer returns a new Server object
func NewServer(db *sql.DB) *Server {
	return &Server{
		db,
	}
}

// API returns an http.Handler implementing the Red String API
func (s *Server) API() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-individual", s.handleGetIndividual)
	mux.HandleFunc("/search-individuals", s.handleSearchIndividuals)
	mux.HandleFunc("/individual-contributions-received", s.handleIndividualContributionsReceived)
	mux.HandleFunc("/individual-contributions-given", s.handleIndividualContributionsGiven)

	mux.HandleFunc("/categories", s.handleGetCategories)
	mux.HandleFunc("/individual-categories", s.handleGetIndividualsByCategory)
	mux.HandleFunc("/category-associations", s.handleGetAssociationsForCategory)
	mux.HandleFunc("/individual-associations", s.handleGetIndividualsByAssociation)
	return http.Handler(mux)
}

func (s *Server) handleGetIndividual(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividual: unmarshaling request body"))
		return
	}
	i, err := s.getIndividual(r.Context(), body.ID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividual: getting individual"))
		return
	}
	respsuccess(w, r, i)
	return
}

func (s *Server) handleSearchIndividuals(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Query string `json:"query"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleSearchIndividuals: unmarshaling request body"))
		return
	}
	resp, err := s.searchIndividuals(r.Context(), body.Query)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleSearchIndividuals: getting search responses"))
		return
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleIndividualContributionsReceived(w http.ResponseWriter, r *http.Request) {
	body := struct {
		IndividualID string `json:"individual_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleIndividualContributionsReceived: unmarshaling request body"))
		return
	}
	resp, err := s.contributionsReceived(r.Context(), body.IndividualID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleIndividualContributionsReceived: getting search responses"))
		return
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleIndividualContributionsGiven(w http.ResponseWriter, r *http.Request) {
	body := struct {
		IndividualID string `json:"individual_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleIndividualContributionsGiven: unmarshaling request body"))
		return
	}
	resp, err := s.contributionsGiven(r.Context(), body.IndividualID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleIndividualContributionsGiven: getting search responses"))
		return
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := s.getCategories(r.Context())
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetCategories: getting categories"))
		return
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleGetIndividualsByCategory(w http.ResponseWriter, r *http.Request) {
	body := struct {
		CategoryID string `json:"category_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividualsByCategory: unmarshaling request body"))
		return
	}
	resp, err := s.getIndividualsByCategory(r.Context(), body.CategoryID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividualsByCategory: getting individuals"))
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleGetAssociationsForCategory(w http.ResponseWriter, r *http.Request) {
	body := struct {
		CategoryID string `json:"category_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetAssociationsForCategory: unmarshaling request body"))
		return
	}
	resp, err := s.getAssociationsForCategory(r.Context(), body.CategoryID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetAssociationsForCategory: getting individuals"))
	}
	respsuccess(w, r, resp)
}

func (s *Server) handleGetIndividualsByAssociation(w http.ResponseWriter, r *http.Request) {
	body := struct {
		AssociationID string `json:"association_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividualsByAssociation: unmarshaling request body"))
		return
	}
	resp, err := s.getIndividualsByAssociation(r.Context(), body.AssociationID)
	if err != nil {
		resperr(w, r, errors.Wrap(err, "handleGetIndividualsByAssociation: getting individuals"))
	}
	respsuccess(w, r, resp)
}

func resperr(w http.ResponseWriter, r *http.Request, err error) {
	// TODO(vicki): better error instrumentation + logging
	log.Println("error: ", err)
	w.WriteHeader(500)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(err)
}

func respsuccess(w http.ResponseWriter, r *http.Request, b interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(b)
}
