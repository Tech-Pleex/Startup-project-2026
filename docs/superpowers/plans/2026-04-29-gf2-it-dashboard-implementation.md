# GF2 IT Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the offline GF2 IT Dashboard as `start.html`, package it as `IT-opstart-GF2.zip`, and include local assets plus Windows/Mac setup assistants.

**Architecture:** The dashboard is a single offline-first HTML file with embedded CSS/JS and one local image dependency at `assets/neg-hero-transition.png`. Setup automation lives beside it in `scripts/` and is launched through double-clickable wrapper files in the zip root. Packaging is handled by a PowerShell build script that validates required files before creating the zip.

**Tech Stack:** HTML, CSS, vanilla JavaScript, browser `localStorage`, PowerShell 5+, Windows batch, POSIX shell, PowerShell `Compress-Archive`.

---

## File Structure

- Create: `start.html`
  - Offline dashboard, embedded CSS/JS, local hero image reference, student/teacher mode, platform selection, local status storage, manual guide, fixed links.
- Keep: `assets/neg-hero-transition.png`
  - Required local hero image. No remote image references are allowed in production UI.
- Create: `scripts/setup-windows.ps1`
  - Windows assistant. Opens selected school/service pages on explicit script launch, checks common system conditions, creates a desktop shortcut to `start.html`.
- Create: `scripts/setup-mac.sh`
  - Mac assistant. Opens selected pages, prints manual fallback guidance, avoids storing credentials.
- Create: `Start Windows setup.cmd`
  - Double-click launcher that runs `scripts/setup-windows.ps1`.
- Create: `Start Mac setup.command`
  - Double-click launcher that runs `scripts/setup-mac.sh`.
- Create: `scripts/build-package.ps1`
  - Validates required files and creates `dist/IT-opstart-GF2.zip`.
- Create: `tests/check-dashboard.ps1`
  - Local static checks for required text, local asset usage, no remote image references, localStorage status keys, and zip contents.
- Create: `dist/`
  - Generated output directory. Do not hand-edit files in this directory.

## Implementation Decisions

- Use a single `start.html` as both the light version and the dashboard inside the zip.
- Use actual first-version URLs in one data object inside `start.html` and matching arrays in scripts:
  - Office 365 / skolemail: `https://www.office.com/`
  - Moodle: `https://online.neg.dk/login/index.php`
  - Lectio: `https://www.lectio.dk/lectio/769/default.aspx`
  - PraxisOnline: `https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery`
  - OneDrive: `https://www.office.com/launch/onedrive`
  - SketchUp / Trimble: `https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg`
  - NEG hjemmeside: `https://www.neg.dk/`
- External pages open only when the user clicks a link or actively starts a setup assistant.
- Store only these local browser values:
  - `gf2-it-dashboard.platform`
  - `gf2-it-dashboard.mode`
  - `gf2-it-dashboard.studentStatus`
  - `gf2-it-dashboard.teacherStatus`
- Never store usernames, passwords, MitID data, or personal login details.

---

### Task 1: Add Static Verification Script

**Files:**
- Create: `tests/check-dashboard.ps1`

- [ ] **Step 1: Create the test directory**

Run:

```powershell
New-Item -ItemType Directory -Force tests
```

Expected: PowerShell prints a `tests` directory path.

- [ ] **Step 2: Create static dashboard checks**

Create `tests/check-dashboard.ps1` with this content:

```powershell
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$StartHtml = Join-Path $Root "start.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$ZipPath = Join-Path $Root "dist\IT-opstart-GF2.zip"

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
    if ($Content -notlike "*$Needle*") {
        throw "Missing required content '$Label': $Needle"
    }
}

function Assert-NotContains {
    param(
        [string]$Content,
        [string]$Needle,
        [string]$Label
    )
    if ($Content -like "*$Needle*") {
        throw "Forbidden content '$Label': $Needle"
    }
}

Assert-File $StartHtml
Assert-File $HeroImage

$Html = Get-Content -Raw -LiteralPath $StartHtml

Assert-Contains $Html "GF2 IT Dashboard" "dashboard title"
Assert-Contains $Html "ET FÆLLESSKAB MED PLADS TIL DIG" "hero kicker"
Assert-Contains $Html "Kom sikkert i gang med IT på GF2" "hero heading"
Assert-Contains $Html "assets/neg-hero-transition.png" "local hero asset"
Assert-Contains $Html "Elev" "student mode"
Assert-Contains $Html "Underviser" "teacher mode"
Assert-Contains $Html "Ikke startet" "status label"
Assert-Contains $Html "I gang" "status label"
Assert-Contains $Html "Færdig" "status label"
Assert-Contains $Html "gf2-it-dashboard.studentStatus" "student localStorage key"
Assert-Contains $Html "gf2-it-dashboard.teacherStatus" "teacher localStorage key"
Assert-Contains $Html "Office 365 / skolemail" "fixed link"
Assert-Contains $Html "Moodle" "fixed link"
Assert-Contains $Html "Lectio" "fixed link"
Assert-Contains $Html "PraxisOnline" "fixed link"
Assert-Contains $Html "OneDrive" "fixed link"
Assert-Contains $Html "SketchUp / Trimble" "fixed link"
Assert-Contains $Html "NEG hjemmeside" "fixed link"
Assert-Contains $Html "+45 72 290 100" "absence phone"
Assert-NotContains $Html "url(`"http" "remote CSS image"
Assert-NotContains $Html "src=`"http" "remote image/script source"
Assert-NotContains $Html "password" "credential wording"
Assert-NotContains $Html "adgangskode gemmes" "stored password wording"

