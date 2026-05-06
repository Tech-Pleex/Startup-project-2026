# Landing Page Platform Download Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move platform choice from the dashboard into a landing page download modal with Windows enabled and Mac visible but disabled.

**Architecture:** Keep the site static. `index.html` owns platform choice and download behavior; `start.html` remains a local status/dashboard page with no platform state. Existing PowerShell tests verify static HTML contracts and package contents.

**Tech Stack:** Static HTML/CSS/JavaScript, PowerShell test scripts, existing package build script.

---

## File Structure

- Modify: `tests/check-setup-delivery.ps1` - add landing page modal assertions.
- Modify: `index.html` - replace direct download link with modal button, Windows download choice, disabled Mac choice, close behavior.
- Modify: `tests/check-dashboard.ps1` - assert dashboard no longer contains platform selection or launcher help.
- Modify: `start.html` - remove platform panel, platform CSS, platform storage key, and platform JavaScript; replace left panel with neutral setup/help navigation.
- Modify: `docs/HANDOFF.md` - record completion of the platform-modal change.
- Rebuild: `dist/GF2-IT-Setup-Windows.zip` - include updated static files.

## Task 1: Add Landing Page Modal Test Coverage

**Files:**
- Modify: `tests/check-setup-delivery.ps1`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Write the failing landing page assertions**

In `tests/check-setup-delivery.ps1`, replace:

```powershell
Assert-Contains $LandingHtml "Download Windows Setup" "download call to action"
Assert-Contains $LandingHtml "dist/GF2-IT-Setup-Windows.zip" "download package path"
```

with:

```powershell
Assert-Contains $LandingHtml "Download setup" "download call to action"
Assert-Contains $LandingHtml "downloadModal" "platform modal"
Assert-Contains $LandingHtml "Vælg computer" "platform modal title"
Assert-Contains $LandingHtml "Windows" "Windows platform choice"
Assert-Contains $LandingHtml "dist/GF2-IT-Setup-Windows.zip" "Windows download package path"
Assert-Contains $LandingHtml "Mac" "Mac platform choice"
Assert-Contains $LandingHtml "disabled" "disabled Mac platform choice"
Assert-Contains $LandingHtml "openDownloadModal" "open modal handler"
Assert-Contains $LandingHtml "closeDownloadModal" "close modal handler"
Assert-Contains $LandingHtml "Escape" "keyboard modal close"
```

- [ ] **Step 2: Run delivery test to verify it fails**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: FAIL with missing content such as `platform modal`.

- [ ] **Step 3: Commit the failing test**

Run:

```powershell
git add tests/check-setup-delivery.ps1
git commit -m "test: require landing page platform modal"
```

## Task 2: Implement Landing Page Download Modal

**Files:**
- Modify: `index.html`
- Test: `tests/check-setup-delivery.ps1`

- [ ] **Step 1: Update button and add modal markup**

In `index.html`, replace the current primary download link:

```html
<a class="primary" href="dist/GF2-IT-Setup-Windows.zip" download>Download Windows Setup</a>
```

with:

```html
<button class="primary" type="button" onclick="openDownloadModal()">Download setup</button>
```

Then add this modal markup immediately before `</body>`:

```html
<div id="downloadModal" class="modal-backdrop" role="presentation" aria-hidden="true" onclick="handleModalBackdropClick(event)">
  <section class="download-modal" role="dialog" aria-modal="true" aria-labelledby="downloadModalTitle">
    <button class="modal-close" type="button" aria-label="Luk" onclick="closeDownloadModal()">×</button>
    <h2 id="downloadModalTitle">Vælg computer</h2>
    <div class="platform-choice">
      <a class="platform-download" href="dist/GF2-IT-Setup-Windows.zip" download>Windows</a>
      <button class="platform-download" type="button" disabled>Mac</button>
    </div>
  </section>
</div>
```

- [ ] **Step 2: Add modal CSS**

In `index.html`, add this CSS before the first `@media` block:

```css
.modal-backdrop {
  position: fixed;
  inset: 0;
  display: none;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(23, 32, 42, .52);
  z-index: 20;
}

.modal-backdrop.is-open {
  display: flex;
}

.download-modal {
  position: relative;
  width: min(420px, 100%);
  padding: 28px;
  border-radius: 8px;
  background: var(--surface);
  box-shadow: 0 26px 70px rgba(23, 32, 42, .30);
}

.download-modal h2 {
  margin: 0 0 20px;
  font-size: 1.45rem;
  line-height: 1.15;
}

.modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border: 1px solid var(--line);
  border-radius: 4px;
  background: #fff;
  color: var(--ink);
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
}

.platform-choice {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.platform-download {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 58px;
  border: 0;
  border-radius: 6px;
  background: var(--blue);
  color: #fff;
  font: inherit;
  font-weight: 900;
  text-decoration: none;
  cursor: pointer;
}

.platform-download:hover {
  background: var(--blue-dark);
}

.platform-download:disabled {
  background: #d7e0e5;
  color: #77838d;
  cursor: not-allowed;
}
```

- [ ] **Step 3: Add modal JavaScript**

In `index.html`, add this script block immediately before `</body>`, after the modal markup:

```html
<script>
  const downloadModal = document.getElementById("downloadModal");

  function openDownloadModal() {
    downloadModal.classList.add("is-open");
    downloadModal.setAttribute("aria-hidden", "false");
  }

  function closeDownloadModal() {
    downloadModal.classList.remove("is-open");
    downloadModal.setAttribute("aria-hidden", "true");
  }

  function handleModalBackdropClick(event) {
    if (event.target === downloadModal) {
      closeDownloadModal();
    }
  }

  document.addEventListener("keydown", event => {
    if (event.key === "Escape") {
      closeDownloadModal();
    }
  });
</script>
```

- [ ] **Step 4: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: PASS with `Setup delivery checks passed.`

- [ ] **Step 5: Commit landing page modal**

Run:

```powershell
git add index.html
git commit -m "feat: add landing page platform download modal"
```

## Task 3: Add Dashboard Platform Removal Test Coverage

**Files:**
- Modify: `tests/check-dashboard.ps1`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Replace old dashboard launcher assertions**

In `tests/check-dashboard.ps1`, replace:

```powershell
Assert-Contains $Html "setupLauncherHelp" "setup launcher help target"
Assert-NotContains $Html 'href="Start Windows setup.cmd"' "browser must not link directly to Windows CMD"
Assert-NotContains $Html 'href="Start Mac setup.command"' "browser must not link directly to Mac command"
```

with:

```powershell
Assert-Contains $Html "Download setup" "dashboard download page link"
Assert-Contains $Html 'href="index.html"' "dashboard links back to landing page"
Assert-NotContains $Html "Vælg computer" "dashboard platform chooser removed"
Assert-NotContains $Html 'data-platform="windows"' "Windows platform button removed"
Assert-NotContains $Html 'data-platform="mac"' "Mac platform button removed"
Assert-NotContains $Html "gf2-it-dashboard.platform" "platform localStorage key removed"
Assert-NotContains $Html "function setPlatform" "platform JavaScript removed"
Assert-NotContains $Html "setupLauncherHelp" "dashboard launcher help removed"
Assert-NotContains $Html 'href="Start Windows setup.cmd"' "browser must not link directly to Windows CMD"
Assert-NotContains $Html 'href="Start Mac setup.command"' "browser must not link directly to Mac command"
```

- [ ] **Step 2: Run dashboard test to verify it fails**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: FAIL with forbidden content such as `Vælg computer` or `gf2-it-dashboard.platform`.

- [ ] **Step 3: Commit the failing test**

Run:

```powershell
git add tests/check-dashboard.ps1
git commit -m "test: require dashboard platform chooser removal"
```

## Task 4: Remove Dashboard Platform Choice

**Files:**
- Modify: `start.html`
- Test: `tests/check-dashboard.ps1`

- [ ] **Step 1: Remove platform CSS**

In `start.html`, delete the full `.platforms`, `.platform`, and `.platform[aria-pressed="true"]` CSS blocks.

- [ ] **Step 2: Replace the left dashboard panel**

In `start.html`, replace the full first `<aside class="panel">...</aside>` inside `<div class="dashboard">` with:

```html
<aside class="panel">
  <div class="panel-header">Kom i gang</div>
  <div class="panel-body">
    <a class="primary-action" href="index.html">Download setup</a>
    <a class="secondary-action" href="#manualGuide">Manuel vejledning</a>
    <div class="notice">Chromebook kan ikke bruges til dette GF2-flow. Windows må ikke være i S-mode, hvis programmer skal installeres.</div>
  </div>
</aside>
```

- [ ] **Step 3: Update hero copy**

In `start.html`, replace:

```html
<p>Vælg Windows eller Mac, følg rækkefølgen første skoledag, og brug dashboardet igen når du skal finde Office, Moodle, Lectio, PraxisOnline, OneDrive eller SketchUp.</p>
```

with:

```html
<p>Følg rækkefølgen første skoledag, og brug dashboardet igen når du skal finde Office, Moodle, Lectio, PraxisOnline, OneDrive eller SketchUp.</p>
```

- [ ] **Step 4: Remove platform storage and JavaScript**

In `start.html`, replace the `storage` object:

```javascript
const storage = {
  platform: "gf2-it-dashboard.platform",
  mode: "gf2-it-dashboard.mode",
  student: "gf2-it-dashboard.studentStatus",
  teacher: "gf2-it-dashboard.teacherStatus"
};
```

with:

```javascript
const storage = {
  mode: "gf2-it-dashboard.mode",
  student: "gf2-it-dashboard.studentStatus",
  teacher: "gf2-it-dashboard.teacherStatus"
};
```

Delete the full `setPlatform(platform)` function and delete this event binding block:

```javascript
document.querySelectorAll("[data-platform]").forEach(button => {
  button.addEventListener("click", () => setPlatform(button.dataset.platform));
});
```

Delete this initialization line:

```javascript
setPlatform(localStorage.getItem(storage.platform) || "windows");
```

- [ ] **Step 5: Run dashboard test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 6: Commit dashboard cleanup**

Run:

```powershell
git add start.html
git commit -m "fix: remove platform choice from dashboard"
```

## Task 5: Update Handoff, Rebuild Package, And Verify

**Files:**
- Modify: `docs/HANDOFF.md`
- Rebuild: `dist/GF2-IT-Setup-Windows.zip`
- Test: all verification commands

- [ ] **Step 1: Update handoff**

In `docs/HANDOFF.md`, add this bullet under `Completed In This Follow-Up Round`:

```markdown
- Moved platform choice to the landing page download flow:
  - `Download setup` now opens a platform modal.
  - Windows downloads `dist/GF2-IT-Setup-Windows.zip`.
  - Mac is visible but disabled until a later Mac package feature.
  - Dashboard no longer contains computer choice or platform launcher logic.
```

- [ ] **Step 2: Rebuild package**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
```

Expected: PASS with `Created ...\dist\GF2-IT-Setup-Windows.zip`.

- [ ] **Step 3: Run dashboard test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
```

Expected: PASS with `Dashboard checks passed.`

- [ ] **Step 4: Run delivery test**

Run:

```powershell
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
```

Expected: PASS with `Setup delivery checks passed.`

- [ ] **Step 5: Commit final package and handoff**

Run:

```powershell
git add docs/HANDOFF.md dist/GF2-IT-Setup-Windows.zip
git commit -m "docs: update handoff for platform download flow"
```

## Final Verification

Run these sequentially:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-package.ps1
powershell -ExecutionPolicy Bypass -File tests/check-dashboard.ps1
powershell -ExecutionPolicy Bypass -File tests/check-setup-delivery.ps1
git status --short --branch
```

Expected:

```text
Created ...\dist\GF2-IT-Setup-Windows.zip
Dashboard checks passed.
Setup delivery checks passed.
## feature/setup-wizard-impl
```

## Self-Review

- Spec coverage: landing modal, Windows enabled download, Mac disabled visible button, dashboard platform removal, and package rebuild are covered.
- Scope: no Mac zip, no installer, no GitHub publishing, no full dashboard redesign.
- Test strategy: delivery test covers landing behavior; dashboard test covers platform removal and preserves existing status behavior.
- Type consistency: modal IDs/functions are `downloadModal`, `openDownloadModal`, and `closeDownloadModal` throughout.
