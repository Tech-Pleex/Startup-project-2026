// Package steps definerer Assistentens trinforløb: de ti trin med al
// elevvendt tekst (på dansk) samt de faste konfigurationsværdier.
// Porteret fra scripts/setup-config.ps1.
package steps

// Kind angiver hvilken slags handling et trin kræver.
type Kind string

const (
	KindManual   Kind = "manual"
	KindWifi     Kind = "wifi"
	KindLink     Kind = "link"
	KindSketchUp Kind = "sketchup"
	KindFinish   Kind = "finish"
)

// Konfigurationsværdier porteret fra $Script:SetupConfig.
const (
	Title             = "GF2 IT Setup"
	TargetWifi        = "NEG"
	GuestWifi         = "NEG Guest"
	SketchUpPackageID = "Trimble.SketchUp.2026"

	SafetyText = "Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login. " +
		"Elever indtaster kun oplysninger på officielle sider og i Windows' egne indstillinger."
)

// URL'er til de officielle sider.
const (
	URLOffice            = "https://www.office.com/"
	URLMoodle            = "https://online.neg.dk/login/index.php"
	URLLectio            = "https://www.lectio.dk/lectio/769/default.aspx"
	URLPraxis            = "https://online.praxis.dk/"
	URLOneDrive          = "https://www.office.com/launch/onedrive"
	URLTrimbleInvitation = "https://www.office.com/"
	URLSketchUpFallback  = "https://sketchup.trimble.com/"
	URLDashboard         = "https://tech-pleex.github.io/Startup-project-2026/start.html"
)

// Step er ét trin i Assistentens forløb. Al tekst er elevvendt og på dansk.
type Step struct {
	ID      string
	Title   string
	Kind    Kind
	Body    string
	Warning string
	URL     string
	Button  string
}

// All returnerer de ti trin i den rækkefølge eleven møder dem.
func All() []Step {
	return []Step{
		{
			ID:     "welcome",
			Title:  "Velkommen",
			Kind:   KindManual,
			Body:   "GF2 IT Setup hjælper dig de rigtige steder hen. Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login.",
			Button: "Start",
		},
		{
			ID:     "wifi",
			Title:  "Wi-Fi",
			Kind:   KindWifi,
			Body:   "Du skal ende på NEG-netværket. Hvis du allerede er på NEG, markerer Assistenten trinnet som gennemført. Brug kun NEG Guest som midlertidigt gæstenet, hvis skolens personale beder om det.",
			Button: "Åbn Wi-Fi-indstillinger",
		},
		{
			ID:     "office",
			Title:  "Office 365 / skolemail",
			Kind:   KindLink,
			Body:   "Log ind på Office og åbn din skolemail. Assistenten ser aldrig dine loginoplysninger.",
			URL:    URLOffice,
			Button: "Åbn Office",
		},
		{
			ID:     "trimble",
			Title:  "Trimble invitation",
			Kind:   KindLink,
			Body:   "Find mailen fra Trimble eller SketchUp i din skolemail, klik invitationslinket, og følg flowet med din skolemail.",
			URL:    URLTrimbleInvitation,
			Button: "Åbn skolemail",
		},
		{
			ID:     "moodle",
			Title:  "Moodle",
			Kind:   KindLink,
			Body:   "Åbn Moodle og kontroller at dine GF2-rum vises.",
			URL:    URLMoodle,
			Button: "Åbn Moodle",
		},
		{
			ID:      "praxis",
			Title:   "PraxisOnline",
			Kind:    KindLink,
			Body:    "Login sker via UNI-Login, hvor du vælger MitID.",
			Warning: "Vigtigt: Brug din skolemail til PraxisOnline. Eksempel: neg04026@edu.neg.dk",
			URL:     URLPraxis,
			Button:  "Åbn PraxisOnline",
		},
		{
			ID:     "lectio",
			Title:  "Lectio",
			Kind:   KindLink,
			Body:   "Log ind med UNI-Login, hvor du vælger MitID. Assistenten åbner kun Lectio og ser aldrig dine oplysninger.",
			URL:    URLLectio,
			Button: "Åbn Lectio",
		},
		{
			ID:     "onedrive",
			Title:  "OneDrive",
			Kind:   KindLink,
			Body:   "Åbn OneDrive via Office 365 og kontroller at du kan se dine filer.",
			URL:    URLOneDrive,
			Button: "Åbn OneDrive",
		},
		{
			ID:     "sketchup",
			Title:  "SketchUp",
			Kind:   KindSketchUp,
			Body:   "Assistenten kan forsøge at installere SketchUp via winget-pakken Trimble.SketchUp.2026. Hvis det ikke virker, bruger du manuel fallback til SketchUp-siden.",
			URL:    URLSketchUpFallback,
			Button: "Installer SketchUp",
		},
		{
			ID:     "finish",
			Title:  "Færdig",
			Kind:   KindFinish,
			Body:   "Du er igennem alle trin. Assistenten kan åbne dashboardet med hurtige links til skolesystemerne, så du nemt finder dem igen.",
			URL:    URLDashboard,
			Button: "Åbn dashboard",
		},
	}
}
