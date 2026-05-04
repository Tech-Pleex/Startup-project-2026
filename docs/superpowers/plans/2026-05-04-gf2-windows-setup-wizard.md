# GF2 Windows Setup Wizard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first Windows-first delivery flow: GitHub Pages landing page, downloadable Windows setup package, and a local Windows GUI setup wizard that never asks for passwords.

**Architecture:** Keep the public landing page static and simple. Keep installation/check logic in Windows PowerShell 5.1 using Windows Forms, split into config, checks, and wizard entrypoint. Keep the existing `start.html` as dashboard/projector guide and final opened dashboard.

**Tech Stack:** Static HTML/CSS, Windows CMD, Windows PowerShell 5.1, Windows Forms, `winget`, existing PowerShell package/test scripts.

---

## File Structure

- Create: `index.html` - GitHub Pages landing page with hero image, safety text, download button, projector guide link, GitHub source link.
- Create: `scripts/setup-config.ps1` - central config for URLs, package labels, SketchUp package id, dashboard path, and status labels.
- Create: `scripts/setup-checks.ps1` - small testable functions for Wi-Fi SSID, Windows S-mode, winget availability, URL opening, and shortcut creation.
- Modify: `scripts/setup-windows.ps1` - replace terminal prompts with Windows Forms wizard UI; dot-source config/check functions.
- Modify: `Start Windows setup.cmd` - start the GUI wizard using standard Windows PowerShell 5.1 and hide the PowerShell process window where possible.
- Modify: `scripts/build-package.ps1` - build `dist/GF2-IT-Setup-Windows.zip` with the Windows wizard files and existing dashboard assets.
- Create: `tests/check-setup-delivery.ps1` - validate landing page, package contents, config values, Windows Forms references, and safety wording.
- Modify: `tests/check-dashboard.ps1` - keep existing dashboard checks and add references for new package name only if needed.
- Modify: `.gitattributes` - ensure `.ps1`, `.cmd`, `.html`, and `.md` text handling is explicit.
- Create or modify: `README.md` - explain purpose, safety model, current private/prototype status, and local test commands.

## Task 1: Add Delivery Test Skeleton

**Files:**
- Create: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Write the failing delivery test**

Create `tests/check-setup-delivery.ps1` with:

```powershell
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$LandingPage = Join-Path $Root "index.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$WindowsLauncher = Join-Path $Root "Start Windows setup.cmd"
$WindowsSetup = Join-Path $Root "scripts\setup-windows.ps1"
$SetupConfig = Join-Path $Root "scripts\setup-config.ps1"
$SetupChecks = Join-Path $Root "scripts\setup-checks.ps1"
$BuildPackage = Join-Path $Root "scripts\build-package.ps1"
$ZipPath = Join-Path $Root "dist\GF2-IT-Setup-Windows.zip"
$Dashboard = Join-Path $Root "start.html"

function Assert-File {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) {
        throw "Missing required file: $Path"
    }
}

function Assert-Contains {
    param(
        [string]$Content,
        [string]$Needle,
        [string]$Label
    )
    if (-not $Content.Contains($Needle)) {
        throw "Missing required content '$Label': $Needle"
    }
}

function Assert-NotContains {
    param(
        [string]$Content,
        [string]$Needle,
        [string]$Label
    )
    if ($Content.Contains($Needle)) {
        throw "Forbidden content '$Label': $Needle"
    }
}

function Assert-ZipContains {
    param(
        [string]$ZipFile,
        [string[]]$RequiredEntries
    )

    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $Zip = [System.IO.Compression.ZipFile]::OpenRead($ZipFile)
    try {
        $Entries = $Zip.Entries.FullName
        foreach ($Required in $RequiredEntries) {
            if ($Entries -notcontains $Required) {
                throw "Zip is missing entry: $Required"
            }
        }
    }
    finally {
        $Zip.Dispose()
    }
}

Assert-File $LandingPage
Assert-File $HeroImage
Assert-File $WindowsLauncher
Assert-File $WindowsSetup
Assert-File $SetupConfig
Assert-File $SetupChecks
Assert-File $BuildPackage
Assert-File $Dashboard

$LandingHtml = Get-Content -Raw -Encoding UTF8 -LiteralPath $LandingPage
$LauncherContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsLauncher
$SetupContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsSetup
$ConfigContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupConfig
$ChecksContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupChecks
$BuildContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $BuildPackage

Assert-Contains $LandingHtml "GF2 IT Setup" "landing title"
Assert-Contains $LandingHtml "assets/neg-hero-transition.png" "landing hero image"
Assert-Contains $LandingHtml "Download Windows Setup" "download button"
Assert-Contains $LandingHtml "dist/GF2-IT-Setup-Windows.zip" "download target"
Assert-Contains $LandingHtml "Vi beder aldrig om adgangskoder" "landing safety message"
Assert-Contains $LandingHtml "start.html" "projector guide link"
Assert-Contains $LandingHtml "GitHub" "source transparency link"

Assert-Contains $ConfigContent "Trimble.SketchUp.2026" "SketchUp package id"
Assert-Contains $ConfigContent "PraxisOnline" "Praxis step"
Assert-Contains $ConfigContent "neg04026@edu.neg.dk" "Praxis school mail example"
Assert-Contains $ConfigContent "UNI-Login" "UNI-Login guidance"
Assert-Contains $ConfigContent "MitID" "MitID guidance"
Assert-Contains $ConfigContent "NEG Guest" "guest Wi-Fi"
Assert-Contains $ConfigContent "`"NEG`"" "target Wi-Fi"

