//go:build linux

package osops

import (
	"os/exec"
)

// Linux er en udviklerstub, så Assistenten kan bygges og afprøves på
// udviklingsmaskinen (WSL). Elever får kun Windows- og Mac-binærer.
type Linux struct{}

var _ OS = Linux{}

func (Linux) OpenWifiSettings() error { return nil }

func (Linux) OpenURL(url string) error {
	return exec.Command("xdg-open", url).Start()
}

func (Linux) SMode() (bool, error) { return false, nil }

// Current returnerer OS-implementeringen for den platform binæren er bygget til.
func Current() OS { return Linux{} }
