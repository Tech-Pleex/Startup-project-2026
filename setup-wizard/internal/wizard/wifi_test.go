package wizard

import (
	"errors"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
)

func TestWifiStatusOnTargetNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "NEG"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOnTarget {
		t.Errorf("State = %q, forventede %q", status.State, WifiOnTarget)
	}
	if status.SSID != "NEG" {
		t.Errorf("SSID = %q, forventede %q", status.SSID, "NEG")
	}
}

func TestWifiStatusOnGuestNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "NEG Guest"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOnGuest {
		t.Errorf("State = %q, forventede %q", status.State, WifiOnGuest)
	}
}

func TestWifiStatusOnUnknownNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: "Naboens Netværk"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOther {
		t.Errorf("State = %q, forventede %q", status.State, WifiOther)
	}
	if status.SSID != "Naboens Netværk" {
		t.Errorf("SSID = %q, forventede %q", status.SSID, "Naboens Netværk")
	}
}

func TestWifiStatusWithNoNetwork(t *testing.T) {
	w := New(&osfake.Fake{SSID: ""})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiNone {
		t.Errorf("State = %q, forventede %q", status.State, WifiNone)
	}
}

func TestWifiStatusPropagatesOSError(t *testing.T) {
	osErr := errors.New("netsh fejlede")
	w := New(&osfake.Fake{SSIDErr: osErr})

	_, err := w.WifiStatus()
	if !errors.Is(err, osErr) {
		t.Errorf("err = %v, forventede at den ombryder %v", err, osErr)
	}
}

func TestOpenWifiSettingsDelegatesToOS(t *testing.T) {
	fake := &osfake.Fake{}
	w := New(fake)

	if err := w.OpenWifiSettings(); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if fake.WifiSettingsOpens != 1 {
		t.Errorf("Wi-Fi-indstillinger åbnet %d gange, forventede 1", fake.WifiSettingsOpens)
	}
}
