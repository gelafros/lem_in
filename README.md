# Lem-in — Ant Colony Max-Flow Pathfinding & Turn Simulator

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/gelafros/lem_in/actions/workflows/ci.yml/badge.svg)](https://github.com/gelafros/lem_in/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A graph pathfinding engine and turn-based discrete simulator written in **Go (zero external dependencies)** that routes $N$ ants from a `##start` room to an `##end` room across an undirected tunnel network in the **minimum number of turns**, without ever exceeding room or tunnel capacity constraints. Includes an interactive **3D WebGL (Three.js) Visualizer**.

---

## Algorithmic Architecture

Naive greedy shortest-path selection fails on graphs with "trap" edges (e.g., a short bridge that blocks two slightly longer independent paths). `lem_in` models the colony as a **unit-capacity directed flow network** and applies the **Edmonds-Karp (BFS Ford-Fulkerson)** algorithm with **vertex splitting**:

1. **Vertex Splitting (`internal/solver`)**:
   Every intermediate room $v \notin \{\text{start}, \text{end}\}$ is split into two virtual nodes $v_{\text{in}} \to v_{\text{out}}$ with unit capacity $\text{cap}(v_{\text{in}}, v_{\text{out}}) = 1$. This guarantees that all extracted paths are strictly **vertex-disjoint** (no two paths share an intermediate room).
2. **Residual Graph & Backward Flow Cancellation**:
   Breadth-First Search (BFS) repeatedly finds shortest augmenting paths from `##start` to `##end`. Augmenting flow decrements forward capacity and increments reverse capacity ($\text{cap}(u, v) \mathrel{-}= 1$, $\text{cap}(v, u) \mathrel{+}= 1$), allowing subsequent BFS passes to "undo" suboptimal bottleneck edges and unlock higher-throughput disjoint path sets.
3. **Optimal Ant Scheduling (`internal/simulator`)**:
   Given extracted disjoint paths $P_1, P_2, \dots, P_k$, each ant $i \in \{1, \dots, N\}$ is greedily assigned to the path minimizing its arrival turn:
   $$\text{arrival}(v) = |P_v| + |\text{antsAssigned}(v)|$$
   Ants then advance simultaneously turn-by-turn, emitting formatted moves (`L<ant>-<room>`).

---

## Project Structure

```text
lem_in/
├── cmd/
│   └── lem-in.go                    # CLI entry point
├── internal/
│   ├── parser/
│   │   ├── parser.go                # Colony file lexer, validator & graph builder
│   │   └── parser_test.go           # Table-driven parser & edge-case tests
│   ├── solver/
│   │   ├── solver.go                # Vertex-split Edmonds-Karp max-flow solver
│   │   └── solver_test.go           # Disjointness & disconnected graph tests
│   └── simulator/
│       ├── simulator.go             # Greedy path scheduler & turn-by-turn engine
│       └── simulator_test.go        # Turn simulation verification tests
├── tests/                           # Benchmark & edge-case colony topology fixtures
├── visualizer/                      # Interactive 3D Three.js colony visualizer
│   ├── index.html
│   ├── app.js
│   ├── app.test.js
│   ├── ui.js
│   └── style.css
├── go.mod
└── LICENSE
```

---

## Getting Started

### Run the Solver

```bash
go run ./cmd tests/example01.txt
```

Example output (`tests/example00.txt`):

```text
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

### Run Unit Tests & Coverage

```bash
go test -v -race -coverpkg=./... ./...
node --test visualizer/app.test.js
```

### Interactive 3D Web Visualizer

Open `visualizer/index.html` in a browser (or serve the `visualizer/` directory with any static file server, e.g., `python3 -m http.server -d visualizer 8080`) and paste the complete solver output to step through ant movements turn-by-turn in an interactive 3D Three.js scene.

---

## License

Released under the [MIT License](LICENSE).
