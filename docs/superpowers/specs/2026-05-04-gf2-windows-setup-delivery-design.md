# GF2 Windows Setup Delivery Design

## Formaal

Projektet skal goere GF2 IT-opstart nemmere for nye elever og undervisere paa en dag med mange nye systemer. Foerste produktionsrettede version skal fokusere paa Windows, fordi Windows kan testes nu, og fordi de fleste tekniske checks og installationer skal koere lokalt paa elevens computer.

Loesningen skal vaere enkel nok til ikke-tekniske elever og kollegaer:

1. Eleven gaar paa `NEG Guest`.
2. Eleven aabner en enkel GitHub Pages landing page.
3. Eleven klikker paa en tydelig Windows download-knap.
4. Eleven starter en lokal Windows setup-assistent.
5. Setup-assistenten guider eleven gennem Wi-Fi, skolemail, Office, Moodle, PraxisOnline, Lectio, OneDrive og SketchUp.

## Overordnede Principper

- Loesningen skal vaere open source og kunne goeres public, naar den er klar.
- Foerste offentlige levering sker via GitHub Pages, ikke et NEG-domaene.
- Et officielt NEG-domaene kan tilfoejes senere, hvis skolen vil goere loesningen officiel.
- Elever skal ikke bruge Git, GitHub CLI, VS Code eller PowerShell 7.
- Bootstrap og setup skal virke paa standard Windows med CMD og Windows PowerShell 5.1.
- Appen maa aldrig bede om, modtage, gemme eller vise adgangskoder, MitID-oplysninger eller UNI-Login.
- Eleven indtaster kun loginoplysninger paa officielle sider eller i Windows' egne indstillinger.
- Setup-assistenten skal vaere gennemsigtig: den viser hvad den goer, og hvorfor naeste trin er noedvendigt.

## Leveringsmodel

Den primaere leveringskanal er en simpel landing page paa GitHub Pages:

```text
https://tech-pleex.github.io/Startup-project-2026/
```

Siden skal ikke ligne GitHub for eleverne. Den skal vaere en rolig, tydelig downloadside med et visuelt NEG/IT-signal, en stor download-knap og sikkerhedstekst.

Foerste landing page kan bruge det nuvaerende hero-billede som midlertidigt visuelt udtryk. Billedet bruges til at skabe genkendelse og signalere IT-opstart, mens resten af siden holdes enkel.

Landing page skal indeholde:

- Hero-billede.
- Titel: `GF2 IT Setup`.
- Kort forklaring af formaalet.
- Primaer knap: `Download Windows Setup`.
- Sekundaer knap/link: `Aabn underviser-/projektorguide`.
- Link: `Se kildekode paa GitHub`.
- Sikkerhedstekst: `Vi beder aldrig om adgangskoder, MitID eller UNI-Login. Du logger kun ind paa officielle sider, som assistenten aabner.`

Download-knappen skal pege paa en GitHub Release zip, naar wizard-pakken findes. Indtil da kan siden pege paa den nuvaerende zip-pakke eller vaere markeret som prototype.

## Underviserflow

Underviseren skal kunne bruge landing page og dashboard som praktisk klasseintro:

1. Vis landing page eller projektorguide paa tavlen.
2. Forklar at assistenten ikke beder om adgangskoder.
3. Bed eleverne gaa paa `NEG Guest`.
4. Bed eleverne aabne landing page-adressen i browseren.
5. Eleverne downloader Windows setup og foelger assistenten.
6. Underviseren hjaelper kun ved trin der viser `Kraever handling`.

Underviser-/projektorguiden maa gerne vaere den eksisterende HTML-dashboardoplevelse, men den skal beskrives som guide og statusoverblik, ikke som installationsmotor.

## Elevflow

Elevens setup skal vaere trin-for-trin i et klassisk Windows-installationsvindue. HTML-dashboardet maa gerne aabnes til sidst, men selve setup-processen skal ikke afhaenge af at browseren kan koere systemchecks.

Planlagt Windows setup-flow:

1. Velkomst og sikkerhed.
2. Wi-Fi check.
3. Office/skolemail.
4. Trimble invitation via skolemail.
5. Moodle.
6. PraxisOnline.
7. Lectio.
8. OneDrive.
9. SketchUp installation eller fallback.
10. Faerdig: opret dashboard-genvej og aabn dashboard.

## Trin og Checks

### Velkomst og Sikkerhed

Wizarden starter med en kort tekst:

```text
GF2 IT Setup hjaelper dig med at komme de rigtige steder hen.
Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login.
Du indtaster kun oplysninger paa officielle sider og i Windows' egne indstillinger.
```

### Wi-Fi

Wizarden skal tjekke aktivt Wi-Fi-navn, hvis Windows tillader det.

- Hvis aktivt netvaerk er `NEG`, markeres trinnet som gennemfoert.
- Hvis aktivt netvaerk er `NEG Guest`, viser wizarden at eleven er paa gaestenet og skal skifte til `NEG`.
- Hvis netvaerket ikke kan laeses, viser wizarden manuel vejledning.

