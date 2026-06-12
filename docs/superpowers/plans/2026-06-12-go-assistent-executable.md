# Kørbar Assistent-binær (Windows + Mac) — implementeringsplan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fra én Go-kodebase bygges en dobbeltklikbar Windows-`.exe` og Mac-binær, der starter et lokalt HTTP-API, serverer en indlejret dansk tringuide og åbner elevens browser.

**Architecture:** Tynd `cmd/assistent/main.go` binder til `127.0.0.1:0` og delegerer alt til en ny `internal/server`-pakke (HTTP-API med stdlib `ServeMux`, trin-status i hukommelsen) og `internal/web` (indlejret HTML/CSS/JS via `go:embed`). Al logik testes med `httptest` og den eksisterende OS-fake, som udtrækkes til en delt `osfake`-pakke.

**Tech Stack:** Go stdlib (net/http, embed, httptest) — ingen dependencies. Rå HTML/CSS/JS uden frameworks eller build-step.

**Spec:** `docs/superpowers/specs/2026-06-12-go-assistent-executable-design.md`

**Arbejdsmappe:** Alle kommandoer køres fra repo-roden `/home/jere/PKA/dev/Startup-project-2026`. Go-kommandoer kører i `setup-wizard/` via `(cd setup-wizard && ...)`.

---

## Filoversigt

| Fil | Ansvar |
|---|---|
| `setup-wizard/internal/osops/osfake/osfake.go` | NY: delt test-fake af `osops.OS` (flyttet fra `wizard/fake_os_test.go`) |
| `setup-wizard/internal/wizard/wifi_test.go` | OMSKRIVES: bruger `osfake.Fake` |
| `setup-wizard/internal/wizard/sketchup_test.go` | OMSKRIVES: bruger `osfake.Fake` |
| `setup-wizard/internal/wizard/fake_os_test.go` | SLETTES: erstattet af `osfake` |
| `setup-wizard/internal/osops/current_windows.go` | NY: `Current()` → `Windows{}` |
| `setup-wizard/internal/osops/current_darwin.go` | NY: `Current()` → `Darwin{}` |
| `setup-wizard/internal/osops/osops_linux.go` | NY: Linux-udviklerstub + `Current()` (så WSL kan bygge/køre) |
| `setup-wizard/internal/steps/steps.go` | ÆNDRES: dashboard-URL på finish-trinnet |
| `setup-wizard/internal/server/state.go` | NY: trin-status i hukommelsen (mutex-beskyttet) |
| `setup-wizard/internal/server/server.go` | NY: HTTP-API + servering af indlejrede assets |
| `setup-wizard/internal/server/server_test.go` | NY: httptest-tests af hele API'et |
| `setup-wizard/internal/web/embed.go` | NY: `go:embed` af `static/` |
| `setup-wizard/internal/web/static/index.html` | NY: tringuide-siden |
| `setup-wizard/internal/web/static/style.css` | NY: stil lånt fra `start.html`-tokens |
| `setup-wizard/internal/web/static/app.js` | NY: fetch-baseret klient |
| `setup-wizard/cmd/assistent/main.go` | NY: tynd main — port, server, åbn browser, vent på quit |
| `setup-wizard/README.md` | NY: build-/test-/kørselsinstruktioner |
| `.gitignore` | ÆNDRES: ignorér `setup-wizard/dist/` |

---

### Task 1: Delt OS-fake (`osfake`)

Server-testene skal bruge samme fake som wizard-testene. Fake'en flyttes fra den private testfil til en delt pakke. Ren refaktorering — adfærden er allerede dækket af de 21 eksisterende tests, som skal blive ved med at bestå.

**Files:**
- Create: `setup-wizard/internal/osops/osfake/osfake.go`
- Delete: `setup-wizard/internal/wizard/fake_os_test.go`
- Rewrite: `setup-wizard/internal/wizard/wifi_test.go`
- Rewrite: `setup-wizard/internal/wizard/sketchup_test.go`

- [ ] **Step 1: Opret osfake-pakken**

`setup-wizard/internal/osops/osfake/osfake.go`:

```go
// Package osfake er en testimplementering af osops.OS. Den simulerer
// Wi-Fi-svar, S-mode og winget-resultater og registrerer hvilke URL'er
// der ville være åbnet — ingen test rører det rigtige OS.
package osfake

import "github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"

// Fake styres via felterne og registrerer alle kald.
type Fake struct {
	SSID       string
	SSIDErr    error
	InSMode    bool
	WingetOK   bool
	InstallErr error

	OpenedURLs        []string
	WifiSettingsOpens int
	InstalledPackages []string
}

var _ osops.OS = (*Fake)(nil)

func (f *Fake) ActiveWifiSSID() (string, error) { return f.SSID, f.SSIDErr }

func (f *Fake) OpenWifiSettings() error {
	f.WifiSettingsOpens++
	return nil
}

func (f *Fake) OpenURL(url string) error {
	f.OpenedURLs = append(f.OpenedURLs, url)
	return nil
}

func (f *Fake) SMode() (bool, error) { return f.InSMode, nil }

func (f *Fake) WingetAvailable() bool { return f.WingetOK }

func (f *Fake) InstallSketchUp(packageID string) error {
	f.InstalledPackages = append(f.InstalledPackages, packageID)
	return f.InstallErr
}
```

