# Design: Kørbar Assistent-binær med tringuide (Windows + Mac)

**Dato:** 2026-06-12
**Udspringer af:** Issue #13 (PRD: Totalt rewrite af setup-leveringen i Go) — andet sub-goal efter porteringen af trinkonfiguration og OS-logik (commit 99f3253).

## Mål

Fra én Go-kodebase bygges en dobbeltklikbar `.exe` til Windows (og samtidig en Mac-binær), der starter en lokal webserver, åbner browseren og viser den danske tringuide. En elev kan gennemføre hele trinforløbet: se trin, markere færdige, tjekke Wi-Fi, åbne officielle sider og starte SketchUp-installationen.

## Afgrænsning (besluttet i brainstorm)

- **Trin-status holdes kun i hukommelsen.** Persistens til fil (PRD user story 14) er et separat, senere sub-goal.
- **Tringuide-siden skrives frisk**, men låner farver/typografi fra dashboardet (`start.html`). Prototypens kode genbruges ikke.
- GitHub Release-artefakter og landing page-links er uden for denne bid.
- JS-adfærd testes ikke automatisk i denne bid; al logik ligger bag det httptest-dækkede API.

## Arkitektur

```
setup-wizard/
├── cmd/assistent/main.go        ← NY: tynd main. Port, server, åbn browser, vent
├── internal/
│   ├── steps/                   ← findes (trinkonfiguration, 10 trin)
│   ├── wizard/                  ← findes (Wi-Fi-klassifikation, SketchUp-flow)
│   ├── osops/                   ← findes (OS-interface, Windows + Mac)
│   ├── server/                  ← NY: HTTP-API + servering af indlejrede assets
│   │   ├── server.go
│   │   ├── server_test.go       ← httptest-tests (PRD'ens højeste søm)
│   │   └── state.go             ← trin-status i hukommelsen (mutex-beskyttet)
│   └── web/                     ← NY: tringuide-siden
│       ├── embed.go             ← //go:embed static
│       └── static/              ← index.html, style.css, app.js
```

**Opstartsflow:** main binder til `127.0.0.1:0` (OS vælger ledig port), starter serveren, åbner `http://127.0.0.1:<port>` via det eksisterende `osops`-interface, og blokerer indtil `POST /api/quit`.

**Build (dokumenteres i `setup-wizard/README.md`):**

```bash
GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o dist/Assistenten.exe ./cmd/assistent
GOOS=darwin  GOARCH=arm64 go build -o dist/Assistenten ./cmd/assistent
```

`-H=windowsgui` forhindrer det sorte konsolvindue ved dobbeltklik på Windows. Ingen CGO, ingen dependencies ud over stdlib — krydskompilering fra enhver maskine.

## HTTP-API

Alle endpoints binder kun til `127.0.0.1` og svarer JSON. Routing med stdlib `http.ServeMux` (Go 1.22+-mønstre, ingen router-dependency).

| Endpoint | Funktion |
|---|---|
| `GET /` | Serverer indlejret tringuide-side |
| `GET /api/steps` | Trinliste: id, titel, brødtekst, advarsel, knaptype, færdig-status |
| `POST /api/steps/{id}/done` | Markér trin færdigt |
| `POST /api/steps/{id}/undo` | Fortryd trin |
| `GET /api/wifi` | SSID + klassifikation (NEG / NEG Guest / ukendt / intet netværk) |
| `POST /api/wifi/settings` | Åbn systemets Wi-Fi-/netværksindstillinger |
| `POST /api/steps/{id}/open` | Åbn trinnets officielle URL (slås op server-side) |
| `POST /api/sketchup/install` | Kør SketchUp-flowet; svarer hvilken vej der blev taget |
| `POST /api/quit` | Luk Assistenten pænt |

**Sikkerhed håndhævet i kode:** Browseren sender aldrig URL'er — `/open` slår URL'en op i den indlejrede trinkonfiguration, så ingen proces kan misbruge Assistenten til at åbne vilkårlige sider. Serveren binder kun til localhost.

**Dataflow:** Ensrettet. Siden henter `GET /api/steps`, tegner, og genhenter efter hver handling. Go-processen er eneste kilde til sandhed; ingen klient-side-tilstand.

**Fejlhåndtering:** Ukendt trin-id → 404. SketchUp-fallback er et forventet udfald, ikke en fejl: 200 med `{"action": "fallback", "reason": ...}`.

## Tringuide-siden

Tre filer uden frameworks eller build-step: `index.html`, `style.css`, `app.js`.

- **Layout:** Venstre kolonne: trinliste med numre, titler, flueben og "trin X af Y". Højre kolonne: aktivt trins brødtekst, evt. advarselsboks (PraxisOnline), handlingsknap, "Markér som færdig"/"Fortryd", frem/tilbage.
- **Wi-Fi-trinnet:** viser SSID + klassifikation, "Tjek igen"-knap.
- **SketchUp-trinnet:** "Installér automatisk"-knap; viser om det blev winget, S-mode-fallback eller fejl-fallback.
- **Sikkerhedsbanner** fast i bunden på alle trin: Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login.
- **Stil:** farver, skrifttype og kort-udseende lånes fra `start.html`; markup skrives frisk.
- Al elevvendt tekst er dansk og bruger ordet "Assistenten".

## Tests (skrives først, TDD)

Højeste søm er HTTP-API'et via `httptest` med den eksisterende `FakeOS` injiceret:

1. `GET /api/steps` returnerer alle 10 trin i rækkefølge med dansk indhold.
2. done/undo ændrer status og afvises for ukendte id'er (404).
3. `GET /api/wifi` klassificerer korrekt for NEG, NEG Guest, ukendt og intet netværk.
4. `POST /api/steps/{id}/open` åbner præcis trinnets konfigurerede URL og afviser ukendte id'er.
5. `POST /api/sketchup/install` vælger winget / S-mode-fallback / fejl-fallback korrekt.
6. `GET /` svarer 200, indeholder sikkerhedsteksten, og ingen indlejret fil indeholder ordet "setup-wizard".
7. `POST /api/quit` signalerer nedlukning.

Krydskompilering verificeres som hidtil: `GOOS=windows go build ./...` og `GOOS=darwin go build ./...` skal begge lykkes.

Manuel verifikation på rigtige maskiner dækkes af issue #12 og er ikke en del af denne bid.

## Tilgange overvejet

1. **JSON-API + fetch-baseret side (valgt):** matcher PRD'ens formulering, giver det stærkeste test-søm (httptest mod JSON).
2. Server-renderet HTML med formular-posts: mindre JS, men side-blink og HTML-parsende tests.
3. Server-Sent Events for live Wi-Fi-status: YAGNI — eleven trykker selv "Tjek igen".
