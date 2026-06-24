# Setup-wizard UI-redesign — "Blueprint → virkelighed"

- **Dato:** 2026-06-24
- **Branch:** `feature/ui-redesign`
- **Status:** Design godkendt, klar til implementeringsplan

## Formål

Assistentens (setup-wizard) UI er funktionelt men visuelt fladt. Dette redesign giver
den et moderne, energisk udtryk, der **strengt** følger NEG's brandguide, og gør
elevens fremdrift til hovedpersonen via et gennemgående koncept: **blueprint →
virkelighed**.

## Koncept

Eleven starter med en *blueprint* af sin NEG-opsætning. For hvert fuldført trin
"renderer" en del af scenen fra blå streg til virkelighed. Når alle trin er klaret,
er hele scenen bygget. Fremdrift = scenen der bliver virkelig. Konceptet er forankret
i brugerens eget hero-billede `assets/neg-hero-transition.png`.

## Brand-overholdelse (hård regel)

Alt skal følge `NEG BRAND/NEG_Brandguide_2024`. Ingen afvigelser uden eksplicit grund.

- **Farver (kun disse + tints 100/80/40/20%):**
  - Mørkeblå `#123c62` (primær — baggrund, tekst)
  - Blå `#86afd8` (primær — streger/blueprint, sekundær tekst, "renderet" flise-fyld)
  - Lysegrå/grøn `#efefef` (sekundær — lyse baggrunde)
  - Orange `#ec8113` (sekundær — **kun** CTA/links og lille aktiv-accent, sparsomt)
- **Typografi:** Space Grotesk (brandfont). Bold=overskrifter, Medium=labels/CTA
  (uppercase + letter-spacing til "tekniske" labels), Regular=brødtekst, Light=footer.
  Bundles lokalt som `.woff2` (OFL-licens) — **ingen CDN**, da Assistenten kører før
  Wi-Fi er sat op.
- **Logo:** brug de officielle logofiler (`NEG BRAND/2023 - NEG logo white.png` på mørk
  baggrund). Må ikke gen-tegnes/forvrænges. Respektafstand = 50% af logohøjde.
- **Bevidst undtagelse:** S-mode-blokeringen vises i **system-rød** (ikke NEG-palet),
  fordi det er en Windows/maskine-fejl ("computeren opfylder ikke kravene"), ikke
  NEG-kommunikation. Holdes visuelt adskilt fra brandets udtryk.

## Scope

- **Ændres:** `setup-wizard/internal/web/static/` (`index.html`, `style.css`, `app.js`)
  + nye lokale assets (Space Grotesk `.woff2`, optimeret hero-billede, logo).
  `embed.FS` dækker allerede hele `static/`, så assets følger automatisk med i binæren.
