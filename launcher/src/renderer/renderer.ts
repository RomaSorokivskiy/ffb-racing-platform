type Car = { id: string; state: "FREE"|"RESERVED"|"BUSY"; assignedTo?: string; updatedAt: string; ttl?: number };

let API_MM = window.FFB?.apiHostMatchmaker ?? "http://localhost:8081";
let API_GW = window.FFB?.apiHostGateway ?? "http://localhost:8080";
let evtSource: EventSource | null = null;

function qs<T extends HTMLElement>(s: string): T { const el = document.querySelector(s); if (!el) throw new Error(`Missing ${s}`); return el as T; }

function setStatus(msg: string, ok = true) { const n = qs<HTMLDivElement>("#status"); n.textContent = msg; n.style.color = ok ? "#9fe2b0" : "#ff9c9c"; }

function badge(state: Car["state"]) {
  const cls = state === "FREE" ? "badge free" : state === "BUSY" ? "badge busy" : "badge reserved";
  return `<span class="${cls}">${state}</span>`;
}

function renderRows(rows: Car[]) {
  const t = qs<HTMLTableElement>("#roomsTable");
  t.innerHTML = `
    <tr><th>ID</th><th>State</th><th>Assigned</th><th>TTL</th><th>Updated</th></tr>
    ${rows.map(r => `
      <tr>
        <td>${r.id}</td>
        <td>${badge(r.state)}</td>
        <td>${r.assignedTo||"-"}</td>
        <td>${r.state==="RESERVED"?(r.ttl??0):"-"}</td>
        <td>${new Date(r.updatedAt).toLocaleString()}</td>
      </tr>`).join("")}
  `;
}

async function loadRooms() {
  try {
    const res = await fetch(`${API_MM}/rooms`);
    if (!res.ok) throw new Error(await res.text());
    const rows: Car[] = await res.json();
    renderRows(rows);
    setStatus(`Rooms loaded (${rows.length})`);
  } catch (e:any) {
    setStatus(`Load rooms failed: ${e.message}`, false);
  }
}

async function claim() {
  const user = qs<HTMLInputElement>("#userId").value || "user-1";
  const ttl = parseInt(qs<HTMLInputElement>("#ttl").value || "120", 10);
  try {
    const res = await fetch(`${API_MM}/claim`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ userId: user, ttlSec: ttl })
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body?.error || JSON.stringify(body));
    qs<HTMLInputElement>("#carId").value = body.id;
    setStatus(`Claimed ${body.id} for ${user}`);
  } catch(e:any) {
    setStatus(`Claim failed: ${e.message}`, false);
  }
}

async function releaseCar() {
  const user = qs<HTMLInputElement>("#userId").value || "user-1";
  const car  = qs<HTMLInputElement>("#carId").value;
  try {
    const res = await fetch(`${API_MM}/release`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ userId: user, carId: car })
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body?.error || JSON.stringify(body));
    setStatus(`Released ${body.id}`);
  } catch(e:any) {
    setStatus(`Release failed: ${e.message}`, false);
  }
}

async function createSession() {
  const user = qs<HTMLInputElement>("#userId").value || "user-1";
  const car  = qs<HTMLInputElement>("#carId").value;
  if (!car) { setStatus("Select car first (Claim)", false); return; }
  try {
    const res = await fetch(`${API_GW}/session/create`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ userId: user, carId: car })
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body?.error || JSON.stringify(body));
    setStatus(`Session token created (len=${(body.token||"").length})`);
    // TODO: pass token to native client.exe on Play
  } catch(e:any) {
    setStatus(`Session create failed: ${e.message}`, false);
  }
}

function connectEvents() {
  if (evtSource) { evtSource.close(); evtSource = null; }
  evtSource = new EventSource(`${API_MM}/events`);
  evtSource.onmessage = (e) => {
    try {
      const ev = JSON.parse((e as MessageEvent).data);
      if (ev.type === "snapshot") renderRows(ev.data);
      if (ev.type === "update") {
        // naive update: reload
        loadRooms();
      }
    } catch {}
  };
  evtSource.onerror = () => setStatus("SSE disconnected (retrying…)", false);
}

function switchTab(target: string) {
  document.querySelectorAll(".item").forEach(el => el.classList.remove("active"));
  document.querySelector(`.item[data-tab="${target}"]`)?.classList.add("active");
  document.querySelectorAll(".tab").forEach(el => (el as HTMLElement).style.display = "none");
  qs<HTMLElement>(`#tab-${target}`).style.display = "";
}

function initNav() {
  document.querySelectorAll(".item").forEach(el => {
    el.addEventListener("click", () => {
      const tab = (el as HTMLElement).dataset.tab!;
      switchTab(tab);
    });
  });
}

function initSettings() {
  const gw = qs<HTMLInputElement>("#cfgGateway");
  const mm = qs<HTMLInputElement>("#cfgMatchmaker");
  gw.value = API_GW; mm.value = API_MM;
  qs<HTMLButtonElement>("#btnSaveCfg").addEventListener("click", () => {
    API_GW = gw.value || API_GW;
    API_MM = mm.value || API_MM;
    connectEvents();
    setStatus("Config saved");
  });
}

window.addEventListener("DOMContentLoaded", () => {
  initNav();
  initSettings();

  qs<HTMLButtonElement>("#btnLoadRooms").addEventListener("click", loadRooms);
  qs<HTMLButtonElement>("#btnClaim").addEventListener("click", claim);
  qs<HTMLButtonElement>("#btnRelease").addEventListener("click", releaseCar);
  qs<HTMLButtonElement>("#btnCreateSession").addEventListener("click", createSession);

  connectEvents();
  loadRooms();

  setStatus(`Ready. MM=${API_MM}, GW=${API_GW}`);
});