- [ ] **Step 2: Slet den gamle fake**

```bash
rm setup-wizard/internal/wizard/fake_os_test.go
```

- [ ] **Step 3: Omskriv wifi_test.go til osfake**

Erstat hele indholdet af `setup-wizard/internal/wizard/wifi_test.go` med:

```go
package wizard

import (
	"errors"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
)

func TestWifiStatusOnTargetNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "NEG"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOnTarget {
		t.Errorf("State = %q, forventede %q", status.State, WifiOnTarget)
	}
	if status.SSID != "NEG" {
		t.Errorf("SSID = %q, forventede %q", status.SSID, "NEG")
	}
}

func TestWifiStatusOnGuestNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "NEG Guest"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOnGuest {
		t.Errorf("State = %q, forventede %q", status.State, WifiOnGuest)
	}
}

func TestWifiStatusOnUnknownNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "Naboens Netværk"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOther {
		t.Errorf("State = %q, forventede %q", status.State, WifiOther)
	}
	if status.SSID != "Naboens Netværk" {
		t.Errorf("SSID = %q, forventede %q", status.SSID, "Naboens Netværk")
	}
}

func TestWifiStatusWithNoNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: ""})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiNone {
		t.Errorf("State = %q, forventede %q", status.State, WifiNone)
	}
}

func TestWifiStatusPropagatesOSError(t *testing.T) {
	osErr := errors.New("netsh fejlede")
	w := New(&osfake.Fake{SSIDErr: osErr})

	_, err := w.WifiStatus()
	if !errors.Is(err, osErr) {
		t.Errorf("err = %v, forventede at den ombryder %v", err, osErr)
	}
}

func TestOpenWifiSettingsDelegatesToOS(t *testing.T) {
	fake := &osfake.Fake{}
	w := New(fake)

	if err := w.OpenWifiSettings(); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if fake.WifiSettingsOpens != 1 {
		t.Errorf("Wi-Fi-indstillinger åbnet %d gange, forventede 1", fake.WifiSettingsOpens)
	}
}
```

- [ ] **Step 4: Omskriv sketchup_test.go til osfake**

Erstat hele indholdet af `setup-wizard/internal/wizard/sketchup_test.go` med (bemærk: `errWingetFailed` flytter med herind):

```go
package wizard

import (
	"errors"
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

var errWingetFailed = errors.New("winget returnerede fejlkode 1")

func TestSketchUpInstallsViaWingetWhenPossible(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpInstalled {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpInstalled)
	}
	if len(fake.InstalledPackages) != 1 || fake.InstalledPackages[0] != steps.SketchUpPackageID {
		t.Errorf("installerede pakker = %v, forventede [%q]", fake.InstalledPackages, steps.SketchUpPackageID)
	}
	if len(fake.OpenedURLs) != 0 {
		t.Errorf("fallback-siden blev åbnet ved vellykket installation: %v", fake.OpenedURLs)
	}
}

func TestSketchUpFallsBackInSModeWithoutTryingWinget(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true, InSMode: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if !strings.Contains(outcome.Reason, "S-mode") {
		t.Errorf("begrundelsen nævner ikke S-mode: %q", outcome.Reason)
	}
	if len(fake.InstalledPackages) != 0 {
		t.Errorf("winget blev forsøgt i S-mode: %v", fake.InstalledPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetMissing(t *testing.T) {
	fake := &osfake.Fake{WingetOK: false}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if len(fake.InstalledPackages) != 0 {
		t.Errorf("winget blev forsøgt selvom det ikke findes: %v", fake.InstalledPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetFails(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true, InstallErr: errWingetFailed}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if !strings.Contains(outcome.Reason, errWingetFailed.Error()) {
		t.Errorf("begrundelsen indeholder ikke winget-fejlen: %q", outcome.Reason)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpReasonsAreDanishStudentFacing(t *testing.T) {
	cases := map[string]*osfake.Fake{
		"s-mode":         {WingetOK: true, InSMode: true},
		"winget mangler": {WingetOK: false},
		"winget fejler":  {WingetOK: true, InstallErr: errWingetFailed},
	}
	for name, fake := range cases {
		outcome := New(fake).InstallSketchUp()
		if outcome.Reason == "" {
			t.Errorf("%s: begrundelsen er tom", name)
		}
		if strings.Contains(strings.ToLower(outcome.Reason), "setup-wizard") {
			t.Errorf("%s: begrundelsen indeholder \"setup-wizard\": %q", name, outcome.Reason)
		}
	}
}

func TestOpenStepLinkDelegatesToOS(t *testing.T) {
	fake := &osfake.Fake{}
	w := New(fake)

	if err := w.OpenURL(steps.URLMoodle); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLMoodle {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLMoodle)
	}
}

func assertOpenedFallback(t *testing.T, fake *osfake.Fake) {
	t.Helper()
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLSketchUpFallback {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLSketchUpFallback)
	}
}
```

