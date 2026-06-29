$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$GitAttributes = Join-Path $Root ".gitattributes"
$StartHtml = Join-Path $Root "start.html"
$HeroImage = Join-Path $Root "assets\neg-hero-transition.png"
$BrandFont = Join-Path $Root "assets\fonts\space-grotesk.woff2"
$BrandLogo = Join-Path $Root "assets\neg-logo.png"
$WindowsSetup = Join-Path $Root "scripts\setup-windows.ps1"
$SetupConfig = Join-Path $Root "scripts\setup-config.ps1"
$SetupChecks = Join-Path $Root "scripts\setup-checks.ps1"
$WindowsLauncher = Join-Path $Root "Start Windows setup.cmd"
$MacSetup = Join-Path $Root "scripts\setup-mac.sh"
$MacLauncher = Join-Path $Root "Start Mac setup.command"
$BuildPackage = Join-Path $Root "scripts\build-package.ps1"
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

function Assert-StartsWithBytes {
    param(
        [string]$Path,
        [byte[]]$Expected,
        [string]$Label
    )

    $Bytes = [System.IO.File]::ReadAllBytes($Path)
    if ($Bytes.Length -lt $Expected.Length) {
        throw "File too short for '$Label': $Path"
    }

    for ($Index = 0; $Index -lt $Expected.Length; $Index++) {
        if ($Bytes[$Index] -ne $Expected[$Index]) {
            throw "Missing expected byte prefix '$Label': $Path"
        }
    }
}

Assert-File $GitAttributes
Assert-File $StartHtml
Assert-File $HeroImage
Assert-File $BrandFont
Assert-File $BrandLogo
Assert-File $WindowsSetup
Assert-File $SetupConfig
Assert-File $SetupChecks
Assert-File $WindowsLauncher
Assert-File $MacSetup
Assert-File $MacLauncher
Assert-File $BuildPackage

$GitAttributesContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $GitAttributes
$Html = Get-Content -Raw -Encoding UTF8 -LiteralPath $StartHtml
$WindowsSetupContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsSetup
$SetupConfigContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $SetupConfig
$WindowsLauncherContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $WindowsLauncher
$MacSetupContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $MacSetup
$MacLauncherContent = Get-Content -Raw -Encoding UTF8 -LiteralPath $MacLauncher

Assert-Contains $Html "GF2 IT Dashboard" "dashboard title"
Assert-Contains $Html "ET FÆLLESSKAB MED PLADS TIL DIG" "hero kicker"
Assert-Contains $Html "Kom sikkert i gang med IT på GF2" "hero heading"
Assert-Contains $Html "assets/neg-hero-transition.png" "local hero asset"
Assert-Contains $Html "Elev" "student mode"
Assert-Contains $Html "Underviser" "teacher mode"
Assert-Contains $Html "Ikke startet" "status label"
Assert-Contains $Html "Færdig" "status label"
Assert-Contains $Html 'type="checkbox"' "compact done-checkbox replaces status dropdown"
Assert-Contains $Html "done-toggle" "done toggle control"
Assert-NotContains $Html "status-control select" "old status dropdown removed"
Assert-Contains $Html "gf2-it-dashboard.studentStatus" "student localStorage key"
Assert-Contains $Html "gf2-it-dashboard.teacherStatus" "teacher localStorage key"
Assert-Contains $Html "setup=complete" "setup completion URL signal"
Assert-Contains $Html "applySetupCompletionFromUrl" "setup completion handler"
Assert-Contains $Html 'fillStudentStatus("Færdig")' "setup completion marks student status done"
Assert-Contains $Html ".step p {" "dashboard step body style"
Assert-Contains $Html "font-size: 13px;" "smaller dashboard step body text"
Assert-Contains $Html "Download setup" "dashboard download page link"
Assert-Contains $Html 'href="index.html"' "dashboard links back to landing page"

# --- NEG brand (#29): palet, font og officielt logo ---
Assert-Contains $Html "Space Grotesk" "NEG brand font Space Grotesk"
Assert-Contains $Html "assets/fonts/space-grotesk.woff2" "locally bundled brand font"
Assert-Contains $Html '@font-face' "font bundled offline, no CDN"
Assert-Contains $Html "#123c62" "NEG navy palette color"
Assert-Contains $Html "#86afd8" "NEG blue palette color"
Assert-Contains $Html "#ec8113" "NEG orange palette color"
Assert-Contains $Html '#efefef' "NEG grey palette color"
Assert-Contains $Html 'src="assets/neg-logo.png"' "official NEG logo image (not redrawn)"
Assert-NotContains $Html "#005aa7" "off-brand blue removed"

