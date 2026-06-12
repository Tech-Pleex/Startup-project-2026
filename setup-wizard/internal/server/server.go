// Package server udstiller Assistentens lokale HTTP-API og serverer den
// indlejrede tringuide-side. API'et er tringuidens eneste vej til
// Go-processen; serveren bindes kun til localhost (se cmd/assistent).
package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/wizard"
)

// Server er Assistentens HTTP-API. Den implementerer http.Handler og
// holder al tilstand, så main kan være tynd.
type Server struct {
	wiz   *wizard.Wizard
	state *state
	mux   *http.ServeMux
}

func New(osImpl osops.OS) *Server {
	s := &Server{
		wiz:   wizard.New(osImpl),
		state: newState(steps.All()),
		mux:   http.NewServeMux(),
	}
	s.mux.HandleFunc("GET /api/steps", s.handleSteps)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleSteps(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{"steps": s.state.list()})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("kunne ikke skrive JSON-svar: %v", err)
	}
}