- [ ] **Step 5: Kør alle tests**

Run: `(cd setup-wizard && go test ./...)`
Expected: `ok` for `internal/steps` og `internal/wizard`, ingen fejl (21 tests består som før).

- [ ] **Step 6: Commit**

```bash
git add setup-wizard/internal/osops/osfake/osfake.go setup-wizard/internal/wizard/
git commit -m "refactor: udtræk OS-fake til delt osfake-pakke"
```

---

### Task 2: `osops.Current()` — platformvalg til main

`main.go` skal kunne vælge OS-implementeringen uden build-tags i main-pakken. Linux-stubben gør at hele modulet kan bygges og køres på udviklingsmaskinen (WSL). Ingen unit-test (én linje pr. platform, platform-tagget); verifikation er at alle tre GOOS bygger.

**Files:**
- Create: `setup-wizard/internal/osops/current_windows.go`
- Create: `setup-wizard/internal/osops/current_darwin.go`
- Create: `setup-wizard/internal/osops/osops_linux.go`

- [ ] **Step 1: Windows-konstruktør**

`setup-wizard/internal/osops/current_windows.go`:

```go
//go:build windows

package osops

// Current returnerer OS-implementeringen for den platform binæren er bygget til.
func Current() OS { return Windows{} }
```

- [ ] **Step 2: Mac-konstruktør**

`setup-wizard/internal/osops/current_darwin.go`:

```go
//go:build darwin

package osops

// Current returnerer OS-implementeringen for den platform binæren er bygget til.
func Current() OS { return Darwin{} }
```

- [ ] **Step 3: Linux-udviklerstub**

`setup-wizard/internal/osops/osops_linux.go`:

```go
//go:build linux

package osops

import (
	"errors"
	"os/exec"
)

// Linux er en udviklerstub, så Assistenten kan bygges og afprøves på
// udviklingsmaskinen (WSL). Elever får kun Windows- og Mac-binærer.
type Linux struct{}

var _ OS = Linux{}

func (Linux) ActiveWifiSSID() (string, error) { return "", nil }

func (Linux) OpenWifiSettings() error { return nil }

func (Linux) OpenURL(url string) error {
	return exec.Command("xdg-open", url).Start()
}

func (Linux) SMode() (bool, error) { return false, nil }

func (Linux) WingetAvailable() bool { return false }

func (Linux) InstallSketchUp(string) error {
	return errors.New("automatisk installation findes ikke på Linux")
}

// Current returnerer OS-implementeringen for den platform binæren er bygget til.
func Current() OS { return Linux{} }
```

- [ ] **Step 4: Verificér at alle tre platforme bygger**

Run: `(cd setup-wizard && go build ./... && GOOS=windows go build ./... && GOOS=darwin go build ./...)`
Expected: ingen output, exit 0 for alle tre.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/osops/
git commit -m "feat: osops.Current() med platformvalg og Linux-udviklerstub"
```

---

### Task 3: Dashboard-URL på finish-trinnet

Finish-trinnet skal kunne åbne dashboardet, men har ingen URL i konfigurationen. Dashboardet er hostet på GitHub Pages. Brødteksten justeres samtidig, så den ikke lover en genvej (desktop-genvejen er ikke en del af denne bid).

**Files:**
- Modify: `setup-wizard/internal/steps/steps.go`
- Test: `setup-wizard/internal/steps/steps_test.go`

- [ ] **Step 1: Skriv den fejlende test**

Tilføj nederst i `setup-wizard/internal/steps/steps_test.go`:

```go
func TestFinishStepOpensDashboard(t *testing.T) {
	for _, s := range All() {
		if s.ID != "finish" {
			continue
		}
		if s.URL != URLDashboard {
			t.Errorf("finish.URL = %q, forventede %q", s.URL, URLDashboard)
		}
		return
	}
	t.Fatal("finish-trinnet findes ikke")
}
```

- [ ] **Step 2: Kør testen — den skal fejle**

Run: `(cd setup-wizard && go test ./internal/steps/ -run TestFinishStepOpensDashboard -v)`
Expected: FAIL — `undefined: URLDashboard` (kompileringsfejl).

- [ ] **Step 3: Implementér**

I `setup-wizard/internal/steps/steps.go` tilføjes i URL-konstantblokken:

```go
	URLDashboard = "https://tech-pleex.github.io/Startup-project-2026/start.html"
