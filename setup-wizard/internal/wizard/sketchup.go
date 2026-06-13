package wizard

import (
	"fmt"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/steps"
)

// SketchUpAction angiver hvordan SketchUp-trinnet endte.
type SketchUpAction string

const (
	SketchUpInstalled SketchUpAction = "installed"
	SketchUpFallback  SketchUpAction = "fallback"
)

// SketchUpOutcome beskriver resultatet af installationsforsøget.
// Reason er en elevvendt dansk forklaring ved fallback.
type SketchUpOutcome struct {
	Action SketchUpAction
	Reason string
}

// InstallSketchUp forsøger at installere SketchUp via winget. Hvis maskinen
// er i S-mode, winget mangler (fx på Mac), eller installationen fejler,
// åbnes fallback-siden i stedet, så eleven altid har en vej videre.
func (w *Wizard) InstallSketchUp() SketchUpOutcome {
	if sMode, _ := w.os.SMode(); sMode {
		return w.sketchUpFallback("Din computer kører Windows S-mode, som blokerer installationen. Assistenten har åbnet SketchUp-siden, hvor du kan hente programmet manuelt.")
	}

	if !w.os.WingetAvailable() {
		return w.sketchUpFallback("Automatisk installation er ikke tilgængelig på denne computer. Assistenten har åbnet SketchUp-siden, hvor du kan hente programmet manuelt.")
	}

	if err := w.os.InstallSketchUp(steps.SketchUpPackageID); err != nil {
		return w.sketchUpFallback(fmt.Sprintf("Den automatiske installation fejlede (%s). Assistenten har åbnet SketchUp-siden, hvor du kan hente programmet manuelt.", err))
	}

	return SketchUpOutcome{Action: SketchUpInstalled}
}

func (w *Wizard) sketchUpFallback(reason string) SketchUpOutcome {
	_ = w.os.OpenURL(steps.URLSketchUpFallback)
	return SketchUpOutcome{Action: SketchUpFallback, Reason: reason}
}

// OpenURL åbner en officiel side i elevens standardbrowser.
func (w *Wizard) OpenURL(url string) error {
	return w.os.OpenURL(url)
}
