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

function Assert-InDirectory {
    param(
        [string]$Path,
        [string]$Directory
    )

    $FullPath = [System.IO.Path]::GetFullPath($Path)
    $FullDirectory = [System.IO.Path]::GetFullPath($Directory)
    if (-not $FullDirectory.EndsWith([System.IO.Path]::DirectorySeparatorChar)) {
        $FullDirectory = "$FullDirectory$([System.IO.Path]::DirectorySeparatorChar)"
    }

    if (-not $FullPath.StartsWith($FullDirectory, [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to operate outside package directory: $FullPath"
    }
}

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
    Assert-InDirectory -Path $Zip -Directory $Dist
    Remove-Item -LiteralPath $Zip -Force
}

Add-Type -AssemblyName System.IO.Compression
Add-Type -AssemblyName System.IO.Compression.FileSystem
$ZipArchive = [System.IO.Compression.ZipFile]::Open($Zip, [System.IO.Compression.ZipArchiveMode]::Create)
try {
    foreach ($Relative in $Required) {
        $Source = Join-Path $Root $Relative
        $EntryName = $Relative.Replace("\", "/")
        [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
            $ZipArchive,
            $Source,
            $EntryName,
            [System.IO.Compression.CompressionLevel]::Optimal
        ) | Out-Null
    }
}
finally {
    $ZipArchive.Dispose()
}

Write-Host "Created $Zip"