Wizarden maa aabne Windows Wi-Fi-indstillinger, men eleven skal selv vaelge netvaerk og indtaste brugernavn/adgangskode.

### Office Og Skolemail

Wizarden aabner Office/skolemail. Eleven logger selv ind med skoleoplysninger. Wizarden kan ikke og maa ikke tjekke loginindhold. Eleven klikker selv `Faerdig`, naar skolemail virker.

### Trimble Invitation

Wizarden forklarer at SketchUp/Trimble invitationen kommer via skolemail. Eleven skal:

1. Aabne skolemail.
2. Finde mailen fra Trimble/SketchUp.
3. Klikke invitationslinket.
4. Oprette eller logge ind hos Trimble med skolemail-flowet.
5. Klikke `Faerdig` i wizarden.

Wizarden aktiverer ikke invitationen automatisk.

### Moodle

Wizarden aabner Moodle. Eleven logger selv ind. Eleven klikker `Faerdig`, naar Moodle virker.

### PraxisOnline

PraxisOnline skal have en tydelig gul advarselsboks.

Tekst:

```text
Vigtigt:
Brug din skolemail til PraxisOnline.
Eksempel: neg04026@edu.neg.dk

Login sker via UNI-Login, hvor du vaelger MitID.
Assistenten beder aldrig om MitID, UNI-Login eller adgangskode.
```

Wizarden aabner PraxisOnline. Eleven opretter eller logger selv ind og klikker `Faerdig`.

### Lectio

Lectio skal forklares separat fra PraxisOnline.

Tekst:

```text
Log ind med UNI-Login, hvor du vaelger MitID.
Assistenten aabner kun Lectio og ser aldrig dine oplysninger.
```

Wizarden aabner Lectio. Eleven logger selv ind og klikker `Faerdig`.

### OneDrive

Wizarden aabner OneDrive via Office 365. Eleven bekraefter selv at OneDrive virker.

### SketchUp

SketchUp installeres foerst via `winget`, hvis `winget` findes og pakken kan installeres.

Foerste stabile konfiguration:

```powershell
$SketchUpPackageId = "Trimble.SketchUp.2026"
```

Wizarden skal vise hvad der installeres, inden installationen startes. Hvis `winget` mangler, installationen fejler, eller maskinen er i Windows S-mode, skal wizarden vise manuel fallback til SketchUp/Trimble-siden.

Pakke-id holdes som en tydelig konfigurationsvaerdi, saa den kan aendres til 2027 eller senere uden at omskrive hele flowet.

## Teknisk Arkitektur

Foerste Windows-version bygges med:

- GitHub Pages landing page.
- GitHub Release zip til Windows setup-pakken.
- CMD/PowerShell bootstrap som fallback.
- Windows PowerShell 5.1 setup-engine.
- Klassisk Windows GUI, sandsynligvis via PowerShell Windows Forms.
- Eksisterende `start.html` som dashboard/projektorguide og slutside.

HTML alene maa ikke bruges som installationsmotor, fordi browseren ikke kan tjekke Wi-Fi, koere `winget`, aabne Windows-indstillinger paa samme kontrollerede maade eller laese lokale systemforhold sikkert.

## Sikkerhed

Loesningen skal vaere nem at forklare:

- Alt kode er open source, naar repoet bliver public.
- Hjemmesiden er statisk.
- Ingen database.
- Ingen login paa hjemmesiden.
- Ingen tracking.
- Ingen elevdata.
- Ingen adgangskoder, MitID eller UNI-Login i appen.
- Login sker kun hos Microsoft, NEG, PraxisOnline, Lectio, Trimble eller Windows' egne indstillinger.

Hvis repoet stadig er privat under udvikling, skal landing page og download kun bruges som test/prototype.

## Scope For Foerste Version

Med i foerste Windows-version:

- GitHub Pages landing page.
- Windows download-knap.
- Sikkerhedstekst.
- Klassisk Windows setup-wizard.
- Wi-Fi check for `NEG` vs `NEG Guest`, hvis muligt.
- Manuelle bekraeftelsestrin for Office/skolemail, Trimble invitation, Moodle, PraxisOnline, Lectio og OneDrive.
- SketchUp installation via konfigureret `winget` package id med manuel fallback.
- Dashboard-genvej.

Ikke med i foerste version:

- Mac setup-wizard.
- Automatisk login.
- Indsamling af elevdata.
- Central status for hele klassen.
- Officielt NEG-domaene.
- Signeret `.exe` installer.
- Microsoft Store eller Intune deployment.

## Aabne Beslutninger Til Implementation Plan

- Om Windows GUI bygges med Windows Forms eller en anden indbygget Windows-teknologi.
- Praecis GitHub Pages struktur: eksisterende repo page eller separat `tech-pleex.github.io` repo.
- Om foerste download peger paa nuvaerende zip eller en ny release-pakke.
- Hvordan branch `feature/setup-wizard` skal opdeles mellem landing page, wizard og packaging.
- Hvilket navn den downloadede zip og lokale mappe skal have.
