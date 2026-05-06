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

## Build

Build the Windows delivery package with:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

The package is created at:

```text
dist/GF2-IT-Setup-Windows.zip
```

## Student delivery

1. Build the package.
2. Share `dist/GF2-IT-Setup-Windows.zip` with the student.
3. Ask the student to extract the zip file to a normal folder.
4. Ask the student to open the extracted folder.
5. Ask the student to run `Start Windows setup.cmd`.
