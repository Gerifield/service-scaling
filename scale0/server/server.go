package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", s.getMessages)
	r.Post("/save", s.saveMessage)

	return r
}

func (s *Server) saveMessage(rw http.ResponseWriter, r *http.Request) {
	type input struct {
		Content string `json:"content"`
	}

	var inp input
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Exec("INSERT INTO messages (id, content) VALUES (?, ?)", id.String(), inp.Content)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(struct {
		ID string `json:"id"`
	}{
		ID: id.String(),
	})
}

func (s *Server) getMessages(rw http.ResponseWriter, r *http.Request) {
	type model struct {
		ID        string    `db:"id" json:"id"`
		Content   string    `db:"content" json:"content"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}

	var resp []model
	err := s.db.Select(&resp, "SELECT id, content, created_at FROM messages LIMIT 100")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(struct {
		Results []model `json:"results"`
	}{
		Results: resp,
	})
}