```

Og finish-trinnet i `All()` ændres til:

```go
		{
			ID:     "finish",
			Title:  "Færdig",
			Kind:   KindFinish,
			Body:   "Du er igennem alle trin. Assistenten kan åbne dashboardet med hurtige links til skolesystemerne, så du nemt finder dem igen.",
			URL:    URLDashboard,
			Button: "Åbn dashboard",
		},
```

- [ ] **Step 4: Kør alle tests — de skal bestå**

Run: `(cd setup-wizard && go test ./...)`
Expected: alle `ok`.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/steps/
git commit -m "feat: dashboard-URL på finish-trinnet"
```

---

### Task 4: Server-skelet + `GET /api/steps`

Første API-endpoint med trin-status i hukommelsen. `Server` implementerer `http.Handler`, så tests kalder den direkte uden netværk.

**Files:**
- Create: `setup-wizard/internal/server/state.go`
- Create: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv den fejlende test**

`setup-wizard/internal/server/server_test.go`:

```go
package server

import (
	"encoding/json"
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
```

- [ ] **Step 2: Kør testen — den skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — pakken kompilerer ikke (`undefined: Server`, `undefined: New`).

- [ ] **Step 3: Implementér state**

`setup-wizard/internal/server/state.go`:

```go
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
```

- [ ] **Step 4: Implementér serveren**

`setup-wizard/internal/server/server.go`:

```go
// Package server udstiller Assistentens lokale HTTP-API og serverer den
// indlejrede tringuide-side. API'et er tringuidens eneste vej til
// Go-processen; serveren bindes kun til localhost (se cmd/assistent).
package server

import (
	"encoding/json"
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

func New(os osops.OS) *Server {
	s := &Server{
		wiz:   wizard.New(os),
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
	_ = json.NewEncoder(w).Encode(v)
}
```

- [ ] **Step 5: Kør testen — den skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS (`TestStepsReturnsAllTenInOrder`).

- [ ] **Step 6: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: lokalt HTTP-API med GET /api/steps"
```

---

### Task 5: `POST /api/steps/{id}/done` og `/undo`

**Files:**
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv de fejlende tests**

Tilføj nederst i `server_test.go`:

```go
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
```

- [ ] **Step 2: Kør testene — de skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — `done: status = 404, forventede 204` (ruten findes ikke endnu).

- [ ] **Step 3: Implementér**

I `New()` i `server.go` tilføjes efter `GET /api/steps`-ruten:

```go
	s.mux.HandleFunc("POST /api/steps/{id}/done", s.handleSetDone(true))
	s.mux.HandleFunc("POST /api/steps/{id}/undo", s.handleSetDone(false))
