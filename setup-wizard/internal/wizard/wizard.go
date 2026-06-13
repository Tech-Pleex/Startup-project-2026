// Package wizard indeholder Assistentens trinlogik: Wi-Fi-status og
// SketchUp-installationsflowet. Al OS-adgang sker gennem osops.OS.
package wizard

import (
	"fmt"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

// Wizard binder trinlogikken til en konkret OS-implementering.
type Wizard struct {
	os osops.OS
}

func New(os osops.OS) *Wizard {
	return &Wizard{os: os}
}

// WifiState klassificerer elevens aktive netværk i forhold til NEG.
type WifiState string

const (
	WifiOnTarget WifiState = "target" // på NEG
	WifiOnGuest  WifiState = "guest"  // på NEG Guest
	WifiOther    WifiState = "other"  // på et andet netværk
	WifiNone     WifiState = "none"   // ikke på noget Wi-Fi-netværk
)

// WifiStatus beskriver elevens aktuelle Wi-Fi-situation.
type WifiStatus struct {
	SSID  string
	State WifiState
}

// WifiStatus aflæser det aktive netværk og klassificerer det.
func (w *Wizard) WifiStatus() (WifiStatus, error) {
	ssid, err := w.os.ActiveWifiSSID()
	if err != nil {
		return WifiStatus{}, fmt.Errorf("kunne ikke aflæse Wi-Fi-status: %w", err)
	}

	status := WifiStatus{SSID: ssid}
	switch ssid {
	case "":
		status.State = WifiNone
	case steps.TargetWifi:
		status.State = WifiOnTarget
	case steps.GuestWifi:
		status.State = WifiOnGuest
	default:
		status.State = WifiOther
	}
	return status, nil
}

// OpenWifiSettings åbner styresystemets Wi-Fi-indstillinger.
func (w *Wizard) OpenWifiSettings() error {
	return w.os.OpenWifiSettings()
}
