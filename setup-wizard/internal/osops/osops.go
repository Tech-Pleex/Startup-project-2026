// Package osops samler alle platformspecifikke kald bag ét interface,
// så trinlogikken kan testes med en fake uden at røre det rigtige OS.
// Windows- og Mac-implementeringerne ligger i hver sin build-taggede fil.
package osops

// OS er sømmen mod styresystemet. Porteret fra scripts/setup-checks.ps1.
type OS interface {
	// OpenWifiSettings åbner styresystemets egne Wi-Fi-indstillinger.
	OpenWifiSettings() error

	// OpenURL åbner en URL i standardbrowseren.
	OpenURL(url string) error

	// SMode rapporterer om Windows kører i S-mode. Altid false på Mac.
	SMode() (bool, error)
}
