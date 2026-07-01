package server

import (
	"sync"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

// state holder trin-status i hukommelsen. Status nulstilles når
// Assistenten lukkes; persistens til fil er et senere sub-goal.
type state struct {
	mu     sync.Mutex
	steps  []steps.Step
	status map[string]string // id -> "done" | "skipped"; fraværende = pending
}

func newState(all []steps.Step) *state {
	return &state{steps: all, status: make(map[string]string)}
}

// stepView er ét trin plus dets status, klar til JSON.
type stepView struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	Body    string `json:"body"`
	Warning string `json:"warning"`
	Button  string `json:"button"`
	Done    bool   `json:"done"`
	Skipped bool   `json:"skipped"`
}

func (st *state) list() []stepView {
	st.mu.Lock()
	defer st.mu.Unlock()
	views := make([]stepView, len(st.steps))
	for i, s := range st.steps {
		status := st.status[s.ID]
		views[i] = stepView{
			ID:      s.ID,
			Title:   s.Title,
			Kind:    string(s.Kind),
			Body:    s.Body,
			Warning: s.Warning,
			Button:  s.Button,
			Done:    status == "done",
			Skipped: status == "skipped",
		}
	}
	return views
}

// setStatus sætter status for et trin. status "" = pending (rydder).
// Returnerer false hvis id'et ikke findes.
func (st *state) setStatus(id, status string) bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	if status != "" && status != "done" && status != "skipped" {
		return false
	}
	if _, ok := st.byIDLocked(id); !ok {
		return false
	}
	if status == "" {
		delete(st.status, id)
	} else {
		st.status[id] = status
	}
	return true
}

// rawStatus returnerer en kopi af trin-status (id -> "done"/"skipped").
// Bruges til at generere dashboardet med elevens fremdrift.
func (st *state) rawStatus() map[string]string {
	st.mu.Lock()
	defer st.mu.Unlock()
	out := make(map[string]string, len(st.status))
	for id, status := range st.status {
		out[id] = status
	}
	return out
}

// byID slår et trin op i konfigurationen.
func (st *state) byID(id string) (steps.Step, bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.byIDLocked(id)
}

func (st *state) byIDLocked(id string) (steps.Step, bool) {
	for _, s := range st.steps {
		if s.ID == id {
			return s, true
		}
	}
	return steps.Step{}, false
}