if (Test-Path -LiteralPath $ZipPath) {
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $Zip = [System.IO.Compression.ZipFile]::OpenRead($ZipPath)
    try {
        $Entries = $Zip.Entries.FullName
        foreach ($Required in @(
            "start.html",
            "assets/neg-hero-transition.png",
            "scripts/setup-windows.ps1",
            "scripts/setup-mac.sh",
            "Start Windows setup.cmd",
            "Start Mac setup.command"
        )) {
            if ($Entries -notcontains $Required) {
                throw "Zip is missing entry: $Required"
            }
        }
    }
    finally {
        $Zip.Dispose()
    }
}

Write-Host "Dashboard checks passed."
```

- [ ] **Step 3: Run the test to verify it fails before implementation**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: FAIL with `Missing required file: ...\start.html`.

- [ ] **Step 4: Commit**

```powershell
git add tests/check-dashboard.ps1
git commit -m "test: add dashboard static checks"
```

---

### Task 2: Build Offline `start.html`

**Files:**
- Create: `start.html`
- Read: `assets/neg-hero-transition.png`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Create the offline dashboard file**

Create `start.html` with this complete first version:

```html
<!DOCTYPE html>
<html lang="da">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GF2 IT Dashboard</title>
  <style>
    :root {
      --ink: #181818;
      --muted: #5c6670;
      --line: #dce4ea;
      --paper: #f4f7f9;
      --surface: #ffffff;
      --neg-blue: #005aa7;
      --neg-blue-dark: #003f73;
      --neg-blue-soft: #e7f2fb;
      --lime: #d8e85a;
      --amber: #efb43f;
      --green: #237a4b;
      --red: #b84232;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      font-family: Arial, Helvetica, sans-serif;
      background: var(--paper);
      color: var(--ink);
      letter-spacing: 0;
    }

    button, a, select { font: inherit; }

    .topbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 18px;
      padding: 16px 28px;
      background: var(--surface);
      border-bottom: 1px solid var(--line);
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 14px;
      min-width: 0;
      font-weight: 900;
    }

    .brand-mark {
      display: grid;
      place-items: center;
      width: 54px;
      height: 54px;
      border-radius: 4px;
      background: var(--neg-blue);
      color: #fff;
      flex: 0 0 auto;
    }

    .brand small {
      display: block;
      margin-top: 3px;
      color: var(--muted);
      font-weight: 700;
    }

    .mode-toggle {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 4px;
      padding: 5px;
      border: 1px solid var(--line);
      border-radius: 999px;
      background: #eef3f6;
    }

    .mode-toggle button {
      border: 0;
      border-radius: 999px;
      padding: 9px 14px;
      background: transparent;
      color: var(--muted);
      font-weight: 900;
      cursor: pointer;
    }

    .mode-toggle button[aria-pressed="true"] {
      background: var(--neg-blue);
      color: #fff;
    }

    .hero {
      min-height: 350px;
      display: grid;
      grid-template-columns: minmax(300px, 0.5fr) minmax(360px, 0.5fr);
      background:
        linear-gradient(90deg, var(--neg-blue) 0%, var(--neg-blue) 41%, rgba(0,90,167,0.82) 58%, rgba(0,90,167,0.1) 80%, rgba(0,90,167,0) 100%),
        url("assets/neg-hero-transition.png") right center / auto 100% no-repeat,
        var(--neg-blue);
      color: #fff;
      overflow: hidden;
    }

    .hero-copy {
      align-self: center;
      padding: 52px 58px;
    }

    .kicker {
      margin-bottom: 14px;
      color: var(--lime);
      font-size: 13px;
      font-weight: 900;
      text-transform: uppercase;
    }

    h1 {
      margin: 0;
      max-width: 680px;
      font-size: 46px;
      line-height: 1.04;
      letter-spacing: 0;
    }

    .hero p {
      margin: 18px 0 0;
      max-width: 620px;
      color: #eef5f8;
      font-size: 16px;
      line-height: 1.45;
    }

    .wrap {
      max-width: 1260px;
      margin: 0 auto;
      padding: 26px;
    }

    .dashboard {
      display: grid;
      grid-template-columns: 300px minmax(0, 1fr) 280px;
      gap: 18px;
      align-items: start;
    }

    .panel {
      background: var(--surface);
      border: 1px solid var(--line);
      border-radius: 6px;
      overflow: hidden;
    }

    .panel-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 15px 16px;
      border-bottom: 1px solid var(--line);
      font-weight: 900;
    }

    .panel-body { padding: 16px; }

    .platforms {
      display: grid;
      gap: 10px;
    }

    .platform {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      width: 100%;
      padding: 13px;
      border: 2px solid var(--line);
      border-radius: 6px;
      background: #fff;
      color: var(--ink);
      font-weight: 900;
      cursor: pointer;
    }

    .platform[aria-pressed="true"] {
      border-color: var(--neg-blue);
      background: var(--neg-blue-soft);
    }

    .primary-action,
    .secondary-action {
      display: block;
      width: 100%;
      margin-top: 12px;
      border-radius: 4px;
      padding: 13px 14px;
      text-align: center;
      text-decoration: none;
      font-weight: 900;
    }

    .primary-action {
      border: 0;
      background: var(--neg-blue);
      color: #fff;
    }

    .secondary-action {
      border: 1px solid var(--ink);
      background: #fff;
      color: var(--ink);
    }

    .notice {
      margin-top: 14px;
      padding: 13px;
      border-left: 5px solid var(--amber);
      background: #fff8e6;
      color: #463914;
      line-height: 1.42;
      font-size: 14px;
    }

    .steps {
      display: grid;
      gap: 12px;
    }

    .step {
      display: grid;
      grid-template-columns: 40px minmax(0, 1fr) 128px;
      gap: 14px;
      align-items: start;
      padding: 15px;
      border: 1px solid var(--line);
      border-radius: 6px;
      background: #fff;
    }

    .num {
      display: grid;
      place-items: center;
      width: 38px;
      height: 38px;
      border-radius: 4px;
      background: var(--neg-blue);
      color: #fff;
      font-weight: 900;
    }

    .step h3 {
      margin: 0 0 6px;
      font-size: 18px;
      line-height: 1.2;
    }

    .step p {
      margin: 0;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.42;
    }

    .status-control {
      display: grid;
      gap: 6px;
    }

    .status-control label {
      font-size: 12px;
      color: var(--muted);
      font-weight: 900;
    }

    .status-control select {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 4px;
      padding: 8px;
      background: #f6f8f9;
      color: var(--ink);
      font-weight: 800;
    }

    .status-control select[data-status="I gang"] {
      border-color: var(--amber);
      background: #fff4d8;
    }

    .status-control select[data-status="Færdig"] {
      border-color: var(--green);
      background: #e7f6ed;
      color: var(--green);
    }

    .quicklinks {
      display: grid;
      gap: 9px;
    }

    .quicklink {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      padding: 12px 13px;
      border: 1px solid var(--line);
      border-radius: 4px;
      background: #fff;
      color: var(--ink);
      text-decoration: none;
      font-weight: 900;
    }

    .quicklink span:last-child {
      color: var(--neg-blue);
    }

    .manual {
      margin-top: 18px;
    }

    .manual summary {
      cursor: pointer;
      font-weight: 900;
    }

    .manual ol {
      margin: 12px 0 0;
      padding-left: 22px;
      color: var(--muted);
      line-height: 1.5;
    }

    .hidden { display: none; }

    @media (max-width: 1000px) {
      .topbar { align-items: flex-start; flex-direction: column; }
      .hero { grid-template-columns: 1fr; }
      .hero-copy { padding: 42px 26px; }
      h1 { font-size: 34px; }
      .dashboard { grid-template-columns: 1fr; }
      .step { grid-template-columns: 40px minmax(0, 1fr); }
      .status-control { grid-column: 2; }
    }
  </style>
