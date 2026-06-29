# ADR-0001: Go som SSOT for links via template-genereret dashboard

**Status:** Accepteret (2026-06-29)
**Kontekst-issues:** #30 (driver), #3, #29, #7, #22, #17, #27

## Kontekst

Links til skolesystemer findes i dag flere steder uden én fælles kilde:

- `setup-wizard/internal/steps/steps.go` — Assistentens links (forbruges af den kompilerede Go-binær).
- `start.html` — dashboardets egen JS-kopi (`const links`), serveret som statisk GitHub Pages-side uden Go på runtime.

De to kopier er allerede begyndt at drive fra hinanden (fx SketchUp-URL'en og Trimble-invitationen). Det giver dobbelt vedligehold og risiko for døde/uenige links.

PRD #3 fastslog tidligere at "dashboardet forbliver uden build-step". Den beslutning stammer fra projektets bash-tid og er ved denne ADR frigjort — den er ikke længere bindende.

## Beslutning

1. **Go er den eneste kilde (SSOT) for dashboard-links.** Link-data defineres ét sted i Go.
2. **Dashboardet genereres fra en template.** Go ejer `start.html.tmpl` (markup/layout/brand) + link-data og renderer `start.html`. Den genererede `start.html` er et artefakt og redigeres ikke i hånden.
3. **Genereringen kobles til CI (#17).** Samme pipeline bygger Assistenten (mod Microsoft Store, #27) og regenererer + deployer dashboardet til GitHub Pages. Workflow for maintainer: ret et link ét sted i Go → build → alt virker igen.

## Konsekvenser

**Positivt:**
- Driften forsvinder — én kilde til alle links.
- Lavt vedligehold: ét sted at rette links.
- Rent ejerskab: #30 ejer link-data + generator; layout/brand-skabelonen ejes af design-arbejdet (#3/#29).

**Negativt / at være opmærksom på:**
- Dashboardet får et build-step og en CI-afhængighed (#17), som projektet ikke havde før.
- `start.html` er fremover genereret — håndredigeringer skal ske i `start.html.tmpl`. Den genererede fil bør bære en "GENERERET — rediger .tmpl"-header.

## Rækkefølge

Det visuelle redesign (#3 layout + #29 brand + #7 lærer-link) laves **først** på `start.html` som ren HTML, for hurtig visuel iteration. Templatiseringen (#30) sker **bagefter** på det stabile, færdige design — så vi ikke templatiserer et mål i bevægelse.

## Alternativer overvejet

- **Delt `links.json` + runtime-`fetch`:** holder dashboardet build-frit, men gør ikke Go til ægte SSOT og giver en runtime-afhængighed (fejlet fetch = ingen links). Fravalgt fordi maintainer alligevel ønsker CI (#17) og ét vedligeholds-sted i Go.
- **Go injicerer kun et fragment i håndholdt `start.html`:** skrøbeligt (markør kan overskrives; to kan redigere samme fil). Fravalgt til fordel for fuld template-rendering.
