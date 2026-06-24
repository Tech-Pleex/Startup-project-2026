# Setup-wizard UI-redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Giv Assistentens web-UI et brand-locked "blueprint → virkelighed"-redesign og indfør en tredje trintilstand ("Spring over", issue #23).

**Architecture:** UI'et er en statisk web-side (`setup-wizard/internal/web/static/`) serveret af Go-serveren og indlejret via `embed.FS`. Backenden holder al trin-status i hukommelsen og udstiller et localhost-HTTP-API. Vi udvider status-modellen fra to til tre tilstande, tilføjer et `skip`-endpoint, og omskriver de tre frontend-filer + lokale assets. Frontend forbliver dependency-frit vanilla JS/CSS.

**Tech Stack:** Go 1.22+ (net/http med metode-mønstre i ServeMux), vanilla HTML/CSS/JS, Space Grotesk (lokal woff2).

## Global Constraints

- **Brand (hård regel):** Kun NEG-paletten — mørkeblå `#123c62`, blå `#86afd8`, lysegrå `#efefef`, orange `#ec8113` (+ tints). Orange **kun** til CTA/links + aktiv-accent. Font: **Space Grotesk** (bundlet lokalt). Officielt logo, ikke gen-tegnet.
- **Offline:** Assistenten kører før Wi-Fi. Ingen CDN/eksterne kald — alle fonte/billeder lokale i `static/`.
- **Bevidst undtagelse:** S-mode-blokering vises i system-rød (uden for paletten) som en Windows/maskine-fejl.
- **Bevar disse strenge i `index.html`** (eksisterende Go-tests kræver dem ordret): `steps.SafetyText`-værdien, ordet `Assistenten`, `Windows S-mode er aktiveret`, `Assistenten kan ikke fortsætte`, `Tjek igen`.
- **Ingen static-fil må indeholde strengen `setup-wizard`** (testet af `TestEmbeddedAssetsContainNoInternalNames`).
- **Filnavne bevares:** `static/style.css` og `static/app.js` skal fortsat svare 200.
- **Arbejdsmappe-rod:** `C:\Users\jere\Documents\Neg_Ai_Stuff\Startup-project-2026`. Go-modulet ligger i `setup-wizard/`.
- **Go-test-kommando (kør fra `setup-wizard/`):** `go test ./...`

---

## File Structure

- `setup-wizard/internal/server/state.go` — **modify**: tre-tilstands-status (pending/done/skipped) i stedet for bool-map; `stepView` får `Skipped`-felt; `setDone` → `setStatus`.
- `setup-wizard/internal/server/server.go` — **modify**: rutér done/undo/skip via fælles `handleSetStatus`; nyt `POST /api/steps/{id}/skip`.
- `setup-wizard/internal/server/server_test.go` — **modify**: tilføj `Skipped` til test-structen + nye skip-tests.
- `setup-wizard/internal/web/static/fonts/space-grotesk.woff2` — **create**: bundlet brandfont.
- `setup-wizard/internal/web/static/img/neg-hero.jpg` — **create**: optimeret hero.
- `setup-wizard/internal/web/static/img/neg-logo-white.png` — **create**: officielt logo (hvid).
- `setup-wizard/internal/web/static/index.html` — **modify (rewrite)**: ny struktur (topbar/scene/rail/kort/banner).
- `setup-wizard/internal/web/static/style.css` — **modify (rewrite)**: brand-tokens + komponenter.
- `setup-wizard/internal/web/static/app.js` — **modify (rewrite)**: status-bevidst, skip, scene-fremdrift, rail, tastatur.

---

### Task 1: Backend — tre trintilstande + skip-endpoint

**Files:**
- Modify: `setup-wizard/internal/server/state.go`
- Modify: `setup-wizard/internal/server/server.go:36-39` (rute-registrering) og `:63-71` (handler)
- Test: `setup-wizard/internal/server/server_test.go`

**Interfaces:**
- Produces:
  - `(*state).setStatus(id, status string) bool` — `status` ∈ `{"", "done", "skipped"}`; `""` = pending; returnerer `false` hvis id ukendt.
  - `stepView.Skipped bool` (JSON `"skipped"`), `stepView.Done bool` (JSON `"done"`, uændret).
  - Endpoints: `POST /api/steps/{id}/done`, `/undo`, `/skip` → 204; ukendt id → 404.

