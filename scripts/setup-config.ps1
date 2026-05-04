$Script:SetupConfig = [ordered]@{
    Title = "GF2 IT Setup"
    DashboardFile = "start.html"
    DesktopShortcutName = "GF2 IT Dashboard.url"
    TargetWifi = "NEG"
    GuestWifi = "NEG Guest"
    SketchUpPackageId = "Trimble.SketchUp.2026"
    SafetyText = "Assistenten beder aldrig om adgangskoder, MitID-koder eller UNI-Login-koder. Elever logger kun ind på officielle sider og i Windows-indstillinger."
    Urls = [ordered]@{
        Office = "https://www.office.com/"
        Moodle = "https://online.neg.dk/login/index.php"
        Lectio = "https://www.lectio.dk/lectio/769/default.aspx"
        Praxis = "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery"
        OneDrive = "https://www.office.com/launch/onedrive"
        Trimble = "https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg"
        Neg = "https://www.neg.dk/"
    }
}

$Script:SetupSteps = @(
    [ordered]@{
        Id = "welcome"
        Title = "Velkommen"
        Text = "Guiden hjelper med at abne skolens systemer og kontrollere de vigtigste Windows-indstillinger. Den gemmer ingen adgangskoder, MitID-oplysninger eller UNI-Login-oplysninger."
    }
    [ordered]@{
        Id = "wifi"
        Title = "Wi-Fi"
        Text = "Forbind til skolens Wi-Fi '$($Script:SetupConfig.TargetWifi)' med dit NEG-login. Brug kun '$($Script:SetupConfig.GuestWifi)' som midlertidigt gaestenet, hvis skolens personale beder om det."
        Settings = "ms-settings:network-wifi"
    }
    [ordered]@{
        Id = "office"
        Title = "Office 365 / skolemail"
        Text = "Abn Office 365, log ind med skolemailen, og bekraeft at Outlook og Office virker."
        Url = $Script:SetupConfig.Urls.Office
    }
    [ordered]@{
        Id = "trimble"
        Title = "Trimble"
        Text = "Abn Trimble-login, og brug den officielle Trimble-side til konto eller skoleadgang."
        Url = $Script:SetupConfig.Urls.Trimble
    }
    [ordered]@{
        Id = "moodle"
        Title = "Moodle"
        Text = "Abn Moodle og log ind med NEG-login. Tjek at GF2-rum og materialer vises."
        Url = $Script:SetupConfig.Urls.Moodle
    }
    [ordered]@{
        Id = "praxis"
        Title = "PraxisOnline"
        Warning = "Brug skolemailen, for eksempel neg04026@edu.neg.dk, hvis PraxisOnline sporger efter mail."
        Text = "Login sker via UNI-Login, hvor eleven vaelger MitID som loginmetode."
        Url = $Script:SetupConfig.Urls.Praxis
    }
    [ordered]@{
        Id = "lectio"
        Title = "Lectio"
        Text = "Abn Lectio. Vaelg UNI-Login, og brug MitID som separat godkendelse, hvis siden beder om det."
        Url = $Script:SetupConfig.Urls.Lectio
    }
    [ordered]@{
        Id = "onedrive"
        Title = "OneDrive"
        Text = "Abn OneDrive fra Office 365, og bekraeft at skolens lagerplads og filer er tilgaengelige."
        Url = $Script:SetupConfig.Urls.OneDrive
    }
    [ordered]@{
        Id = "sketchup"
        Title = "SketchUp"
        Text = "Installer SketchUp med winget-pakken Trimble.SketchUp.2026, og log derefter ind via Trimble."
        PackageId = $Script:SetupConfig.SketchUpPackageId
        Url = $Script:SetupConfig.Urls.Trimble
    }
    [ordered]@{
        Id = "finish"
        Title = "Faerdig"
        Text = "Afslut med at abne dashboardet og kontroller, at de vigtigste systemer er markeret som klaret."
        File = $Script:SetupConfig.DashboardFile
    }
)