```

Og nederst i filen:

```go
func (s *Server) handleSetDone(done bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.state.setDone(r.PathValue("id"), done) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
```

- [ ] **Step 4: Kør testene — de skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS for alle.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: markér trin færdigt og fortryd via API"
```

---

### Task 6: Wi-Fi-endpoints

`GET /api/wifi` klassificerer netværket via den eksisterende wizard-logik; `POST /api/wifi/settings` åbner systemindstillingerne.

**Files:**
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv de fejlende tests**

Tilføj nederst i `server_test.go` (og tilføj `"errors"` til importblokken):

```go
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
```

- [ ] **Step 2: Kør testene — de skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — 404 i stedet for 200/204 (ruterne findes ikke endnu).

- [ ] **Step 3: Implementér**

I `New()` tilføjes:

```go
	s.mux.HandleFunc("GET /api/wifi", s.handleWifi)
	s.mux.HandleFunc("POST /api/wifi/settings", s.handleWifiSettings)
```

Og nederst i `server.go`:

```go
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
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
```

- [ ] **Step 4: Kør testene — de skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS for alle.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: Wi-Fi-status og -indstillinger via API"
```

---

### Task 7: `POST /api/steps/{id}/open` — åbn officiel side

Sikkerhedskritisk detalje: URL'en slås op server-side i trinkonfigurationen. Browseren kan ikke sende en URL med, så API'et kan ikke misbruges til at åbne vilkårlige sider.

**Files:**
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv de fejlende tests**

Tilføj nederst i `server_test.go`:

```go
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
```

- [ ] **Step 2: Kør testene — de skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — 404 for `/api/steps/moodle/open`.

- [ ] **Step 3: Implementér**

I `New()` tilføjes:

```go
	s.mux.HandleFunc("POST /api/steps/{id}/open", s.handleOpen)
```

Og nederst i `server.go`:

```go
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
```

- [ ] **Step 4: Kør testene — de skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS for alle.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: åbn trinnets officielle side via API (URL slås op server-side)"
```

---

### Task 8: `POST /api/sketchup/install`

**Files:**
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv den fejlende test**

Tilføj nederst i `server_test.go`:

```go
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
		})
	}
}
```

- [ ] **Step 2: Kør testen — den skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — 404 for `/api/sketchup/install`.

- [ ] **Step 3: Implementér**

I `New()` tilføjes:

```go
	s.mux.HandleFunc("POST /api/sketchup/install", s.handleSketchUp)
```

Og nederst i `server.go`:

```go
// handleSketchUp kører installationsflowet. Fallback er et forventet
// udfald (S-mode, manglende winget, fejlet installation) — ikke en
// serverfejl — så svaret er altid 200 med action og evt. begrundelse.
func (s *Server) handleSketchUp(w http.ResponseWriter, r *http.Request) {
	outcome := s.wiz.InstallSketchUp()
	writeJSON(w, map[string]string{
		"action": string(outcome.Action),
		"reason": outcome.Reason,
	})
}
```

- [ ] **Step 4: Kør testene — de skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS for alle.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: SketchUp-installation via API med fallback-svar"
```

---

### Task 9: `POST /api/quit`

Serveren signalerer nedlukning via en kanal, som main venter på. `sync.Once` beskytter mod dobbelt-close hvis eleven trykker to gange.

**Files:**
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv den fejlende test**

Tilføj nederst i `server_test.go`:

```go
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
```

- [ ] **Step 2: Kør testen — den skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — kompileringsfejl `srv.Quit undefined`.

- [ ] **Step 3: Implementér**

I `server.go`: tilføj `"sync"` til importblokken, udvid `Server`-structen:

```go
type Server struct {
	wiz      *wizard.Wizard
	state    *state
	mux      *http.ServeMux
	quit     chan struct{}
	quitOnce sync.Once
}
```

I `New()` initialiseres kanalen og ruten registreres:

```go
	s := &Server{
		wiz:   wizard.New(os),
		state: newState(steps.All()),
		mux:   http.NewServeMux(),
		quit:  make(chan struct{}),
	}
```

```go
	s.mux.HandleFunc("POST /api/quit", s.handleQuit)
```

Og nederst i filen:

```go
// Quit lukkes når eleven afslutter Assistenten; main venter på kanalen.
func (s *Server) Quit() <-chan struct{} { return s.quit }

func (s *Server) handleQuit(w http.ResponseWriter, r *http.Request) {
	s.quitOnce.Do(func() { close(s.quit) })
	w.WriteHeader(http.StatusNoContent)
}
```

- [ ] **Step 4: Kør testene — de skal bestå**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: PASS for alle.

- [ ] **Step 5: Commit**

```bash
git add setup-wizard/internal/server/
git commit -m "feat: afslut Assistenten via API"
```

---

### Task 10: Indlejret tringuide-side + `GET /`

Siden indlejres med `go:embed` og serveres fra rod-URL'en. Indholdstests i samme ånd som de gamle PowerShell-tjek: sikkerhedsteksten skal være til stede, og ingen indlejret fil må indeholde ordet "setup-wizard".

**Files:**
- Create: `setup-wizard/internal/web/embed.go`
- Create: `setup-wizard/internal/web/static/index.html`
- Create: `setup-wizard/internal/web/static/style.css`
- Create: `setup-wizard/internal/web/static/app.js`
- Modify: `setup-wizard/internal/server/server.go`
- Test: `setup-wizard/internal/server/server_test.go`

- [ ] **Step 1: Skriv de fejlende tests**

Tilføj nederst i `server_test.go` (og tilføj `"io/fs"`, `"strings"` og web-pakken til importblokken):

```go
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
```

Importblokken i `server_test.go` skal nu samlet indeholde:

```go
import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/web"
)
```

- [ ] **Step 2: Kør testene — de skal fejle**

Run: `(cd setup-wizard && go test ./internal/server/ -v)`
Expected: FAIL — kompileringsfejl, `internal/web` findes ikke.

- [ ] **Step 3: Opret embed-pakken**

`setup-wizard/internal/web/embed.go`:

```go
// Package web indlejrer tringuide-siden i binæren, så Assistenten
// leveres som én fil uden løse assets.
package web

import "embed"

