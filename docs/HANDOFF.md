# Handoff

Current working branch: `feature/setup-wizard-impl`

Worktree:

```text
C:\Users\jere\.config\superpowers\worktrees\Startup-project-2026\setup-wizard
```

Base worktree:

```text
C:\Users\jere\Documents\Neg_Ai_Stuff\Startup-project-2026
```

Base branch: `feature/setup-wizard`

## Current State

- Work continues in the implementation worktree on `feature/setup-wizard-impl`.
- The branch has not been merged into `feature/setup-wizard` yet.
- The implementation branch has not intentionally been published to GitHub yet.
- Latest confirmed clean commit: `0b5fdf8 feat: add developer credit`.

## Completed In This Follow-Up Round

- Added `ROADMAP.md`.
- Added future roadmap item for packaging Windows setup as a single executable or installer.
- Added a setup-wizard back button:
  - Hidden on step 1.
  - Visible from step 2 onward.
  - Works on the final step before closing/opening the dashboard.
- Added tests for back-button delivery.
- Added developer credit:
  - Landing page footer: `Udviklet af Jesper Reenberg`.
  - Dashboard right column below the local-status notice.
- Rebuilt `dist/GF2-IT-Setup-Windows.zip` after code/content changes.

## Confirmed / Closed Discussion Points

- Point 1 is approved: setup wizard back button.
- Point 2 is clarified: landing page and setup wizard opened quickly because local test commands opened both manually. This is not a product bug.
- Point 3 is approved: discreet developer credit.
- Packaging as a single installer/executable is a future feature, not part of the current fix round.

## Remaining Original Points

4. Add or document a walkthrough for exactly what happens when users click `Download Windows Setup`.
5. Make the landing page `Kom i gang på tre minutter` box about 90% transparent so the image is more visible.
6. Investigate dashboard/setup-wizard status sync. Current state: wizard finishing does not automatically set dashboard status to `Færdig`.
7. Reduce dashboard explanatory text size in the middle column.
8. Fix or redesign the dashboard `Start setup-assistent` link. Current browser behavior can show the `.cmd` file text instead of running it.
9. Keep for later: move computer choice to landing page, dashboard prototypes, and SketchUp tutorial links.

## Verification Commands

Run these sequentially, not in parallel. The build script deletes and recreates the ZIP, while the delivery test reads it.

```powershell
cd C:\Users\jere\.config\superpowers\worktrees\Startup-project-2026\setup-wizard
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

## Git Checks For A New Chat

```powershell
cd C:\Users\jere\.config\superpowers\worktrees\Startup-project-2026\setup-wizard
git status --short --branch
git log --oneline --decorate -8
git worktree list
```

Expected branch:

```text
feature/setup-wizard-impl
```

Expected clean state after the last completed commit:

```text
## feature/setup-wizard-impl
```

## GitHub Learning Notes

- `Changes` in VS Code means files are not committed yet.
- `Commit` saves changes locally on the current branch.
- `Publish Branch` creates the current local branch on GitHub for the first time.
- `Push` sends new commits to GitHub after the branch already exists there.
- `Pull Request` asks GitHub to merge one branch into another.

## Existing Project Docs

- `README.md`
- `ROADMAP.md`
- `docs/superpowers/specs/2026-04-29-gf2-it-dashboard-design.md`
- `docs/superpowers/specs/2026-05-04-gf2-windows-setup-delivery-design.md`
- `docs/superpowers/plans/2026-04-29-gf2-it-dashboard-implementation.md`
- `docs/superpowers/plans/2026-05-04-gf2-windows-setup-wizard.md`
