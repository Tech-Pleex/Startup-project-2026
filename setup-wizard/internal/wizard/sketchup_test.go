package wizard

import (
	"errors"
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops/osfake"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

var errWingetFailed = errors.New("winget returnerede fejlkode 1")

func TestSketchUpInstallsViaWingetWhenPossible(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpInstalled {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpInstalled)
	}
	if len(fake.InstalledPackages) != 1 || fake.InstalledPackages[0] != steps.SketchUpPackageID {
		t.Errorf("installerede pakker = %v, forventede [%q]", fake.InstalledPackages, steps.SketchUpPackageID)
	}
	if len(fake.OpenedURLs) != 0 {
		t.Errorf("fallback-siden blev åbnet ved vellykket installation: %v", fake.OpenedURLs)
	}
}

func TestSketchUpFallsBackInSModeWithoutTryingWinget(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true, InSMode: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if !strings.Contains(outcome.Reason, "S-mode") {
		t.Errorf("begrundelsen nævner ikke S-mode: %q", outcome.Reason)
	}
	if len(fake.InstalledPackages) != 0 {
		t.Errorf("winget blev forsøgt i S-mode: %v", fake.InstalledPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetMissing(t *testing.T) {
	fake := &osfake.Fake{WingetOK: false}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if len(fake.InstalledPackages) != 0 {
		t.Errorf("winget blev forsøgt selvom det ikke findes: %v", fake.InstalledPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetFails(t *testing.T) {
	fake := &osfake.Fake{WingetOK: true, InstallErr: errWingetFailed}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if !strings.Contains(outcome.Reason, errWingetFailed.Error()) {
		t.Errorf("begrundelsen indeholder ikke winget-fejlen: %q", outcome.Reason)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpReasonsAreDanishStudentFacing(t *testing.T) {
	cases := map[string]*osfake.Fake{
		"s-mode":         {WingetOK: true, InSMode: true},
		"winget mangler": {WingetOK: false},
		"winget fejler":  {WingetOK: true, InstallErr: errWingetFailed},
	}
	for name, fake := range cases {
		outcome := New(fake).InstallSketchUp()
		if outcome.Reason == "" {
			t.Errorf("%s: begrundelsen er tom", name)
		}
		if strings.Contains(strings.ToLower(outcome.Reason), "setup-wizard") {
			t.Errorf("%s: begrundelsen indeholder \"setup-wizard\": %q", name, outcome.Reason)
		}
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

func assertOpenedFallback(t *testing.T, fake *osfake.Fake) {
	t.Helper()
	if len(fake.OpenedURLs) != 1 || fake.OpenedURLs[0] != steps.URLSketchUpFallback {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.OpenedURLs, steps.URLSketchUpFallback)
	}
}