</head>
<body>
  <header class="topbar">
    <div class="brand">
      <div class="brand-mark" aria-hidden="true">NEG</div>
      <div>
        GF2 IT Dashboard
        <small>Opstart, links og lokal status</small>
      </div>
    </div>
    <div class="mode-toggle" aria-label="Skift mellem elev og underviser">
      <button id="studentMode" type="button" aria-pressed="true">Elev</button>
      <button id="teacherMode" type="button" aria-pressed="false">Underviser</button>
    </div>
  </header>

  <section class="hero">
    <div class="hero-copy">
      <div class="kicker">ET FÆLLESSKAB MED PLADS TIL DIG</div>
      <h1>Kom sikkert i gang med IT på GF2</h1>
      <p>Vælg Windows eller Mac, følg rækkefølgen første skoledag, og brug dashboardet igen når du skal finde Office, Moodle, Lectio, PraxisOnline, OneDrive eller SketchUp.</p>
    </div>
  </section>

  <main class="wrap">
    <div class="dashboard">
      <aside class="panel">
        <div class="panel-header">Vælg computer</div>
        <div class="panel-body">
          <div class="platforms">
            <button class="platform" type="button" data-platform="windows" aria-pressed="true">Windows <span>Valgt</span></button>
            <button class="platform" type="button" data-platform="mac" aria-pressed="false">Mac <span>Skift</span></button>
          </div>
          <a id="setupLink" class="primary-action" href="Start Windows setup.cmd">Start setup-assistent</a>
          <a class="secondary-action" href="#manualGuide">Manuel vejledning</a>
          <div class="notice">Chromebook kan ikke bruges til dette GF2-flow. Windows må ikke være i S-mode, hvis programmer skal installeres.</div>
        </div>
      </aside>

      <section class="panel" aria-live="polite">
        <div class="panel-header">
          <span id="mainTitle">Opstart første dag</span>
          <span>Status gemmes lokalt</span>
        </div>
        <div class="panel-body">
          <div id="studentSteps" class="steps"></div>
          <div id="teacherSteps" class="steps hidden"></div>
          <details id="manualGuide" class="manual">
            <summary>Manuel vejledning</summary>
            <ol>
              <li>Få dit NEG-brugernavn og din adgangskode udleveret af skolen.</li>
              <li>Log på skolens Wi-Fi med dit NEG-login.</li>
              <li>Åbn Office 365, tjek skolemail, og find SketchUp/Trimble-invitationen.</li>
              <li>Åbn Moodle med samme NEG-login.</li>
              <li>Åbn Lectio og PraxisOnline med MitID.</li>
              <li>Gem dine opgaver i OneDrive via Office 365.</li>
            </ol>
          </details>
        </div>
      </section>

      <aside class="panel">
        <div class="panel-header">Faste links</div>
        <div class="panel-body">
          <div id="quickLinks" class="quicklinks"></div>
          <div class="notice"><strong>Syg eller fravær?</strong><br>Ring til +45 72 290 100 mellem kl. 8.00-9.00.</div>
          <div class="notice">Dashboardet gemmer kun status lokalt i browseren. Det gemmer aldrig brugernavne, adgangskoder eller MitID-oplysninger.</div>
        </div>
      </aside>
    </div>
  </main>

  <script>
    const links = [
      ["Office 365 / skolemail", "https://www.office.com/"],
      ["Moodle", "https://online.neg.dk/login/index.php"],
      ["Lectio", "https://www.lectio.dk/lectio/769/default.aspx"],
      ["PraxisOnline", "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery"],
      ["OneDrive", "https://www.office.com/launch/onedrive"],
      ["SketchUp / Trimble", "https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg"],
      ["NEG hjemmeside", "https://www.neg.dk/"]
    ];

    const studentSteps = [
      ["Få NEG-login", "Få udleveret NEG-brugernavn og adgangskode. Skriv ikke oplysningerne i dashboardet."],
      ["Log på skolens Wi-Fi", "Brug dit NEG-login til skolens netværk før du åbner de øvrige systemer."],
      ["Åbn Office 365 og skolemail", "Log på Office.com og kontroller at du kan læse skolemail."],
      ["Find SketchUp/Trimble-invitationen", "Åbn skolemailen og find invitationen til SketchUp/Trimble-flowet."],
      ["Log på Moodle", "Brug samme NEG-login til Moodle og tjek at dine GF2-rum vises."],
      ["Log på PraxisOnline og Lectio", "Brug MitID til PraxisOnline og Lectio. Dashboardet kan kun åbne siderne."],
      ["Installer eller åbn SketchUp", "Installer SketchUp hvis nødvendigt, og opret eller log ind på Trimble-profil via skolemail-flowet."],
      ["Gem opgaver i OneDrive", "Brug OneDrive via Office 365 som fast sted til opgaver og filer."]
    ];

    const teacherSteps = [
      ["Tjek Wi-Fi", "Bekræft at eleven er på skolens Wi-Fi med NEG-login."],
      ["Bekræft Office og skolemail", "Kontroller at Office.com virker, og at eleven kan åbne skolemail."],
      ["Bekræft Moodle", "Kontroller at Moodle virker med elevens NEG-login."],
      ["Forklar SketchUp og Trimble", "Scriptet kan hjælpe med installation, men Trimble/licens kræver elevens skolemail-flow."],
      ["Forklar MitID-sporet", "Lectio og PraxisOnline kræver MitID. Scriptet kan kun åbne de rigtige sider."],
      ["Bekræft OneDrive", "Eleven skal bruge OneDrive til filer. FileCloud/P-drev er ikke en del af elevflowet."],
      ["Tjek typiske fejl", "Windows S-mode, Chromebook, forkert login, manglende skolemail og manglende MitID er de hyppigste blokeringer."]
    ];

    const storage = {
      platform: "gf2-it-dashboard.platform",
      mode: "gf2-it-dashboard.mode",
      student: "gf2-it-dashboard.studentStatus",
      teacher: "gf2-it-dashboard.teacherStatus"
    };

    const statusLabels = ["Ikke startet", "I gang", "Færdig"];

    function readJson(key, fallback) {
      try {
        const value = JSON.parse(localStorage.getItem(key));
        return value || fallback;
      } catch {
        return fallback;
      }
    }

    function writeJson(key, value) {
      localStorage.setItem(key, JSON.stringify(value));
    }

    function renderLinks() {
      const holder = document.getElementById("quickLinks");
      holder.innerHTML = "";
      links.forEach(([label, url]) => {
        const link = document.createElement("a");
        link.className = "quicklink";
        link.href = url;
        link.target = "_blank";
        link.rel = "noopener noreferrer";
        link.innerHTML = `<span>${label}</span><span>Åbn</span>`;
        holder.appendChild(link);
      });
    }

    function renderSteps(targetId, items, key) {
      const holder = document.getElementById(targetId);
      const statuses = readJson(key, {});
      holder.innerHTML = "";
      items.forEach(([title, body], index) => {
        const article = document.createElement("article");
        article.className = "step";
        const current = statuses[index] || "Ikke startet";
        article.innerHTML = `
          <div class="num">${index + 1}</div>
          <div>
            <h3>${title}</h3>
            <p>${body}</p>
          </div>
          <div class="status-control">
            <label for="${targetId}-${index}">Status</label>
            <select id="${targetId}-${index}" data-status="${current}">
              ${statusLabels.map(label => `<option value="${label}" ${label === current ? "selected" : ""}>${label}</option>`).join("")}
            </select>
          </div>
        `;
        const select = article.querySelector("select");
        select.addEventListener("change", () => {
          statuses[index] = select.value;
          select.dataset.status = select.value;
          writeJson(key, statuses);
        });
        holder.appendChild(article);
      });
    }

    function setPlatform(platform) {
      localStorage.setItem(storage.platform, platform);
      document.querySelectorAll("[data-platform]").forEach(button => {
        const active = button.dataset.platform === platform;
        button.setAttribute("aria-pressed", String(active));
        button.querySelector("span").textContent = active ? "Valgt" : "Skift";
      });
      document.getElementById("setupLink").href = platform === "mac" ? "Start Mac setup.command" : "Start Windows setup.cmd";
    }

    function setMode(mode) {
      localStorage.setItem(storage.mode, mode);
      document.getElementById("studentMode").setAttribute("aria-pressed", String(mode === "student"));
      document.getElementById("teacherMode").setAttribute("aria-pressed", String(mode === "teacher"));
      document.getElementById("studentSteps").classList.toggle("hidden", mode !== "student");
      document.getElementById("teacherSteps").classList.toggle("hidden", mode !== "teacher");
      document.getElementById("mainTitle").textContent = mode === "teacher" ? "Underviser-tjekliste" : "Opstart første dag";
    }

    document.querySelectorAll("[data-platform]").forEach(button => {
      button.addEventListener("click", () => setPlatform(button.dataset.platform));
    });

    document.getElementById("studentMode").addEventListener("click", () => setMode("student"));
    document.getElementById("teacherMode").addEventListener("click", () => setMode("teacher"));

    renderLinks();
    renderSteps("studentSteps", studentSteps, storage.student);
    renderSteps("teacherSteps", teacherSteps, storage.teacher);
    setPlatform(localStorage.getItem(storage.platform) || "windows");
    setMode(localStorage.getItem(storage.mode) || "student");
  </script>
