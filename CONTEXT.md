# Context: GF2 IT Setup

## Formål

GF2 IT Setup er et onboarding-værktøj til nye GF2-elever på NEG. Det hjælper elever med at komme i gang med Windows, skolemail, og de centrale systemer de bruger i undervisningen — uden at de behøver teknisk erfaring.

Projektet består af:
- **Landing page** (`index.html`) — GitHub Pages-side hvor elever downloader setup-pakken.
- **Dashboard** (`start.html`) — lokal tjekliste med trin-for-trin opsætning og hurtige links til skolesystemer. Har elev- og underviser-tilstand.
- **Windows setup-wizard** — en `.cmd`-baseret guide der åbner de rigtige sider og hjælper med installation (leveret som ZIP-pakke via `dist/`).
- **Mac setup** — planlagt men ikke implementeret endnu.

## Domæneordliste

Brug disse termer konsekvent. Undgå synonymer der ikke er listet her.

- **GF2** — Grundforløb 2. Den del af en erhvervsuddannelse hvor eleverne specialiserer sig. Projektets målgruppe.
- **NEG** — Nordjyllands Erhvervs Gymnasium. Skolen der bruger dette værktøj.
- **Elev** — en GF2-studerende. Projektets primære bruger. Brug ikke "bruger" eller "user" i eleven's kontekst.
- **Underviser** — en lærer på NEG der hjælper elever med opsætning. Sekundær bruger af dashboardet.
- **NEG-login** — skolens brugernavn + adgangskode (ikke MitID). Bruges til Wi-Fi, Office 365, Moodle.
- **UNI-Login** — national login-gateway for uddannelsessystemer. Eleven møder UNI-Login-portalen når de logger på PraxisOnline og Lectio, og vælger derfra MitID som autentificeringsmetode. Projektet håndterer aldrig UNI-Login-credentials.
- **MitID** — Danmarks nationale digitale identitet. Den autentificeringsmetode eleven vælger inde i UNI-Login-portalen for at tilgå Lectio og PraxisOnline. Projektet åbner kun siderne — det håndterer aldrig MitID-credentials.
- **Assistent** — det brugervendte navn for setup-flowet. Al UI-tekst og elevkommunikation bruger "Assistenten". _Undgå_: "setup-wizard" i elevvendt tekst.
- **Setup-wizard** — det teknisk-interne navn for `.cmd`-entry-pointet + PowerShell-flowet (`Start Windows setup.cmd` → `setup-windows.ps1`). Bruges kun i kode, docs og udviklersammenhæng. Eleven ser aldrig dette ord.
- **Dashboard** — `start.html`. Lokal tjekliste med status (gemmes i `localStorage`), hurtige links, og vejledning. Har elev-tilstand og underviser-tilstand.
- **Landing page** — `index.html`. Den offentlige GitHub Pages-side med download-knap og projektbeskrivelse.
- **Setup-pakke** — ZIP-filen (`dist/GF2-IT-Setup-Windows.zip`) der leveres til elever. Indeholder setup-wizard og dashboard.
- **S-mode** — Windows S-mode. En begrænsning der forhindrer installation af programmer udenfor Microsoft Store. En kendt blokering for elever.

## Centrale systemer eleverne skal bruge

- **Office 365** — skolemail, Word, OneDrive. Tilgås via NEG-login.
- **Moodle** — læringsplatform. Tilgås via NEG-login.
- **Lectio** — skema og fravær. Tilgås via UNI-Login-portalen, hvor eleven vælger MitID.
- **PraxisOnline** — opgaveplatform. Tilgås via UNI-Login-portalen, hvor eleven vælger MitID.
- **OneDrive** — fillagring via Office 365.
- **SketchUp / Trimble** — 3D-modelleringsprogram. Kræver invitation via skolemail.

## Sikkerhedsprincip

Projektet beder **aldrig** om adgangskoder, MitID eller UNI-Login. Elever indtaster kun loginoplysninger på officielle sider eller i Windows' egne indstillinger. Dette princip skal overholdes i al kode, UI-tekst og dokumentation.

## Teknisk kontekst

- **Windows-first** — primær platform. Mac-support er planlagt men ikke bygget.
- **Ingen backend** — alt kører lokalt eller som statiske sider (GitHub Pages).
- **Ingen build-step for frontend** — ren HTML/CSS/JS uden frameworks.
- **PowerShell scripts** — bruges til test (`tests/`) og build (`scripts/build-package.ps1`).
- **Lokal status** — dashboardet gemmer trin-status i `localStorage`. Ingen data sendes nogen steder.
