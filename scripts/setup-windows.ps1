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

$DashboardPath = Join-Path $Root $SetupConfig.DashboardFile
$ShortcutPath = Join-Path ([Environment]::GetFolderPath("Desktop")) $SetupConfig.DesktopShortcutName
$SafetyReminder = "Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login."
$CurrentStepIndex = 0

function Set-StatusText {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message
    )

    $StatusLabel.Text = $Message
}

function Get-StepProgressText {
    $displayIndex = $CurrentStepIndex + 1
    return "Trin $displayIndex af $($SetupSteps.Count)"
}

function Get-WifiStatusText {
    $ssid = Get-ActiveWifiSsid
    if ($ssid -eq $SetupConfig.TargetWifi) {
        return "Wi-Fi: Du er på $($SetupConfig.TargetWifi)."
    }

    if ($ssid -eq $SetupConfig.GuestWifi) {
        return "Wi-Fi: Du er på $($SetupConfig.GuestWifi). Brug $($SetupConfig.TargetWifi), hvis skolens personale beder om det."
    }

    if ([string]::IsNullOrWhiteSpace($ssid)) {
        return "Wi-Fi: Aktivt netværk er ukendt."
    }

    return "Wi-Fi: Aktivt netværk er $ssid. Målet er $($SetupConfig.TargetWifi)."
}

function Render-Step {
    $step = $SetupSteps[$CurrentStepIndex]
    $Form.Text = $SetupConfig.Title
    $ProgressLabel.Text = Get-StepProgressText
    $TitleLabel.Text = $step.Title
    $BodyTextBox.Text = $step.Body
    $PrimaryButton.Text = $step.Button
    $PrimaryButton.Enabled = $true

    if ($step.Contains("Warning") -and -not [string]::IsNullOrWhiteSpace($step.Warning)) {
        $WarningLabel.Text = $step.Warning
        $WarningPanel.Visible = $true
    }
    else {
        $WarningLabel.Text = ""
        $WarningPanel.Visible = $false
    }

    if ($step.Kind -eq "wifi") {
        Set-StatusText (Get-WifiStatusText)
    }
    elseif ($step.Kind -eq "finish") {
        Set-StatusText "Klar til at oprette genvej og åbne dashboardet."
    }
    elseif ($step.Kind -eq "manual") {
        Set-StatusText $SafetyReminder
    }
    else {
        Set-StatusText "Klar."
    }

    if ($CurrentStepIndex -ge ($SetupSteps.Count - 1)) {
        $NextButton.Text = "Luk"
    }
    else {
        $NextButton.Text = "Færdig - næste"
    }
}

function Invoke-StepAction {
    $step = $SetupSteps[$CurrentStepIndex]

    try {
        switch ($step.Kind) {
            "manual" {
                [System.Windows.Forms.MessageBox]::Show(
                    $SetupConfig.SafetyText,
                    $SetupConfig.Title,
                    [System.Windows.Forms.MessageBoxButtons]::OK,
                    [System.Windows.Forms.MessageBoxIcon]::Information
                ) | Out-Null
                Set-StatusText "Sikkerhedsteksten er vist."
            }
            "wifi" {
                Open-WifiSettings
                Set-StatusText "$(Get-WifiStatusText) Windows Wi-Fi-indstillinger er åbnet."
            }
            "link" {
                Open-SetupLink $step.Url
                Set-StatusText "Log ind på den officielle side, og klik Færdig her i assistenten, når trinnet er klaret."
            }
            "sketchup" {
                if (Test-WindowsSMode) {
                    Open-SetupLink $step.Url
                    Set-StatusText "Windows S-mode kan blokere installation. SketchUp fallback-siden er åbnet."
                    return
                }

                $PrimaryButton.Enabled = $false
                Set-StatusText "Starter SketchUp-installation via winget..."
                [System.Windows.Forms.Application]::DoEvents()
                # Equivalent command: winget install --id $($SetupConfig.SketchUpPackageId) -e --source winget
                $result = Install-SketchUpPackage $SetupConfig.SketchUpPackageId
                if ($result.Success) {
                    Set-StatusText $result.Message
                }
                else {
                    Open-SetupLink $step.Url
                    Set-StatusText "$($result.Message) SketchUp fallback-siden er åbnet."
                }
            }
            "finish" {
                if (New-DashboardShortcut -DashboardPath $DashboardPath -ShortcutPath $ShortcutPath) {
                    Set-StatusText "Skrivebordsgenvej oprettet. Dashboardet åbnes."
                }
                else {
                    Set-StatusText "Dashboardet blev ikke fundet ved $DashboardPath."
                }

                if (Test-Path -LiteralPath $DashboardPath) {
                    Start-Process $DashboardPath
                }
            }
            default {
                Set-StatusText "Ukendt trin: $($step.Kind)"
            }
        }
    }
    catch {
        Set-StatusText "Fejl: $($_.Exception.Message)"
    }
    finally {
        if (-not $Form.IsDisposed) {
            $PrimaryButton.Enabled = $true
        }
    }
}