- [ ] **Step 1: Tilføj `Skipped` til test-structen og skriv de fejlende skip-tests**

I `server_test.go`, tilføj feltet til `stepJSON` (efter `Done`):

```go
	Done    bool `json:"done"`
	Skipped bool `json:"skipped"`
```

Tilføj disse tests sidst i filen:

```go
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
	do(t, srv, http.MethodPost, "/api/steps/sketchup/skip")
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
```

- [ ] **Step 2: Kør testene og bekræft at de fejler**

Run (fra `setup-wizard/`): `go test ./internal/server/ -run 'Skip' -v`
Expected: FAIL — `/api/steps/sketchup/skip` svarer 404 (ruten findes ikke endnu), og `Skipped` er altid false.

- [ ] **Step 3: Omskriv `state.go` til tre tilstande**

Erstat hele indholdet af `state.go` med:

```go
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

- [ ] **Step 4: Opdatér `server.go` — fælles statushandler + skip-rute**

I `server.go`, erstat de tre rute-linjer (`POST /api/steps/{id}/done` og `/undo`) i `New` med:

```go
	s.mux.HandleFunc("POST /api/steps/{id}/done", s.handleSetStatus("done"))
	s.mux.HandleFunc("POST /api/steps/{id}/undo", s.handleSetStatus(""))
	s.mux.HandleFunc("POST /api/steps/{id}/skip", s.handleSetStatus("skipped"))
```

Erstat hele `handleSetDone`-funktionen med:

```go
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
```

- [ ] **Step 5: Kør hele server-testpakken og bekræft grøn**

Run (fra `setup-wizard/`): `go test ./internal/server/ -v`
Expected: PASS — inkl. de nye Skip-tests og de eksisterende (`TestMarkStepDoneAndUndo`, `TestStepsReturnsAllTenInOrder` m.fl.).

- [ ] **Step 6: Commit**

```bash
git add setup-wizard/internal/server/state.go setup-wizard/internal/server/server.go setup-wizard/internal/server/server_test.go
git commit -m "feat: tilføj 'sprunget over'-tilstand og skip-endpoint (issue #23)"
```

---

### Task 2: Bundl lokale assets (font, hero, logo)

**Files:**
- Create: `setup-wizard/internal/web/static/fonts/space-grotesk.woff2`
- Create: `setup-wizard/internal/web/static/img/neg-hero.jpg`
- Create: `setup-wizard/internal/web/static/img/neg-logo-white.png`
- Test: `setup-wizard/internal/server/server_test.go`

**Interfaces:**
- Produces: assets serveres på `/static/fonts/space-grotesk.woff2`, `/static/img/neg-hero.jpg`, `/static/img/neg-logo-white.png` (alle 200). `embed.FS` (`//go:embed static`) indlejrer undermapper automatisk — ingen Go-ændring nødvendig.

- [ ] **Step 1: Skriv den fejlende asset-test**

I `server_test.go`, tilføj:

```go
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
```

- [ ] **Step 2: Kør testen og bekræft at den fejler**

Run (fra `setup-wizard/`): `go test ./internal/server/ -run TestBrandAssetsAreServed -v`
Expected: FAIL — filerne findes ikke endnu (404).

- [ ] **Step 3: Hent Space Grotesk (variabel woff2) lokalt**

Opret mappen og hent fonten (latin, vægt 300–700). Fra arbejdsmappe-roden:

```bash
mkdir -p setup-wizard/internal/web/static/fonts
curl -L -o setup-wizard/internal/web/static/fonts/space-grotesk.woff2 \
  "https://cdn.jsdelivr.net/fontsource/fonts/space-grotesk:vf@latest/latin-wght-normal.woff2"
```

