$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Readme = Join-Path $Root "README.md"
$LandingPage = Join-Path $Root "index.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$WindowsLauncher = Join-Path $Root "Start Windows setup.cmd"
$SetupScript = Join-Path $Root "scripts\setup-windows.ps1"
$SetupConfig = Join-Path $Root "scripts\setup-config.ps1"
$SetupChecks = Join-Path $Root "scripts\setup-checks.ps1"
$BuildScript = Join-Path $Root "scripts\build-package.ps1"
$StartHtml = Join-Path $Root "start.html"
$WindowsReleaseUrl = "https://github.com/Tech-Pleex/Startup-project-2026/releases/latest/download/Assistenten-Windows.exe"
$MacAppleSiliconReleaseUrl = "https://github.com/Tech-Pleex/Startup-project-2026/releases/latest/download/Assistenten-Mac-Apple-Silicon"
$MacIntelReleaseUrl = "https://github.com/Tech-Pleex/Startup-project-2026/releases/latest/download/Assistenten-Mac-Intel"

function Assert-File {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
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

Assert-File $Readme
Assert-File $LandingPage
Assert-File $HeroImage
Assert-File $WindowsLauncher
Assert-File $SetupScript
Assert-File $SetupConfig
Assert-File $SetupChecks
Assert-File $BuildScript
Assert-File $StartHtml

$ReadmeContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $Readme
$LandingHtml = Get-Content -Raw -Encoding UTF8 -LiteralPath $LandingPage
$LauncherContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsLauncher
$SetupContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupScript
$ConfigContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupConfig
$ChecksContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupChecks
$BuildContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $BuildScript

Assert-Contains $ReadmeContent "never asks for passwords, MitID, or UNI-Login" "README credential safety"
Assert-Contains $ReadmeContent "Assistenten-Windows.exe" "README Windows release asset"
Assert-Contains $ReadmeContent "Assistenten-Mac-Apple-Silicon" "README Apple Silicon release asset"
Assert-Contains $ReadmeContent "Assistenten-Mac-Intel" "README Intel Mac release asset"

Assert-Contains $LandingHtml "GF2 IT Setup" "landing title"
Assert-Contains $LandingHtml "assets/neg-hero-transition.png" "landing hero asset"
Assert-Contains $LandingHtml "Download setup" "download call to action"
Assert-Contains $LandingHtml "downloadModal" "platform modal"
Assert-Contains $LandingHtml "Vælg computer" "platform modal title"
Assert-Contains $LandingHtml "Windows" "Windows platform choice"
Assert-Contains $LandingHtml $WindowsReleaseUrl "Windows release download URL"
Assert-Contains $LandingHtml "Mac" "Mac platform choice"
Assert-Contains $LandingHtml "Vælg Mac-type" "Mac type selection title"
Assert-Contains $LandingHtml "Vælg Apple Silicon for M1/M2/M3/M4. Vælg Intel for ældre Macs." "Mac type guidance"
Assert-Contains $LandingHtml "Apple Silicon" "Apple Silicon platform choice"
Assert-Contains $LandingHtml "Intel Mac" "Intel Mac platform choice"
Assert-Contains $LandingHtml $MacAppleSiliconReleaseUrl "Apple Silicon release download URL"
Assert-Contains $LandingHtml $MacIntelReleaseUrl "Intel Mac release download URL"
Assert-Contains $LandingHtml "Tilbage" "Mac selection back button"
Assert-Contains $LandingHtml "openDownloadModal" "open modal handler"
Assert-Contains $LandingHtml "closeDownloadModal" "close modal handler"
Assert-Contains $LandingHtml "showMacChoices" "Mac choice handler"
Assert-Contains $LandingHtml "showPlatformChoices" "platform choice handler"
Assert-Contains $LandingHtml "Escape" "keyboard modal close"
Assert-Contains $LandingHtml "Vi beder aldrig om adgangskoder" "credential safety"
Assert-Contains $LandingHtml "start.html" "dashboard entry"
Assert-Contains $LandingHtml "GitHub" "GitHub fallback"
Assert-Contains $LandingHtml "Udviklet af Jesper Reenberg" "developer credit"
Assert-NotContains $LandingHtml "dist/GF2-IT-Setup-Windows.zip" "old ZIP landing download"

Assert-Contains $ConfigContent "Trimble.SketchUp.2026" "SketchUp package"
Assert-Contains $ConfigContent "PraxisOnline" "PraxisOnline service"
Assert-Contains $ConfigContent "neg04026@edu.neg.dk" "school email"
Assert-Contains $ConfigContent "UNI-Login" "UNI-Login"
Assert-Contains $ConfigContent "MitID" "MitID"
Assert-Contains $ConfigContent "NEG Guest" "guest Wi-Fi"
Assert-Contains $ConfigContent "TargetWifi = `"NEG`"" "target Wi-Fi"

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
Assert-Contains $SetupContent "function Move-PreviousStep" "previous step navigation"
Assert-Contains $SetupContent '$BackButton.Text = "Tilbage"' "back button label"
Assert-Contains $SetupContent '$BackButton.Visible = $CurrentStepIndex -gt 0' "hide back button on first step"
Assert-Contains $SetupContent '$BackButton.Add_Click({ Move-PreviousStep })' "back button click handler"

Assert-Contains $LauncherContent "powershell.exe" "PowerShell launcher"
Assert-Contains $LauncherContent "setup-windows.ps1" "setup target"
Assert-Contains $LauncherContent "chcp 65001" "UTF-8 codepage"

Assert-Contains $BuildContent "GF2-IT-Setup-Windows.zip" "legacy package name"

Assert-NotContains $LandingHtml "password" "credential wording"
Assert-NotContains $SetupContent "Read-Host `"Åbn" "old terminal prompt"

Write-Host "Setup delivery checks passed."
