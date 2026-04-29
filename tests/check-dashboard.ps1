$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$GitAttributes = Join-Path $Root ".gitattributes"
$StartHtml = Join-Path $Root "start.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$WindowsSetup = Join-Path $Root "scripts\setup-windows.ps1"
$WindowsLauncher = Join-Path $Root "Start Windows setup.cmd"
$MacSetup = Join-Path $Root "scripts\setup-mac.sh"
$MacLauncher = Join-Path $Root "Start Mac setup.command"
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

Assert-File $GitAttributes
Assert-File $StartHtml
Assert-File $HeroImage
Assert-File $WindowsSetup
Assert-File $WindowsLauncher
Assert-File $MacSetup
Assert-File $MacLauncher

$GitAttributesContent = Get-Content -Raw -LiteralPath $GitAttributes
$Html = Get-Content -Raw -LiteralPath $StartHtml
$WindowsSetupContent = Get-Content -Raw -LiteralPath $WindowsSetup
$WindowsLauncherContent = Get-Content -Raw -LiteralPath $WindowsLauncher
$MacSetupContent = Get-Content -Raw -LiteralPath $MacSetup
$MacLauncherContent = Get-Content -Raw -LiteralPath $MacLauncher

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

foreach ($ExpectedLink in @(
    "https://online.neg.dk/login/index.php",
    "https://www.lectio.dk/lectio/769/default.aspx",
    "https://authentication.praxis.dk/Account/Login?ReturnUrl=",
    "https://id.trimble.com/ui/sign_in.html?state="
)) {
    Assert-Contains $Html $ExpectedLink "correct service URL"
    Assert-Contains $WindowsSetupContent $ExpectedLink "correct Windows setup URL"
    Assert-Contains $MacSetupContent $ExpectedLink "correct Mac setup URL"
}

Assert-Contains $WindowsSetupContent "Assistenten gemmer ingen brugernavne" "credential safety"
Assert-Contains $WindowsSetupContent "GF2 IT Dashboard.url" "desktop shortcut"
Assert-Contains $WindowsSetupContent "Read-Host `"Åbn" "asks before opening links"
Assert-Contains $WindowsLauncherContent "scripts\setup-windows.ps1" "Windows launcher target"
Assert-Contains $MacSetupContent "GF2 IT setup-assistent til Mac" "Mac setup heading"
Assert-Contains $MacSetupContent "Assistenten gemmer ingen brugernavne" "Mac credential safety"
Assert-Contains $MacSetupContent "Skriv j for ja" "Mac asks before opening links"
Assert-Contains $MacSetupContent "Dashboard-genvej på Mac holdes manuel" "Mac manual dashboard shortcut"
Assert-Contains $MacLauncherContent "scripts/setup-mac.sh" "Mac launcher target"
Assert-Contains $GitAttributesContent "*.sh text eol=lf" "shell script LF endings"
Assert-Contains $GitAttributesContent "*.command text eol=lf" "Mac launcher LF endings"

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
