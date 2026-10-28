// Connect DOM controls and the Three.js scene to the standalone playback API.
(function () {
  "use strict";
  const $ = (id) => document.getElementById(id);
  const elements = {
    input: $("output"), load: $("load"), error: $("error"), scene: $("scene"),
    play: $("play"), pause: $("pause"), step: $("step"), back: $("back"), reset: $("reset"),
    speed: $("speed"), speedValue: $("speed-value"), turn: $("turn"), summary: $("summary")
  };
  let parsed = null;
  let playback = null;
  let view = null;

  // A playback snapshot drives both control availability and ant positions.
  function update(state) {
    if (!state) return;
    elements.turn.textContent = `Turn ${state.currentTurn} / ${state.totalTurns}`;
    elements.play.disabled = state.isPlaying || state.currentTurn >= state.totalTurns;
    elements.pause.disabled = !state.isPlaying;
    elements.step.disabled = state.currentTurn >= state.totalTurns;
    elements.back.disabled = state.currentTurn === 0;
    elements.reset.disabled = state.currentTurn === 0;
    if (view) view.updateAnts(state.antRooms);
  }

  function setControls(enabled) {
    for (const key of ["play", "pause", "step", "back", "reset"]) elements[key].disabled = !enabled;
  }

  elements.load.addEventListener("click", () => {
    elements.error.textContent = "";
    try {
      // Stop the previous replay before replacing its data and scene.
      if (playback) playback.pause();
      parsed = LemInVisualizer.parseLemInOutput(elements.input.value);
      if (view) view.dispose();
      view = createView(elements.scene, parsed);
      playback = LemInVisualizer.createPlayback(parsed, update);
      playback.setSpeed(elements.speed.value);
      elements.summary.textContent = `${parsed.ants} ants · ${Object.keys(parsed.rooms).length} rooms · ${parsed.links.length} tunnels`;
      setControls(true);
      update(playback.getState());
    } catch (error) {
      elements.error.textContent = error.message;
      setControls(false);
    }
  });
  elements.play.addEventListener("click", () => playback && playback.play());
  elements.pause.addEventListener("click", () => playback && playback.pause());
  elements.step.addEventListener("click", () => playback && playback.stepForward());
  elements.back.addEventListener("click", () => playback && playback.stepBackward());
  elements.reset.addEventListener("click", () => playback && playback.reset());
  elements.speed.addEventListener("input", () => {
    elements.speedValue.textContent = `${elements.speed.value} ms`;
    if (playback) playback.setSpeed(elements.speed.value);
  });
  setControls(false);

  // Own the renderer and its resources for one loaded colony.
  function createView(container, data) {
    if (typeof THREE === "undefined") throw new Error("Three.js could not load. Check your internet connection and reload.");
    container.innerHTML = "";
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x07111f);
    scene.fog = new THREE.Fog(0x07111f, 32, 90);
    const camera = new THREE.PerspectiveCamera(52, 1, 0.1, 200);
    camera.position.set(0, 13, 24);
    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    container.appendChild(renderer.domElement);
    const controls = THREE.OrbitControls ? new THREE.OrbitControls(camera, renderer.domElement) : null;
    if (controls) { controls.enableDamping = true; controls.dampingFactor = 0.06; }
    scene.add(new THREE.HemisphereLight(0xb9ddff, 0x18202a, 1.4));
    const light = new THREE.DirectionalLight(0xffffff, 1.1);
    light.position.set(12, 18, 10);
    scene.add(light);

    // Center and scale map coordinates to fit the view. Small depth offsets
    // add visual separation without changing the map's tunnel relationships.
    const roomEntries = Object.values(data.rooms);
    const centerX = roomEntries.reduce((sum, room) => sum + room.x, 0) / roomEntries.length;
    const centerY = roomEntries.reduce((sum, room) => sum + room.y, 0) / roomEntries.length;
    const extent = Math.max(1, ...roomEntries.map((room) => Math.max(Math.abs(room.x - centerX), Math.abs(room.y - centerY))));
    const scale = Math.min(2.2, 11 / extent);
    const positions = {};
    roomEntries.forEach((room, index) => {
      const z = ((index % 3) - 1) * 0.55;
      positions[room.name] = new THREE.Vector3((room.x - centerX) * scale, (room.y - centerY) * scale, z);
    });

    const tunnelMaterial = new THREE.LineBasicMaterial({ color: 0x476985, transparent: true, opacity: 0.78 });
    for (const [from, to] of data.links) {
      const geometry = new THREE.BufferGeometry().setFromPoints([positions[from], positions[to]]);
      scene.add(new THREE.Line(geometry, tunnelMaterial));
    }
    const roomGeometry = new THREE.SphereGeometry(0.42, 24, 16);
    const colors = { start: 0x48d5a7, end: 0xff6f91, normal: 0x62a9e8 };
    for (const room of roomEntries) {
      const sphere = new THREE.Mesh(roomGeometry, new THREE.MeshStandardMaterial({ color: colors[room.role], roughness: 0.45 }));
      sphere.position.copy(positions[room.name]);
      sphere.userData.room = room.name;
      scene.add(sphere);
    }
    const antGeometry = new THREE.SphereGeometry(0.2, 18, 12);
    const antMeshes = {};

    // Sprites keep ant labels facing the camera as the user orbits the scene.
    function createAntLabel(ant) {
      const canvas = document.createElement("canvas");
      canvas.width = 128;
      canvas.height = 64;
      const context = canvas.getContext("2d");
      context.fillStyle = "rgba(5, 12, 22, 0.88)";
      context.beginPath();
      context.roundRect(4, 4, 120, 56, 16);
      context.fill();
      context.strokeStyle = "#ffd166";
      context.lineWidth = 4;
      context.stroke();
      context.fillStyle = "#ffffff";
      context.font = "700 36px ui-sans-serif, sans-serif";
      context.textAlign = "center";
      context.textBaseline = "middle";
      context.fillText(`L${ant}`, 64, 34);

      const texture = new THREE.CanvasTexture(canvas);
      texture.minFilter = THREE.LinearFilter;
      const label = new THREE.Sprite(new THREE.SpriteMaterial({ map: texture, transparent: true, depthTest: false }));
      label.position.set(0, 0.5, 0);
      label.scale.set(0.9, 0.45, 1);
      label.renderOrder = 10;
      return label;
    }

    for (let ant = 1; ant <= data.ants; ant++) {
      const antGroup = new THREE.Group();
      const hue = ((ant - 1) * 137.508) % 360;
      const antMaterial = new THREE.MeshStandardMaterial({
        color: new THREE.Color(`hsl(${hue}, 78%, 58%)`),
        roughness: 0.38
      });
      const ball = new THREE.Mesh(antGeometry, antMaterial);
      ball.userData.ant = ant;
      antGroup.add(ball);
      antGroup.add(createAntLabel(ant));
      antGroup.userData.ant = ant;
      antMeshes[ant] = antGroup;
      scene.add(antGroup);
    }

    const legend = document.createElement("div");
    legend.className = "legend";
    legend.innerHTML = '<span><i class="dot start"></i>Start</span><span><i class="dot end"></i>End</span><span><i class="dot ant"></i>Ant</span>';
    container.appendChild(legend);

    // Offset ants around shared rooms so start/end crowds remain distinguishable.
    function updateAnts(antRooms) {
      const occupants = {};
      for (const ant of Object.keys(antRooms).sort((a, b) => Number(a) - Number(b))) {
        const room = antRooms[ant];
        const slot = occupants[room] || 0;
        occupants[room] = slot + 1;
        const angle = slot * 2.399;
        const radius = slot === 0 ? 0.72 : 0.72 + Math.floor(slot / 6) * 0.32;
        antMeshes[ant].position.copy(positions[room]).add(new THREE.Vector3(Math.cos(angle) * radius, Math.sin(angle) * radius, 0.55));
      }
    }
    updateAnts(Object.fromEntries(Array.from({ length: data.ants }, (_, i) => [i + 1, data.start])));

    function resize() {
      const width = container.clientWidth || 1;
      const height = container.clientHeight || 1;
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height, false);
    }
    const resizeObserver = new ResizeObserver(resize);
    resizeObserver.observe(container);
    resize();
    let frame;
    // Render continuously for camera damping, even while playback is paused.
    function animate() {
      frame = requestAnimationFrame(animate);
      if (controls) controls.update();
      renderer.render(scene, camera);
    }
    animate();

    // Release animation, observers, and GPU resources when loading another map.
    function dispose() {
      cancelAnimationFrame(frame);
      resizeObserver.disconnect();
      if (controls) controls.dispose();
      scene.traverse((object) => {
        if (object.geometry) object.geometry.dispose();
        if (object.material) {
          const materials = Array.isArray(object.material) ? object.material : [object.material];
          materials.forEach((material) => {
            if (material.map) material.map.dispose();
            material.dispose();
          });
        }
      });
      renderer.dispose();
      container.innerHTML = "";
    }
    return { updateAnts, dispose };
  }
})();
