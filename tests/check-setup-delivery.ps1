$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$LandingPage = Join-Path $Root "index.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$WindowsLauncher = Join-Path $Root "Start Windows setup.cmd"
$SetupScript = Join-Path $Root "scripts\setup-windows.ps1"
$SetupConfig = Join-Path $Root "scripts\setup-config.ps1"
$SetupChecks = Join-Path $Root "scripts\setup-checks.ps1"
$BuildScript = Join-Path $Root "scripts\build-package.ps1"
$ZipPath = Join-Path $Root "dist\GF2-IT-Setup-Windows.zip"
$StartHtml = Join-Path $Root "start.html"

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
        [string]$Path,
        [string]$Entry
    )

    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $Zip = [System.IO.Compression.ZipFile]::OpenRead($Path)
    try {
        $Entries = $Zip.Entries.FullName
        if ($Entries -notcontains $Entry) {
            throw "Zip is missing entry: $Entry"
        }
    }
    finally {
        $Zip.Dispose()
    }
}

Assert-File $LandingPage
Assert-File $HeroImage
Assert-File $WindowsLauncher
Assert-File $SetupScript
Assert-File $SetupConfig
Assert-File $SetupChecks
Assert-File $BuildScript
Assert-File $StartHtml

$LandingHtml = Get-Content -Raw -Encoding UTF8 -LiteralPath $LandingPage
$LauncherContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsLauncher
$SetupContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupScript
$ConfigContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupConfig
$ChecksContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupChecks
$BuildContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $BuildScript

Assert-Contains $LandingHtml "GF2 IT Setup" "landing title"
Assert-Contains $LandingHtml "assets/neg-hero-transition.png" "landing hero asset"
Assert-Contains $LandingHtml "Download Windows Setup" "download call to action"
Assert-Contains $LandingHtml "dist/GF2-IT-Setup-Windows.zip" "download package path"
Assert-Contains $LandingHtml "Vi beder aldrig om adgangskoder" "credential safety"
Assert-Contains $LandingHtml "start.html" "dashboard entry"
Assert-Contains $LandingHtml "GitHub" "GitHub fallback"

Assert-Contains $ConfigContent "Trimble.SketchUp.2026" "SketchUp package"
Assert-Contains $ConfigContent "PraxisOnline" "PraxisOnline service"
Assert-Contains $ConfigContent "neg04026@edu.neg.dk" "school email"
Assert-Contains $ConfigContent "UNI-Login" "UNI-Login"
Assert-Contains $ConfigContent "MitID" "MitID"
Assert-Contains $ConfigContent "NEG Guest" "guest Wi-Fi"
Assert-Contains $ConfigContent "`"NEG" "NEG network prefix"

Assert-Contains $ChecksContent "Get-ActiveWifiSsid" "Wi-Fi helper"
Assert-Contains $ChecksContent "Test-WingetAvailable" "winget helper"
Assert-Contains $ChecksContent "Test-WindowsSMode" "S mode helper"
Assert-Contains $ChecksContent "New-DashboardShortcut" "dashboard shortcut helper"

Assert-Contains $SetupContent "System.Windows.Forms" "Windows Forms"
Assert-Contains $SetupContent "GF2 IT Setup" "setup title"
Assert-Contains $SetupContent "Assistenten beder aldrig om adgangskoder" "credential safety"
Assert-Contains $SetupContent "Get-ActiveWifiSsid" "Wi-Fi helper"
Assert-Contains $SetupContent "winget install" "winget install"
Assert-Contains $SetupContent "Start-Process" "process launcher"

Assert-Contains $LauncherContent "powershell.exe" "PowerShell launcher"
Assert-Contains $LauncherContent "setup-windows.ps1" "setup target"
Assert-Contains $LauncherContent "chcp 65001" "UTF-8 codepage"

Assert-Contains $BuildContent "GF2-IT-Setup-Windows.zip" "package name"

Assert-NotContains $LandingHtml "password" "credential wording"
Assert-NotContains $SetupContent "Read-Host `"Åbn" "old terminal prompt"

if (Test-Path -LiteralPath $ZipPath) {
    Assert-ZipContains $ZipPath "start.html"
    Assert-ZipContains $ZipPath "index.html"
    Assert-ZipContains $ZipPath "assets/neg-hero-transition.png"
    Assert-ZipContains $ZipPath "scripts/setup-windows.ps1"
    Assert-ZipContains $ZipPath "scripts/setup-config.ps1"
    Assert-ZipContains $ZipPath "scripts/setup-checks.ps1"
    Assert-ZipContains $ZipPath "Start Windows setup.cmd"
}

Write-Host "Setup delivery checks passed."
