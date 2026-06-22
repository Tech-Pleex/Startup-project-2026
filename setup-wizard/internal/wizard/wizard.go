// Package wizard indeholder Assistentens trinlogik for Wi-Fi-status,
// S-mode og åbning af officielle sider. Al OS-adgang sker gennem osops.OS.
package wizard

import "github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"

// Wizard binder trinlogikken til en konkret OS-implementering.
type Wizard struct {
	os osops.OS
}

func New(os osops.OS) *Wizard {
	return &Wizard{os: os}
}

// OpenWifiSettings åbner styresystemets Wi-Fi-indstillinger.
func (w *Wizard) OpenWifiSettings() error {
	return w.os.OpenWifiSettings()
}

// SMode rapporterer om Windows kører i S-mode.
func (w *Wizard) SMode() (bool, error) {
	return w.os.SMode()
}

// OpenURL åbner en officiel side i elevens standardbrowser.
func (w *Wizard) OpenURL(url string) error {
	return w.os.OpenURL(url)
}