//go:embed static
var Static embed.FS
```

- [ ] **Step 4: Skriv tringuide-siden**

`setup-wizard/internal/web/static/index.html`:

```html
<!DOCTYPE html>
<html lang="da">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GF2 IT Setup — Assistenten</title>
  <link rel="stylesheet" href="/static/style.css">
</head>
<body>
  <header>
    <h1>GF2 IT Setup</h1>
    <p id="progress">Indlæser …</p>
  </header>
  <main>
    <nav id="step-list" aria-label="Trinliste"></nav>
    <section id="step-detail">
      <h2 id="step-title"></h2>
      <p id="step-body"></p>
      <div id="step-warning" class="warning" hidden></div>
      <div id="wifi-panel" hidden>
        <p id="wifi-status">Wi-Fi-status er ikke tjekket endnu.</p>
        <button id="wifi-check">Tjek igen</button>
        <button id="wifi-settings">Åbn Wi-Fi-indstillinger</button>
      </div>
      <p id="sketchup-result" hidden></p>
      <button id="action" hidden></button>
      <button id="quit" hidden>Afslut Assistenten</button>
      <div class="step-controls">
        <button id="toggle-done"></button>
      </div>
      <div class="nav-controls">
        <button id="prev">Tilbage</button>
        <button id="next">Frem</button>
      </div>
    </section>
  </main>
  <footer>Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login. Elever indtaster kun oplysninger på officielle sider og i Windows' egne indstillinger.</footer>
  <script src="/static/app.js"></script>
</body>
</html>
```

(Footer-teksten skal matche `steps.SafetyText` ordret — testen sammenligner med konstanten.)

- [ ] **Step 5: Skriv stilarket**

`setup-wizard/internal/web/static/style.css` (design-tokens lånt fra `start.html`):

```css
:root {
  --ink: #181818;
  --muted: #5c6670;
  --line: #dce4ea;
  --paper: #f4f7f9;
  --surface: #ffffff;
  --neg-blue: #005aa7;
  --neg-blue-dark: #003f73;
  --neg-blue-soft: #e7f2fb;
  --amber: #efb43f;
  --green: #237a4b;
}

* { box-sizing: border-box; }

body {
  margin: 0;
  font-family: Arial, Helvetica, sans-serif;
  color: var(--ink);
  background: var(--paper);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

header {
  background: var(--neg-blue);
  color: #fff;
  padding: 16px 24px;
}

header h1 { margin: 0; font-size: 1.4rem; }

#progress { margin: 4px 0 0; color: var(--neg-blue-soft); }

main {
  flex: 1;
  display: flex;
  gap: 24px;
  padding: 24px;
  max-width: 960px;
  width: 100%;
  margin: 0 auto;
}

#step-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 260px;
  flex-shrink: 0;
}

.step-item {
  text-align: left;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--surface);
  cursor: pointer;
  font: inherit;
  font-size: 0.95rem;
}

.step-item.active {
  border-color: var(--neg-blue);
  background: var(--neg-blue-soft);
  font-weight: bold;
}

#step-detail {
  flex: 1;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 24px;
}

.warning {
  background: #fdf3df;
  border: 1px solid var(--amber);
  border-radius: 8px;
  padding: 12px;
  margin: 12px 0;
}

#action, #toggle-done, #wifi-check, #wifi-settings, #quit {
  font: inherit;
  background: var(--neg-blue);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 16px;
  margin: 8px 8px 0 0;
  cursor: pointer;
}

#action:hover, #wifi-check:hover, #wifi-settings:hover, #quit:hover {
  background: var(--neg-blue-dark);
}

#toggle-done { background: var(--green); }

.nav-controls {
  margin-top: 24px;
  display: flex;
  justify-content: space-between;
}

.nav-controls button {
  font: inherit;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 8px 16px;
  cursor: pointer;
}

.nav-controls button:disabled { opacity: 0.4; cursor: default; }

footer {
  background: var(--surface);
  border-top: 1px solid var(--line);
  color: var(--muted);
  padding: 12px 24px;
  font-size: 0.85rem;
  text-align: center;
}
```

- [ ] **Step 6: Skriv klienten**

`setup-wizard/internal/web/static/app.js`:

```js
// Tringuidens klient. Al tilstand bor i Go-processen; siden henter
// trinlisten efter hver handling og tegner forfra.
const api = {
  steps: () => fetch("/api/steps").then(r => r.json()),
  setDone: (id, done) => fetch(`/api/steps/${id}/${done ? "done" : "undo"}`, { method: "POST" }),
  open: id => fetch(`/api/steps/${id}/open`, { method: "POST" }),
  wifi: () => fetch("/api/wifi").then(r => (r.ok ? r.json() : null)),
  wifiSettings: () => fetch("/api/wifi/settings", { method: "POST" }),
  sketchup: () => fetch("/api/sketchup/install", { method: "POST" }).then(r => r.json()),
  quit: () => fetch("/api/quit", { method: "POST" }),
};

