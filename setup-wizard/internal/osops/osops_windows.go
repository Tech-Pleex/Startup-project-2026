//go:build windows

package osops

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Windows implementerer OS for Windows. Porteret fra scripts/setup-checks.ps1.
type Windows struct{}

var _ OS = Windows{}

// ssidLine matcher "SSID : <navn>" men ikke "BSSID : ..." — samme
// udtryk som PowerShell-originalen.
var ssidLine = regexp.MustCompile(`(?m)^\s*SSID\s*:\s*(.+)$`)

func (Windows) ActiveWifiSSID() (string, error) {
	out, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
	if err != nil {
		// netsh fejler bl.a. når der ikke er noget WLAN-interface;
		// det behandles som "ikke på Wi-Fi", ikke som en fejl.
		return "", nil
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "BSSID") {
			continue
		}
		if m := ssidLine.FindStringSubmatch(line); m != nil {
			if ssid := strings.TrimSpace(m[1]); ssid != "" {
				return ssid, nil
			}
		}
	}
	return "", nil
}

func (Windows) OpenWifiSettings() error {
	return openWithShell("ms-settings:network-wifi")
}

func (Windows) OpenURL(url string) error {
	return openWithShell(url)
}

// openWithShell åbner en URL eller URI med dens standardprogram.
// rundll32-varianten undgår cmd's specialtegns-fortolkning.
func openWithShell(target string) error {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
}

func (Windows) SMode() (bool, error) {
	out, err := exec.Command("reg", "query",
		`HKLM\SYSTEM\CurrentControlSet\Control\CI\Policy`,
		"/v", "SkuPolicyRequired").Output()
	if err != nil {
		// Nøglen findes ikke på maskiner uden S-mode-politik.
		return false, nil
	}
	return strings.Contains(string(out), "0x1"), nil
}

func (Windows) WingetAvailable() bool {
	_, err := exec.LookPath("winget")
	return err == nil
}

func (Windows) InstallSketchUp(packageID string) error {
	cmd := exec.Command("winget", "install",
		"--id", packageID,
		"-e",
		"--source", "winget",
		"--accept-source-agreements",
		"--accept-package-agreements",
	)
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("winget returnerede fejlkode %d", exit.ExitCode())
		}
		return fmt.Errorf("winget kunne ikke startes: %w", err)
	}
	return nil
}