function Move-NextStep {
    if ($CurrentStepIndex -ge ($SetupSteps.Count - 1)) {
        $Form.Close()
        return
    }

    $script:CurrentStepIndex += 1
    Render-Step
}

$Form = New-Object System.Windows.Forms.Form
$Form.Text = "GF2 IT Setup"
$Form.StartPosition = [System.Windows.Forms.FormStartPosition]::CenterScreen
$Form.Size = New-Object System.Drawing.Size(760, 520)
$Form.MinimumSize = New-Object System.Drawing.Size(720, 500)
$Form.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$Form.BackColor = [System.Drawing.Color]::FromArgb(248, 249, 251)

$HeaderPanel = New-Object System.Windows.Forms.Panel
$HeaderPanel.Dock = [System.Windows.Forms.DockStyle]::Top
$HeaderPanel.Height = 96
$HeaderPanel.BackColor = [System.Drawing.Color]::White
$Form.Controls.Add($HeaderPanel)

$ProgressLabel = New-Object System.Windows.Forms.Label
$ProgressLabel.AutoSize = $true
$ProgressLabel.Location = New-Object System.Drawing.Point(28, 18)
$ProgressLabel.ForeColor = [System.Drawing.Color]::FromArgb(92, 101, 112)
$HeaderPanel.Controls.Add($ProgressLabel)

$TitleLabel = New-Object System.Windows.Forms.Label
$TitleLabel.AutoSize = $false
$TitleLabel.Location = New-Object System.Drawing.Point(24, 40)
$TitleLabel.Size = New-Object System.Drawing.Size(690, 40)
$TitleLabel.Font = New-Object System.Drawing.Font("Segoe UI Semibold", 20)
$TitleLabel.ForeColor = [System.Drawing.Color]::FromArgb(20, 29, 38)
$HeaderPanel.Controls.Add($TitleLabel)

$BodyTextBox = New-Object System.Windows.Forms.TextBox
$BodyTextBox.Multiline = $true
$BodyTextBox.ReadOnly = $true
$BodyTextBox.BorderStyle = [System.Windows.Forms.BorderStyle]::None
$BodyTextBox.BackColor = $Form.BackColor
$BodyTextBox.Location = New-Object System.Drawing.Point(28, 126)
$BodyTextBox.Size = New-Object System.Drawing.Size(690, 150)
$BodyTextBox.Font = New-Object System.Drawing.Font("Segoe UI", 12)
$BodyTextBox.ForeColor = [System.Drawing.Color]::FromArgb(36, 45, 55)
$BodyTextBox.TabStop = $false
$Form.Controls.Add($BodyTextBox)

$WarningPanel = New-Object System.Windows.Forms.Panel
$WarningPanel.Location = New-Object System.Drawing.Point(28, 292)
$WarningPanel.Size = New-Object System.Drawing.Size(690, 62)
$WarningPanel.BackColor = [System.Drawing.Color]::FromArgb(255, 247, 224)
$WarningPanel.Visible = $false
$Form.Controls.Add($WarningPanel)

$WarningLabel = New-Object System.Windows.Forms.Label
$WarningLabel.AutoSize = $false
$WarningLabel.Location = New-Object System.Drawing.Point(14, 11)
$WarningLabel.Size = New-Object System.Drawing.Size(660, 42)
$WarningLabel.ForeColor = [System.Drawing.Color]::FromArgb(102, 73, 0)
$WarningLabel.Font = New-Object System.Drawing.Font("Segoe UI Semibold", 10)
$WarningPanel.Controls.Add($WarningLabel)

$StatusLabel = New-Object System.Windows.Forms.Label
$StatusLabel.AutoSize = $false
$StatusLabel.Location = New-Object System.Drawing.Point(28, 370)
$StatusLabel.Size = New-Object System.Drawing.Size(690, 44)
$StatusLabel.ForeColor = [System.Drawing.Color]::FromArgb(67, 76, 87)
$Form.Controls.Add($StatusLabel)

$ButtonPanel = New-Object System.Windows.Forms.Panel
$ButtonPanel.Dock = [System.Windows.Forms.DockStyle]::Bottom
$ButtonPanel.Height = 78
$ButtonPanel.BackColor = [System.Drawing.Color]::White
$Form.Controls.Add($ButtonPanel)

$PrimaryButton = New-Object System.Windows.Forms.Button
$PrimaryButton.Size = New-Object System.Drawing.Size(210, 38)
$PrimaryButton.Location = New-Object System.Drawing.Point(278, 20)
$PrimaryButton.UseVisualStyleBackColor = $true
$PrimaryButton.Add_Click({ Invoke-StepAction })
$ButtonPanel.Controls.Add($PrimaryButton)

$NextButton = New-Object System.Windows.Forms.Button
$NextButton.Size = New-Object System.Drawing.Size(110, 38)
$NextButton.Location = New-Object System.Drawing.Point(608, 20)
$NextButton.UseVisualStyleBackColor = $true
$NextButton.Add_Click({ Move-NextStep })
$ButtonPanel.Controls.Add($NextButton)

Render-Step
[System.Windows.Forms.Application]::Run($Form)