</body>
</html>
```

- [ ] **Step 2: Run static checks**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 3: Open the light version locally**

Run:

```powershell
Start-Process .\start.html
```

Expected: Browser opens the offline dashboard from the local filesystem. Do not click external links during this verification.

- [ ] **Step 4: Manual UI check**

Verify:

- Hero uses the local NEG image on the right.
- Elev mode is active by default.
- Underviser mode swaps the center content to the teacher checklist.
- Windows/Mac buttons change the setup-assistant link.
- Status dropdowns retain values after browser refresh.
- Right-side links only open when clicked.

- [ ] **Step 5: Commit**

```powershell
git add start.html
git commit -m "feat: add offline GF2 IT dashboard"
```

---

### Task 3: Add Windows Setup Assistant

**Files:**
- Create: `scripts/setup-windows.ps1`
- Create: `Start Windows setup.cmd`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Create scripts directory**

Run:

```powershell
New-Item -ItemType Directory -Force scripts
```

Expected: PowerShell prints a `scripts` directory path.

- [ ] **Step 2: Create `scripts/setup-windows.ps1`**

Create `scripts/setup-windows.ps1`:

```powershell
$ErrorActionPreference = "Continue"

$Root = Split-Path -Parent $PSScriptRoot
$Dashboard = Join-Path $Root "start.html"
$Desktop = [Environment]::GetFolderPath("Desktop")
$ShortcutPath = Join-Path $Desktop "GF2 IT Dashboard.url"

