async function loadHealth() {
  const h = await fetch("/api/health").then(r => r.json());
  document.getElementById("health").textContent =
    "node=" + h.NodeID + " ok=" + h.OK;
}

async function loadLayers() {
  const layers = await fetch("/api/layers").then(r => r.json());
  const ol = document.getElementById("layers");
  ol.innerHTML = "";
  layers.forEach(l => {
    const li = document.createElement("li");
    li.textContent = l.ID + " prio=" + l.Priority + " keys=" +
      Object.keys(l.Values || {}).length + " src=" + (l.Source || "");
    ol.appendChild(li);
  });
}

async function mergePreview() {
  const merged = await fetch("/api/merge").then(r => r.json());
  document.getElementById("merged").textContent =
    JSON.stringify(merged, null, 2);
}

async function explainKey() {
  const key = document.getElementById("keyInput").value.trim();
  if (!key) return;
  const steps = await fetch("/api/explain/" + encodeURIComponent(key))
    .then(r => r.json());
  document.getElementById("explainOut").textContent =
    JSON.stringify(steps, null, 2);
}

document.getElementById("refresh").onclick = async () => {
  await loadHealth();
  await loadLayers();
};
document.getElementById("merge").onclick = mergePreview;
document.getElementById("explain").onclick = explainKey;
document.getElementById("pop").onclick = async () => {
  await fetch("/api/layers", { method: "DELETE" });
  await loadLayers();
};

loadHealth();
loadLayers();
