# GF2 IT Setup

GF2 IT Setup is a Windows-first helper for new GF2 students. It helps students find the right setup flow, run the Windows setup helper, and open local setup guidance for GF2 IT onboarding.

## Safety

The assistant never asks for passwords, MitID, or UNI-Login.

Students only enter login information on official pages or Windows' own settings. If a page or prompt asks for login details, students should confirm that it is an official service before continuing.

## Current status

This is a prototype and the first version is Windows first.

Not in the first version:

- Mac setup
- Signed installer
- Central status
- Intune deployment
- Official NEG hosting

## Local checks

Run these commands from the repository root:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Run the Go tests from the Go module:

```powershell
cd setup-wizard
go test ./...
```

## Release build

The landing page downloads the latest approved Assistenten binaries from
GitHub Releases. A release should only be marked as `latest` when all three
assets are built and uploaded:

- `Assistenten-Windows.exe`
- `Assistenten-Mac-Apple-Silicon`
- `Assistenten-Mac-Intel`

Build the release binaries from `setup-wizard`:

```powershell
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -ldflags "-H=windowsgui" -o ..\dist\Assistenten-Windows.exe .\cmd\assistent
$env:GOOS = "darwin"; $env:GOARCH = "arm64"; go build -o ..\dist\Assistenten-Mac-Apple-Silicon .\cmd\assistent
$env:GOOS = "darwin"; $env:GOARCH = "amd64"; go build -o ..\dist\Assistenten-Mac-Intel .\cmd\assistent
```

Release checklist:

1. Run `go test ./...` from `setup-wizard`.
2. Build the three Assistenten binaries.
3. Upload all three files as assets on the GitHub Release.
4. Mark the release as latest only after all three assets are present.
5. Run `powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1`.

## Student delivery

1. Ask the student to open the landing page.
2. Ask the student to choose Windows or Mac.
3. Mac students choose Apple Silicon for M1/M2/M3/M4, or Intel for older Macs.
4. Ask the student to open the downloaded Assistenten file.
5. The Assistenten starts a local browser guide and never asks for passwords,
   MitID, or UNI-Login.

## License

The source code in this repository is released under the **MIT License** (see
[LICENSE](LICENSE)). You are free to use, modify, and distribute the code,
including for NEG to take over and maintain the project at any time.

## Brand and trademarks

The MIT License covers the **source code only** — not brands or trademarks.

- **NEG** (Nordvestsjællands Erhvervs- og Gymnasieuddannelser): the NEG name,
  logo, brand guide, colors, and fonts are the property of NEG and are used in
  this project **with NEG's permission**. They are not part of the MIT license
  grant and may not be reused without NEG's consent.
- **Trimble / SketchUp**, **Microsoft / Windows / Microsoft Store**, and
  **Apple / App Store / macOS** are trademarks of their respective owners.

This project is a student onboarding helper and is **not** an official product
of, or endorsed by, any of the above parties.
