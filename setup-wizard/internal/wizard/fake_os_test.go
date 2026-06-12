package wizard

import "errors"

// fakeOS er en testimplementering af osops.OS. Den simulerer Wi-Fi-svar,
// S-mode og winget-resultater og registrerer hvilke URL'er der ville være
// åbnet — ingen test rører det rigtige OS.
type fakeOS struct {
	ssid            string
	ssidErr         error
	sMode           bool
	wingetAvailable bool
	installErr      error

	openedURLs         []string
	openedWifiSettings int
	installedPackages  []string
}

func (f *fakeOS) ActiveWifiSSID() (string, error) {
	return f.ssid, f.ssidErr
}

func (f *fakeOS) OpenWifiSettings() error {
	f.openedWifiSettings++
	return nil
}

func (f *fakeOS) OpenURL(url string) error {
	f.openedURLs = append(f.openedURLs, url)
	return nil
}

func (f *fakeOS) SMode() (bool, error) {
	return f.sMode, nil
}

func (f *fakeOS) WingetAvailable() bool {
	return f.wingetAvailable
}

func (f *fakeOS) InstallSketchUp(packageID string) error {
	f.installedPackages = append(f.installedPackages, packageID)
	return f.installErr
}

var errWingetFailed = errors.New("winget returnerede fejlkode 1")
