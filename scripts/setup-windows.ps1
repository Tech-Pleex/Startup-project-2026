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