Bekræft at filen er > 10 KB: `ls -l setup-wizard/internal/web/static/fonts/space-grotesk.woff2`
(Hvis URL'en ikke virker: hent `space-grotesk` woff2 fra https://gwfh.mranftl.com/fonts/space-grotesk og omdøb til `space-grotesk.woff2`.)

- [ ] **Step 4: Optimér hero og kopiér logo**

Fra arbejdsmappe-roden (kræver ImageMagick `magick`):

```bash
mkdir -p setup-wizard/internal/web/static/img
magick "assets/neg-hero-transition.png" -resize 1600x -quality 82 \
  setup-wizard/internal/web/static/img/neg-hero.jpg
magick "NEG BRAND/2023 - NEG logo white.png" -trim +repage -resize 260x \
  setup-wizard/internal/web/static/img/neg-logo-white.png
```

Hvis `magick` ikke er installeret: skalér `assets/neg-hero-transition.png` til ~1600px bredde og gem som `neg-hero.jpg` (kvalitet ~82), og kopiér logoet manuelt til `neg-logo-white.png`. Verificér begge filer findes: `ls -l setup-wizard/internal/web/static/img/`

- [ ] **Step 5: Kør asset-testen + hele pakken og bekræft grøn**

Run (fra `setup-wizard/`): `go test ./...`
Expected: PASS — inkl. `TestBrandAssetsAreServed` og `TestEmbeddedAssetsContainNoInternalNames` (binære assets indeholder ikke "setup-wizard").

- [ ] **Step 6: Commit**

```bash
git add setup-wizard/internal/web/static/fonts setup-wizard/internal/web/static/img setup-wizard/internal/server/server_test.go
git commit -m "feat: bundl Space Grotesk, optimeret hero og NEG-logo i static"
```

---

### Task 3: Omskriv `index.html` til ny struktur

**Files:**
- Modify (rewrite): `setup-wizard/internal/web/static/index.html`

**Interfaces:**
- Produces (element-id'er som `app.js` i Task 5 forbruger): `progress`, `progress-bar`, `scene-full`, `scene-blueprint`, `step-rail`, `step-kicker`, `step-title`, `step-body`, `smode-warning`, `smode-retry`, `step-warning`, `wifi-panel`, `wifi-settings`, `action`, `toggle-done`, `skip`, `prev`, `next`, `quit`.
- Bevarer ordrette strenge krævet af Go-tests (se Global Constraints).

- [ ] **Step 1: Erstat hele `index.html`**

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
  <header class="topbar">
    <div class="brand">
      <img src="/static/img/neg-logo-white.png" alt="NEG" class="logo">
      <span class="brand-tag">IT-Opsætning · Assistenten</span>
    </div>
    <div class="progress-wrap">
      <span id="progress" class="progress-label">Indlæser …</span>
      <span class="bar"><i id="progress-bar"></i></span>
    </div>
  </header>

  <section class="scene" aria-hidden="true">
    <img id="scene-full" class="scene-img" src="/static/img/neg-hero.jpg" alt="">
    <img id="scene-blueprint" class="scene-img scene-bp" src="/static/img/neg-hero.jpg" alt="">
    <div class="scene-grid"></div>
    <div class="scene-cap">Blueprint → virkelighed · scenen bygges trin for trin</div>
  </section>

  <nav id="step-rail" class="rail" aria-label="Trinliste"></nav>

  <main class="card">
    <p id="step-kicker" class="kicker"></p>
    <h2 id="step-title"></h2>
    <p id="step-body"></p>

    <div id="smode-warning" class="smode-warning" role="alert" hidden>
      <h3>Windows S-mode er aktiveret</h3>
      <p>Assistenten kan ikke fortsætte, mens computeren kører i S-mode. Kontakt din lærer eller IT, slå S-mode fra, og prøv igen.</p>
      <button id="smode-retry">Tjek igen</button>
    </div>

    <div id="step-warning" class="warning" hidden></div>

    <div id="wifi-panel" hidden>
      <button id="wifi-settings" class="btn btn-done">Åbn Wi-Fi-indstillinger</button>
    </div>

    <div class="actions">
      <button id="action" class="btn btn-primary" hidden></button>
      <button id="toggle-done" class="btn btn-done"></button>
      <button id="skip" class="btn btn-skip"></button>
    </div>

    <button id="quit" class="btn btn-primary" hidden>Afslut Assistenten</button>

    <div class="nav-controls">
      <button id="prev">‹ Tilbage</button>
      <button id="next">Frem ›</button>
    </div>
  </main>

  <footer class="safety">
    <span class="shield" aria-hidden="true">🛡️</span>
    <p>Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login. Elever indtaster kun oplysninger på officielle sider og i Windows' egne indstillinger.</p>
  </footer>

  <script src="/static/app.js"></script>
</body>
</html>
```

- [ ] **Step 2: Bekræft at de krævede strenge stadig serveres**

Run (fra `setup-wizard/`): `go test ./internal/server/ -run 'TestIndexServesTringuideWithSafetyText|TestEmbeddedAssetsContainNoInternalNames' -v`
Expected: PASS — siden indeholder `steps.SafetyText`, "Assistenten", "Windows S-mode er aktiveret", "Assistenten kan ikke fortsætte", "Tjek igen", og ingen "setup-wizard".

- [ ] **Step 3: Commit**

```bash
git add setup-wizard/internal/web/static/index.html
git commit -m "feat: ny index.html-struktur (topbar/scene/rail/kort/banner)"
```

---

### Task 4: Omskriv `style.css` (brand-tokens + komponenter)

**Files:**
- Modify (rewrite): `setup-wizard/internal/web/static/style.css`

**Interfaces:**
- Consumes klasser/id'er fra `index.html` (Task 3).
- Produces visuelle tilstande som `app.js` (Task 5) skifter via klasser: `.tile`, `.tile.rendered`, `.tile.skipped`, `.tile.is-current`; samt inline `opacity` på `#scene-blueprint` og `width` på `#progress-bar`.

- [ ] **Step 1: Erstat hele `style.css`**

```css
/* NEG brand-locked · Space Grotesk (lokal) · blueprint → virkelighed */
@font-face {
  font-family: "Space Grotesk";
  src: url("/static/fonts/space-grotesk.woff2") format("woff2");
  font-weight: 300 700;
  font-display: swap;
}

:root {
  --navy: #123c62;
  --navy-deep: #0d2f4f;
  --blue: #86afd8;
  --grey: #efefef;
  --orange: #ec8113;
  --ink: #eaf1f8;
  --muted: #9dbdde;
  --line: #2b5680;
}

* { box-sizing: border-box; }

body {
  margin: 0;
  font-family: "Space Grotesk", system-ui, Arial, sans-serif;
  font-weight: 400;
  color: var(--ink);
  background: var(--navy-deep);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.lbl, .brand-tag, .kicker, .progress-label, .scene-cap {
  font-weight: 500;
  letter-spacing: .14em;
  text-transform: uppercase;
}

/* topbar */
.topbar {
  display: flex; align-items: center; justify-content: space-between;
  gap: 16px; padding: 14px 22px;
  background: linear-gradient(#15426c, var(--navy));
  border-bottom: 1px solid var(--line);
}
.brand { display: flex; align-items: center; gap: 14px; }
.brand .logo { height: 28px; display: block; }
.brand-tag { color: var(--blue); font-size: 11px; }
.progress-wrap { display: flex; align-items: center; gap: 11px; }
.progress-label { color: var(--blue); font-size: 11px; }
.bar { width: 140px; height: 8px; border-radius: 5px; background: #0c2840;
  overflow: hidden; border: 1px solid var(--line); }
.bar > i { display: block; height: 100%; width: 0; background: var(--blue);
  transition: width .5s ease; }

/* scene: hero med blueprint-krydstoning */
.scene { position: relative; height: 190px; overflow: hidden;
  background: var(--navy-deep); border-bottom: 1px solid var(--line); }
.scene-img { position: absolute; inset: 0; width: 100%; height: 100%;
  object-fit: cover; object-position: center; }
.scene-bp {
  filter: brightness(.82) grayscale(.45) sepia(1) hue-rotate(172deg) saturate(2.7);
  opacity: 1; transition: opacity .5s ease;
}
.scene-grid { position: absolute; inset: 0;
  background-image:
    linear-gradient(#86afd816 1px, transparent 1px),
    linear-gradient(90deg, #86afd816 1px, transparent 1px);
  background-size: 22px 22px; }
.scene-cap { position: absolute; left: 0; right: 0; bottom: 0;
  padding: 18px 16px 9px; font-size: 10px; color: #cfe0f2;
  background: linear-gradient(transparent, #0d2f4fcc); }

/* flise-rail */
.rail { display: grid; grid-template-columns: repeat(4, 1fr); gap: 9px;
  padding: 14px 18px; background: #0b2840; border-bottom: 1px solid var(--line); }
.tile { display: flex; align-items: center; gap: 8px; padding: 9px 11px;
  border-radius: 9px; font: inherit; font-size: 12px; font-weight: 500;
  text-align: left; cursor: pointer;
  border: 1px dashed #4d7eaa; color: var(--blue); background: #0e305088; }
.tile .ic { width: 9px; height: 9px; border-radius: 3px;
  border: 1px dashed var(--blue); flex: none; }
.tile.rendered { border: 1px solid var(--blue); color: var(--navy-deep);
  background: var(--blue); font-weight: 700; }
.tile.rendered .ic { background: var(--navy-deep); border: none; }
.tile.skipped { border: 1px solid var(--line); color: #7c9bbb; background: #0c2236aa; }
.tile.skipped .ic { border: none; }
.tile.skipped .ic::before { content: "↷"; color: #7c9bbb; font-size: 13px; line-height: 1; }
.tile.is-current { box-shadow: 0 0 0 2px var(--orange); }
.tile:disabled { opacity: .5; cursor: default; }

/* trin-kort */
.card { flex: 1; padding: 22px; max-width: 880px; width: 100%; margin: 0 auto; }
.kicker { font-size: 11px; color: var(--blue); margin: 0; }
.card h2 { margin: 8px 0 10px; font-size: 27px; font-weight: 700; letter-spacing: -.01em; }
.card > p { margin: 0; color: var(--muted); line-height: 1.55; font-size: 16px; }

.warning { margin-top: 14px; padding: 12px 14px; border-radius: 9px; font-size: 14px;
  background: #11334f; border: 1px solid var(--orange);
  border-left: 5px solid var(--orange); color: #ffd9a6; }

/* S-mode = system/Windows-fejl, bevidst rød uden for paletten */
.smode-warning { margin: 16px 0; padding: 16px; border-radius: 10px;
  background: #3a1212; border: 1px solid #a33; border-left: 5px solid #d34; color: #ffd4d4; }
.smode-warning h3 { margin: 0 0 8px; color: #ffecec; }
#smode-retry { font: inherit; background: #fff; color: #8b1d1d;
  border: 1px solid #c62828; border-radius: 8px; padding: 8px 14px; cursor: pointer; }

/* knapper */
.actions { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-top: 20px; }
.btn { font: inherit; font-size: 14px; font-weight: 700; border: none; border-radius: 10px;
  padding: 14px 10px; cursor: pointer; text-align: center; }
.btn-primary { background: var(--orange); color: #2a1804; }
.btn-done { background: var(--blue); color: var(--navy-deep); }
.btn-skip { background: transparent; color: var(--ink); border: 1.5px solid var(--line); }
.btn:disabled { opacity: .45; cursor: default; }
#wifi-panel { margin-top: 16px; }
#quit { margin-top: 16px; }

.nav-controls { margin-top: 16px; display: flex; justify-content: space-between; }
.nav-controls button { font: inherit; font-size: 14px; background: transparent;
  border: 1px solid var(--line); color: var(--ink); border-radius: 10px;
  padding: 11px 18px; cursor: pointer; }
.nav-controls button:disabled { opacity: .4; cursor: default; }

/* sikkerheds-banner */
.safety { display: flex; align-items: center; gap: 14px; padding: 15px 22px;
  background: #0e3050; border-top: 1px solid var(--line); border-left: 5px solid var(--blue); }
.safety .shield { font-size: 22px; flex: none; }
.safety p { margin: 0; color: var(--muted); font-size: 13.5px; line-height: 1.5; }

@media (prefers-reduced-motion: reduce) {
  .bar > i, .scene-bp { transition: none; }
}
```

- [ ] **Step 2: Bekræft at CSS'en stadig serveres + ingen interne navne**

Run (fra `setup-wizard/`): `go test ./internal/server/ -run 'TestStaticAssetsAreServed|TestEmbeddedAssetsContainNoInternalNames' -v`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add setup-wizard/internal/web/static/style.css
git commit -m "feat: brand-locked stylesheet (NEG-palet, scene, rail, knapper)"
```

---

### Task 5: Omskriv `app.js` (status-bevidst, skip, scene-fremdrift)

**Files:**
- Modify (rewrite): `setup-wizard/internal/web/static/app.js`

**Interfaces:**
- Consumes element-id'er fra Task 3 og API-endpoints fra Task 1 (`/done`, `/undo`, `/skip`).
- Produces: ren funktion `sceneProgress()` (andel done blandt ikke-finish-trin); klasse-skift på fliser; opdatering af scene-opacity og fremdriftsbjælke.

- [ ] **Step 1: Erstat hele `app.js`**

```js
// Tringuidens klient. Al tilstand bor i Go-processen; siden henter
// trinlisten efter hver handling og tegner forfra.

async function postOK(path) {
  const r = await fetch(path, { method: "POST" });
  if (!r.ok) throw new Error(`${path} svarede ${r.status}`);
  return r;
}

const api = {
  steps: () => fetch("/api/steps").then(r => r.json()),
  system: () => fetch("/api/system").then(r => (r.ok ? r.json() : null)),
  done: id => postOK(`/api/steps/${id}/done`),
  undo: id => postOK(`/api/steps/${id}/undo`),
  skip: id => postOK(`/api/steps/${id}/skip`),
  open: id => postOK(`/api/steps/${id}/open`),
  wifiSettings: () => postOK("/api/wifi/settings"),
  quit: () => postOK("/api/quit"),
};

let allSteps = [];
let current = 0;
let sModeBlocked = false;

const $ = id => document.getElementById(id);

async function refresh() {
  allSteps = (await api.steps()).steps;
  render();
}

// Fremdrift mod virkelighed: andel af ikke-finish-trin der er markeret done.
function sceneProgress() {
  const renderable = allSteps.filter(s => s.kind !== "finish");
  if (renderable.length === 0) return 0;
  return renderable.filter(s => s.done).length / renderable.length;
}

function renderRail() {
  const rail = $("step-rail");
  rail.innerHTML = "";
  allSteps.forEach((s, i) => {
    if (s.id === "welcome" || s.kind === "finish") return;
    const tile = document.createElement("button");
    let cls = "tile";
    if (s.done) cls += " rendered";
    else if (s.skipped) cls += " skipped";
    if (i === current) cls += " is-current";
    tile.className = cls;
    tile.disabled = sModeBlocked;
    tile.innerHTML = `<span class="ic"></span>${s.title}`;
    tile.addEventListener("click", () => { current = i; render(); });
    rail.appendChild(tile);
  });
}

function render() {
  if (allSteps.length === 0) {
    $("progress").textContent = "Ingen trin at vise.";
    return;
  }
  const step = allSteps[current];
  const progress = sceneProgress();

  $("progress").textContent = `Trin ${current + 1} / ${allSteps.length}`;
  $("progress-bar").style.width = `${Math.round(progress * 100)}%`;
  $("scene-blueprint").style.opacity = String(1 - progress);

  renderRail();

  $("step-kicker").textContent = `Trin ${current + 1} / ${allSteps.length} · ${step.kind}`;
  $("step-title").textContent = step.title;
  $("step-body").textContent = step.body;

  const welcomeBlocked = sModeBlocked && current === 0;
  $("smode-warning").hidden = !welcomeBlocked;

  const warning = $("step-warning");
  warning.hidden = !step.warning;
  warning.textContent = step.warning || "";

  $("wifi-panel").hidden = step.kind !== "wifi";
  $("quit").hidden = step.kind !== "finish";

  const action = $("action");
  const hasAction = step.kind === "link" || step.kind === "finish";
  action.hidden = !hasAction;
  action.textContent = step.button;

  const toggleDone = $("toggle-done");
  toggleDone.textContent = step.done ? "Fortryd" : "Markér som færdig";
  toggleDone.disabled = sModeBlocked;

  const skip = $("skip");
  skip.hidden = step.kind === "finish";
  skip.textContent = step.skipped ? "Fortryd spring over" : "Spring over";
  skip.disabled = sModeBlocked;

  $("prev").disabled = current === 0;
  $("next").disabled = current === allSteps.length - 1 || sModeBlocked;
}

async function checkSystem() {
  let status = null;
  try { status = await api.system(); } catch (err) { /* ukendt status må ikke blokere */ }
  sModeBlocked = Boolean(status?.sMode);
  if (sModeBlocked) current = 0;
  render();
}

async function start() {
  allSteps = (await api.steps()).steps;
  await checkSystem();
}

$("toggle-done").addEventListener("click", async () => {
  const step = allSteps[current];
  try {
    await (step.done ? api.undo(step.id) : api.done(step.id));
    await refresh();
  } catch (err) { alert("Kunne ikke opdatere trinnet. Prøv igen."); }
});

$("skip").addEventListener("click", async () => {
  const step = allSteps[current];
  const wasSkipped = step.skipped;
  try {
    await (wasSkipped ? api.undo(step.id) : api.skip(step.id));
    await refresh();
    if (!wasSkipped && current < allSteps.length - 1) { current++; render(); }
  } catch (err) { alert("Kunne ikke springe trinnet over. Prøv igen."); }
});

$("prev").addEventListener("click", () => { if (current > 0) { current--; render(); } });
$("next").addEventListener("click", () => { if (current < allSteps.length - 1) { current++; render(); } });

document.addEventListener("keydown", e => {
  if (e.key === "ArrowLeft") $("prev").click();
  else if (e.key === "ArrowRight") $("next").click();
});

$("action").addEventListener("click", async () => {
  const step = allSteps[current];
  try { await api.open(step.id); }
  catch (err) { alert("Siden kunne ikke åbnes. Prøv at åbne den manuelt i din browser."); }
});

$("smode-retry").addEventListener("click", checkSystem);

$("wifi-settings").addEventListener("click", async () => {
  try { await api.wifiSettings(); }
  catch (err) { alert("Wi-Fi-indstillingerne kunne ikke åbnes. Åbn dem selv via proceslinjen."); }
});

$("quit").addEventListener("click", async () => {
  try { await api.quit(); } catch (err) { /* sandsynligvis allerede lukket */ }
  document.body.innerHTML = "<p style='padding:40px;text-align:center'>Assistenten er lukket. Du kan lukke denne fane.</p>";
});

start();
```

- [ ] **Step 2: Bekræft at hele Go-pakken stadig er grøn**

Run (fra `setup-wizard/`): `go test ./...`
Expected: PASS (frontend-ændringen rører ikke server-adfærd).

- [ ] **Step 3: Commit**

```bash
git add setup-wizard/internal/web/static/app.js
git commit -m "feat: status-bevidst klient med skip, scene-fremdrift og tastatur"
```

---

### Task 6: Manuel verifikation af hele flowet

**Files:** (ingen ændringer — verifikation)

- [ ] **Step 1: Byg og kør Assistenten**

Fra `setup-wizard/`:

```bash
go run ./cmd/assistent
```

Åbn den localhost-URL den udskriver i en browser.

- [ ] **Step 2: Gennemgå tjeklisten visuelt**

Bekræft hver:
- Topbar viser det hvide NEG-logo + "Trin 1 / 10" og en fremdriftsbjælke.
- Hero-scenen er stærkt blåtonet (blueprint) ved start.
- Markér flere trin som færdige → deres fliser bliver blå-fyldte, bjælken vokser, og hero-scenen toner gradvist over i fuld farve.
- "Spring over" på et trin → flisen bliver dæmpet grå med ↷, og visningen går til næste trin; "Fortryd spring over" rydder den igen.
- Tre knapper (Åbn / Markér som færdig / Spring over) har samme størrelse; orange bruges kun til Åbn-CTA og den aktive flises ring.
- Sikkerheds-banneren nederst er tydelig og indeholder den fulde sikkerhedstekst.
- Tastatur ←/→ skifter trin.
- Alt tekst er på dansk med korrekte tegn (æ/ø/å).

- [ ] **Step 3: Kør hele testpakken en sidste gang**

Run (fra `setup-wizard/`): `go test ./...`
Expected: PASS.

- [ ] **Step 4: Afslut udviklingsgrenen**

Brug `superpowers:finishing-a-development-branch` til at vælge merge/PR/oprydning. Henvis evt. til issue #23 i PR-teksten (skip-tilstanden er implementeret).

---

## Self-Review-noter (udført under planlægning)

- **Spec-dækning:** brand-palet/font/logo → Task 2+3+4; blueprint→virkelighed-scene → Task 4+5 (`sceneProgress`, scene-opacity); flise-rail tre tilstande → Task 4+5; tre trintilstande/skip + tests → Task 1; S-mode-rød → Task 3+4; tilgængelighed (reduced-motion, tastatur) → Task 4+5; offline assets → Task 2; bevarede testede strenge → Task 3; manuel test → Task 6.
- **Typekonsistens:** `setStatus(id, status)` bruges ens i state.go (Task 1, Step 3) og server.go (Task 1, Step 4). Frontend `api.done/undo/skip` matcher ruterne. Element-id'er i Task 3 matcher `$()`-opslag i Task 5.
- **Ingen placeholders:** al kode er fuldt udfyldt.
