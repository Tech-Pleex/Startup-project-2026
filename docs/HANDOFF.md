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
- Latest confirmed clean commit: `8f1f779 style: refine landing page quick steps`.

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
- Refined the landing page quick-start area:
  - `Kom i gang på tre minutter` is now 90% transparent.
  - The quick-start box no longer stretches too far to the right.
  - The landing page lead text now matches the quick-start body text size.
- Fixed dashboard/setup-wizard status sync:
  - Setup wizard now opens `start.html?setup=complete` on finish.
  - Dashboard reads the completion signal and marks the student checklist `Færdig` in local browser storage.
  - Dashboard removes the completion query from the visible URL after applying it.
- Reduced dashboard middle-column explanatory text size.
- Redesigned the dashboard setup-assistant action so the browser no longer links directly to `.cmd` or `.command` files.
- Moved platform choice to the landing page download flow:
  - `Download setup` now opens a platform modal.
  - Windows downloads `dist/GF2-IT-Setup-Windows.zip`.
  - Mac is visible but disabled until a later Mac package feature.
  - Dashboard no longer contains computer choice or platform launcher logic.
- Rebuilt `dist/GF2-IT-Setup-Windows.zip` after code/content changes.

## Confirmed / Closed Discussion Points

- Point 1 is approved: setup wizard back button.
- Point 2 is clarified: landing page and setup wizard opened quickly because local test commands opened both manually. This is not a product bug.
- Point 3 is approved: discreet developer credit.
- Point 4 was explored and reverted. The extra `Efter download` block was removed because it duplicated the existing quick-start box.
- Point 5 is approved: landing page quick-start transparency, box width, and lead text size.
- Packaging as a single installer/executable is a future feature, not part of the current fix round.

## Remaining Original Points

9. Keep for later: dashboard prototypes and SketchUp tutorial links.

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
