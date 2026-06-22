package wizard

import (
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

func TestOpenWifiSettingsDelegatesToOS(t *testing.T) {
	fake := &osfake.Fake{}
	w := New(fake)

	if err := w.OpenWifiSettings(); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if fake.WifiSettingsOpens != 1 {
		t.Errorf("Wi-Fi-indstillinger åbnet %d gange, forventede 1", fake.WifiSettingsOpens)
	}
}

func TestOpenStepLinkDelegatesToOS(t *testing.T) {
	fake := &osfake.Fake{}
	w := New(fake)

	if err := w.OpenURL(steps.URLMoodle); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLMoodle {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLMoodle)
	}
}
