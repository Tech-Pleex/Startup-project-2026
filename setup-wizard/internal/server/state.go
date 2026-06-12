package server

import (
	"sync"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

// state holder trin-status i hukommelsen. Status nulstilles når
// Assistenten lukkes; persistens til fil er et senere sub-goal.
type state struct {
	mu    sync.Mutex
	steps []steps.Step
	done  map[string]bool
}

func newState(all []steps.Step) *state {
	return &state{steps: all, done: make(map[string]bool)}
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
}

func (st *state) list() []stepView {
	st.mu.Lock()
	defer st.mu.Unlock()
	views := make([]stepView, len(st.steps))
	for i, s := range st.steps {
		views[i] = stepView{
			ID:      s.ID,
			Title:   s.Title,
			Kind:    string(s.Kind),
			Body:    s.Body,
			Warning: s.Warning,
			Button:  s.Button,
			Done:    st.done[s.ID],
		}
	}
	return views
}

// setDone sætter status for et trin; false hvis id'et ikke findes.
func (st *state) setDone(id string, done bool) bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.byIDLocked(id); !ok {
		return false
	}
	st.done[id] = done
	return true
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
