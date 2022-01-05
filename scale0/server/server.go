package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gerifield/service-scaling/scale0/model"
	"github.com/go-chi/chi/v5"
)

type appLogic interface {
	Save(ctx context.Context, content string) (string, error)
	GetAll(ctx context.Context) ([]model.Message, error)
}

type Server struct {
	app appLogic
}

func New(app appLogic) *Server {
	return &Server{
		app: app,
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

	id, err := s.app.Save(r.Context(), inp.Content)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(struct {
		ID string `json:"id"`
	}{
		ID: id,
	})
}

func (s *Server) getMessages(rw http.ResponseWriter, r *http.Request) {
	resp, err := s.app.GetAll(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(struct {
		Results []model.Message `json:"results"`
	}{
		Results: resp,
	})
}