Assert-Contains $ChecksContent "Get-ActiveWifiSsid" "Wi-Fi check function"
Assert-Contains $ChecksContent "Test-WingetAvailable" "winget check function"
Assert-Contains $ChecksContent "Test-WindowsSMode" "S-mode check function"
Assert-Contains $ChecksContent "New-DashboardShortcut" "dashboard shortcut function"

Assert-Contains $SetupContent "System.Windows.Forms" "Windows Forms GUI"
Assert-Contains $SetupContent "GF2 IT Setup" "wizard title"
Assert-Contains $SetupContent "Assistenten beder aldrig om adgangskoder" "wizard safety text"
Assert-Contains $SetupContent "Get-ActiveWifiSsid" "wizard uses Wi-Fi check"
Assert-Contains $SetupContent "winget install" "wizard can install SketchUp"
Assert-Contains $SetupContent "Start-Process" "wizard opens official links"

Assert-Contains $LauncherContent "powershell.exe" "standard Windows PowerShell launcher"
Assert-Contains $LauncherContent "setup-windows.ps1" "launcher target"
Assert-Contains $LauncherContent "chcp 65001" "launcher UTF-8 codepage"

Assert-Contains $BuildContent "GF2-IT-Setup-Windows.zip" "new Windows package name"
Assert-NotContains $LandingHtml "password" "English password wording"
Assert-NotContains $SetupContent "Read-Host `"Åbn" "old terminal prompt flow"

if (Test-Path -LiteralPath $ZipPath) {
    Assert-ZipContains $ZipPath @(
        "start.html",
        "index.html",
        "assets/neg-hero-transition.png",
        "scripts/setup-windows.ps1",
        "scripts/setup-config.ps1",
        "scripts/setup-checks.ps1",
        "Start Windows setup.cmd"
    )
}

