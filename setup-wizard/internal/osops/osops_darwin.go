//go:build darwin

package osops

import (
	"os/exec"
)

// Darwin implementerer OS for Mac. S-mode findes ikke på Mac.
type Darwin struct{}

var _ OS = Darwin{}

func (Darwin) OpenWifiSettings() error {
	return exec.Command("open", "x-apple.systempreferences:com.apple.preference.network").Start()
}

func (Darwin) OpenURL(url string) error {
	return exec.Command("open", url).Start()
}

func (Darwin) SMode() (bool, error) {
	return false, nil
}
