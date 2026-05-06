function Get-ActiveWifiSsid {
    $Output = netsh wlan show interfaces 2>$null
    if (-not $Output) {
        return $null
    }

    foreach ($Line in $Output) {
        if ($Line -match "^\s*SSID\s*:\s*(.+)$" -and $Line -notmatch "BSSID") {
            $Ssid = $Matches[1].Trim()
            if ($Ssid.Length -gt 0) {
                return $Ssid
            }
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
    return [bool]($Policy -and $Policy.SkuPolicyRequired -eq 1)
}

function Open-SetupLink {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Url
    )

    Start-Process $Url
}

function New-DashboardShortcut {
    param(
        [Parameter(Mandatory = $true)]
        [string]$DashboardPath,

        [Parameter(Mandatory = $true)]
        [string]$ShortcutPath
    )

    if (-not (Test-Path -LiteralPath $DashboardPath -PathType Leaf)) {
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
    param(
        [Parameter(Mandatory = $true)]
        [string]$PackageId
    )

    if (-not (Test-WingetAvailable)) {
        return [ordered]@{
            Success = $false
            Message = "winget blev ikke fundet."
        }
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
        return [ordered]@{
            Success = $true
            Message = "SketchUp-installationen blev startet eller gennemfort."
        }
    }

    return [ordered]@{
        Success = $false
        Message = "winget returnerede fejlkode $($Process.ExitCode)."
    }
}
