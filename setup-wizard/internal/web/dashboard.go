package web

import (
	"encoding/base64"
	"encoding/json"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strconv"
)

// dashboardTmpl er den redesignede start.html gjort selvstændig: Assistenten
// indlejrer font og logo som data-URI'er, så den genererede fil virker offline
// på elevens skrivebord uden løse assets. Se ADR-0001 (Go ejer templaten).
//
//go:embed dashboard.gohtml
var dashboardTmpl string

// studentStepWizardIDs kobler dashboardets studentSteps (i rækkefølge) til de
// trin-ID'er i Assistenten (steps.All) som skal være "done" for at trinnet
// vises med flueben. Et dashboard-trin er færdigt når ALLE dets ID'er er done.
//
// Indeks svarer 1:1 til studentSteps-arrayet i dashboard.gohtml. Ændrer man
// rækkefølgen dér, skal denne liste opdateres (dækket af TestStudentStepMapping).
var studentStepWizardIDs = [][]string{
	{"wifi"},              // 0: Få NEG-login (du har login når du er på Wi-Fi)
	{"wifi"},              // 1: Log på skolens Wi-Fi
	{"office"},            // 2: Åbn Office 365 og skolemail
	{"trimble"},           // 3: Find SketchUp/Trimble-invitationen
	{"moodle"},            // 4: Log på Moodle
	{"praxis", "lectio"},  // 5: Log på PraxisOnline og Lectio
	{"sketchup"},          // 6: Installer eller åbn SketchUp
	{"onedrive"},          // 7: Gem opgaver i OneDrive
}

// StudentStatusFromWizard oversætter Assistentens trin-status (id -> "done"/
// "skipped") til dashboardets fluebens-map (studentSteps-indeks -> "Færdig").
// Kun "done" tæller som flueben; "skipped" og manglende trin gør ikke.
func StudentStatusFromWizard(status map[string]string) map[string]string {
	out := make(map[string]string)
	for index, ids := range studentStepWizardIDs {
		allDone := true
		for _, id := range ids {
			if status[id] != "done" {
				allDone = false
				break
			}
		}
		if allDone {
			out[strconv.Itoa(index)] = "Færdig"
		}
	}
	return out
}

type dashboardData struct {
	FontDataURI       template.URL
	LogoDataURI       template.URL
	StudentStatusJSON template.JS
}

// RenderDashboard skriver en selvstændig, personlig dashboard-HTML til w.
// studentStatus er dashboardets fluebens-map (fra StudentStatusFromWizard).
func RenderDashboard(w io.Writer, studentStatus map[string]string) error {
	font, err := Static.ReadFile("static/fonts/space-grotesk.woff2")
	if err != nil {
		return fmt.Errorf("kunne ikke læse font: %w", err)
	}
	logo, err := Static.ReadFile("static/img/neg-logo.png")
	if err != nil {
		return fmt.Errorf("kunne ikke læse logo: %w", err)
	}
	// map[string]string marshaller til deterministisk (sorteret) JSON, så
	// output er stabilt og testbart.
	statusJSON, err := json.Marshal(studentStatus)
	if err != nil {
		return fmt.Errorf("kunne ikke serialisere trin-status: %w", err)
	}

	tmpl, err := template.New("dashboard").Parse(dashboardTmpl)
	if err != nil {
		return fmt.Errorf("kunne ikke parse dashboard-template: %w", err)
	}
	return tmpl.Execute(w, dashboardData{
		FontDataURI:       dataURI("font/woff2", font),
		LogoDataURI:       dataURI("image/png", logo),
		StudentStatusJSON: template.JS(statusJSON),
	})
}

func dataURI(mime string, b []byte) template.URL {
	return template.URL("data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b))
}