- **Lille backend-ændring (pga. issue #23, "Spring over"):** trin får en tredje
  tilstand. Berører `internal/server/state.go` (tilstandsmodel) og `internal/server/server.go`
  (nyt endpoint). API-kontrakten udvides bagudkompatibelt.
- **Hero-asset:** den nuværende `neg-hero-transition.png` er 2,3 MB. Den skal
  optimeres/skaleres (fx bredde ~1600px, komprimeret) før indlejring, så binæren ikke
  vokser unødigt.
- **Uden for scope:** ændringer i trinforløbets indhold/tekster, øvrige GH-issues,
  dashboard, release-pipeline.

## Layout (valgt: "A — scene øverst")

```
┌─────────────────────────────────────────────┐
│  [NEG-logo] · IT-OPSÆTNING        Trin 3/10 ▓░│  topbar: officielt logo + fremdrift
├─────────────────────────────────────────────┤
│   HERO-SCENE (blueprint-krydstoning)          │  renderer efter samlet fremdrift
├─────────────────────────────────────────────┤
│  [Wi-Fi✓][Office✓][Skolemail•][Moodle]…       │  flise-rail: fliser = trin, klik = hop
├─────────────────────────────────────────────┤
│  TRIN 3/10 · LINK                             │
│  Office 365 / skolemail                       │  aktivt trin-kort
│  Log ind på Office …                          │
│  [ Warning-felt v. behov ]                    │
│  [Åbn ↗]  [Marker færdig ✓]  [Spring over ↷]  │  tre ens-store knapper
│  ‹ Tilbage                          Frem ›    │
├─────────────────────────────────────────────┤
│  🛡️ Assistenten beder aldrig om adgangskoder… │  tydelig sikkerheds-banner
└─────────────────────────────────────────────┘
(S-mode: blokerende rød systemfejl-banner når aktiv)
```

## Komponenter (frontend, vanilla JS i klare enheder)

UI'et forbliver dependency-frit (vanilla). `app.js` opdeles i fokuserede enheder med
klare grænseflader:

- **`api`** — tynd wrapper om de eksisterende endpoints (`GET /api/steps`, `/api/system`,
  `POST …/done`, `…/undo`, `…/open`, `/api/wifi/settings`, `/api/quit`) + det nye
  skip-endpoint. Returnerer rene data; ingen DOM.
- **`state`** — afledt model fra `/api/steps`: `{ steps[], currentIndex }`, hvor hvert
  trin har en status (`pending` | `done` | `skipped`). Render-tilstand er **ren
  afledning** af status — ingen dobbelt sandhed.
- **`scene`** — styrer hero-krydstoningen: blueprint-lagets opacity = `1 − fremdrift`,
  hvor `fremdrift = (antal trin med status done) / (antal trin eksklusive finish)`.
  `prefers-reduced-motion` → instant skift uden transition. Bemærk: dette er adskilt fra
  topbarens "Trin 3/10", som blot er elevens nuværende position i forløbet.
- **`rail`** — tegner fliserne (1:1 med trin) i fire tilstande (todo/blueprint, rendered,
  active, skipped); klik på flise sætter `currentIndex`.
- **`stepCard`** — aktivt trin: kicker-label, titel, body, valgfri warning, knaprække,
  wifi-panel, S-mode-blokering, Tilbage/Frem.
- **Rene funktioner** udtrækkes hvor det giver mening (fx `sceneProgress(steps)`,
  `tileState(step)`) så logik kan ræsonneres om isoleret.

## Trintilstande & "Spring over" (issue #23)

Hvert trin har tre praktiske udfald: **ikke færdig**, **færdig**, **sprunget over**.

- Hver trin-visning har en synlig **"Spring over ↷"**-handling.
- Klik på "Spring over" → gem status `skipped` og gå til næste trin.
- `skipped` lagres **adskilt** fra `done` i serverens tilstand.
- Visuelt adskilt: færdig = blå-fyldt flise; sprunget over = dæmpet grå flise med ↷;
  ikke-færdig = stiplet blueprint-blå.
- Et sprunget-over-trin kan **genåbnes** og ændres til færdig/ikke-færdig igen.
- **Fremdrift:** kun `done` fylder scenen mod virkelighed; `skipped` tæller ikke som
  "bygget" (efterlader bevidst en markeret flise), men blokerer ikke for at nå "Færdig".

## Dataflow

1. Load → `GET /api/steps` + `GET /api/system` (S-mode).
2. Render scene (krydstoning efter fremdrift) + rail + aktivt kort.
3. Handling (færdig/undo/skip/åbn) → `POST` → opdatér `state` → re-render afledt:
   flise skifter tilstand, scenen renderer et hak, fremdriftsbjælke opdateres.
4. Ingen dobbelt-tilstand; alt udledes af serverens trin-status.

## Fejlhåndtering & robusthed

- API-fejl (fx `open`/`wifi`) vises i trin-kortets warning-felt (eksisterende mønster).
- **S-mode:** blokerende rød systemfejl-banner; scenen holdes i blueprint-tilstand indtil
  løst. Bevarer eksisterende "Tjek igen"-knap.
- **Offline er normaltilstanden:** nul eksterne kald; fonte/billeder lokale.

## Tilgængelighed

- `prefers-reduced-motion` slår scene-/flise-animationer fra (instant tilstandsskift).
- Fuld tastaturnavigation (←/→ for Tilbage/Frem, fokuserbare fliser/knapper).
- Bevar `aria`-attributter og semantik. UI'et skal være fuldt brugbart uden effekterne.
- Tilstrækkelig kontrast: lys tekst på mørkeblå, navy tekst på blå fliser.

## Test

- Eksisterende Go-server-tests (`server_test.go`, `steps_test.go`, `wifi_test.go`) skal
  forblive grønne.
- Nye Go-tests for skip-tilstanden: skip et normalt trin, skip sidste trin, og skift et
  sprunget-over-trin tilbage til færdig/ikke-færdig (jf. issue #23 acceptkriterier).
- Manuel verifikation: byg binæren, kør Assistenten, gennemgå alle trin og se scenen
  rendere fuldt + skip/genåbn fungere.
- Ingen JS-testharness (YAGNI) — render-logik holdes i rene funktioner og verificeres
  manuelt + via Go-laget.

## Assets der skal tilføjes

- Space Grotesk `.woff2` (de nødvendige vægte: 300/400/500/700), OFL — lægges i `static/`.
- Optimeret hero-billede (afledt af `assets/neg-hero-transition.png`).
- Officielt NEG-logo (hvid variant) til topbaren.

## Referencer

- Issue #23 — "Add Spring over action to every Assistent step"
- `NEG BRAND/NEG_Brandguide_2024_Holbæk_Audebo.pdf` (lokal, gitignored)
- Godkendt mockup: `.superpowers/brainstorm/1929-1782286794/content/final-mockup-v5.html`
