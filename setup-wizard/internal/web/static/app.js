// Tringuidens klient. Al tilstand bor i Go-processen; siden henter
// trinlisten efter hver handling og tegner forfra.
const api = {
  steps: () => fetch("/api/steps").then(r => r.json()),
  setDone: (id, done) => fetch(`/api/steps/${id}/${done ? "done" : "undo"}`, { method: "POST" }),
  open: id => fetch(`/api/steps/${id}/open`, { method: "POST" }),
  wifi: () => fetch("/api/wifi").then(r => (r.ok ? r.json() : null)),
  wifiSettings: () => fetch("/api/wifi/settings", { method: "POST" }),
  sketchup: () => fetch("/api/sketchup/install", { method: "POST" }).then(r => r.json()),
  quit: () => fetch("/api/quit", { method: "POST" }),
};

let allSteps = [];
let current = 0;

async function refresh() {
  allSteps = (await api.steps()).steps;
  render();
}

function render() {
  const step = allSteps[current];
  document.getElementById("progress").textContent = `Trin ${current + 1} af ${allSteps.length}`;

  const list = document.getElementById("step-list");
  list.innerHTML = "";
  allSteps.forEach((s, i) => {
    const item = document.createElement("button");
    item.className = "step-item" + (i === current ? " active" : "");
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
  await api.setDone(step.id, !step.done);
  await refresh();
});

document.getElementById("prev").addEventListener("click", () => { current--; render(); });
document.getElementById("next").addEventListener("click", () => { current++; render(); });

document.getElementById("action").addEventListener("click", async () => {
  const step = allSteps[current];
  if (step.kind === "sketchup") {
    const result = document.getElementById("sketchup-result");
    result.hidden = false;
    result.textContent = "Installerer … det kan tage nogle minutter.";
    const outcome = await api.sketchup();
    result.textContent = outcome.action === "installed"
      ? "SketchUp er installeret."
      : outcome.reason;
  } else {
    await api.open(step.id);
  }
});

document.getElementById("wifi-check").addEventListener("click", checkWifi);
document.getElementById("wifi-settings").addEventListener("click", () => api.wifiSettings());

document.getElementById("quit").addEventListener("click", async () => {
  await api.quit();
  document.body.innerHTML = "<p style='padding:40px;text-align:center'>Assistenten er lukket. Du kan lukke denne fane.</p>";
});

refresh();
