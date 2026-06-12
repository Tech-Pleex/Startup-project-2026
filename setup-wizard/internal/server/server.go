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
	s.mux.HandleFunc("POST /api/steps/{id}/done", s.handleSetDone(true))
	s.mux.HandleFunc("POST /api/steps/{id}/undo", s.handleSetDone(false))
	s.mux.HandleFunc("POST /api/steps/{id}/open", s.handleOpen)
	s.mux.HandleFunc("GET /api/wifi", s.handleWifi)
	s.mux.HandleFunc("POST /api/wifi/settings", s.handleWifiSettings)
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

func (s *Server) handleSetDone(done bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.state.setDone(r.PathValue("id"), done) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleWifi(w http.ResponseWriter, r *http.Request) {
	status, err := s.wiz.WifiStatus()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, map[string]string{"ssid": status.SSID, "state": string(status.State)})
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
		http.Error(w, "trinnet har ingen tilknyttet side", http.StatusBadRequest)
		return
	}
	if err := s.wiz.OpenURL(step.URL); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
