// Package server udstiller Assistentens lokale HTTP-API og serverer den
// indlejrede tringuide-side. API'et er tringuidens eneste vej til
// Go-processen; serveren bindes kun til localhost (se cmd/assistent).
package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/web"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/wizard"
)

// Server er Assistentens HTTP-API. Den implementerer http.Handler og
// holder al tilstand, så main kan være tynd.
type Server struct {
	wiz      *wizard.Wizard
	state    *state
	mux      *http.ServeMux
	quit     chan struct{}
	quitOnce sync.Once
}

func New(osImpl osops.OS) *Server {
	s := &Server{
		wiz:   wizard.New(osImpl),
		state: newState(steps.All()),
		mux:   http.NewServeMux(),
		quit:  make(chan struct{}),
	}
	s.mux.HandleFunc("GET /api/steps", s.handleSteps)
	s.mux.HandleFunc("POST /api/steps/{id}/done", s.handleSetStatus("done"))
	s.mux.HandleFunc("POST /api/steps/{id}/undo", s.handleSetStatus(""))
	s.mux.HandleFunc("POST /api/steps/{id}/skip", s.handleSetStatus("skipped"))
	s.mux.HandleFunc("POST /api/steps/{id}/open", s.handleOpen)
	s.mux.HandleFunc("GET /api/system", s.handleSystem)
	s.mux.HandleFunc("POST /api/wifi/settings", s.handleWifiSettings)
	s.mux.HandleFunc("POST /api/quit", s.handleQuit)
	s.mux.Handle("GET /static/", http.FileServerFS(web.Static))
	s.mux.HandleFunc("GET /", s.handleIndex)
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

// handleSetStatus opdaterer et trins status (done/skipped/pending).
func (s *Server) handleSetStatus(status string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.state.setStatus(r.PathValue("id"), status) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleSystem(w http.ResponseWriter, r *http.Request) {
	sMode, err := s.wiz.SMode()
	if err != nil {
		sMode = false
	}
	writeJSON(w, map[string]bool{"sMode": sMode})
}

func (s *Server) handleWifiSettings(w http.ResponseWriter, r *http.Request) {
	if err := s.wiz.OpenWifiSettings(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
		log.Printf("kunne ikke skrive JSON-fejlsvar: %v", encErr)
	}
}

// handleOpen åbner trinnets officielle URL. URL'en slås altid op i den
// indlejrede trinkonfiguration, så API'et kun kan åbne de officielle sider.
func (s *Server) handleOpen(w http.ResponseWriter, r *http.Request) {
	step, ok := s.state.byID(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	if step.URL == "" {
		writeError(w, http.StatusBadRequest, errors.New("trinnet har ingen tilknyttet side"))
		return
	}
	if err := s.wiz.OpenURL(step.URL); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Quit lukkes når eleven afslutter Assistenten; main venter på kanalen.
func (s *Server) Quit() <-chan struct{} { return s.quit }

func (s *Server) handleQuit(w http.ResponseWriter, r *http.Request) {
	s.quitOnce.Do(func() { close(s.quit) })
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFileFS(w, r, web.Static, "static/index.html")
}
