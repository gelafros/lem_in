// Share the pure parser/playback API between browser scripts and Node tests.
(function (root, factory) {
  const api = factory();
  if (typeof module === "object" && module.exports) module.exports = api;
  else root.LemInVisualizer = api;
})(typeof globalThis !== "undefined" ? globalThis : this, function () {
  "use strict";

  // Read the CLI's map + blank line + turns format for replay. This parser
  // checks movement syntax, but does not enforce all solver movement rules.
  function parseLemInOutput(raw) {
    if (typeof raw !== "string" || !raw.trim()) throw new Error("Output is empty.");
    const normalized = raw.replace(/\r\n?/g, "\n").trimEnd();
    // The first blank line marks the boundary between map and movement data.
    const separator = normalized.search(/\n[ \t]*\n/);
    if (separator < 0) throw new Error("A blank line must separate the map from the moves.");

    const mapLines = normalized.slice(0, separator).split("\n");
    const moveText = normalized.slice(separator).replace(/^\n[ \t]*\n/, "");
    const antText = (mapLines.shift() || "").trim();
    if (!/^\d+$/.test(antText) || Number(antText) < 1) throw new Error("Invalid ant count.");
    const ants = Number(antText);
    const rooms = Object.create(null);
    const links = [];
    let pendingRole = null;
    let startName = null;
    let endName = null;

    for (const rawLine of mapLines) {
      const line = rawLine.trim();
      if (!line) continue;
      if (line === "##start" || line === "##end") {
        pendingRole = line.slice(2);
        continue;
      }
      if (line.startsWith("#")) continue;
      const fields = line.split(/\s+/);
      if (fields.length === 3 && /^-?\d+$/.test(fields[1]) && /^-?\d+$/.test(fields[2])) {
        const name = fields[0];
        if (rooms[name]) throw new Error(`Duplicate room: ${name}.`);
        const role = pendingRole || "normal";
        rooms[name] = { name, x: Number(fields[1]), y: Number(fields[2]), role };
        if (role === "start") startName = name;
        if (role === "end") endName = name;
        pendingRole = null;
        continue;
      }
      if (/^\S+-\S+$/.test(line)) {
        const splitAt = line.indexOf("-");
        links.push([line.slice(0, splitAt), line.slice(splitAt + 1)]);
      }
      // Unknown commands and map lines are intentionally ignored.
    }

    if (!startName) throw new Error("The map has no start room.");
    if (!endName) throw new Error("The map has no end room.");
    for (const [from, to] of links) {
      if (!rooms[from] || !rooms[to]) throw new Error(`Link references an unknown room: ${from}-${to}.`);
    }

    const turns = [];
    const moveLines = moveText.split("\n").filter((line) => line.trim());
    if (!moveLines.length) throw new Error("No movement turns were found.");
    for (const line of moveLines) {
      const moves = [];
      const seen = new Set();
      for (const token of line.trim().split(/\s+/)) {
        const match = /^L([1-9]\d*)-(\S+)$/.exec(token);
        if (!match) throw new Error(`Malformed move: ${token}.`);
        const ant = Number(match[1]);
        const room = match[2];
        if (ant > ants) throw new Error(`Move references unknown ant L${ant}.`);
        if (!rooms[room]) throw new Error(`Move references unknown room: ${room}.`);
        if (seen.has(ant)) throw new Error(`Ant L${ant} moves twice in one turn.`);
        seen.add(ant);
        moves.push({ ant, room });
      }
      turns.push(moves);
    }
    return { ants, rooms, links, turns, start: startName, end: endName };
  }

  // Keep playback independent of the DOM; notify the caller with state snapshots.
  function createPlayback(parsed, onChange) {
    if (!parsed || !parsed.start || !Array.isArray(parsed.turns)) throw new Error("Invalid parsed output.");
    let currentTurn = 0;
    let speedMs = 600;
    let timer = null;
    let antRooms = initialPositions();

    function initialPositions() {
      const result = {};
      for (let ant = 1; ant <= parsed.ants; ant++) result[ant] = parsed.start;
      return result;
    }
    // Copy positions so consumers cannot mutate the playback's internal state.
    function snapshot() {
      return { currentTurn, totalTurns: parsed.turns.length, isPlaying: timer !== null, speedMs, antRooms: { ...antRooms } };
    }
    function notify() {
      const state = snapshot();
      if (typeof onChange === "function") onChange(state);
      return state;
    }
    // Rewind by replaying from the start instead of storing every prior state.
    function rebuild(turnCount) {
      antRooms = initialPositions();
      for (let index = 0; index < turnCount; index++) {
        for (const move of parsed.turns[index]) antRooms[move.ant] = move.room;
      }
      currentTurn = turnCount;
    }
    function stepForward() {
      if (currentTurn < parsed.turns.length) {
        for (const move of parsed.turns[currentTurn]) antRooms[move.ant] = move.room;
        currentTurn++;
      }
      if (currentTurn >= parsed.turns.length && timer !== null) pause();
      return notify();
    }
    function stepBackward() {
      if (currentTurn > 0) rebuild(currentTurn - 1);
      return notify();
    }
    function reset() {
      pause();
      rebuild(0);
      return notify();
    }
    function pause() {
      if (timer !== null) clearInterval(timer);
      timer = null;
      return notify();
    }
    function play() {
      if (timer !== null || currentTurn >= parsed.turns.length) return notify();
      timer = setInterval(stepForward, speedMs);
      return notify();
    }
    function setSpeed(value) {
      const next = Number(value);
      if (!Number.isFinite(next) || next <= 0) throw new Error("Playback speed must be positive.");
      const wasPlaying = timer !== null;
      // Replace the active timer so the new speed applies during playback too.
      if (wasPlaying) clearInterval(timer);
      speedMs = next;
      if (wasPlaying) timer = setInterval(stepForward, speedMs);
      return notify();
    }
    return { stepForward, stepBackward, reset, play, pause, setSpeed, getState: snapshot };
  }

  return { parseLemInOutput, createPlayback };
});
