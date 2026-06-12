package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

type stepJSON struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	Body    string `json:"body"`
	Warning string `json:"warning"`
	Button  string `json:"button"`
	Done    bool   `json:"done"`
}

type stepsResponse struct {
	Steps []stepJSON `json:"steps"`
}

// do sender ét request gennem serveren uden netværk.
func do(t *testing.T, srv *Server, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func getSteps(t *testing.T, srv *Server) stepsResponse {
	t.Helper()
	rec := do(t, srv, http.MethodGet, "/api/steps")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/steps: status = %d, forventede 200", rec.Code)
	}
	var resp stepsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ugyldig JSON: %v", err)
	}
	return resp
}

func findStep(t *testing.T, resp stepsResponse, id string) stepJSON {
	t.Helper()
	for _, s := range resp.Steps {
		if s.ID == id {
			return s
		}
	}
	t.Fatalf("trin %q findes ikke i svaret", id)
	return stepJSON{}
}

func TestStepsReturnsAllTenInOrder(t *testing.T) {
	srv := New(&osfake.Fake{})
	resp := getSteps(t, srv)

	all := steps.All()
	if len(resp.Steps) != len(all) {
		t.Fatalf("antal trin = %d, forventede %d", len(resp.Steps), len(all))
	}
	for i, want := range all {
		got := resp.Steps[i]
		if got.ID != want.ID {
			t.Errorf("trin %d: id = %q, forventede %q", i, got.ID, want.ID)
		}
		if got.Title != want.Title || got.Body != want.Body || got.Button != want.Button || got.Warning != want.Warning {
			t.Errorf("trin %q: tekstfelter matcher ikke konfigurationen", want.ID)
		}
		if got.Kind != string(want.Kind) {
			t.Errorf("trin %q: kind = %q, forventede %q", want.ID, got.Kind, want.Kind)
		}
		if got.Done {
			t.Errorf("trin %q: nyt trin er allerede markeret færdigt", want.ID)
		}
	}
}

func TestMarkStepDoneAndUndo(t *testing.T) {
	srv := New(&osfake.Fake{})

	if rec := do(t, srv, http.MethodPost, "/api/steps/wifi/done"); rec.Code != http.StatusNoContent {
		t.Fatalf("done: status = %d, forventede 204", rec.Code)
	}
	if !findStep(t, getSteps(t, srv), "wifi").Done {
		t.Errorf("wifi-trinnet er ikke markeret færdigt efter done")
	}

	if rec := do(t, srv, http.MethodPost, "/api/steps/wifi/undo"); rec.Code != http.StatusNoContent {
		t.Fatalf("undo: status = %d, forventede 204", rec.Code)
	}
	if findStep(t, getSteps(t, srv), "wifi").Done {
		t.Errorf("wifi-trinnet er stadig færdigt efter undo")
	}
}

func TestMarkUnknownStepReturnsNotFound(t *testing.T) {
	srv := New(&osfake.Fake{})
	for _, path := range []string{"/api/steps/findes-ikke/done", "/api/steps/findes-ikke/undo"} {
		if rec := do(t, srv, http.MethodPost, path); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, forventede 404", path, rec.Code)
		}
	}
}

func TestWifiStatusClassifiesNetworks(t *testing.T) {
	cases := []struct {
		name, ssid, wantState string
	}{
		{"på NEG", "NEG", "target"},
		{"på NEG Guest", "NEG Guest", "guest"},
		{"ukendt netværk", "Naboens Netværk", "other"},
		{"intet netværk", "", "none"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := New(&osfake.Fake{SSID: tc.ssid})
			rec := do(t, srv, http.MethodGet, "/api/wifi")
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, forventede 200", rec.Code)
			}
			var resp struct {
				SSID  string `json:"ssid"`
				State string `json:"state"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("ugyldig JSON: %v", err)
			}
			if resp.State != tc.wantState || resp.SSID != tc.ssid {
				t.Errorf("svar = %+v, forventede ssid=%q state=%q", resp, tc.ssid, tc.wantState)
			}
		})
	}
}

func TestWifiStatusReportsOSError(t *testing.T) {
	srv := New(&osfake.Fake{SSIDErr: errors.New("netsh fejlede")})
	if rec := do(t, srv, http.MethodGet, "/api/wifi"); rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, forventede 500", rec.Code)
	}
}

func TestOpenWifiSettings(t *testing.T) {
	fake := &osfake.Fake{}
	srv := New(fake)
	if rec := do(t, srv, http.MethodPost, "/api/wifi/settings"); rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, forventede 204", rec.Code)
	}
	if fake.WifiSettingsOpens != 1 {
		t.Errorf("Wi-Fi-indstillinger åbnet %d gange, forventede 1", fake.WifiSettingsOpens)
	}
}