Write-Host "Setup delivery checks passed."
```

- [ ] **Step 2: Run the test to verify it fails**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with `Missing required file: ...\index.html`.

- [ ] **Step 3: Commit the failing test**

Run:

```powershell
git add tests/check-setup-delivery.ps1
git commit -m "test: add setup delivery checks"
```

## Task 2: Add GitHub Pages Landing Page

**Files:**
- Create: `index.html`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Create the static landing page**

Create `index.html` with:

```html
<!DOCTYPE html>
<html lang="da">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GF2 IT Setup</title>
  <style>
    :root {
      --ink: #17202a;
      --muted: #506070;
      --line: #dce4ea;
      --paper: #f5f7f9;
      --surface: #ffffff;
      --neg-blue: #005aa7;
      --neg-blue-dark: #003f73;
      --lime: #d8e85a;
      --warning: #fff5d6;
      --warning-line: #efb43f;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      font-family: Arial, Helvetica, sans-serif;
      background: var(--paper);
      color: var(--ink);
      letter-spacing: 0;
    }

    a, button { font: inherit; }

    .hero {
      min-height: 520px;
      display: grid;
      align-items: end;
      background:
        linear-gradient(90deg, rgba(0, 63, 115, 0.96) 0%, rgba(0, 90, 167, 0.78) 42%, rgba(0, 90, 167, 0.16) 74%),
        url("assets/neg-hero-transition.png") center / cover no-repeat,
        var(--neg-blue);
      color: #fff;
    }

    .hero-inner {
      width: min(1120px, calc(100% - 40px));
      margin: 0 auto;
      padding: 56px 0;
    }

    .kicker {
      margin-bottom: 12px;
      color: var(--lime);
      font-size: 13px;
      font-weight: 900;
      text-transform: uppercase;
    }

    h1 {
      margin: 0;
      max-width: 720px;
      font-size: clamp(34px, 7vw, 64px);
      line-height: 1;
      letter-spacing: 0;
    }

    .hero p {
      max-width: 680px;
      margin: 18px 0 0;
      color: #eef6fb;
      font-size: 18px;
      line-height: 1.45;
    }

    main {
      width: min(1120px, calc(100% - 40px));
      margin: 0 auto;
      padding: 30px 0 44px;
    }

    .actions {
      display: grid;
      grid-template-columns: minmax(280px, 420px) minmax(240px, 1fr);
      gap: 18px;
      align-items: stretch;
    }

    .download-panel,
    .info-panel {
      background: var(--surface);
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 22px;
    }

    .download-button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      min-height: 58px;
      border-radius: 6px;
      background: var(--neg-blue);
      color: #fff;
      text-decoration: none;
      font-size: 19px;
      font-weight: 900;
    }

    .download-button:hover { background: var(--neg-blue-dark); }

    .secondary-links {
      display: grid;
      gap: 10px;
      margin-top: 14px;
    }

    .secondary-links a {
      color: var(--neg-blue);
      font-weight: 900;
    }

    .warning {
      margin: 16px 0 0;
      padding: 14px;
      border-left: 5px solid var(--warning-line);
      background: var(--warning);
      line-height: 1.45;
    }

    .steps {
      margin: 0;
      padding-left: 22px;
      color: var(--muted);
      line-height: 1.55;
    }

    .steps strong { color: var(--ink); }

    @media (max-width: 760px) {
      .hero { min-height: 440px; }
      .actions { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <section class="hero">
    <div class="hero-inner">
      <div class="kicker">GF2 IT-opstart</div>
      <h1>GF2 IT Setup</h1>
      <p>Kom sikkert i gang med Wi-Fi, Office, Moodle, PraxisOnline, Lectio, OneDrive og SketchUp.</p>
    </div>
  </section>

  <main>
    <section class="actions" aria-label="Start setup">
      <div class="download-panel">
        <a class="download-button" href="dist/GF2-IT-Setup-Windows.zip" download>Download Windows Setup</a>
        <div class="secondary-links">
          <a href="start.html">Åbn underviser-/projektorguide</a>
          <a href="https://github.com/Tech-Pleex/Startup-project-2026">Se kildekode på GitHub</a>
        </div>
        <div class="warning">
          Vi beder aldrig om adgangskoder, MitID eller UNI-Login. Du logger kun ind på officielle sider, som assistenten åbner.
        </div>
      </div>

      <div class="info-panel">
        <ol class="steps">
          <li><strong>Gå på NEG Guest.</strong> Brug gæstenettet til at hente setup-pakken.</li>
          <li><strong>Download Windows Setup.</strong> Åbn zip-filen og start setup-assistenten.</li>
          <li><strong>Følg trinnene.</strong> Assistenten sender dig til de rigtige sider og indstillinger.</li>
          <li><strong>Indtast kun login på officielle sider.</strong> Setup-assistenten gemmer ingen personlige oplysninger.</li>
        </ol>
      </div>
    </section>
  </main>
</body>
</html>
```

- [ ] **Step 2: Run the delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with `Missing required file: ...\scripts\setup-config.ps1`.

- [ ] **Step 3: Commit the landing page**

Run:

```powershell
git add index.html
git commit -m "feat: add setup landing page"
```

## Task 3: Add Setup Config

**Files:**
- Create: `scripts/setup-config.ps1`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Create central setup config**

Create `scripts/setup-config.ps1` as UTF-8 with BOM:

```powershell
$Script:SetupConfig = [ordered]@{
    Title = "GF2 IT Setup"
    DashboardFile = "start.html"
    DesktopShortcutName = "GF2 IT Dashboard.url"
    TargetWifi = "NEG"
    GuestWifi = "NEG Guest"
    SketchUpPackageId = "Trimble.SketchUp.2026"
    SafetyText = "Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login. Du indtaster kun oplysninger på officielle sider og i Windows' egne indstillinger."
}

$Script:SetupSteps = @(
    [ordered]@{
        Id = "welcome"
        Title = "Velkomst og sikkerhed"
        Kind = "manual"
        Body = "GF2 IT Setup hjælper dig de rigtige steder hen. Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login."
        Button = "Start"
    },
    [ordered]@{
        Id = "wifi"
        Title = "Wi-Fi"
        Kind = "wifi"
        Body = "Du skal ende på NEG-netværket. Hvis du allerede er på NEG, markerer assistenten trinnet som gennemført."
        Button = "Åbn Wi-Fi-indstillinger"
    },
    [ordered]@{
        Id = "office"
        Title = "Office og skolemail"
        Kind = "link"
        Body = "Log ind på Office og åbn din skolemail. Assistenten ser aldrig dine loginoplysninger."
        Url = "https://www.office.com/"
        Button = "Åbn Office"
    },
    [ordered]@{
        Id = "trimble"
        Title = "Trimble invitation"
        Kind = "link"
        Body = "Find mailen fra Trimble eller SketchUp i din skolemail, klik invitationslinket, og følg flowet med din skolemail."
        Url = "https://www.office.com/"
        Button = "Åbn skolemail"
    },
    [ordered]@{
        Id = "moodle"
        Title = "Moodle"
        Kind = "link"
        Body = "Åbn Moodle og kontroller at dine GF2-rum vises."
        Url = "https://online.neg.dk/login/index.php"
        Button = "Åbn Moodle"
    },
    [ordered]@{
        Id = "praxis"
        Title = "PraxisOnline"
        Kind = "link"
        Body = "Login sker via UNI-Login, hvor du vælger MitID."
        Warning = "Vigtigt: Brug din skolemail til PraxisOnline. Eksempel: neg04026@edu.neg.dk"
        Url = "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery"
        Button = "Åbn PraxisOnline"
    },
    [ordered]@{
        Id = "lectio"
        Title = "Lectio"
        Kind = "link"
        Body = "Log ind med UNI-Login, hvor du vælger MitID. Assistenten åbner kun Lectio og ser aldrig dine oplysninger."
        Url = "https://www.lectio.dk/lectio/769/default.aspx"
        Button = "Åbn Lectio"
    },
    [ordered]@{
        Id = "onedrive"
        Title = "OneDrive"
        Kind = "link"
        Body = "Åbn OneDrive via Office 365 og kontroller at du kan se dine filer."
        Url = "https://www.office.com/launch/onedrive"
        Button = "Åbn OneDrive"
    },
    [ordered]@{
        Id = "sketchup"
        Title = "SketchUp"
        Kind = "sketchup"
        Body = "Assistenten kan forsøge at installere SketchUp via winget. Hvis det ikke virker, bruger du manuel fallback."
        Url = "https://www.sketchup.com/"
        Button = "Installer SketchUp"
    },
    [ordered]@{
        Id = "finish"
        Title = "Færdig"
        Kind = "finish"
        Body = "Assistenten opretter en genvej til dashboardet og åbner dashboardet."
        Button = "Åbn dashboard"
    }
)
```

- [ ] **Step 2: Run the delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with `Missing required file: ...\scripts\setup-checks.ps1`.

- [ ] **Step 3: Commit config**

Run:

```powershell
git add scripts/setup-config.ps1
git commit -m "feat: add setup wizard config"
```

## Task 4: Add Setup Check Functions

**Files:**
- Create: `scripts/setup-checks.ps1`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Create check helper functions**

Create `scripts/setup-checks.ps1` as UTF-8 with BOM:

```powershell
function Get-ActiveWifiSsid {
    $Output = netsh wlan show interfaces 2>$null
    if (-not $Output) {
        return $null
    }

    foreach ($Line in $Output) {
        if ($Line -match "^\s*SSID\s*:\s*(.+)$" -and $Line -notmatch "BSSID") {
            return $Matches[1].Trim()
        }
    }

    return $null
}

function Open-WifiSettings {
    Start-Process "ms-settings:network-wifi"
}

function Test-WingetAvailable {
    return [bool](Get-Command winget -ErrorAction SilentlyContinue)
}

function Test-WindowsSMode {
    $Policy = Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\CI\Policy" -ErrorAction SilentlyContinue
    return ($Policy -and $Policy.SkuPolicyRequired -eq 1)
}

function Open-SetupLink {
    param([Parameter(Mandatory = $true)][string]$Url)
    Start-Process $Url
}

function New-DashboardShortcut {
    param(
        [Parameter(Mandatory = $true)][string]$DashboardPath,
        [Parameter(Mandatory = $true)][string]$ShortcutPath
    )

    if (-not (Test-Path -LiteralPath $DashboardPath)) {
        return $false
    }

    $DashboardUri = (Get-Item -LiteralPath $DashboardPath).FullName.Replace("\", "/")
    Set-Content -LiteralPath $ShortcutPath -Encoding ASCII -Value @(
        "[InternetShortcut]",
        "URL=file:///$DashboardUri"
    )
    return $true
}

function Install-SketchUpPackage {
    param([Parameter(Mandatory = $true)][string]$PackageId)

    if (-not (Test-WingetAvailable)) {
        return [ordered]@{ Success = $false; Message = "winget blev ikke fundet." }
    }

    $Arguments = @(
        "install",
        "--id", $PackageId,
        "-e",
        "--source", "winget",
        "--accept-source-agreements",
        "--accept-package-agreements"
    )

    $Process = Start-Process -FilePath "winget" -ArgumentList $Arguments -Wait -PassThru -WindowStyle Normal
    if ($Process.ExitCode -eq 0) {
        return [ordered]@{ Success = $true; Message = "SketchUp-installationen blev startet eller gennemført." }
    }

    return [ordered]@{ Success = $false; Message = "winget returnerede fejlkode $($Process.ExitCode)." }
}
```

- [ ] **Step 2: Run parser check**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "$Errors = $null; $null = [System.Management.Automation.PSParser]::Tokenize((Get-Content -Raw -Encoding UTF8 scripts/setup-checks.ps1), [ref]$Errors); if ($Errors) { $Errors; exit 1 } else { 'Setup checks parse passed' }"
```

Expected: PASS with `Setup checks parse passed`.

- [ ] **Step 3: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with missing Windows Forms GUI content in `scripts/setup-windows.ps1`.

- [ ] **Step 4: Commit helper functions**

Run:

```powershell
git add scripts/setup-checks.ps1
git commit -m "feat: add setup check helpers"
```

## Task 5: Replace Terminal Setup With Windows Forms Wizard

**Files:**
- Modify: `scripts/setup-windows.ps1`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Replace `scripts/setup-windows.ps1` with GUI wizard**

Replace the file with this structure, saved as UTF-8 with BOM:

```powershell
$ErrorActionPreference = "Continue"

$Utf8Encoding = New-Object System.Text.UTF8Encoding $false
[Console]::InputEncoding = $Utf8Encoding
[Console]::OutputEncoding = $Utf8Encoding
$OutputEncoding = $Utf8Encoding

$Root = Split-Path -Parent $PSScriptRoot
. (Join-Path $PSScriptRoot "setup-config.ps1")
. (Join-Path $PSScriptRoot "setup-checks.ps1")

Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing
[System.Windows.Forms.Application]::EnableVisualStyles()

$Dashboard = Join-Path $Root $SetupConfig.DashboardFile
$Desktop = [Environment]::GetFolderPath("Desktop")
$ShortcutPath = Join-Path $Desktop $SetupConfig.DesktopShortcutName
$CurrentStepIndex = 0

$Form = New-Object System.Windows.Forms.Form
$Form.Text = $SetupConfig.Title
$Form.StartPosition = "CenterScreen"
$Form.Size = New-Object System.Drawing.Size(760, 520)
$Form.MinimumSize = New-Object System.Drawing.Size(680, 460)
$Form.BackColor = [System.Drawing.Color]::White

$TitleLabel = New-Object System.Windows.Forms.Label
$TitleLabel.Font = New-Object System.Drawing.Font("Segoe UI", 18, [System.Drawing.FontStyle]::Bold)
$TitleLabel.Location = New-Object System.Drawing.Point(26, 22)
$TitleLabel.Size = New-Object System.Drawing.Size(690, 42)
$Form.Controls.Add($TitleLabel)

$ProgressLabel = New-Object System.Windows.Forms.Label
$ProgressLabel.Font = New-Object System.Drawing.Font("Segoe UI", 9, [System.Drawing.FontStyle]::Regular)
$ProgressLabel.ForeColor = [System.Drawing.Color]::FromArgb(80, 96, 112)
$ProgressLabel.Location = New-Object System.Drawing.Point(30, 68)
$ProgressLabel.Size = New-Object System.Drawing.Size(690, 24)
$Form.Controls.Add($ProgressLabel)

$BodyBox = New-Object System.Windows.Forms.TextBox
$BodyBox.Multiline = $true
$BodyBox.ReadOnly = $true
$BodyBox.BorderStyle = "None"
$BodyBox.BackColor = [System.Drawing.Color]::White
$BodyBox.Font = New-Object System.Drawing.Font("Segoe UI", 11)
$BodyBox.Location = New-Object System.Drawing.Point(30, 110)
$BodyBox.Size = New-Object System.Drawing.Size(680, 110)
$Form.Controls.Add($BodyBox)

$WarningBox = New-Object System.Windows.Forms.TextBox
$WarningBox.Multiline = $true
$WarningBox.ReadOnly = $true
$WarningBox.BorderStyle = "FixedSingle"
$WarningBox.BackColor = [System.Drawing.Color]::FromArgb(255, 245, 214)
$WarningBox.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$WarningBox.Location = New-Object System.Drawing.Point(30, 238)
$WarningBox.Size = New-Object System.Drawing.Size(680, 78)
$Form.Controls.Add($WarningBox)

$StatusLabel = New-Object System.Windows.Forms.Label
$StatusLabel.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$StatusLabel.ForeColor = [System.Drawing.Color]::FromArgb(0, 90, 167)
$StatusLabel.Location = New-Object System.Drawing.Point(30, 334)
$StatusLabel.Size = New-Object System.Drawing.Size(680, 34)
$Form.Controls.Add($StatusLabel)

$PrimaryButton = New-Object System.Windows.Forms.Button
$PrimaryButton.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$PrimaryButton.Location = New-Object System.Drawing.Point(30, 396)
$PrimaryButton.Size = New-Object System.Drawing.Size(210, 42)
$Form.Controls.Add($PrimaryButton)

$DoneButton = New-Object System.Windows.Forms.Button
$DoneButton.Text = "Færdig - næste"
$DoneButton.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$DoneButton.Location = New-Object System.Drawing.Point(500, 396)
$DoneButton.Size = New-Object System.Drawing.Size(210, 42)
$Form.Controls.Add($DoneButton)

function Get-CurrentStep {
    return $SetupSteps[$CurrentStepIndex]
}

function Set-StatusText {
    param([string]$Text)
    $StatusLabel.Text = $Text
}

function Render-Step {
    $Step = Get-CurrentStep
    $TitleLabel.Text = $Step.Title
    $ProgressLabel.Text = "Trin $($CurrentStepIndex + 1) af $($SetupSteps.Count)"
    $BodyBox.Text = $Step.Body
    $PrimaryButton.Text = $Step.Button
    $WarningBox.Visible = [bool]$Step.Warning
    $WarningBox.Text = if ($Step.Warning) { $Step.Warning } else { "" }
    Set-StatusText ""

    if ($Step.Kind -eq "wifi") {
        $Ssid = Get-ActiveWifiSsid
        if ($Ssid -eq $SetupConfig.TargetWifi) {
            Set-StatusText "Du er allerede på $($SetupConfig.TargetWifi). Trinnet er gennemført."
        }
        elseif ($Ssid -eq $SetupConfig.GuestWifi) {
            Set-StatusText "Du er på $($SetupConfig.GuestWifi). Skift til $($SetupConfig.TargetWifi), og klik derefter Færdig."
        }
        elseif ($Ssid) {
            Set-StatusText "Aktivt Wi-Fi: $Ssid. Skift til $($SetupConfig.TargetWifi), hvis det ikke er korrekt."
        }
        else {
            Set-StatusText "Assistenten kunne ikke læse aktivt Wi-Fi. Brug Windows Wi-Fi-indstillinger."
        }
    }

    if ($Step.Kind -eq "sketchup") {
        if (Test-WindowsSMode) {
            Set-StatusText "Mulig Windows S-mode fundet. SketchUp-installation kan være blokeret."
        }
        elseif (Test-WingetAvailable) {
            Set-StatusText "winget er fundet. Assistenten kan forsøge at installere $($SetupConfig.SketchUpPackageId)."
        }
        else {
            Set-StatusText "winget blev ikke fundet. Brug manuel SketchUp-fallback."
        }
    }

    if ($Step.Kind -eq "finish") {
        $PrimaryButton.Text = "Åbn dashboard"
        $DoneButton.Text = "Luk"
    }
}

function Invoke-PrimaryAction {
    $Step = Get-CurrentStep

    if ($Step.Kind -eq "wifi") {
        Open-WifiSettings
        Set-StatusText "Windows Wi-Fi-indstillinger er åbnet. Vælg NEG og indtast selv dine oplysninger."
        return
    }

    if ($Step.Kind -eq "link") {
        Open-SetupLink -Url $Step.Url
        Set-StatusText "Siden er åbnet i browseren. Log selv ind på den officielle side og klik Færdig bagefter."
        return
    }

    if ($Step.Kind -eq "sketchup") {
        if (Test-WindowsSMode) {
            Open-SetupLink -Url $Step.Url
            Set-StatusText "Manuel SketchUp-side er åbnet, fordi Windows S-mode kan blokere installation."
            return
        }

        $Result = Install-SketchUpPackage -PackageId $SetupConfig.SketchUpPackageId
        if (-not $Result.Success) {
            Open-SetupLink -Url $Step.Url
        }
        Set-StatusText $Result.Message
        return
    }

    if ($Step.Kind -eq "finish") {
        New-DashboardShortcut -DashboardPath $Dashboard -ShortcutPath $ShortcutPath | Out-Null
        if (Test-Path -LiteralPath $Dashboard) {
            Start-Process $Dashboard
        }
        Set-StatusText "Dashboardet er åbnet, og setup er færdigt."
        return
    }

    Set-StatusText $SetupConfig.SafetyText
}

$PrimaryButton.Add_Click({ Invoke-PrimaryAction })

$DoneButton.Add_Click({
    if ($CurrentStepIndex -lt ($SetupSteps.Count - 1)) {
        $Script:CurrentStepIndex++
        Render-Step
    }
    else {
        $Form.Close()
    }
})

Render-Step
[void]$Form.ShowDialog()
```

- [ ] **Step 2: Run parser check**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "$Errors = $null; $null = [System.Management.Automation.PSParser]::Tokenize((Get-Content -Raw -Encoding UTF8 scripts/setup-windows.ps1), [ref]$Errors); if ($Errors) { $Errors; exit 1 } else { 'Setup wizard parse passed' }"
```

Expected: PASS with `Setup wizard parse passed`.

- [ ] **Step 3: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with package-name checks in `scripts/build-package.ps1`.

- [ ] **Step 4: Commit wizard**

Run:

```powershell
git add scripts/setup-windows.ps1
git commit -m "feat: add Windows setup wizard UI"
```

## Task 6: Update Windows Launcher

**Files:**
- Modify: `Start Windows setup.cmd`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Update CMD launcher**

Replace `Start Windows setup.cmd` with:

```cmd
@echo off
setlocal
chcp 65001 >nul
cd /d "%~dp0"
start "" powershell.exe -NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File "%~dp0scripts\setup-windows.ps1"
endlocal
```

- [ ] **Step 2: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with package-name checks in `scripts/build-package.ps1`.

- [ ] **Step 3: Commit launcher**

Run:

```powershell
git add "Start Windows setup.cmd"
git commit -m "fix: launch setup wizard from CMD"
```

## Task 7: Update Package Builder

**Files:**
- Modify: `scripts/build-package.ps1`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Update package file list and zip name**

Modify the top of `scripts/build-package.ps1` so it uses:

```powershell
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Dist = Join-Path $Root "dist"
$Zip = Join-Path $Dist "GF2-IT-Setup-Windows.zip"

$Required = @(
    "index.html",
    "start.html",
    "assets\neg-hero-transition.png",
    "scripts\setup-windows.ps1",
    "scripts\setup-config.ps1",
    "scripts\setup-checks.ps1",
    "Start Windows setup.cmd"
)
```

Keep the existing `Assert-InDirectory`, required-file loop, zip creation code, and `Write-Host "Created $Zip"` behavior unchanged.

- [ ] **Step 2: Build the package**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

Expected: PASS with `Created ...\dist\GF2-IT-Setup-Windows.zip`.

- [ ] **Step 3: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: PASS with `Setup delivery checks passed.`

- [ ] **Step 4: Commit package builder and zip**

Run:

```powershell
git add scripts/build-package.ps1 dist/GF2-IT-Setup-Windows.zip
git commit -m "build: package Windows setup wizard"
```

## Task 8: Update Dashboard Test Coverage

**Files:**
- Modify: `tests/check-dashboard.ps1`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Keep dashboard test compatible with new Windows setup**

In `tests/check-dashboard.ps1`, replace the old terminal prompt assertion:

```powershell
Assert-Contains $WindowsSetupContent "Read-Host `"Åbn" "asks before opening links"
```

with:

```powershell
Assert-Contains $WindowsSetupContent "System.Windows.Forms" "Windows setup GUI"
Assert-Contains $WindowsSetupContent "Assistenten beder aldrig om adgangskoder" "Windows setup safety"
```

Add new file paths near the existing Windows setup variables:

```powershell
$SetupConfig = Join-Path $Root "scripts\setup-config.ps1"
$SetupChecks = Join-Path $Root "scripts\setup-checks.ps1"
```

Add required file checks after `$WindowsSetup`:

```powershell
Assert-File $SetupConfig
Assert-File $SetupChecks
```

- [ ] **Step 2: Run dashboard test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 3: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: PASS with `Setup delivery checks passed.`

- [ ] **Step 4: Commit test compatibility**

Run:

```powershell
git add tests/check-dashboard.ps1 tests/check-setup-delivery.ps1
git commit -m "test: cover setup wizard delivery"
```

## Task 9: Add Text Handling Rules

**Files:**
- Modify: `.gitattributes`
- Test: both PowerShell tests

- [ ] **Step 1: Update `.gitattributes`**

Replace `.gitattributes` with:

```gitattributes
*.cmd text eol=crlf
*.html text eol=lf
*.md text eol=lf
*.ps1 text eol=crlf
*.sh text eol=lf
*.command text eol=lf
```

- [ ] **Step 2: Run tests**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected:

```text
Dashboard checks passed.
Setup delivery checks passed.
```

- [ ] **Step 3: Commit attributes**

Run:

```powershell
git add .gitattributes
git commit -m "chore: define setup file line endings"
```

## Task 10: Add README For Colleagues And Future Public Repo

**Files:**
- Create: `README.md`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Create README**

Create `README.md`:

```markdown
# GF2 IT Setup

GF2 IT Setup is a Windows-first helper for new GF2 students. It guides students through Wi-Fi, Office, school mail, Moodle, PraxisOnline, Lectio, OneDrive, Trimble/SketchUp, and the local dashboard.

## Safety

The setup assistant never asks for passwords, MitID, or UNI-Login.

Students only enter login information on official pages or in Windows' own settings. The assistant opens the right places and asks the student to confirm when a step is done.

## Current Status

This repository is still a prototype. The first supported platform is Windows.

Mac support, signed installers, central classroom status, Intune deployment, and official NEG hosting are not part of the first version.

## Local Checks

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

## Build Windows Package

Run:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

The package is created at:

```text
dist/GF2-IT-Setup-Windows.zip
```

## Student Delivery

The intended student flow is:

1. Connect to NEG Guest.
2. Open the GitHub Pages landing page.
3. Download Windows Setup.
4. Open the downloaded zip.
5. Start `Start Windows setup.cmd`.
6. Follow the setup assistant.
```

- [ ] **Step 2: Add README assertions to delivery test**

In `tests/check-setup-delivery.ps1`, add after `$Dashboard`:

```powershell
$Readme = Join-Path $Root "README.md"
```

Add after `Assert-File $Dashboard`:

```powershell
Assert-File $Readme
```

Add after content reads:

```powershell
$ReadmeContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $Readme
```

Add assertions:

```powershell
Assert-Contains $ReadmeContent "never asks for passwords, MitID, or UNI-Login" "README safety"
Assert-Contains $ReadmeContent "dist/GF2-IT-Setup-Windows.zip" "README package path"
Assert-Contains $ReadmeContent "Start Windows setup.cmd" "README launcher"
```

- [ ] **Step 3: Run tests**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected:

```text
Dashboard checks passed.
Setup delivery checks passed.
```

- [ ] **Step 4: Commit README**

Run:

```powershell
git add README.md tests/check-setup-delivery.ps1
git commit -m "docs: explain setup safety and delivery"
```

## Task 11: Manual GUI Smoke Test

**Files:**
- No code changes unless the smoke test exposes a defect.

- [ ] **Step 1: Start the wizard from PowerShell**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/setup-windows.ps1
```

Expected:

```text
No terminal prompt flow appears.
A Windows Forms window titled GF2 IT Setup opens.
```

- [ ] **Step 2: Click through non-destructive steps**

Manual expected behavior:

```text
Welcome shows safety text.
Wi-Fi step shows active Wi-Fi status or manual fallback.
Office button opens office.com.
PraxisOnline step shows yellow schoolmail warning.
Lectio step mentions UNI-Login and MitID.
SketchUp step does not run winget until the Install SketchUp button is clicked.
Finish opens start.html and creates GF2 IT Dashboard.url on Desktop.
```

- [ ] **Step 3: Re-run tests after smoke test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected:

```text
Dashboard checks passed.
Setup delivery checks passed.
```

- [ ] **Step 4: Commit any smoke-test fixes**

If no code changes were needed, do not commit. If a defect was fixed, run:

```powershell
git add scripts/setup-windows.ps1 scripts/setup-checks.ps1 tests/check-setup-delivery.ps1
git commit -m "fix: polish setup wizard smoke test issues"
```

## Task 12: Final Verification

**Files:**
- No code changes expected.

- [ ] **Step 1: Rebuild package**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

Expected:

```text
Created C:\Users\jere\Documents\Neg_Ai_Stuff\Startup-project-2026\dist\GF2-IT-Setup-Windows.zip
```

- [ ] **Step 2: Run all checks**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected:

```text
Dashboard checks passed.
Setup delivery checks passed.
```

- [ ] **Step 3: Check Git status**

Run:

```powershell
git status --short --branch
```

Expected:

```text
## feature/setup-wizard
```

or a clean branch with upstream information if the branch has been pushed.

## Self-Review

- Spec coverage: landing page, GitHub Pages delivery, Windows-only first version, no-password safety model, Wi-Fi, PraxisOnline schoolmail warning, Lectio UNI-Login/MitID text, SketchUp `winget`, packaging, and dashboard fallback are covered.
- Scope: Mac, signed `.exe`, official NEG domain, Intune, central class status, and automatic login remain out of first version.
- Test strategy: delivery test covers landing page/package/wizard text and existing dashboard test remains active.
- Type consistency: config keys used by wizard are `SetupConfig`, `SetupSteps`, `Kind`, `Title`, `Body`, `Warning`, `Url`, and `Button` throughout.