$Links = @(
    @{ Name = "Office 365 / skolemail"; Url = "https://www.office.com/" },
    @{ Name = "Moodle"; Url = "https://online.neg.dk/login/index.php" },
    @{ Name = "Lectio"; Url = "https://www.lectio.dk/lectio/769/default.aspx" },
    @{ Name = "PraxisOnline"; Url = "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery" },
    @{ Name = "OneDrive"; Url = "https://www.office.com/launch/onedrive" },
    @{ Name = "SketchUp / Trimble"; Url = "https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg" },
    @{ Name = "NEG hjemmeside"; Url = "https://www.neg.dk/" }
)

function Write-Section {
    param([string]$Text)
    Write-Host ""
    Write-Host "== $Text ==" -ForegroundColor Cyan
}

function Open-Link {
    param([string]$Url)
    Start-Process $Url
    Start-Sleep -Milliseconds 600
}

Write-Section "GF2 IT setup-assistent til Windows"
Write-Host "Assistenten gemmer ingen brugernavne, adgangskoder eller MitID-oplysninger."

if (Test-Path -LiteralPath $Dashboard) {
    $DashboardUri = (Get-Item -LiteralPath $Dashboard).FullName
    Set-Content -LiteralPath $ShortcutPath -Encoding ASCII -Value @(
        "[InternetShortcut]",
        "URL=file:///$($DashboardUri.Replace('\','/'))"
    )
    Write-Host "Skrivebordsgenvej oprettet: $ShortcutPath"
}
else {
    Write-Host "Dashboardet blev ikke fundet ved: $Dashboard" -ForegroundColor Yellow
}

