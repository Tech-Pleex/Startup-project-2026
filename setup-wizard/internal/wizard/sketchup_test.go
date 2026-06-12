package wizard

import (
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

func TestSketchUpInstallsViaWingetWhenPossible(t *testing.T) {
	fake := &fakeOS{wingetAvailable: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpInstalled {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpInstalled)
	}
	if len(fake.installedPackages) != 1 || fake.installedPackages[0] != steps.SketchUpPackageID {
		t.Errorf("installerede pakker = %v, forventede [%q]", fake.installedPackages, steps.SketchUpPackageID)
	}
	if len(fake.openedURLs) != 0 {
		t.Errorf("fallback-siden blev åbnet ved vellykket installation: %v", fake.openedURLs)
	}
}

func TestSketchUpFallsBackInSModeWithoutTryingWinget(t *testing.T) {
	fake := &fakeOS{wingetAvailable: true, sMode: true}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if !strings.Contains(outcome.Reason, "S-mode") {
		t.Errorf("begrundelsen nævner ikke S-mode: %q", outcome.Reason)
	}
	if len(fake.installedPackages) != 0 {
		t.Errorf("winget blev forsøgt i S-mode: %v", fake.installedPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetMissing(t *testing.T) {
	fake := &fakeOS{wingetAvailable: false}
	w := New(fake)

	outcome := w.InstallSketchUp()

	if outcome.Action != SketchUpFallback {
		t.Errorf("Action = %q, forventede %q", outcome.Action, SketchUpFallback)
	}
	if len(fake.installedPackages) != 0 {
		t.Errorf("winget blev forsøgt selvom det ikke findes: %v", fake.installedPackages)
	}
	assertOpenedFallback(t, fake)
}

func TestSketchUpFallsBackWhenWingetFails(t *testing.T) {
	fake := &fakeOS{wingetAvailable: true, installErr: errWingetFailed}
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
	cases := map[string]*fakeOS{
		"s-mode":         {wingetAvailable: true, sMode: true},
		"winget mangler": {wingetAvailable: false},
		"winget fejler":  {wingetAvailable: true, installErr: errWingetFailed},
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
	fake := &fakeOS{}
	w := New(fake)

	if err := w.OpenURL(steps.URLMoodle); err != nil {
		t.Fatalf("uventet fejl: %v", err)
	}
	if len(fake.openedURLs) != 1 || fake.openedURLs[0] != steps.URLMoodle {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.openedURLs, steps.URLMoodle)
	}
}

func assertOpenedFallback(t *testing.T, fake *fakeOS) {
	t.Helper()
	if len(fake.openedURLs) != 1 || fake.openedURLs[0] != steps.URLSketchUpFallback {
		t.Errorf("åbnede URL'er = %v, forventede [%q]", fake.openedURLs, steps.URLSketchUpFallback)
	}
}
