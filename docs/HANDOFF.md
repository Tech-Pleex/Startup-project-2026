# Handoff

Current main project worktree:

```text
C:\Users\jere\Documents\Neg_Ai_Stuff\Startup-project-2026
```

Current branch:

```text
master
```

Implementation worktree still exists:

```text
C:\Users\jere\.config\superpowers\worktrees\Startup-project-2026\setup-wizard
```

## Current State

- GitHub PR `#1` titled `Add GF2 setup wizard and platform download flow` was merged into `master`.
- Local `master` was updated with `git pull origin master`.
- The user ran a quick manual test and said it looked like it worked.
- The user is learning Git and wants to make commits themselves. Do not commit, push, merge, or open PRs automatically unless explicitly asked.
- Explain Git steps one at a time.

## Local Git Status At Handoff

Last observed in the main project worktree:

```text
## master...origin/master
 M dist/GF2-IT-Setup-Windows.zip
```

Meaning:

- `master` is tracking `origin/master`.
- `dist/GF2-IT-Setup-Windows.zip` was modified locally after running the build script.
- This zip change has not been committed.
- In the next chat, start by checking whether the user wants to keep and commit the rebuilt zip, or restore it if it is identical/unwanted.

## Completed Work Now On Master

- Added Windows setup wizard delivery flow.
- Added `index.html` landing page.
- Added Windows setup package at `dist/GF2-IT-Setup-Windows.zip`.
- Added setup wizard config/check scripts:
  - `scripts/setup-config.ps1`
  - `scripts/setup-checks.ps1`
  - `scripts/setup-windows.ps1`
- Added setup delivery tests:
  - `tests/check-setup-delivery.ps1`
  - updated `tests/check-dashboard.ps1`
- Added `ROADMAP.md`.
- Added developer credit:
  - Landing page footer: `Udviklet af Jesper Reenberg`
  - Dashboard right column below local-status notice.
- Added setup-wizard back button:
  - Hidden on step 1.
  - Visible from step 2 onward.
  - Works on final step before closing/opening dashboard.
- Fixed dashboard/setup-wizard status sync:
  - Setup wizard opens `start.html?setup=complete` on finish.
  - Dashboard marks student checklist `Færdig` in local browser storage.
  - Dashboard removes completion query from visible URL after applying it.
- Reduced dashboard middle-column explanatory text size.
- Moved platform choice to landing page:
  - `Download setup` opens platform modal.
  - `Windows` downloads `dist/GF2-IT-Setup-Windows.zip`.
  - `Mac` is visible but disabled until a later Mac-package feature.
  - Dashboard no longer contains computer choice or platform launcher logic.

## Git Learning Notes For User

- `Commit` saves local changes on the current branch.
- `Publish Branch` in VS Code creates a local branch on GitHub for the first time and pushes it.
- `Push` sends local commits to GitHub after the branch already exists there.
- `Pull Request` asks GitHub to merge one branch into another.
- `Merge pull request` on GitHub applies the branch changes to the base branch.
- After a GitHub merge, local `master` needs `git pull origin master`.
- If merge conflicts happen, resolve files first, then `git add <file>` marks them resolved, then `git commit` completes the merge.

## Verification To Run In Next Chat

Run from:

```powershell
cd C:\Users\jere\Documents\Neg_Ai_Stuff\Startup-project-2026
```

Start with:

```powershell
git status --short --branch
```

Then run verification sequentially:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

If build modifies `dist/GF2-IT-Setup-Windows.zip`, ask the user whether they want to commit that rebuilt zip themselves.

## Remaining Later Work

- Dashboard prototypes.
- SketchUp tutorial links.
- Mac package feature.
- Future installer or single executable for Windows setup.

## Existing Project Docs

- `README.md`
- `ROADMAP.md`
- `docs/superpowers/specs/2026-04-29-gf2-it-dashboard-design.md`
- `docs/superpowers/specs/2026-05-04-gf2-windows-setup-delivery-design.md`
- `docs/superpowers/specs/2026-05-06-landing-page-platform-download-design.md`
- `docs/superpowers/plans/2026-04-29-gf2-it-dashboard-implementation.md`
- `docs/superpowers/plans/2026-05-04-gf2-windows-setup-wizard.md`
- `docs/superpowers/plans/2026-05-06-landing-page-platform-download.md`