Write-Section "Systemtjek"
$Os = Get-CimInstance -ClassName Win32_OperatingSystem
Write-Host "Windows: $($Os.Caption) $($Os.Version)"

$SModeSignals = Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\CI\Policy" -ErrorAction SilentlyContinue
if ($SModeSignals -and $SModeSignals.SkuPolicyRequired -eq 1) {
    Write-Host "Mulig Windows S-mode fundet. Installation af programmer kan være blokeret." -ForegroundColor Yellow
}
else {
    Write-Host "Der blev ikke fundet et klart S-mode signal."
}

Write-Section "Åbn dashboard"
if (Test-Path -LiteralPath $Dashboard) {
    Start-Process $Dashboard
}

Write-Section "Åbn skole- og programlinks"
foreach ($Link in $Links) {
    $Answer = Read-Host "Åbn $($Link.Name)? Skriv j for ja"
    if ($Answer -match "^[jJ]$") {
        Open-Link $Link.Url
    }
}

Write-Section "SketchUp"
$Winget = Get-Command winget -ErrorAction SilentlyContinue
if ($Winget) {
    $Answer = Read-Host "Forsøg at søge efter SketchUp med winget? Skriv j for ja"
    if ($Answer -match "^[jJ]$") {
        winget search SketchUp
        Write-Host "Hvis den rigtige SketchUp-version vises, kan installation køres manuelt fra winget eller SketchUp-siden."
    }
}
else {
    Write-Host "winget blev ikke fundet. Brug SketchUp-linket som manuel fallback." -ForegroundColor Yellow
}

Write-Section "Færdig"
Write-Host "Luk vinduet når eleven er videre. Brug dashboardet som manuel fallback."
Read-Host "Tryk Enter for at lukke"
```

- [ ] **Step 3: Create double-click Windows launcher**

Create `Start Windows setup.cmd`:

```bat
@echo off
setlocal
cd /d "%~dp0"
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\setup-windows.ps1"
endlocal
```

- [ ] **Step 4: Run syntax and static checks**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "$null = [System.Management.Automation.PSParser]::Tokenize((Get-Content -Raw scripts/setup-windows.ps1), [ref]$null); 'PowerShell parse passed'"
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected:

```text
PowerShell parse passed
Dashboard checks passed.
```

- [ ] **Step 5: Manual assistant check**

Run:

```powershell
.\scripts\setup-windows.ps1
```

Expected:

- The script explains that it stores no credentials.
- It creates `GF2 IT Dashboard.url` on the desktop.
- It opens `start.html`.
- It asks before each external school/service link.
- It does not ask for or store login details.

- [ ] **Step 6: Commit**

```powershell
git add scripts/setup-windows.ps1 "Start Windows setup.cmd"
git commit -m "feat: add Windows setup assistant"
```

---

### Task 4: Add Mac Setup Assistant

**Files:**
- Create: `scripts/setup-mac.sh`
- Create: `Start Mac setup.command`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Create `scripts/setup-mac.sh`**

Create `scripts/setup-mac.sh`:

```sh
#!/bin/sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
DASHBOARD="$ROOT_DIR/start.html"

