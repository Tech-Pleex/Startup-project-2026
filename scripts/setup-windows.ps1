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

$WindowTitle = if ($SetupConfig.Title) { $SetupConfig.Title } else { "GF2 IT Setup" }
$SafetyText = if ($SetupConfig.SafetyText) {
    $SetupConfig.SafetyText
}
else {
    "Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login."
}

$Dashboard = Join-Path $Root $SetupConfig.DashboardFile
$Desktop = [Environment]::GetFolderPath("Desktop")
$ShortcutPath = Join-Path $Desktop $SetupConfig.DesktopShortcutName
$CurrentStepIndex = 0

# SketchUp helper runs winget install through Install-SketchUpPackage.

$Form = New-Object System.Windows.Forms.Form
$Form.Text = $WindowTitle
$Form.StartPosition = "CenterScreen"
$Form.Size = New-Object System.Drawing.Size(760, 520)
$Form.MinimumSize = New-Object System.Drawing.Size(680, 460)
$Form.BackColor = [System.Drawing.Color]::White

$TitleLabel = New-Object System.Windows.Forms.Label
$TitleLabel.Font = New-Object System.Drawing.Font("Segoe UI", 18, [System.Drawing.FontStyle]::Bold)
$TitleLabel.Location = New-Object System.Drawing.Point(26, 22)
$TitleLabel.Size = New-Object System.Drawing.Size(690, 42)
$TitleLabel.Anchor = "Top, Left, Right"
$Form.Controls.Add($TitleLabel)

$ProgressLabel = New-Object System.Windows.Forms.Label
$ProgressLabel.Font = New-Object System.Drawing.Font("Segoe UI", 9, [System.Drawing.FontStyle]::Regular)
$ProgressLabel.ForeColor = [System.Drawing.Color]::FromArgb(80, 96, 112)
$ProgressLabel.Location = New-Object System.Drawing.Point(30, 68)
$ProgressLabel.Size = New-Object System.Drawing.Size(690, 24)
$ProgressLabel.Anchor = "Top, Left, Right"
$Form.Controls.Add($ProgressLabel)

$BodyBox = New-Object System.Windows.Forms.TextBox
$BodyBox.Multiline = $true
$BodyBox.ReadOnly = $true
$BodyBox.BorderStyle = "None"
$BodyBox.BackColor = [System.Drawing.Color]::White
$BodyBox.Font = New-Object System.Drawing.Font("Segoe UI", 11)
$BodyBox.Location = New-Object System.Drawing.Point(30, 110)
$BodyBox.Size = New-Object System.Drawing.Size(680, 110)
$BodyBox.Anchor = "Top, Left, Right"
$Form.Controls.Add($BodyBox)

$WarningBox = New-Object System.Windows.Forms.TextBox
$WarningBox.Multiline = $true
$WarningBox.ReadOnly = $true
$WarningBox.BorderStyle = "FixedSingle"
$WarningBox.BackColor = [System.Drawing.Color]::FromArgb(255, 245, 214)
$WarningBox.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$WarningBox.Location = New-Object System.Drawing.Point(30, 238)
$WarningBox.Size = New-Object System.Drawing.Size(680, 78)
$WarningBox.Anchor = "Top, Left, Right"
$Form.Controls.Add($WarningBox)

$StatusLabel = New-Object System.Windows.Forms.Label
$StatusLabel.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$StatusLabel.ForeColor = [System.Drawing.Color]::FromArgb(0, 90, 167)
$StatusLabel.Location = New-Object System.Drawing.Point(30, 334)
$StatusLabel.Size = New-Object System.Drawing.Size(680, 38)
$StatusLabel.Anchor = "Top, Left, Right"
$Form.Controls.Add($StatusLabel)

$PrimaryButton = New-Object System.Windows.Forms.Button
$PrimaryButton.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$PrimaryButton.Location = New-Object System.Drawing.Point(30, 396)
$PrimaryButton.Size = New-Object System.Drawing.Size(240, 42)
$PrimaryButton.Anchor = "Bottom, Left"
$Form.Controls.Add($PrimaryButton)

