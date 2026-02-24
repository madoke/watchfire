<p align="center">
  <img src="assets/watchfire_banner-black.png" alt="Watchfire" width="600" />
</p>

<p align="center">
  <strong>Remote control for AI coding agents</strong>
</p>

---

Watchfire orchestrates coding agent sessions (starting with Claude Code) based on task files. It manages multiple projects in parallel, spawning agents in sandboxed PTYs with git worktree isolation. A daemon handles all orchestration while thin clients (CLI/TUI and GUI) connect over gRPC.

## Components

| Component | Binary | Description |
|-----------|--------|-------------|
| **Daemon** | `watchfired` | Orchestration, PTY management, git workflows, gRPC server, system tray |
| **CLI/TUI** | `watchfire` | Project-scoped CLI commands + interactive TUI mode |
| **GUI** | `Watchfire.app` | Electron multi-project client |

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+ (for GUI)
- macOS (sandbox support uses `sandbox-exec`)

### Build

```bash
# Install dev tools (golangci-lint, air, protoc plugins)
make install-tools

# Build daemon + CLI
make build

# Install to /usr/local/bin
make install
```

### Run

```bash
# Initialize a project
cd your-project
watchfire init

# Add tasks
watchfire task add

# Launch the TUI
watchfire

# Or start the GUI
make dev-gui
```

### Development

```bash
# Daemon with hot reload
make dev-daemon

# TUI (build + run)
make dev-tui

# GUI dev mode
make dev-gui
```

## Agent Modes

| Mode | Description |
|------|-------------|
| **Chat** | Interactive session with the coding agent |
| **Task** | Execute a specific task from the task list |
| **Start All** | Run all ready tasks sequentially |
| **Wildfire** | Autonomous loop: execute ready tasks, refine drafts, generate new tasks |
| **Generate Definition** | Auto-generate a project definition |
| **Generate Tasks** | Auto-generate tasks from the project definition |

## How It Works

1. **Define** your project with `watchfire init` and a project definition
2. **Create tasks** describing what you want built
3. **Start agents** that work in isolated git worktrees (one branch per task)
4. **Monitor** progress through the TUI or GUI with live terminal output
5. **Review and merge** completed work back to your default branch

The daemon watches task files for changes. When an agent marks a task as done, Watchfire automatically stops the agent, merges the worktree (if auto-merge is enabled), and chains to the next task.

## Project Structure

```
watchfire/
├── cmd/                    # Entry points
│   ├── watchfire/          # CLI/TUI binary
│   └── watchfired/         # Daemon binary
├── internal/               # Go packages
│   ├── cli/                # CLI commands
│   ├── config/             # Config loading, paths, logging
│   ├── daemon/             # Daemon internals
│   │   ├── agent/          # Agent manager, process, worktree, sandbox
│   │   ├── server/         # gRPC server + services
│   │   ├── task/           # Task manager
│   │   ├── project/        # Project manager
│   │   └── watcher/        # File watcher
│   ├── models/             # Data structures
│   └── tui/                # Bubbletea TUI
├── proto/                  # Protobuf definitions
├── gui/                    # Electron GUI
└── assets/                 # Icons, logos, brand assets
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build daemon + CLI (native arch) |
| `make build-universal` | Build universal (fat) binaries for macOS |
| `make install` | Build and install to `/usr/local/bin` |
| `make install-all` | Install CLI, daemon, and GUI app |
| `make test` | Run tests with race detector |
| `make lint` | Run golangci-lint |
| `make proto` | Regenerate protobuf code |
| `make clean` | Remove build artifacts |
| `make dev-daemon` | Run daemon with hot reload (air) |
| `make dev-tui` | Build and run TUI |
| `make dev-gui` | Run Electron GUI in dev mode |

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full design document.

## License

Proprietary. All rights reserved.