open_link() {
  url="$1"
  open "$url"
  sleep 1
}

ask_open() {
  name="$1"
  url="$2"
  printf "Åbn %s? Skriv j for ja: " "$name"
  read answer
  case "$answer" in
    j|J) open_link "$url" ;;
    *) printf "Springer over: %s\n" "$name" ;;
  esac
}

printf "\n== GF2 IT setup-assistent til Mac ==\n"
printf "Assistenten gemmer ingen brugernavne, adgangskoder eller MitID-oplysninger.\n"

if [ -f "$DASHBOARD" ]; then
  printf "\nÅbner dashboardet lokalt.\n"
  open "$DASHBOARD"
else
  printf "\nDashboardet blev ikke fundet ved: %s\n" "$DASHBOARD"
fi

printf "\n== Åbn skole- og programlinks ==\n"
ask_open "Office 365 / skolemail" "https://www.office.com/"
ask_open "Moodle" "https://online.neg.dk/login/index.php"
ask_open "Lectio" "https://www.lectio.dk/lectio/769/default.aspx"
ask_open "PraxisOnline" "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery"
ask_open "OneDrive" "https://www.office.com/launch/onedrive"
ask_open "SketchUp / Trimble" "https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg"
ask_open "NEG hjemmeside" "https://www.neg.dk/"

printf "\n== SketchUp ==\n"
printf "Hvis SketchUp ikke er installeret, brug SketchUp-linket og følg skolens Trimble-flow fra skolemailen.\n"
printf "Dashboard-genvej på Mac holdes manuel i første version: læg mappen et fast sted og åbn start.html.\n"

printf "\nFærdig. Tryk Enter for at lukke.\n"
read unused
```

- [ ] **Step 2: Create double-click Mac launcher**

Create `Start Mac setup.command`:

```sh
#!/bin/sh
DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
chmod +x "$DIR/scripts/setup-mac.sh"
"$DIR/scripts/setup-mac.sh"
```

- [ ] **Step 3: Mark Mac scripts executable where supported**

Run:

```powershell
git update-index --chmod=+x scripts/setup-mac.sh
git update-index --chmod=+x "Start Mac setup.command"
```

Expected: no output.

- [ ] **Step 4: Run static checks**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 5: Manual Mac check on macOS**

Run on a Mac after transferring the folder:

```sh
sh -n scripts/setup-mac.sh
sh -n "Start Mac setup.command"
./scripts/setup-mac.sh
```

Expected:

- Shell syntax checks pass silently.
- Script opens local `start.html`.
- Script asks before each external link.
- Script does not ask for or store login details.

- [ ] **Step 6: Commit**

```powershell
git add scripts/setup-mac.sh "Start Mac setup.command"
git commit -m "feat: add Mac setup assistant"
```

---

### Task 5: Add Zip Packaging Script

**Files:**
- Create: `scripts/build-package.ps1`
- Generate: `dist/IT-opstart-GF2.zip`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Create `scripts/build-package.ps1`**

Create `scripts/build-package.ps1`:

```powershell
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Dist = Join-Path $Root "dist"
$Zip = Join-Path $Dist "IT-opstart-GF2.zip"

$Required = @(
    "start.html",
    "assets\neg-hero-transition.png",
    "scripts\setup-windows.ps1",
    "scripts\setup-mac.sh",
    "Start Windows setup.cmd",
    "Start Mac setup.command"
)

foreach ($Relative in $Required) {
    $Path = Join-Path $Root $Relative
    if (-not (Test-Path -LiteralPath $Path)) {
        throw "Missing required package file: $Relative"
    }
}

if (-not (Test-Path -LiteralPath $Dist)) {
    New-Item -ItemType Directory -Path $Dist | Out-Null
}

if (Test-Path -LiteralPath $Zip) {
    Remove-Item -LiteralPath $Zip -Force
}

$PackageRoot = Join-Path $Dist "IT-opstart-GF2"
if (Test-Path -LiteralPath $PackageRoot) {
    Remove-Item -LiteralPath $PackageRoot -Recurse -Force
}

New-Item -ItemType Directory -Path $PackageRoot | Out-Null
New-Item -ItemType Directory -Path (Join-Path $PackageRoot "assets") | Out-Null
New-Item -ItemType Directory -Path (Join-Path $PackageRoot "scripts") | Out-Null

foreach ($Relative in $Required) {
    $Source = Join-Path $Root $Relative
    $Target = Join-Path $PackageRoot $Relative
    Copy-Item -LiteralPath $Source -Destination $Target -Force
}

Compress-Archive -Path (Join-Path $PackageRoot "*") -DestinationPath $Zip -Force
Remove-Item -LiteralPath $PackageRoot -Recurse -Force

