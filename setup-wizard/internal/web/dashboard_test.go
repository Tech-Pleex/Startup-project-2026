package web

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

// TestStudentStepMapping sikrer at hver dashboard-trin-mapning peger på ID'er
// der faktisk findes i Assistentens trin — så mapningen ikke stille forældes.
func TestStudentStepMapping(t *testing.T) {
	valid := map[string]bool{}
	for _, s := range steps.All() {
		valid[s.ID] = true
	}
	for index, ids := range studentStepWizardIDs {
		if len(ids) == 0 {
			t.Errorf("studentStep %d har ingen trin-ID'er", index)
		}
		for _, id := range ids {
			if !valid[id] {
				t.Errorf("studentStep %d peger på ukendt trin-ID %q", index, id)
			}
		}
	}
}

func TestStudentStatusFromWizard(t *testing.T) {
	// Wi-Fi + Office done → studentSteps 0,1,2 færdige. Praxis uden Lectio
	// tæller ikke (kombitrin 5 kræver begge). Skipped tæller ikke.
	status := map[string]string{
		"wifi":   "done",
		"office": "done",
		"praxis": "done",
		"moodle": "skipped",
	}
	got := StudentStatusFromWizard(status)

	wantDone := []string{"0", "1", "2"}
	for _, k := range wantDone {
		if got[k] != "Færdig" {
			t.Errorf("forventede studentStep %s = Færdig, fik %q", k, got[k])
		}
	}
	if _, ok := got["4"]; ok {
		t.Error("moodle var skipped — studentStep 4 må ikke være færdig")
	}
	if _, ok := got["5"]; ok {
		t.Error("kun praxis (ikke lectio) done — kombitrin 5 må ikke være færdig")
	}
}

func TestStudentStatusCombinedStep(t *testing.T) {
	got := StudentStatusFromWizard(map[string]string{"praxis": "done", "lectio": "done"})
	if got["5"] != "Færdig" {
		t.Errorf("praxis+lectio begge done → studentStep 5 skal være Færdig, fik %q", got["5"])
	}
}

func TestRenderDashboardSelfContained(t *testing.T) {
	var buf bytes.Buffer
	status := StudentStatusFromWizard(map[string]string{"wifi": "done", "office": "done"})
	if err := RenderDashboard(&buf, status); err != nil {
		t.Fatalf("RenderDashboard fejlede: %v", err)
	}
	html := buf.String()

	// Selvstændig: font + logo skal være indlejret som data-URI'er, ingen
	// løse asset-stier tilbage.
	for _, needle := range []string{"data:font/woff2;base64,", "data:image/png;base64,", "data:image/jpeg;base64,"} {
		if !strings.Contains(html, needle) {
			t.Errorf("output mangler indlejret asset %q", needle)
		}
	}
	for _, leftover := range []string{"assets/fonts/", "assets/neg-logo.png", "assets/neg-hero"} {
		if strings.Contains(html, leftover) {
			t.Errorf("output har stadig løs asset-sti %q", leftover)
		}
	}

	// Elevens fremdrift skal være bagt ind i JS-objektet.
	if !strings.Contains(html, `const assistentStudentStatus = {`) {
		t.Error("output mangler injiceret assistentStudentStatus")
	}
	if !strings.Contains(html, `"0":"Færdig"`) {
		t.Errorf("output mangler flueben for studentStep 0; fik status %v", status)
	}
}
