// Package osfake er en testimplementering af osops.OS. Den simulerer
// Wi-Fi-svar, S-mode og winget-resultater og registrerer hvilke URL'er
// der ville være åbnet — ingen test rører det rigtige OS.
package osfake

import "github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"

// Fake styres via felterne og registrerer alle kald.
type Fake struct {
	SSID       string
	SSIDErr    error
	InSMode    bool
	WingetOK   bool
	InstallErr error
	OpenURLErr error

	OpenedURLs        []string
	WifiSettingsOpens int
	InstalledPackages []string
}

var _ osops.OS = (*Fake)(nil)

func (f *Fake) ActiveWifiSSID() (string, error) { return f.SSID, f.SSIDErr }

func (f *Fake) OpenWifiSettings() error {
	f.WifiSettingsOpens++
	return nil
}

func (f *Fake) OpenURL(url string) error {
	f.OpenedURLs = append(f.OpenedURLs, url)
	return f.OpenURLErr
}

func (f *Fake) SMode() (bool, error) { return f.InSMode, nil }

func (f *Fake) WingetAvailable() bool { return f.WingetOK }

func (f *Fake) InstallSketchUp(packageID string) error {
	f.InstalledPackages = append(f.InstalledPackages, packageID)
	return f.InstallErr
}