Write-Host "Created $Zip"
```

- [ ] **Step 2: Build the zip**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

Expected: PASS with `Created ...\dist\IT-opstart-GF2.zip`.

- [ ] **Step 3: Verify zip contents**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 4: Inspect zip entries**

Run:

```powershell
Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip = [System.IO.Compression.ZipFile]::OpenRead((Resolve-Path dist\IT-opstart-GF2.zip))
$zip.Entries.FullName
$zip.Dispose()
```

Expected entries include:

```text
start.html
assets/neg-hero-transition.png
scripts/setup-windows.ps1
scripts/setup-mac.sh
Start Windows setup.cmd
Start Mac setup.command
```

- [ ] **Step 5: Commit**

```powershell
git add scripts/build-package.ps1 dist/IT-opstart-GF2.zip
git commit -m "build: add GF2 startup zip package"
```

---

### Task 6: Final Manual QA

**Files:**
- Read: `start.html`
- Read: `dist/IT-opstart-GF2.zip`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Run static verification**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 2: Test extracted zip flow**

Run:

```powershell
New-Item -ItemType Directory -Force dist\qa-extract
Expand-Archive -LiteralPath dist\IT-opstart-GF2.zip -DestinationPath dist\qa-extract -Force
Start-Process .\dist\qa-extract\start.html
```

Expected:

- Dashboard opens from the extracted folder.
- Hero image loads from `assets/neg-hero-transition.png`.
- Windows setup link points to `Start Windows setup.cmd`.
- Mac setup link appears after choosing Mac and points to `Start Mac setup.command`.
- No network request happens until a fixed link or setup assistant link is clicked.

- [ ] **Step 3: Test local status persistence**

In the browser:

- Change one student step to `I gang`.
- Change another student step to `Færdig`.
- Refresh the page.
- Confirm both statuses remain visible.
- Toggle `Underviser`.
- Change one teacher step status.
- Refresh the page.
- Confirm teacher status remains visible.

- [ ] **Step 4: Test credential safety**

Verify:

- There is no text input for username, password, MitID, CPR, phone number, or email.
- Scripts ask only yes/no questions for opening pages.
- `start.html` uses `localStorage` only for mode, platform, and status values.

- [ ] **Step 5: Remove temporary QA extraction**

Run:

```powershell
Remove-Item -LiteralPath dist\qa-extract -Recurse -Force
```

Expected: `dist\qa-extract` is removed.

- [ ] **Step 6: Commit any final verification updates**

If Task 6 caused changes to tests or scripts, commit them:

```powershell
git status --short
git add tests/check-dashboard.ps1 scripts/build-package.ps1 scripts/setup-windows.ps1 scripts/setup-mac.sh start.html
git commit -m "test: tighten GF2 dashboard verification"
```

If `git status --short` shows no tracked file changes, skip this commit.

---

## Self-Review Against Spec

- `start.html` light version: Task 2 creates it as a local-file dashboard.
- `IT-opstart-GF2.zip`: Task 5 builds `dist/IT-opstart-GF2.zip`.
- Local asset: Task 2 references `assets/neg-hero-transition.png`; Task 1 checks no remote image/script references.
- Windows setup assistant: Task 3 creates `scripts/setup-windows.ps1` and `Start Windows setup.cmd`.
- Mac setup assistant: Task 4 creates `scripts/setup-mac.sh` and `Start Mac setup.command`.
- Offline dashboard: Task 2 uses embedded CSS/JS and local image only.
- External links only on click: Task 2 renders normal anchor links and setup scripts ask before opening each URL.
- No credential storage: Task 1 checks forbidden wording and Task 6 verifies no credential fields exist.
- Local status only: Task 2 uses browser `localStorage` for status, platform, and mode.
- OneDrive as official file location: Task 2 includes OneDrive in the ordered student flow and fixed links.
- Printer help not a main link: No printer link is added.
- NEG visual direction: Task 2 uses `#005aa7`, topbar, NEG mark, and hero copy from the spec.
- Student view: Task 2 makes it the default mode and includes the required eight-step order.
- Teacher view: Task 2 adds the topbar slider and teacher checklist items from the spec.
- Fixed links and absence phone: Task 2 includes all required links and `+45 72 290 100`.
- Setup assistant boundaries: Tasks 3 and 4 open pages, offer SketchUp guidance, and avoid automatic login.
- Excluded first-version items: The plan does not add central status collection, teacher login, admin panel, saved credentials, automatic login, FileCloud/P-drev, or remote production images.

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-04-29-gf2-it-dashboard-implementation.md`. Two execution options:

**1. Subagent-Driven (recommended)** - dispatch a fresh subagent per task, review between tasks, fast iteration.

**2. Inline Execution** - execute tasks in this session using executing-plans, batch execution with checkpoints.

Which approach?
