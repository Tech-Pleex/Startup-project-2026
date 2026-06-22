//go:build windows

package osops

import (
	"os/exec"
	"strings"
	"syscall"
)

// Windows implementerer OS for Windows. Porteret fra scripts/setup-checks.ps1.
type Windows struct{}

var _ OS = Windows{}

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
	cmd := exec.Command("reg", "query",
		`HKLM\SYSTEM\CurrentControlSet\Control\CI\Policy`,
		"/v", "SkuPolicyRequired")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		// Nøglen findes ikke på maskiner uden S-mode-politik.
		return false, nil
	}
	return strings.Contains(string(out), "0x1"), nil
}