let allSteps = [];
let current = 0;

async function refresh() {
  allSteps = (await api.steps()).steps;
  render();
}

function render() {
  const step = allSteps[current];
  document.getElementById("progress").textContent = `Trin ${current + 1} af ${allSteps.length}`;

  const list = document.getElementById("step-list");
  list.innerHTML = "";
  allSteps.forEach((s, i) => {
    const item = document.createElement("button");
    item.className = "step-item" + (i === current ? " active" : "");
    item.textContent = `${i + 1}. ${s.title}` + (s.done ? " ✓" : "");
    item.addEventListener("click", () => { current = i; render(); });
    list.appendChild(item);
  });

  document.getElementById("step-title").textContent = step.title;
  document.getElementById("step-body").textContent = step.body;

  const warning = document.getElementById("step-warning");
  warning.hidden = !step.warning;
  warning.textContent = step.warning || "";

  document.getElementById("wifi-panel").hidden = step.kind !== "wifi";
  document.getElementById("sketchup-result").hidden = true;
  document.getElementById("quit").hidden = step.kind !== "finish";

  const action = document.getElementById("action");
  const hasAction = step.kind === "link" || step.kind === "sketchup" || step.kind === "finish";
  action.hidden = !hasAction;
  action.textContent = step.button;

  document.getElementById("toggle-done").textContent = step.done ? "Fortryd" : "Markér som færdig";

  document.getElementById("prev").disabled = current === 0;
  document.getElementById("next").disabled = current === allSteps.length - 1;

  if (step.kind === "wifi") {
    checkWifi();
  }
}

async function checkWifi() {
  const el = document.getElementById("wifi-status");
  const status = await api.wifi();
  if (!status) {
    el.textContent = "Wi-Fi-status kunne ikke aflæses. Tjek selv i dine netværksindstillinger.";
    return;
  }
  const texts = {
    target: `Du er på ${status.ssid} — det rigtige netværk. Trinnet er klaret!`,
    guest: `Du er på ${status.ssid}. Det er kun til midlertidig gæsteadgang — skift til NEG.`,
    other: `Du er på "${status.ssid}", ikke NEG. Åbn Wi-Fi-indstillinger og skift til NEG.`,
    none: "Du er ikke på et Wi-Fi-netværk. Åbn Wi-Fi-indstillinger og vælg NEG.",
  };
  el.textContent = texts[status.state];
}

document.getElementById("toggle-done").addEventListener("click", async () => {
  const step = allSteps[current];
  await api.setDone(step.id, !step.done);
  await refresh();
});

document.getElementById("prev").addEventListener("click", () => { current--; render(); });
document.getElementById("next").addEventListener("click", () => { current++; render(); });

document.getElementById("action").addEventListener("click", async () => {
  const step = allSteps[current];
  if (step.kind === "sketchup") {
    const result = document.getElementById("sketchup-result");
    result.hidden = false;
    result.textContent = "Installerer … det kan tage nogle minutter.";
    const outcome = await api.sketchup();
    result.textContent = outcome.action === "installed"
      ? "SketchUp er installeret."
      : outcome.reason;
  } else {
    await api.open(step.id);
  }
});

document.getElementById("wifi-check").addEventListener("click", checkWifi);
document.getElementById("wifi-settings").addEventListener("click", () => api.wifiSettings());

document.getElementById("quit").addEventListener("click", async () => {
  await api.quit();
  document.body.innerHTML = "<p style='padding:40px;text-align:center'>Assistenten er lukket. Du kan lukke denne fane.</p>";
});

refresh();
```

- [ ] **Step 7: Servér siden fra serveren**

I `server.go`: tilføj web-pakken til importblokken:

```go
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/web"
```

I `New()` tilføjes (efter API-ruterne):

```go
	s.mux.Handle("GET /static/", http.FileServerFS(web.Static))
	s.mux.HandleFunc("GET /", s.handleIndex)
```

Og nederst i filen:

```go
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFileFS(w, r, web.Static, "static/index.html")
}
```

- [ ] **Step 8: Kør alle tests — de skal bestå**

Run: `(cd setup-wizard && go test ./...)`
Expected: alle `ok`.

- [ ] **Step 9: Commit**

```bash
git add setup-wizard/internal/web/ setup-wizard/internal/server/
git commit -m "feat: indlejret dansk tringuide-side serveret fra binæren"
```

---

### Task 11: `main.go`, README og krydskompilering

Den tynde main: ledig port på localhost, server op, åbn browser, vent på quit. Verifikation: begge elev-binærer bygger, og Assistenten kan startes manuelt på WSL.

**Files:**
- Create: `setup-wizard/cmd/assistent/main.go`
- Create: `setup-wizard/README.md`
- Modify: `.gitignore` (opret hvis den ikke findes)

- [ ] **Step 1: Skriv main**

`setup-wizard/cmd/assistent/main.go`:

```go
// Assistenten: starter det lokale API, serverer tringuiden og åbner
// elevens browser. Al logik ligger i internal/server og er testet dér.
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/server"
)

