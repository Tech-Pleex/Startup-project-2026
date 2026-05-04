function Get-ActiveWifiSsid {
    $interfaces = netsh wlan show interfaces 2>$null

    foreach ($line in $interfaces) {
        if ($line -match '^\s*SSID\s+:\s*(.*)\s*$') {
            $ssid = $Matches[1].Trim()
            if ($ssid.Length -gt 0) {
                return $ssid
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
    $policy = Get-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\CI\Policy' -ErrorAction SilentlyContinue
    return ($null -ne $policy -and $policy.SkuPolicyRequired -eq 1)
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

    if (-not (Test-Path -LiteralPath $DashboardPath)) {
        return $false
    }

    $resolvedDashboardPath = (Resolve-Path -LiteralPath $DashboardPath).ProviderPath
    $shortcutUrl = [System.Uri]::new($resolvedDashboardPath).AbsoluteUri
    $content = @(
        '[InternetShortcut]'
        "URL=$shortcutUrl"
    )

    Set-Content -LiteralPath $ShortcutPath -Value $content -Encoding UTF8
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

    $arguments = @(
        'install'
        '--id'
        $PackageId
        '-e'
        '--source'
        'winget'
        '--accept-source-agreements'
        '--accept-package-agreements'
    )

    $process = Start-Process -FilePath 'winget' -ArgumentList $arguments -Wait -PassThru -WindowStyle Normal
    if ($process.ExitCode -eq 0) {
        return [ordered]@{
            Success = $true
            Message = "SketchUp blev installeret."
        }
    }

    return [ordered]@{
        Success = $false
        Message = "winget afsluttede med kode $($process.ExitCode)."
    }
}
