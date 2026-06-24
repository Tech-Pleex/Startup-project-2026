// Tringuidens klient. Al tilstand bor i Go-processen; siden henter
// trinlisten efter hver handling og tegner forfra.

async function postOK(path) {
  const r = await fetch(path, { method: "POST" });
  if (!r.ok) throw new Error(`${path} svarede ${r.status}`);
  return r;
}

const api = {
  steps: () => fetch("/api/steps").then(r => r.json()),
  system: () => fetch("/api/system").then(r => (r.ok ? r.json() : null)),
  done: id => postOK(`/api/steps/${id}/done`),
  undo: id => postOK(`/api/steps/${id}/undo`),
  skip: id => postOK(`/api/steps/${id}/skip`),
  open: id => postOK(`/api/steps/${id}/open`),
  wifiSettings: () => postOK("/api/wifi/settings"),
  quit: () => postOK("/api/quit"),
};

let allSteps = [];
let current = 0;
let sModeBlocked = false;

const $ = id => document.getElementById(id);

async function refresh() {
  allSteps = (await api.steps()).steps;
  render();
}

// Fremdrift mod virkelighed: andel af ikke-finish-trin der er markeret done.
function sceneProgress() {
  const renderable = allSteps.filter(s => s.kind !== "finish" && s.id !== "welcome");
  if (renderable.length === 0) return 0;
  return renderable.filter(s => s.done).length / renderable.length;
}

function renderRail() {
  const rail = $("step-rail");
  rail.innerHTML = "";
  allSteps.forEach((s, i) => {
    if (s.id === "welcome" || s.kind === "finish") return;
    const tile = document.createElement("button");
    let cls = "tile";
    if (s.done) cls += " rendered";
    else if (s.skipped) cls += " skipped";
    if (i === current) cls += " is-current";
    tile.className = cls;
    tile.disabled = sModeBlocked;
    tile.innerHTML = `<span class="ic"></span>${s.title}`;
    tile.addEventListener("click", () => { current = i; render(); });
    rail.appendChild(tile);
  });
}

function render() {
  if (allSteps.length === 0) {
    $("progress").textContent = "Ingen trin at vise.";
    return;
  }
  const step = allSteps[current];
  const progress = sceneProgress();

  $("progress").textContent = `Trin ${current + 1} / ${allSteps.length}`;
  $("progress-bar").style.width = `${Math.round(progress * 100)}%`;
  $("scene-blueprint").style.opacity = String(1 - progress);

  renderRail();

  $("step-kicker").textContent = `Trin ${current + 1} / ${allSteps.length} · ${step.kind}`;
  $("step-title").textContent = step.title;
  $("step-body").textContent = step.body;

  const welcomeBlocked = sModeBlocked && current === 0;
  $("smode-warning").hidden = !welcomeBlocked;

  const warning = $("step-warning");
  warning.hidden = !step.warning;
  warning.textContent = step.warning || "";

  $("wifi-panel").hidden = step.kind !== "wifi";
  $("quit").hidden = step.kind !== "finish";

  const action = $("action");
  const hasAction = step.kind === "link" || step.kind === "finish";
  action.hidden = !hasAction;
  action.textContent = step.button;

  const isWelcome = step.id === "welcome";

  const toggleDone = $("toggle-done");
  toggleDone.hidden = isWelcome;
  toggleDone.textContent = step.done ? "Fortryd" : "Markér som færdig";
  toggleDone.disabled = sModeBlocked;

  const skip = $("skip");
  skip.hidden = step.kind === "finish" || isWelcome;
  skip.textContent = step.skipped ? "Fortryd spring over" : "Spring over";
  skip.disabled = sModeBlocked;

  $("prev").disabled = current === 0;
  $("next").disabled = current === allSteps.length - 1 || sModeBlocked;
}

async function checkSystem() {
  let status = null;
  try { status = await api.system(); } catch (err) { /* ukendt status må ikke blokere */ }
  sModeBlocked = Boolean(status?.sMode);
  if (sModeBlocked) current = 0;
  render();
}

async function start() {
  allSteps = (await api.steps()).steps;
  await checkSystem();
}

$("toggle-done").addEventListener("click", async () => {
  const step = allSteps[current];
  try {
    await (step.done ? api.undo(step.id) : api.done(step.id));
    await refresh();
  } catch (err) { alert("Kunne ikke opdatere trinnet. Prøv igen."); }
});

$("skip").addEventListener("click", async () => {
  const step = allSteps[current];
  const wasSkipped = step.skipped;
  try {
    await (wasSkipped ? api.undo(step.id) : api.skip(step.id));
    await refresh();
    if (!wasSkipped && current < allSteps.length - 1) { current++; render(); }
  } catch (err) { alert("Kunne ikke springe trinnet over. Prøv igen."); }
});

$("prev").addEventListener("click", () => { if (current > 0) { current--; render(); } });
$("next").addEventListener("click", () => { if (current < allSteps.length - 1) { current++; render(); } });

document.addEventListener("keydown", e => {
  if (e.key === "ArrowLeft") $("prev").click();
  else if (e.key === "ArrowRight") $("next").click();
});

$("action").addEventListener("click", async () => {
  const step = allSteps[current];
  try { await api.open(step.id); }
  catch (err) { alert("Siden kunne ikke åbnes. Prøv at åbne den manuelt i din browser."); }
});

$("smode-retry").addEventListener("click", checkSystem);

$("wifi-settings").addEventListener("click", async () => {
  try { await api.wifiSettings(); }
  catch (err) { alert("Wi-Fi-indstillingerne kunne ikke åbnes. Åbn dem selv via proceslinjen."); }
});

$("quit").addEventListener("click", async () => {
  try { await api.quit(); } catch (err) { /* sandsynligvis allerede lukket */ }
  document.body.innerHTML = "<p style='padding:40px;text-align:center'>Assistenten er lukket. Du kan lukke denne fane.</p>";
});

start();
