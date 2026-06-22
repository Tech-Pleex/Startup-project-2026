// Package osfake er en testimplementering af osops.OS. Den simulerer
// S-mode og registrerer hvilke URL'er
// der ville være åbnet — ingen test rører det rigtige OS.
package osfake

import "github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"

// Fake styres via felterne og registrerer alle kald.
type Fake struct {
	InSMode    bool
	SModeErr   error
	OpenURLErr error

	OpenedURLs        []string
	WifiSettingsOpens int
}

var _ osops.OS = (*Fake)(nil)

func (f *Fake) OpenWifiSettings() error {
	f.WifiSettingsOpens++
	return nil
}

func (f *Fake) OpenURL(url string) error {
	f.OpenedURLs = append(f.OpenedURLs, url)
	return f.OpenURLErr
}

func (f *Fake) SMode() (bool, error) { return f.InSMode, f.SModeErr }
