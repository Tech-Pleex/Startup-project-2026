package wizard

import (
	"errors"
	"testing"
)

func TestWifiStatusOnTargetNetwork(t *testing.T) {
	w := New(&fakeOS{ssid: "NEG"})

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
	w := New(&fakeOS{ssid: "NEG Guest"})

	status, err := w.WifiStatus()
	if err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if status.State != WifiOnGuest {
		t.Errorf("State = %q, forventede %q", status.State, WifiOnGuest)
	}
}

func TestWifiStatusOnUnknownNetwork(t *testing.T) {
	w := New(&fakeOS{ssid: "Naboens Netværk"})

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
	w := New(&fakeOS{ssid: ""})

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
	w := New(&fakeOS{ssidErr: osErr})

	_, err := w.WifiStatus()
	if !errors.Is(err, osErr) {
		t.Errorf("err = %v, forventede at den ombryder %v", err, osErr)
	}
}

func TestOpenWifiSettingsDelegatesToOS(t *testing.T) {
	fake := &fakeOS{}
	w := New(fake)

	if err := w.OpenWifiSettings(); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if fake.openedWifiSettings != 1 {
		t.Errorf("Wi-Fi-indstillinger åbnet %d gange, forventede 1", fake.openedWifiSettings)
	}
}
