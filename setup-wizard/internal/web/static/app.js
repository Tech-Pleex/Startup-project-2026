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
  setDone: (id, done) => postOK(`/api/steps/${id}/${done ? "done" : "undo"}`),
  open: id => postOK(`/api/steps/${id}/open`),
  wifi: () => fetch("/api/wifi").then(r => (r.ok ? r.json() : null)),
  wifiSettings: () => postOK("/api/wifi/settings"),
  sketchup: () => postOK("/api/sketchup/install").then(r => r.json()),
  quit: () => postOK("/api/quit"),
};

let allSteps = [];
let current = 0;

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
    item.textContent = `${i + 1}. ${s.title}` + (s.done ? " ✓" : "");
    item.addEventListener("click", () => { current = i; render(); });
    list.appendChild(item);
  });

  document.getElementById("step-title").textContent = step.title;
  document.getElementById("step-body").textContent = step.body;

  const warning = document.getElementById("step-warning");
  warning.hidden = !step.warning;
  warning.textContent = step.warning || "";

  document.getElementById("wifi-panel").hidden = step.kind !== "wifi";
  document.getElementById("sketchup-result").hidden = true;
  document.getElementById("quit").hidden = step.kind !== "finish";

  const action = document.getElementById("action");
  const hasAction = step.kind === "link" || step.kind === "sketchup" || step.kind === "finish";
  action.hidden = !hasAction;
  action.textContent = step.button;

  document.getElementById("toggle-done").textContent = step.done ? "Fortryd" : "Markér som færdig";

  document.getElementById("prev").disabled = current === 0;
  document.getElementById("next").disabled = current === allSteps.length - 1;

  if (step.kind === "wifi") {
    checkWifi();
  }
}

async function checkWifi() {
  const el = document.getElementById("wifi-status");
  const status = await api.wifi();
  if (!status) {
    el.textContent = "Wi-Fi-status kunne ikke aflæses. Tjek selv i dine netværksindstillinger.";
    return;
  }
  const texts = {
    target: `Du er på ${status.ssid} — det rigtige netværk. Trinnet er klaret!`,
    guest: `Du er på ${status.ssid}. Det er kun til midlertidig gæsteadgang — skift til NEG.`,
    other: `Du er på "${status.ssid}", ikke NEG. Åbn Wi-Fi-indstillinger og skift til NEG.`,
    none: "Du er ikke på et Wi-Fi-netværk. Åbn Wi-Fi-indstillinger og vælg NEG.",
  };
  el.textContent = texts[status.state];
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
  if (step.kind === "sketchup") {
    const result = document.getElementById("sketchup-result");
    result.hidden = false;
    result.textContent = "Installerer … det kan tage nogle minutter.";
    try {
      const outcome = await api.sketchup();
      result.textContent = outcome.action === "installed"
        ? "SketchUp er installeret."
        : outcome.reason;
    } catch (err) {
      result.textContent = "Installationen kunne ikke startes. Åbn SketchUp-siden manuelt med knappen ovenfor.";
    }
  } else {
    try {
      await api.open(step.id);
    } catch (err) {
      alert("Siden kunne ikke åbnes. Prøv at åbne den manuelt i din browser.");
    }
  }
});

document.getElementById("wifi-check").addEventListener("click", checkWifi);
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

refresh();
