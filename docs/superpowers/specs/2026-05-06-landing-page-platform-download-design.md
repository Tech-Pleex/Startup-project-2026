# Landing Page Platform Download Design

## Goal

Move the computer choice out of the dashboard and into the landing page download flow.

The landing page should be the only place where a student chooses platform. The dashboard should focus on status, fixed links, and manual guidance after setup.

## User Flow

1. Student opens the landing page.
2. Student clicks `Download setup`.
3. A modal opens with the title `Vælg computer`.
4. Student chooses one of two visible platform buttons:
   - `Windows`
   - `Mac`
5. `Windows` downloads `dist/GF2-IT-Setup-Windows.zip`.
6. `Mac` is visible but disabled until a later Mac package feature exists.
7. Dashboard remains available as a guide/status page after setup.

## Landing Page Behavior

- The primary download control changes from a direct zip link to a button.
- The button opens a modal instead of downloading immediately.
- The modal contains only the platform choice and close behavior.
- The Windows choice is the only enabled download action.
- The Mac choice is visible and disabled.
- No explanatory body text is shown under the Windows or Mac choices.
- The modal must be keyboard-closeable with `Escape` and by a visible close button.
- Clicking outside the modal should close it.

## Dashboard Behavior

- Remove the dashboard computer-choice UI.
- Remove the Windows/Mac platform toggle behavior from dashboard JavaScript.
- Remove direct dashboard launcher behavior for `.cmd` and `.command` files.
- Keep the dashboard's student/teacher mode toggle.
- Keep local status controls, fixed links, manual guide, absence notice, local-storage safety notice, and developer credit.
- The left dashboard panel should become neutral setup/help navigation, not platform selection.

## Package Scope

- Existing Windows package stays at `dist/GF2-IT-Setup-Windows.zip`.
- No Mac package is created in this feature.
- The Mac button is disabled rather than hidden, so the future Mac feature has an obvious destination.

## Testing

Update the existing PowerShell tests to cover:

- Landing page includes a platform modal opened from `Download setup`.
- Windows choice points to `dist/GF2-IT-Setup-Windows.zip`.
- Mac choice exists and is disabled.
- Landing page still includes safety wording, dashboard link, GitHub link, hero asset, and developer credit.
- Dashboard no longer contains `Vælg computer`.
- Dashboard no longer contains the Windows/Mac platform buttons.
- Dashboard no longer contains direct `.cmd` or `.command` launcher links.
- Existing dashboard status and setup completion behavior still exists.

## Out Of Scope

- Creating a Mac zip package.
- Building a signed installer or single executable.
- Publishing or merging the implementation branch.
- Redesigning the whole dashboard layout beyond removing platform selection and preserving a useful left panel.
