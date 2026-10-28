const test = require("node:test");
const assert = require("node:assert/strict");
const { parseLemInOutput, createPlayback } = require("./app.js");

const direct = `2
##start
start 0 0
##end
end 10 0
start-end

L1-end L2-end`;

const multi = `3
##start
s 0 0
a 3 0
##end
e 6 0
s-a
a-e

L1-a
L1-e L2-a
L2-e L3-a
L3-e`;

test("parseValidDirectLink", () => {
  const result = parseLemInOutput(direct);
  assert.equal(result.ants, 2);
  assert.equal(result.rooms.start.role, "start");
  assert.equal(result.rooms.end.role, "end");
  assert.deepEqual(result.links, [["start", "end"]]);
  assert.deepEqual(result.turns[0], [{ ant: 1, room: "end" }, { ant: 2, room: "end" }]);
});

test("parseValidMultiTurn", () => {
  const result = parseLemInOutput(multi);
  assert.equal(Object.keys(result.rooms).length, 3);
  assert.equal(result.links.length, 2);
  assert.equal(result.turns.length, 4);
  assert.deepEqual(result.turns[1], [{ ant: 1, room: "e" }, { ant: 2, room: "a" }]);
});

test("parseIgnoresComments", () => {
  const result = parseLemInOutput(direct.replace("##start", "# note\n##start"));
  assert.deepEqual(Object.keys(result.rooms).sort(), ["end", "start"]);
});

test("parseRejectsMissingSeparator", () => {
  assert.throws(() => parseLemInOutput(direct.replace("start-end\n\n", "start-end\n")), /blank line/i);
});

test("parseRejectsUnknownMoveRoom", () => {
  assert.throws(() => parseLemInOutput(direct.replace("L1-end", "L1-missing")), /unknown room/i);
});

test("parseRejectsMalformedMove", () => {
  assert.throws(() => parseLemInOutput(direct.replace("L1-end", "Lx-end")), /malformed move/i);
});

test("playbackStartsAtStart", () => {
  const state = createPlayback(parseLemInOutput(multi)).getState();
  assert.equal(state.currentTurn, 0);
  assert.deepEqual(state.antRooms, { 1: "s", 2: "s", 3: "s" });
});

test("stepForwardAppliesNextTurn", () => {
  const playback = createPlayback(parseLemInOutput(multi));
  playback.stepForward();
  const state = playback.getState();
  assert.equal(state.currentTurn, 1);
  assert.equal(state.antRooms[1], "a");
});

test("stepForwardStopsAtLastTurn", () => {
  const playback = createPlayback(parseLemInOutput(direct));
  playback.stepForward();
  const final = playback.stepForward();
  assert.equal(final.currentTurn, 1);
});

test("stepBackwardRestoresPreviousState", () => {
  const playback = createPlayback(parseLemInOutput(multi));
  playback.stepForward();
  playback.stepForward();
  const state = playback.stepBackward();
  assert.equal(state.currentTurn, 1);
  assert.equal(state.antRooms[1], "a");
  assert.equal(state.antRooms[2], "s");
});

test("resetMovesAllAntsToStart", () => {
  const playback = createPlayback(parseLemInOutput(multi));
  playback.stepForward();
  const state = playback.reset();
  assert.equal(state.currentTurn, 0);
  assert.deepEqual(state.antRooms, { 1: "s", 2: "s", 3: "s" });
});

test("setSpeedStoresPositiveSpeed and rejects invalid values", () => {
  const playback = createPlayback(parseLemInOutput(multi));
  playback.setSpeed(250);
  assert.equal(playback.getState().speedMs, 250);
  assert.throws(() => playback.setSpeed(0), /positive/i);
});

test("getState returns an independent snapshot", () => {
  const playback = createPlayback(parseLemInOutput(multi));
  const snapshot = playback.getState();
  snapshot.antRooms[1] = "e";
  assert.equal(playback.getState().antRooms[1], "s");
});
