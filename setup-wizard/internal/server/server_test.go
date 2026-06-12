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
	rec := do(t, srv, http.MethodGet, "/api/wifi")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, forventede 500", rec.Code)
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ugyldig JSON: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("fejlsvaret mangler en error-besked")
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

func TestOpenStepOpensConfiguredURL(t *testing.T) {
	fake := &osfake.Fake{}
	srv := New(fake)

	if rec := do(t, srv, http.MethodPost, "/api/steps/moodle/open"); rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, forventede 204", rec.Code)
	}
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLMoodle {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLMoodle)
	}
}

func TestOpenRejectsUnknownAndURLLessSteps(t *testing.T) {
	fake := &osfake.Fake{}
	srv := New(fake)

	if rec := do(t, srv, http.MethodPost, "/api/steps/findes-ikke/open"); rec.Code != http.StatusNotFound {
		t.Errorf("ukendt trin: status = %d, forventede 404", rec.Code)
	}
	// welcome-trinnet har ingen URL i konfigurationen
	if rec := do(t, srv, http.MethodPost, "/api/steps/welcome/open"); rec.Code != http.StatusBadRequest {
		t.Errorf("trin uden URL: status = %d, forventede 400", rec.Code)
	}
	if len(fake.OpenedURLs) != 0 {
		t.Errorf("der blev åbnet URL'er: %v", fake.OpenedURLs)
	}
}

func TestOpenStepReportsOSError(t *testing.T) {
	srv := New(&osfake.Fake{OpenURLErr: errors.New("ingen browser fundet")})
	rec := do(t, srv, http.MethodPost, "/api/steps/moodle/open")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, forventede 500", rec.Code)
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ugyldig JSON: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("fejlsvaret mangler en error-besked")
	}
}

func TestOpenURLLessStepReturnsJSONError(t *testing.T) {
	srv := New(&osfake.Fake{})
	rec := do(t, srv, http.MethodPost, "/api/steps/welcome/open")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, forventede 400", rec.Code)
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ugyldig JSON: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("fejlsvaret mangler en error-besked")
	}
}

func TestSketchUpInstallReportsOutcome(t *testing.T) {
	cases := []struct {
		name       string
		fake       *osfake.Fake
		wantAction string
		wantReason bool
	}{
		{"winget virker", &osfake.Fake{WingetOK: true}, "installed", false},
		{"S-mode", &osfake.Fake{WingetOK: true, InSMode: true}, "fallback", true},
		{"winget mangler", &osfake.Fake{}, "fallback", true},
		{"winget fejler", &osfake.Fake{WingetOK: true, InstallErr: errors.New("fejlkode 1")}, "fallback", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := New(tc.fake)
			rec := do(t, srv, http.MethodPost, "/api/sketchup/install")
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, forventede 200", rec.Code)
			}
			var resp struct {
				Action string `json:"action"`
				Reason string `json:"reason"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("ugyldig JSON: %v", err)
			}
			if resp.Action != tc.wantAction {
				t.Errorf("action = %q, forventede %q", resp.Action, tc.wantAction)
			}
			if tc.wantReason && resp.Reason == "" {
				t.Errorf("fallback uden elevvendt begrundelse")
			}
			if !tc.wantReason && resp.Reason != "" {
				t.Errorf("reason = %q ved vellykket installation, forventede tom streng", resp.Reason)
			}
		})
	}
}
