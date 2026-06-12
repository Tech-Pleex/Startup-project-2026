package steps

import (
	"strings"
	"testing"
)

func TestAllReturnsTenStepsInOrder(t *testing.T) {
	want := []string{
		"welcome", "wifi", "office", "trimble", "moodle",
		"praxis", "lectio", "onedrive", "sketchup", "finish",
	}

	all := All()
	if len(all) != len(want) {
		t.Fatalf("All() returnerede %d trin, forventede %d", len(all), len(want))
	}
	for i, step := range all {
		if step.ID != want[i] {
			t.Errorf("trin %d har id %q, forventede %q", i, step.ID, want[i])
		}
	}
}

func TestEveryStepHasStudentFacingText(t *testing.T) {
	for _, step := range All() {
		if step.Title == "" {
			t.Errorf("trin %q mangler titel", step.ID)
		}
		if step.Body == "" {
			t.Errorf("trin %q mangler brødtekst", step.ID)
		}
		if step.Button == "" {
			t.Errorf("trin %q mangler knaptekst", step.ID)
		}
	}
}

func TestNoStudentFacingTextSaysSetupWizard(t *testing.T) {
	for _, step := range All() {
		for _, text := range []string{step.Title, step.Body, step.Warning, step.Button} {
			if strings.Contains(strings.ToLower(text), "setup-wizard") {
				t.Errorf("trin %q indeholder det forbudte ord \"setup-wizard\": %q", step.ID, text)
			}
		}
	}
}

func TestStepKinds(t *testing.T) {
	wantKinds := map[string]Kind{
		"welcome":  KindManual,
		"wifi":     KindWifi,
		"office":   KindLink,
		"trimble":  KindLink,
		"moodle":   KindLink,
		"praxis":   KindLink,
		"lectio":   KindLink,
		"onedrive": KindLink,
		"sketchup": KindSketchUp,
		"finish":   KindFinish,
	}
	for _, step := range All() {
		if step.Kind != wantKinds[step.ID] {
			t.Errorf("trin %q har type %q, forventede %q", step.ID, step.Kind, wantKinds[step.ID])
		}
	}
}

func TestLinkStepsHaveURLs(t *testing.T) {
	wantURLs := map[string]string{
		"office":   "https://www.office.com/",
		"trimble":  "https://www.office.com/",
		"moodle":   "https://online.neg.dk/login/index.php",
		"praxis":   "https://online.praxis.dk/",
		"lectio":   "https://www.lectio.dk/lectio/769/default.aspx",
		"onedrive": "https://www.office.com/launch/onedrive",
		"sketchup": "https://sketchup.trimble.com/",
	}
	for _, step := range All() {
		want, ok := wantURLs[step.ID]
		if !ok {
			continue
		}
		if step.URL != want {
			t.Errorf("trin %q har URL %q, forventede %q", step.ID, step.URL, want)
		}
	}
}

func TestPraxisStepWarnsAboutSchoolMail(t *testing.T) {
	step := mustFind(t, "praxis")
	if !strings.Contains(step.Warning, "skolemail") {
		t.Errorf("praxis-trinnets advarsel nævner ikke skolemail: %q", step.Warning)
	}
}

func TestSafetyTextMentionsCredentialsMitIDAndUNILogin(t *testing.T) {
	for _, want := range []string{"adgangskoder", "MitID", "UNI-Login"} {
		if !strings.Contains(SafetyText, want) {
			t.Errorf("sikkerhedsteksten mangler %q: %q", want, SafetyText)
		}
	}
}

func TestWelcomeStepRepeatsSafetyPrinciple(t *testing.T) {
	step := mustFind(t, "welcome")
	for _, want := range []string{"adgangskoder", "MitID", "UNI-Login"} {
		if !strings.Contains(step.Body, want) {
			t.Errorf("velkomsttrinnet nævner ikke %q i brødteksten: %q", want, step.Body)
		}
	}
}

func TestConfigConstants(t *testing.T) {
	if TargetWifi != "NEG" {
		t.Errorf("TargetWifi = %q, forventede %q", TargetWifi, "NEG")
	}
	if GuestWifi != "NEG Guest" {
		t.Errorf("GuestWifi = %q, forventede %q", GuestWifi, "NEG Guest")
	}
	if SketchUpPackageID != "Trimble.SketchUp.2026" {
		t.Errorf("SketchUpPackageID = %q, forventede %q", SketchUpPackageID, "Trimble.SketchUp.2026")
	}
}

func TestFinishStepOpensDashboard(t *testing.T) {
	for _, s := range All() {
		if s.ID != "finish" {
			continue
		}
		if s.URL != URLDashboard {
			t.Errorf("finish.URL = %q, forventede %q", s.URL, URLDashboard)
		}
		return
	}
	t.Fatal("finish-trinnet findes ikke")
}

func mustFind(t *testing.T, id string) Step {
	t.Helper()
	for _, step := range All() {
		if step.ID == id {
			return step
		}
	}
	t.Fatalf("trin %q findes ikke", id)
	return Step{}
}
