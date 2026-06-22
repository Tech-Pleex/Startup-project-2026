// Tringuidens klient. Al tilstand bor i Go-processen; siden henter
// trinlisten efter hver handling og tegner forfra.

// postOK sender et POST-kald og fejler hvis serveren svarer med en fejlkode,
// så eleven aldrig efterlades i en misvisende tilstand uden besked.
async function postOK(path) {
  const r = await fetch(path, { method: "POST" });
  if (!r.ok) {
    throw new Error(`${path} svarede ${r.status}`);
  }
  return r;
}

const api = {
  steps: () => fetch("/api/steps").then(r => r.json()),
  system: () => fetch("/api/system").then(r => (r.ok ? r.json() : null)),
  setDone: (id, done) => postOK(`/api/steps/${id}/${done ? "done" : "undo"}`),
  open: id => postOK(`/api/steps/${id}/open`),
  wifiSettings: () => postOK("/api/wifi/settings"),
  quit: () => postOK("/api/quit"),
};

let allSteps = [];
let current = 0;
let sModeBlocked = false;

async function refresh() {
  allSteps = (await api.steps()).steps;
  render();
}

function render() {
  if (allSteps.length === 0) {
    document.getElementById("progress").textContent = "Ingen trin at vise.";
    return;
  }
  const step = allSteps[current];
  document.getElementById("progress").textContent = `Trin ${current + 1} af ${allSteps.length}`;

  const list = document.getElementById("step-list");
  list.innerHTML = "";
  allSteps.forEach((s, i) => {
    const item = document.createElement("button");
    item.className = "step-item" + (i === current ? " active" : "");
    if (i === current) {
      item.setAttribute("aria-current", "step");
    }
    item.disabled = sModeBlocked && i > 0;
    item.textContent = `${i + 1}. ${s.title}` + (s.done ? " ✓" : "");
    item.addEventListener("click", () => { current = i; render(); });
    list.appendChild(item);
  });

  document.getElementById("step-title").textContent = step.title;
  document.getElementById("step-body").textContent = step.body;

  const welcomeBlocked = sModeBlocked && current === 0;
  document.getElementById("smode-warning").hidden = !welcomeBlocked;

  const warning = document.getElementById("step-warning");
  warning.hidden = !step.warning;
  warning.textContent = step.warning || "";

  document.getElementById("wifi-panel").hidden = step.kind !== "wifi";
  document.getElementById("quit").hidden = step.kind !== "finish";

  const action = document.getElementById("action");
  const hasAction = step.kind === "link" || step.kind === "finish";
  action.hidden = !hasAction;
  action.textContent = step.button;

  const toggleDone = document.getElementById("toggle-done");
  toggleDone.textContent = step.done ? "Fortryd" : "Markér som færdig";
  toggleDone.disabled = sModeBlocked;

  document.getElementById("prev").disabled = current === 0;
  document.getElementById("next").disabled = current === allSteps.length - 1 || sModeBlocked;
}

async function checkSystem() {
  let status = null;
  try {
    status = await api.system();
  } catch (err) {
    // En ukendt status må ikke blokere eleven.
  }
  sModeBlocked = Boolean(status?.sMode);
  if (sModeBlocked) {
    current = 0;
  }
  render();
}

async function start() {
  allSteps = (await api.steps()).steps;
  await checkSystem();
}

document.getElementById("toggle-done").addEventListener("click", async () => {
  const step = allSteps[current];
  try {
    await api.setDone(step.id, !step.done);
    await refresh();
  } catch (err) {
    alert("Kunne ikke opdatere trinnet. Prøv igen.");
  }
});

document.getElementById("prev").addEventListener("click", () => { current--; render(); });
document.getElementById("next").addEventListener("click", () => { current++; render(); });

document.getElementById("action").addEventListener("click", async () => {
  const step = allSteps[current];
  try {
    await api.open(step.id);
  } catch (err) {
    alert("Siden kunne ikke åbnes. Prøv at åbne den manuelt i din browser.");
  }
});

document.getElementById("smode-retry").addEventListener("click", checkSystem);
document.getElementById("wifi-settings").addEventListener("click", async () => {
  try {
    await api.wifiSettings();
  } catch (err) {
    alert("Wi-Fi-indstillingerne kunne ikke åbnes. Åbn dem selv via proceslinjen.");
  }
});

document.getElementById("quit").addEventListener("click", async () => {
  try {
    await api.quit();
  } catch (err) {
    // Assistenten er sandsynligvis allerede lukket; vis afslutningsbeskeden uanset.
  }
  document.body.innerHTML = "<p style='padding:40px;text-align:center'>Assistenten er lukket. Du kan lukke denne fane.</p>";
});

start();