# --- Layout (#3): tjekliste-domineret, fremgang og kontekstuelle links ---
Assert-Contains $Html "panel-main" "checklist is the dominant panel"
Assert-Contains $Html "grid-template-columns: 3fr 1fr;" "checklist 3/4 vs sidebar 1/4 layout"
Assert-Contains $Html "grid-template-columns: repeat(3, 1fr);" "steps render in three columns"
Assert-Contains $Html 'class="shell"' "page split into checklist + right column"
Assert-Contains $Html "hero-media" "hero image shares the right column with the sidebar"
Assert-Contains $Html "col-side" "sidebar column same width as hero image"
Assert-Contains $Html 'id="progressBar"' "progress bar element"
Assert-Contains $Html 'id="progressCount"' "progress count element"
Assert-Contains $Html "af 8 færdige" "student progress wording"
Assert-Contains $Html "step-links" "contextual quicklinks at steps"
Assert-Contains $Html "function makeChip" "contextual link chip rendering"
Assert-Contains $Html "linkByKey" "links referenced by key, not duplicated"

# --- Lærer-reminder + selvbetjening (#7) ---
Assert-Contains $Html "https://selvbetjening.neg.dk/UMSLogin/weblogin.aspx?ReturnUrl=%2f" "NEG selvbetjening link"
Assert-Contains $Html "udskriv bruger" "teacher self-service workflow reminder"
Assert-Contains $Html "glemt mail eller login" "teacher forgotten-login step"
# Selvbetjening er tilgængelig for elever (skift egen adgangskode) såvel som undervisere.
Assert-Contains $Html "skifte din adgangskode" "student self-service password wording"
Assert-NotContains $Html "Vælg computer" "dashboard platform chooser removed"
Assert-NotContains $Html 'data-platform="windows"' "Windows platform button removed"
Assert-NotContains $Html 'data-platform="mac"' "Mac platform button removed"
Assert-NotContains $Html "gf2-it-dashboard.platform" "platform localStorage key removed"
Assert-NotContains $Html "function setPlatform" "platform JavaScript removed"
Assert-NotContains $Html "setupLauncherHelp" "dashboard launcher help removed"
Assert-NotContains $Html 'href="Start Windows setup.cmd"' "browser must not link directly to Windows CMD"
Assert-NotContains $Html 'href="Start Mac setup.command"' "browser must not link directly to Mac command"
Assert-Contains $Html "Office 365 / skolemail" "fixed link"
Assert-Contains $Html "Moodle" "fixed link"
Assert-Contains $Html "Lectio" "fixed link"
Assert-Contains $Html "PraxisOnline" "fixed link"
Assert-Contains $Html "OneDrive" "fixed link"
Assert-Contains $Html "SketchUp / Trimble" "fixed link"
Assert-Contains $Html "NEG hjemmeside" "fixed link"
Assert-Contains $Html "+45 72 290 100" "absence phone"
Assert-Contains $Html "Udviklet af Jesper Reenberg" "developer credit"

foreach ($ExpectedLink in @(
    "https://online.neg.dk/login/index.php",
    "https://www.lectio.dk/lectio/769/default.aspx",
    "https://online.praxis.dk/",
    "https://sketchup.trimble.com/"
)) {
    Assert-Contains $Html $ExpectedLink "correct service URL"
    Assert-Contains $MacSetupContent $ExpectedLink "correct Mac setup URL"
}

foreach ($ExpectedLink in @(
    "https://online.neg.dk/login/index.php",
    "https://www.lectio.dk/lectio/769/default.aspx",
    "https://online.praxis.dk/",
    "https://sketchup.trimble.com/"
)) {
    Assert-Contains $SetupConfigContent $ExpectedLink "correct Windows setup URL"
}

Assert-Contains $SetupConfigContent "GF2 IT Dashboard.url" "desktop shortcut"
Assert-Contains $WindowsSetupContent "[Console]::OutputEncoding" "PowerShell UTF-8 output"
Assert-Contains $WindowsSetupContent "System.Windows.Forms" "Windows setup GUI"
Assert-Contains $WindowsSetupContent "Assistenten beder aldrig om adgangskoder" "Windows setup safety"
Assert-Contains $WindowsSetupContent "setup=complete" "Windows wizard opens dashboard with completion signal"
Assert-Contains $WindowsLauncherContent "chcp 65001" "Windows launcher UTF-8 codepage"
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
# Selvbetjenings-login'et (#7) bruger ReturnUrl=%2f (tilbage til roden) og er bevidst tilladt.
# Sessionfulde service-URL'er er stadig forbudt:
Assert-NotContains $Html "online.praxis.dk/?ReturnUrl" "sessionful Praxis URL"
Assert-Contains $Html "https://online.praxis.dk/" "Praxis uses clean root URL"
Assert-NotContains $Html "id.trimble.com/ui/sign_in.html?state=" "sessionful Trimble URL"
Assert-NotContains $Html "password" "credential wording"
Assert-NotContains $Html "adgangskode gemmes" "stored password wording"

Assert-StartsWithBytes $WindowsSetup ([byte[]](0xEF, 0xBB, 0xBF)) "Windows PowerShell UTF-8 BOM"

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