func main() {
	osImpl := osops.Current()
	srv := server.New(osImpl)

	// Port 0: OS'et vælger en ledig port, så Assistenten aldrig
	// kolliderer med noget andet på elevens maskine.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("kunne ikke starte Assistenten: %v", err)
	}

	go func() {
		if err := http.Serve(ln, srv); err != nil {
			log.Fatalf("serveren stoppede uventet: %v", err)
		}
	}()

	url := "http://" + ln.Addr().String()
	log.Printf("Assistenten kører på %s", url)
	if err := osImpl.OpenURL(url); err != nil {
		log.Printf("kunne ikke åbne browseren automatisk — åbn selv %s", url)
	}

	<-srv.Quit()
}
```

- [ ] **Step 2: Skriv README**

`setup-wizard/README.md`:

```markdown
# Assistenten (setup-wizard)

Go-koden bag GF2 IT Setup-Assistenten: én kodebase, der kompileres til en
Windows-`.exe` og en Mac-binær. Binæren starter en lokal webserver på
localhost, åbner elevens browser og viser den danske tringuide.

## Tests

    go test ./...

Alle tests kører uden at røre det rigtige OS (OS-laget fakes via
`internal/osops/osfake`).

## Byg elev-binærerne

    GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o dist/Assistenten.exe ./cmd/assistent
    GOOS=darwin  GOARCH=arm64 go build -o dist/Assistenten ./cmd/assistent

`-H=windowsgui` forhindrer et sort konsolvindue bag browseren, når eleven
dobbeltklikker på .exe'en. Ingen CGO og ingen dependencies — bygget virker
fra enhver maskine med Go installeret.

## Kør lokalt under udvikling (WSL/Linux)

    go run ./cmd/assistent

Linux-implementeringen er en udviklerstub: Wi-Fi-status svarer "intet
netværk", og SketchUp-installation går altid til fallback. Hvis browseren
ikke åbner automatisk, står URL'en i terminalen.
```

- [ ] **Step 3: Ignorér dist-mappen**

Tilføj til `.gitignore` i repo-roden (opret filen hvis den ikke findes):

```
setup-wizard/dist/
```

- [ ] **Step 4: Kør alle tests og vet**

Run: `(cd setup-wizard && go vet ./... && go test ./...)`
Expected: ingen vet-fejl, alle tests `ok`.

- [ ] **Step 5: Byg begge elev-binærer**

Run:
```bash
(cd setup-wizard && GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o dist/Assistenten.exe ./cmd/assistent)
(cd setup-wizard && GOOS=darwin GOARCH=arm64 go build -o dist/Assistenten ./cmd/assistent)
ls -la setup-wizard/dist/
```
Expected: `Assistenten.exe` og `Assistenten` findes i `setup-wizard/dist/` (typisk 7–10 MB hver).

- [ ] **Step 6: Manuel røgtest på WSL**

Run: `(cd setup-wizard && go run ./cmd/assistent &) ; sleep 1`

Find URL'en i outputtet (eller test med curl hvis browseren ikke åbner):

```bash
# erstat PORT med porten fra log-outputtet
curl -s http://127.0.0.1:PORT/api/steps | head -c 300
curl -s -X POST http://127.0.0.1:PORT/api/quit -o /dev/null -w "%{http_code}\n"
```
Expected: JSON med trinlisten; quit svarer 204, og processen afslutter.

- [ ] **Step 7: Commit**

```bash
git add setup-wizard/cmd/ setup-wizard/README.md .gitignore
git commit -m "feat: kørbar Assistent-binær — main, README og krydskompilering"
```

---

## Færdigkriterier (fra spec'en)

- [ ] Alle tests består: `(cd setup-wizard && go test ./...)`
- [ ] `GOOS=windows go build ./...` og `GOOS=darwin go build ./...` lykkes begge
- [ ] `dist/Assistenten.exe` bygges med `-H=windowsgui` (intet konsolvindue)
- [ ] `GET /` serverer tringuiden med sikkerhedsteksten; ingen indlejret fil indeholder "setup-wizard"
- [ ] Manuel røgtest: binæren starter, API'et svarer, `/api/quit` lukker processen
