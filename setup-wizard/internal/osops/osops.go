// Package osops samler alle platformspecifikke kald bag ét interface,
// så trinlogikken kan testes med en fake uden at røre det rigtige OS.
// Windows- og Mac-implementeringerne ligger i hver sin build-taggede fil.
package osops

import (
	"os"
	"path/filepath"
)

// OS er sømmen mod styresystemet. Porteret fra scripts/setup-checks.ps1.
type OS interface {
	// OpenWifiSettings åbner styresystemets egne Wi-Fi-indstillinger.
	OpenWifiSettings() error

	// OpenURL åbner en URL (eller en lokal filsti) i standardprogrammet.
	OpenURL(url string) error

	// SMode rapporterer om Windows kører i S-mode. Altid false på Mac.
	SMode() (bool, error)

	// DesktopDir returnerer stien til elevens skrivebord, hvor Assistenten
	// gemmer det genererede dashboard.
	DesktopDir() (string, error)
}

// desktopDir er den fælles skrivebords-opslagslogik for alle platforme:
// hjemmemappen + "Desktop". Bruges af hver OS-implementering.
func desktopDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Desktop"), nil
}
