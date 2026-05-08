const map = document.querySelector("#reactor-map");
const table = document.querySelector("#events-table");
const form = document.querySelector("#controls");
const playButton = document.querySelector("#play");
const resetButton = document.querySelector("#reset");
const speedInput = document.querySelector("#speed");

const fmtKg = new Intl.NumberFormat("pt-BR", {
  maximumFractionDigits: 0,
});

let demo = null;
let frame = 0;
let timer = null;
let playing = false;

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  await loadDemo();
  startPlayback();
});

playButton.addEventListener("click", () => {
  if (playing) {
    pausePlayback();
  } else {
    startPlayback();
  }
});

resetButton.addEventListener("click", () => {
  pausePlayback();
  frame = 0;
  renderFrame();
});

speedInput.addEventListener("input", () => {
  if (playing) {
    pausePlayback();
    startPlayback();
  }
});

loadDemo();

async function loadDemo() {
  pausePlayback();
  const events = document.querySelector("#events").value;
  const seed = document.querySelector("#seed").value;
  const response = await fetch(`/api/demo?events=${events}&seed=${seed}`);
  demo = await response.json();
  frame = 0;
  renderFrame();
}

function startPlayback() {
  if (!demo || demo.amostras.length === 0) return;
  playing = true;
  playButton.textContent = "Pause";
  timer = window.setInterval(nextFrame, frameDelay());
}

function pausePlayback() {
  playing = false;
  playButton.textContent = "Play";
  if (timer) {
    window.clearInterval(timer);
    timer = null;
  }
}

function nextFrame() {
  if (!demo) return;
  if (frame >= demo.amostras.length - 1) {
    pausePlayback();
    return;
  }
  frame += 1;
  renderFrame();
}

function frameDelay() {
  return 850 / Number(speedInput.value || 1);
}

function renderFrame() {
  if (!demo || demo.amostras.length === 0) return;

  const sample = demo.amostras[frame];
  const currentTime = new Date(sample.timestamp);
  const events = filterByTime(demo.eventos_inferidos, currentTime);
  const alerts = filterByTime(demo.alertas, currentTime);
  const reactors = reactorsAt(events);

  renderMetrics(sample, alerts);
  renderPlayback(sample);
  renderMap(sample, reactors);
  renderEvents(events);
  renderAlerts(alerts);
  renderLoadChart(reactors);
}

function filterByTime(items, currentTime) {
  return items.filter((item) => new Date(item.timestamp) <= currentTime);
}

function reactorsAt(events) {
  const reactors = demo.reatores.map((reactor) => ({
    ...reactor,
    carga_kg: 0,
    ciclos: 0,
    ultima_carga_kg: 0,
  }));
  const index = new Map(reactors.map((reactor) => [reactor.id, reactor]));

  for (const event of events) {
    const reactor = index.get(event.reator_id);
    if (!reactor) continue;
    reactor.carga_kg += event.quantidade_kg;
    reactor.ciclos += 1;
    reactor.ultima_carga_kg = event.quantidade_kg;
  }

  return reactors;
}

function renderMetrics(sample, alerts) {
  document.querySelector("#metric-state").textContent = sample.estado;
  document.querySelector("#metric-load").textContent = `${fmtKg.format(sample.peso_kg)} kg`;
  document.querySelector("#metric-limit").textContent = `${fmtKg.format(sample.limite_operacional_kg)} kg`;
  document.querySelector("#metric-alerts").textContent = alerts.length;
  document.querySelector("#crane-radius").textContent = `${sample.raio.toFixed(1)} m`;
  document.querySelector("#crane-height").textContent = `${sample.altura.toFixed(1)} m`;
  document.querySelector("#crane-speed").textContent = sample.velocidade.toFixed(2);
  document.querySelector("#crane-target").textContent = `R-${String(sample.reator_mais_proximo_id).padStart(2, "0")}`;
}