$DoneButton = New-Object System.Windows.Forms.Button
$DoneButton.Text = "Done / Next"
$DoneButton.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$DoneButton.Location = New-Object System.Drawing.Point(470, 396)
$DoneButton.Size = New-Object System.Drawing.Size(240, 42)
$DoneButton.Anchor = "Bottom, Right"
$Form.Controls.Add($DoneButton)

function Get-CurrentStep {
    return $SetupSteps[$CurrentStepIndex]
}

function Set-StatusText {
    param([string]$Text)

    $StatusLabel.Text = $Text
}

function Set-WifiStatus {
    $Ssid = Get-ActiveWifiSsid

    if ($Ssid -eq $SetupConfig.TargetWifi) {
        Set-StatusText "Du er allerede på $($SetupConfig.TargetWifi). Trinnet er gennemført."
    }
    elseif ($Ssid -eq $SetupConfig.GuestWifi) {
        Set-StatusText "Du er på $($SetupConfig.GuestWifi). Skift til $($SetupConfig.TargetWifi), og klik Done / Next."
    }
    elseif ($Ssid) {
        Set-StatusText "Aktivt Wi-Fi: $Ssid. Skift til $($SetupConfig.TargetWifi), hvis det ikke er korrekt."
    }
    else {
        Set-StatusText "Assistenten kunne ikke læse aktivt Wi-Fi. Brug Windows Wi-Fi-indstillinger."
    }
}

function Set-SketchUpStatus {
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

function Render-Step {
    $Step = Get-CurrentStep

    $TitleLabel.Text = $Step.Title
    $ProgressLabel.Text = "Trin $($CurrentStepIndex + 1) af $($SetupSteps.Count)"
    $BodyBox.Text = $Step.Body
    $PrimaryButton.Text = $Step.Button
    $PrimaryButton.Enabled = ($Step.Kind -ne "manual")
    $WarningBox.Visible = [bool]$Step.Warning
    $WarningBox.Text = if ($Step.Warning) { $Step.Warning } else { "" }
    $DoneButton.Text = if ($Step.Kind -eq "finish") { "Luk" } else { "Done / Next" }
    Set-StatusText ""

    if ($Step.Kind -eq "manual") {
        Set-StatusText $SafetyText
    }
    elseif ($Step.Kind -eq "wifi") {
        Set-WifiStatus
    }
    elseif ($Step.Kind -eq "sketchup") {
        Set-SketchUpStatus
    }
    elseif ($Step.Kind -eq "finish") {
        $PrimaryButton.Text = $Step.Button
        $PrimaryButton.Enabled = $true
        Set-StatusText "Klik for at oprette genvej og åbne dashboardet."
    }
}

function Invoke-FallbackLink {
    param([string]$Url)

    if ($Url) {
        Open-SetupLink -Url $Url
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
        Set-StatusText "Siden er åbnet i browseren. Log selv ind på den officielle side og klik Done / Next bagefter."
        return
    }

    if ($Step.Kind -eq "sketchup") {
        if (Test-WindowsSMode) {
            Invoke-FallbackLink -Url $Step.Url
            Set-StatusText "Manuel SketchUp-side er åbnet, fordi Windows S-mode kan blokere installation."
            return
        }

        if (-not (Test-WingetAvailable)) {
            Invoke-FallbackLink -Url $Step.Url
            Set-StatusText "winget blev ikke fundet. Manuel SketchUp-side er åbnet."
            return
        }

        $Result = Install-SketchUpPackage -PackageId $SetupConfig.SketchUpPackageId
        if (-not $Result.Success) {
            Invoke-FallbackLink -Url $Step.Url
            Set-StatusText "$($Result.Message) Manuel SketchUp-side er åbnet."
            return
        }

        Set-StatusText $Result.Message
        return
    }

    if ($Step.Kind -eq "finish") {
        $ShortcutCreated = New-DashboardShortcut -DashboardPath $Dashboard -ShortcutPath $ShortcutPath
        if (Test-Path -LiteralPath $Dashboard) {
            Start-Process $Dashboard
        }

        if ($ShortcutCreated) {
            Set-StatusText "Dashboardet er åbnet, og skrivebordsgenvejen er oprettet."
        }
        else {
            Set-StatusText "Dashboardet kunne ikke findes ved $Dashboard."
        }
        return
    }

    Set-StatusText $SafetyText
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
