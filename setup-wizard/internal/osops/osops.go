// Package osops samler alle platformspecifikke kald bag ét interface,
// så trinlogikken kan testes med en fake uden at røre det rigtige OS.
// Windows- og Mac-implementeringerne ligger i hver sin build-taggede fil.
package osops

// OS er sømmen mod styresystemet. Porteret fra scripts/setup-checks.ps1.
type OS interface {
	// ActiveWifiSSID returnerer SSID'et for det aktive Wi-Fi-netværk,
	// eller "" hvis maskinen ikke er på et Wi-Fi-netværk.
	ActiveWifiSSID() (string, error)

	// OpenWifiSettings åbner styresystemets egne Wi-Fi-indstillinger.
	OpenWifiSettings() error

	// OpenURL åbner en URL i standardbrowseren.
	OpenURL(url string) error

	// SMode rapporterer om Windows kører i S-mode. Altid false på Mac.
	SMode() (bool, error)

	// WingetAvailable rapporterer om winget findes. Altid false på Mac.
	WingetAvailable() bool

	// InstallSketchUp kører winget-installationen af den angivne pakke
	// og returnerer en fejl hvis winget afslutter med en fejlkode.
	InstallSketchUp(packageID string) error
}
