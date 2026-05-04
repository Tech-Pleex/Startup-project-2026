$Script:SetupConfig = [ordered]@{
    Title = "GF2 IT Setup"
    DashboardFile = "start.html"
    DesktopShortcutName = "GF2 IT Dashboard.url"
    TargetWifi = "NEG"
    GuestWifi = "NEG Guest"
    SketchUpPackageId = "Trimble.SketchUp.2026"
    SafetyText = "Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login. Elever indtaster kun oplysninger på officielle sider og i Windows' egne indstillinger."
    Urls = [ordered]@{
        Office = "https://www.office.com/"
        Moodle = "https://online.neg.dk/login/index.php"
        Lectio = "https://www.lectio.dk/lectio/769/default.aspx"
        Praxis = "https://online.praxis.dk/"
        OneDrive = "https://www.office.com/launch/onedrive"
        TrimbleInvitation = "https://www.office.com/"
        SketchUpFallback = "https://sketchup.trimble.com/"
        Neg = "https://www.neg.dk/"
    }
}

$Script:SetupSteps = @(
    [ordered]@{
        Id = "welcome"
        Title = "Velkommen"
        Kind = "manual"
        Body = "GF2 IT Setup hjælper dig de rigtige steder hen. Assistenten beder aldrig om adgangskoder, MitID eller UNI-Login."
        Button = "Start"
    }
    [ordered]@{
        Id = "wifi"
        Title = "Wi-Fi"
        Kind = "wifi"
        Body = "Du skal ende på NEG-netværket. Hvis du allerede er på NEG, markerer assistenten trinnet som gennemført. Brug kun NEG Guest som midlertidigt gæstenet, hvis skolens personale beder om det."
        Button = "Åbn Wi-Fi-indstillinger"
    }
    [ordered]@{
        Id = "office"
        Title = "Office 365 / skolemail"
        Kind = "link"
        Body = "Log ind på Office og åbn din skolemail. Assistenten ser aldrig dine loginoplysninger."
        Url = $Script:SetupConfig.Urls.Office
        Button = "Åbn Office"
    }
    [ordered]@{
        Id = "trimble"
        Title = "Trimble invitation"
        Kind = "link"
        Body = "Find mailen fra Trimble eller SketchUp i din skolemail, klik invitationslinket, og følg flowet med din skolemail."
        Url = $Script:SetupConfig.Urls.TrimbleInvitation
        Button = "Åbn skolemail"
    }
    [ordered]@{
        Id = "moodle"
        Title = "Moodle"
        Kind = "link"
        Body = "Åbn Moodle og kontroller at dine GF2-rum vises."
        Url = $Script:SetupConfig.Urls.Moodle
        Button = "Åbn Moodle"
    }
    [ordered]@{
        Id = "praxis"
        Title = "PraxisOnline"
        Kind = "link"
        Body = "Login sker via UNI-Login, hvor du vælger MitID."
        Warning = "Vigtigt: Brug din skolemail til PraxisOnline. Eksempel: neg04026@edu.neg.dk"
        Url = $Script:SetupConfig.Urls.Praxis
        Button = "Åbn PraxisOnline"
    }
    [ordered]@{
        Id = "lectio"
        Title = "Lectio"
        Kind = "link"
        Body = "Log ind med UNI-Login, hvor du vælger MitID. Assistenten åbner kun Lectio og ser aldrig dine oplysninger."
        Url = $Script:SetupConfig.Urls.Lectio
        Button = "Åbn Lectio"
    }
    [ordered]@{
        Id = "onedrive"
        Title = "OneDrive"
        Kind = "link"
        Body = "Åbn OneDrive via Office 365 og kontroller at du kan se dine filer."
        Url = $Script:SetupConfig.Urls.OneDrive
        Button = "Åbn OneDrive"
    }
    [ordered]@{
        Id = "sketchup"
        Title = "SketchUp"
        Kind = "sketchup"
        Body = "Assistenten kan forsøge at installere SketchUp via winget-pakken Trimble.SketchUp.2026. Hvis det ikke virker, bruger du manuel fallback til SketchUp-siden."
        Url = $Script:SetupConfig.Urls.SketchUpFallback
        Button = "Installer SketchUp"
    }
    [ordered]@{
        Id = "finish"
        Title = "Færdig"
        Kind = "finish"
        Body = "Assistenten opretter en genvej til dashboardet og åbner dashboardet."
        Button = "Åbn dashboard"
    }
)