function renderPlayback(sample) {
  const date = new Date(sample.timestamp);
  const total = demo.amostras.length;
  const pct = total <= 1 ? 0 : (frame / (total - 1)) * 100;
  document.querySelector("#clock").textContent = date.toLocaleTimeString("pt-BR", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
  document.querySelector("#progress-label").textContent = `${frame + 1} / ${total} amostras`;
  document.querySelector("#progress-fill").style.width = `${pct}%`;
  document.querySelector("#map-subtitle").textContent = `Angulo ${sample.angulo.toFixed(1)} graus, raio ${sample.raio.toFixed(1)} m`;
}

function renderMap(sample, reactors) {
  const center = 310;
  const craneRadius = 232;
  const maxLoad = Math.max(...reactors.map((r) => r.carga_kg), 1);
  const crane = polar(center, craneRadius * (sample.raio / 8), sample.angulo);
  const people = sample.pessoas || [];
  const isLoaded = sample.carga_suspensa && sample.peso_kg > 50;

  map.innerHTML = `
    <circle cx="${center}" cy="${center}" r="${craneRadius * (3 / 8)}" fill="none" stroke="#2f8f46" stroke-width="1.5" stroke-dasharray="6 7" />
    <circle cx="${center}" cy="${center}" r="${craneRadius * (6 / 8)}" fill="none" stroke="#e2a72e" stroke-width="1.5" stroke-dasharray="6 7" />
    <circle cx="${center}" cy="${center}" r="${craneRadius}" fill="none" stroke="#d9dee7" stroke-width="2" />
    ${isLoaded ? `<circle cx="${crane.x}" cy="${crane.y}" r="39" fill="#d1495b" opacity="0.12" />` : ""}
    <circle cx="${center}" cy="${center}" r="38" fill="#edf2f4" stroke="#b8c1cc" />
    <line x1="${center}" y1="${center}" x2="${crane.x}" y2="${crane.y}" stroke="#1f2933" stroke-width="4" stroke-linecap="round" />
    <circle cx="${crane.x}" cy="${crane.y}" r="${isLoaded ? 9 : 6}" fill="${isLoaded ? "#d1495b" : "#667085"}" />
    ${reactors.map((reactor) => renderReactor(reactor, center, craneRadius, maxLoad)).join("")}
    ${people.map((person) => renderPerson(person, center, craneRadius)).join("")}
  `;
}

function renderReactor(reactor, center, craneRadius, maxLoad) {
  const point = polar(center, craneRadius * (reactor.raio / 8), reactor.angulo);
  const fillRatio = reactor.carga_kg / maxLoad;
  const size = 11 + fillRatio * 19;
  const fill = fillRatio > 0.66 ? "#d1495b" : fillRatio > 0.33 ? "#e2a72e" : "#2f8f46";

  return `
    <g>
      <circle cx="${point.x}" cy="${point.y}" r="${size}" fill="${fill}" opacity="0.9" />
      <text x="${point.x}" y="${point.y + 4}" text-anchor="middle" font-size="11" font-weight="800" fill="#fff">${reactor.id}</text>
      <title>${reactor.nome}: ${fmtKg.format(reactor.carga_kg)} kg, ${reactor.ciclos} ciclos</title>
    </g>
  `;
}

function renderPerson(person, center, craneRadius) {
  const scale = craneRadius / 8;
  const x = center + person.x * scale;
  const y = center + person.y * scale;
  const fill = person.em_zona_risco ? "#d1495b" : "#176b87";
  return `
    <g>
      <circle cx="${x}" cy="${y}" r="8" fill="${fill}" stroke="#fff" stroke-width="2" />
      <title>${person.id}</title>
    </g>
  `;
}

function renderEvents(events) {
  table.innerHTML = events
    .slice(-12)
    .reverse()
    .map((event) => {
      const date = new Date(event.timestamp);
      return `
        <tr>
          <td>${date.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</td>
          <td>R-${String(event.reator_id).padStart(2, "0")}</td>
          <td>${fmtKg.format(event.quantidade_kg)} kg</td>
          <td><span class="confidence">${Math.round(event.confianca * 100)}%</span></td>
        </tr>
      `;
    })
    .join("");
}

function renderAlerts(alerts) {
  const list = document.querySelector("#alerts-list");
  if (alerts.length === 0) {
    list.innerHTML = `<li class="empty">Sem alertas ate este instante.</li>`;
    return;
  }
  list.innerHTML = alerts
    .slice(-8)
    .reverse()
    .map((alert) => `
      <li class="${alert.nivel}">
        <strong>${alert.nivel}</strong>
        <span>${alert.mensagem}</span>
      </li>
    `)
    .join("");
}

function renderLoadChart(reactors) {
  const chart = document.querySelector("#load-chart");
  const max = Math.max(...reactors.map((reactor) => reactor.capacidade_kg), 1);
  chart.innerHTML = reactors
    .map((reactor) => {
      const pct = Math.min(100, (reactor.carga_kg / max) * 100);
      return `
        <div class="bar-row">
          <strong>R-${String(reactor.id).padStart(2, "0")}</strong>
          <div class="bar-track"><div class="bar-fill" style="width:${pct}%"></div></div>
          <span>${fmtKg.format(reactor.carga_kg)} kg</span>
        </div>
      `;
    })
    .join("");
}

function polar(center, radius, angleDeg) {
  const radians = ((angleDeg - 90) * Math.PI) / 180;
  return {
    x: center + Math.cos(radians) * radius,
    y: center + Math.sin(radians) * radius,
  };
}
