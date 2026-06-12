//go:build darwin

package osops

import (
	"errors"
	"os/exec"
	"strings"
)

// Darwin implementerer OS for Mac. S-mode og winget findes ikke på Mac,
// så SketchUp-trinnet fører altid til fallback-siden (jf. PRD'en i #13).
type Darwin struct{}

var _ OS = Darwin{}

func (Darwin) ActiveWifiSSID() (string, error) {
	device, err := wifiDevice()
	if err != nil || device == "" {
		return "", nil
	}

	out, err := exec.Command("networksetup", "-getairportnetwork", device).Output()
	if err != nil {
		return "", nil
	}

	// Forventet format: "Current Wi-Fi Network: <navn>".
	// Uden tilknyttet netværk svarer kommandoen med en fejlbesked i stedet.
	text := strings.TrimSpace(string(out))
	if _, ssid, found := strings.Cut(text, ": "); found {
		return strings.TrimSpace(ssid), nil
	}
	return "", nil
}

// wifiDevice finder Wi-Fi-enhedens navn (typisk "en0") via networksetup.
func wifiDevice() (string, error) {
	out, err := exec.Command("networksetup", "-listallhardwareports").Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "Hardware Port: Wi-Fi") && i+1 < len(lines) {
			if _, device, found := strings.Cut(lines[i+1], "Device: "); found {
				return strings.TrimSpace(device), nil
			}
		}
	}
	return "", nil
}

func (Darwin) OpenWifiSettings() error {
	return exec.Command("open", "x-apple.systempreferences:com.apple.preference.network").Start()
}

func (Darwin) OpenURL(url string) error {
	return exec.Command("open", url).Start()
}

func (Darwin) SMode() (bool, error) {
	return false, nil
}

func (Darwin) WingetAvailable() bool {
	return false
}

func (Darwin) InstallSketchUp(string) error {
	return errors.New("automatisk installation findes ikke på Mac")
}
