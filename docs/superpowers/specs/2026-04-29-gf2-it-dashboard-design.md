# GF2 IT Dashboard Design

## Formål

Projektet skal give nye GF2-elever en enkel, synlig og fast IT-indgang, som kan hentes fra Lectio og bruges både første skoledag og senere i forløbet. Dashboardet skal være forståeligt for ikke-tekniske elever og lærere.

## Leverancer

- `start.html`: light-version som kan åbnes direkte og fungerer som manuel guide.
- `IT-opstart-GF2.zip`: fuld elevpakke med dashboard, lokale assets, Windows/Mac setup-assistenter og manuel fallback.
- `assets/neg-hero-transition.png`: lokalt hero-billede med NEG fra blueprint/AutoCAD-stil til 3D-renderet skole-/IT-miljø.
- `scripts/setup-windows.ps1`: Windows setup-assistent.
- `scripts/setup-mac.sh`: Mac setup-assistent.

## Grundprincipper

- Dashboardet skal kunne åbnes 100% offline som lokal fil.
- Eksterne links må kun åbnes, når eleven aktivt klikker.
- Der må ikke gemmes brugernavne, adgangskoder, MitID-oplysninger eller andre personlige loginoplysninger.
- Status gemmes kun lokalt i browseren på elevens computer.
- OneDrive er elevens officielle sted til opgaver og filer. FileCloud/P-drev er ikke del af elevflowet.
- Printerhjælp ligger under Windows/Mac-vejledningen, ikke som fast hovedlink.

## Visuel Retning

- Professionelt NEG-blåt dashboard med `#005aa7` som primær blå.
- Topbar med NEG-markering og teksten `GF2 IT Dashboard`.
- Hero-bjælke med teksten:
  - `ET FÆLLESSKAB MED PLADS TIL DIG`
  - `Kom sikkert i gang med IT på GF2`
  - Kort underrubrik i mindre tekst.
- Hero-billedet beskæres ind i den blå bjælke mod højre, så teksten står roligt til venstre.
- Ingen eksterne billedreferencer i produktionsfiler.

## Elevvisning

Elevvisningen skal være den normale tilstand.

Hovedstruktur:

- Venstre panel: vælg computer, Windows eller Mac.
- Venstre panel: `Start setup-assistent` og `Manuel vejledning`.
- Midtersektion: opstart første dag i korrekt rækkefølge.
- Højre panel: faste links og praktisk info.

Opstartsrækkefølge:

1. Få udleveret NEG-brugernavn og adgangskode.
2. Log på skolens Wi-Fi med NEG-login.
3. Log på Office.com og skolemail med NEG-login.
4. Find SketchUp/Trimble-invitationen i skolemailen.
5. Log på Moodle med samme NEG-login.
6. Log på PraxisOnline og Lectio med MitID.
7. Installer eller åbn SketchUp og opret/log ind på Trimble-profil.
8. Gem opgaver i OneDrive via Office 365.

Status pr. trin:

- `Ikke startet`
- `I gang`
- `Færdig`

Statussen skal være lokal og synlig på skærmen, så underviseren hurtigt kan se, hvor langt eleven er.

## Underviservisning

Dashboardet skal have en Moodle-inspireret `Elev / Underviser` slider-toggle i topbaren.

Når `Underviser` er aktiv, skifter midtersektionen til en underviser-tjekliste. Venstre platformspanel og højre linkpanel kan blive stående.

Underviser-tjeklisten skal dække:

- Tjek at eleven er på Wi-Fi.
- Bekræft at Office.com og skolemail virker.
- Bekræft at Moodle virker med NEG-login.
- Forklar at SketchUp-installation kan hjælpes af scriptet, men Trimble/licens kræver elevens skolemail-flow.
- Forklar at Lectio og PraxisOnline er MitID-spor, hvor scriptet kun kan åbne de rigtige sider.
- Bekræft at eleven bruger OneDrive til filer.
- Typiske fejl: Windows S-mode, Chromebook, forkert login, manglende adgang til skolemail, manglende MitID.

Underviser-toggle er ikke loginbeskyttet i første version. Den er en praktisk visning, ikke et administrationssystem.

## Faste Links

Højre side skal indeholde disse faste links:

- Office 365 / skolemail
- Moodle
- Lectio
- PraxisOnline
- OneDrive
- SketchUp / Trimble
- NEG hjemmeside

Derudover vises en lille praktisk info-blok:

`Syg eller fravær? Ring til +45 72 290 100 mellem kl. 8.00-9.00.`

## Setup-assistenter

Scripts er assistenter, ikke fulde auto-installere.

Første version skal:

- Tjekke om eleven bruger Windows eller Mac.
- Åbne relevante sider i korrekt rækkefølge.
- Hjælpe eleven med at komme til Office, Moodle, Lectio, PraxisOnline, OneDrive, SketchUp og NEG.
- Forsøge at hjælpe med nyeste SketchUp-installation, hvor det er stabilt.
- Bruge manuel fallback, hvis automatisk SketchUp-installation ikke virker.
- Ikke gemme loginoplysninger.

Windows:

- Zip-pakken kan indeholde en dobbeltklikbar `Start Windows setup.cmd`, der kalder `scripts/setup-windows.ps1`.
- Scriptet må gerne tjekke Windows S-mode og relevante systemforhold.
- Scriptet kan oprette en skrivebordsgenvej til dashboardet.

Mac:

- Zip-pakken kan indeholde en dobbeltklikbar `Start Mac setup.command`, der kalder `scripts/setup-mac.sh`.
- Scriptet kan åbne relevante sider og hjælpe med SketchUp-download.
- Dashboard-genvej på Mac skal holdes enkel og ikke afhænge af browserbogmærker eller Dock-ændringer i første version.

## Pakkeform

Der laves to spor:

- Light: enkelt `start.html` til manuel guide.
- Fuld: zip-pakke med dashboard, assets og scripts.

Light-versionen er den sikre fallback. Zip-versionen er den anbefalede version til holdopstart.

## Ikke Med I Første Version

- Central indsamling af elevstatus.
- Lærer-login eller administrationspanel.
- Gemte elev-loginoplysninger.
- Automatisk login på Office, Moodle, Lectio, PraxisOnline eller MitID.
- FileCloud/P-drev som elevflow.
- Eksternt hentede billeder i dashboardets lokale UI.

## Åbne Punkter Før Implementation

- Endelige URL’er til Office, Moodle, Lectio, PraxisOnline, SketchUp/Trimble, OneDrive og NEG.
- Om SketchUp kan installeres stabilt via `winget` på Windows og en tilsvarende metode på Mac.
- Om Windows S-mode kan tjekkes robust uden administratorrettigheder.
- Om dashboardet skal bygges som én selvstændig HTML-fil med indlejret CSS/JS eller som flere lokale filer i zip-pakken.
