function q<T extends HTMLElement>(id: string): T {
  const el = document.getElementById(id);
  if (!el) throw new Error(`Missing element #${id}`);
  return el as T;
}

function setStatus(msg: string, ok = true) {
  const statusEl = q<HTMLParagraphElement>("status");
  statusEl.textContent = msg;
  statusEl.style.color = ok ? "green" : "red";
}

async function claim() {
  if (!window.FFB) return setStatus("Preload not ready (FFB undefined)", false);
  const userInput = q<HTMLInputElement>("userId");
  const carInput = q<HTMLInputElement>("carId");
  try {
    const res = await fetch(`${window.FFB.apiHostMatchmaker}/claim`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ userId: userInput.value || "user-1" }),
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body?.error || JSON.stringify(body));
    carInput.value = body.id;
    setStatus(`Claimed ${body.id} for ${body.assignedTo}`);
  } catch (e: any) {
    setStatus(`Claim failed: ${e.message}`, false);
  }
}

async function releaseCar() {
  if (!window.FFB) return setStatus("Preload not ready (FFB undefined)", false);
  const userInput = q<HTMLInputElement>("userId");
  const carInput = q<HTMLInputElement>("carId");
  try {
    const res = await fetch(`${window.FFB.apiHostMatchmaker}/release`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ userId: userInput.value, carId: carInput.value }),
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body?.error || JSON.stringify(body));
    setStatus(`Released ${body.id}`);
  } catch (e: any) {
    setStatus(`Release failed: ${e.message}`, false);
  }
}

async function loadRooms() {
  if (!window.FFB) return setStatus("Preload not ready (FFB undefined)", false);
  const roomsTable = q<HTMLTableElement>("roomsTable");
  try {
    const res = await fetch(`${window.FFB.apiHostMatchmaker}/rooms`);
    if (!res.ok) throw new Error(await res.text());
    const rows = await res.json();
    roomsTable.innerHTML = `
      <tr><th>ID</th><th>State</th><th>AssignedTo</th><th>Updated</th></tr>
      ${rows.map((r: any)=>`<tr>
        <td>${r.id}</td>
        <td>${r.state}</td>
        <td>${r.assignedTo||"-"}</td>
        <td>${new Date(r.updatedAt).toLocaleString()}</td>
      </tr>`).join("")}
    `;
    setStatus("Rooms loaded");
  } catch (e: any) {
    roomsTable.innerHTML = `<tr><td colspan="4" style="color:red">Load failed: ${e.message}</td></tr>`;
    setStatus(`Load rooms failed: ${e.message}`, false);
  }
}

window.addEventListener("DOMContentLoaded", () => {
  q<HTMLButtonElement>("btnClaim").addEventListener("click", claim);
  q<HTMLButtonElement>("btnRelease").addEventListener("click", releaseCar);
  q<HTMLButtonElement>("btnLoadRooms").addEventListener("click", loadRooms);
  if (window.FFB) {
    setStatus(`Ready. MM=${window.FFB.apiHostMatchmaker}, GW=${window.FFB.apiHostGateway}`);
  } else {
    setStatus("Preload not ready (FFB undefined). Check preload.js path.", false);
  }
});
