package server

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/web"
)

type stepJSON struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	Body    string `json:"body"`
	Warning string `json:"warning"`
	Button  string `json:"button"`
	Done    bool   `json:"done"`
	Skipped bool   `json:"skipped"`
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

func TestOpenSketchUpStepOpensOfficialDownload(t *testing.T) {
	fake := &osfake.Fake{}
	srv := New(fake)

	rec := do(t, srv, http.MethodPost, "/api/steps/sketchup/open")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, forventede 204", rec.Code)
	}
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLSketchUpDownload {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLSketchUpDownload)
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

func TestSystemStatusReportsSMode(t *testing.T) {
	for _, tc := range []struct {
		name string
		fake *osfake.Fake
		want bool
	}{
		{"inaktiv", &osfake.Fake{}, false},
		{"aktiv", &osfake.Fake{InSMode: true}, true},
		{"ukendt", &osfake.Fake{SModeErr: errors.New("reg fejlede")}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := New(tc.fake)
			rec := do(t, srv, http.MethodGet, "/api/system")
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, forventede 200", rec.Code)
			}
			var resp struct {
				SMode bool `json:"sMode"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("ugyldig JSON: %v", err)
			}
			if resp.SMode != tc.want {
				t.Errorf("sMode = %v, forventede %v", resp.SMode, tc.want)
			}
		})
	}
}

func TestQuitSignalsShutdown(t *testing.T) {
	srv := New(&osfake.Fake{})
	select {
	case <-srv.Quit():
		t.Fatal("Quit-kanalen er lukket før /api/quit er kaldt")
	default:
	}

	if rec := do(t, srv, http.MethodPost, "/api/quit"); rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, forventede 204", rec.Code)
	}
	select {
	case <-srv.Quit():
	default:
		t.Error("Quit-kanalen er ikke lukket efter /api/quit")
	}

	// Endnu et kald må ikke panikke (dobbelt close).
	if rec := do(t, srv, http.MethodPost, "/api/quit"); rec.Code != http.StatusNoContent {
		t.Errorf("andet kald: status = %d, forventede 204", rec.Code)
	}
}

func TestIndexServesTringuideWithSafetyText(t *testing.T) {
	srv := New(&osfake.Fake{})
	rec := do(t, srv, http.MethodGet, "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /: status = %d, forventede 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, steps.SafetyText) {
		t.Errorf("siden mangler sikkerhedsteksten om adgangskoder/MitID/UNI-Login")
	}
	if !strings.Contains(body, "Assistenten") {
		t.Errorf("siden bruger ikke navnet Assistenten")
	}
	for _, want := range []string{
		"Windows S-mode er aktiveret",
		"Assistenten kan ikke fortsætte",
		"Tjek igen",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("siden mangler %q", want)
		}
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	srv := New(&osfake.Fake{})
	for _, path := range []string{"/static/style.css", "/static/app.js"} {
		if rec := do(t, srv, http.MethodGet, path); rec.Code != http.StatusOK {
			t.Errorf("GET %s: status = %d, forventede 200", path, rec.Code)
		}
	}
}

func TestEmbeddedAssetsContainNoInternalNames(t *testing.T) {
	err := fs.WalkDir(web.Static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(web.Static, path)
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(string(data)), "setup-wizard") {
			t.Errorf("%s indeholder det interne navn \"setup-wizard\"", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("kunne ikke gennemgå indlejrede filer: %v", err)
	}
}

func TestSkipStepReopenAndComplete(t *testing.T) {
	srv := New(&osfake.Fake{})

	// Spring sketchup over.
	if rec := do(t, srv, http.MethodPost, "/api/steps/sketchup/skip"); rec.Code != http.StatusNoContent {
		t.Fatalf("skip: status = %d, forventede 204", rec.Code)
	}
	s := findStep(t, getSteps(t, srv), "sketchup")
	if !s.Skipped || s.Done {
		t.Errorf("efter skip: skipped=%v done=%v, forventede skipped=true done=false", s.Skipped, s.Done)
	}

	// Genåbn (undo) -> hverken done eller skipped.
	if rec := do(t, srv, http.MethodPost, "/api/steps/sketchup/undo"); rec.Code != http.StatusNoContent {
		t.Fatalf("undo: status = %d, forventede 204", rec.Code)
	}
	s = findStep(t, getSteps(t, srv), "sketchup")
	if s.Skipped || s.Done {
		t.Errorf("efter undo: skipped=%v done=%v, forventede begge false", s.Skipped, s.Done)
	}

	// Skip og dernæst markér færdig -> done rydder skipped.
	if rec := do(t, srv, http.MethodPost, "/api/steps/sketchup/skip"); rec.Code != http.StatusNoContent {
		t.Fatalf("skip: status = %d, forventede 204", rec.Code)
	}
	if rec := do(t, srv, http.MethodPost, "/api/steps/sketchup/done"); rec.Code != http.StatusNoContent {
		t.Fatalf("done: status = %d, forventede 204", rec.Code)
	}
	s = findStep(t, getSteps(t, srv), "sketchup")
	if !s.Done || s.Skipped {
		t.Errorf("efter done: done=%v skipped=%v, forventede done=true skipped=false", s.Done, s.Skipped)
	}
}

func TestSkipLastStep(t *testing.T) {
	srv := New(&osfake.Fake{})
	if rec := do(t, srv, http.MethodPost, "/api/steps/finish/skip"); rec.Code != http.StatusNoContent {
		t.Fatalf("skip finish: status = %d, forventede 204", rec.Code)
	}
	if !findStep(t, getSteps(t, srv), "finish").Skipped {
		t.Errorf("finish-trinnet er ikke markeret sprunget over")
	}
}

func TestSkipUnknownStepReturnsNotFound(t *testing.T) {
	srv := New(&osfake.Fake{})
	if rec := do(t, srv, http.MethodPost, "/api/steps/findes-ikke/skip"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, forventede 404", rec.Code)
	}
}

func TestBrandAssetsAreServed(t *testing.T) {
	srv := New(&osfake.Fake{})
	for _, path := range []string{
		"/static/fonts/space-grotesk.woff2",
		"/static/img/neg-hero.jpg",
		"/static/img/neg-logo-white.png",
	} {
		if rec := do(t, srv, http.MethodGet, path); rec.Code != http.StatusOK {
			t.Errorf("GET %s: status = %d, forventede 200", path, rec.Code)
		}
	}
}

func TestDashboardWritesToDesktop(t *testing.T) {
	desktop := t.TempDir()
	fake := &osfake.Fake{Desktop: desktop}
	srv := New(fake)

	// Marker wifi + office færdige; det giver flueben på dashboardets
	// studentSteps 0, 1 og 2.
	for _, id := range []string{"wifi", "office"} {
		if rec := do(t, srv, http.MethodPost, "/api/steps/"+id+"/done"); rec.Code != http.StatusNoContent {
			t.Fatalf("done %s: status = %d", id, rec.Code)
		}
	}

	rec := do(t, srv, http.MethodPost, "/api/dashboard")
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/dashboard: status = %d, forventede 200 (%s)", rec.Code, rec.Body.String())
	}

	var resp struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ugyldig JSON: %v", err)
	}
	if filepath.Dir(resp.Path) != desktop {
		t.Errorf("filen blev gemt i %q, forventede skrivebordet %q", filepath.Dir(resp.Path), desktop)
	}

	data, err := os.ReadFile(resp.Path)
	if err != nil {
		t.Fatalf("dashboard-filen blev ikke skrevet: %v", err)
	}
	html := string(data)
	if !strings.Contains(html, "data:font/woff2;base64,") {
		t.Error("den gemte fil er ikke selvstændig (mangler indlejret font)")
	}
	if !strings.Contains(html, `"0":"Færdig"`) {
		t.Error("den gemte fil mangler elevens flueben")
	}

	// Filen skal være forsøgt åbnet for eleven.
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != resp.Path {
		t.Errorf("forventede at filen blev åbnet (%s), fik %v", resp.Path, fake.OpenedURLs)
	}
}
