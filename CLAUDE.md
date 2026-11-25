# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
just build                          # builds to $(go env GOPATH)/bin/bark
$(go env GOPATH)/bin/bark --help
LOG_LEVEL=debug $(go env GOPATH)/bin/bark clone
```

No test suite yet. When tests are added, the convention will be `go test ./...`.

## Architecture

**bark** is a Git worktree manager CLI. The intended workflow: clone a repo as a bare git repo, then manage multiple worktrees off it (one per branch), with optional tmux session integration.

Dependency flow: `main.go` → `internal/cli/` → `internal/config/`, `internal/{git,shell,tmux,zoxide}/`

- `main.go` — wires up structured JSON logging (`log/slog`) from `$LOG_LEVEL`, then delegates entirely to `internal/cli/`
- `internal/cli/root.go` — owns startup: calls `config.Load()`, constructs all subcommands via closures (dependency injection), returns `(*cobra.Command, error)`
- `internal/cli/clone.go` — stub; will shell out to `git clone --bare` then `git worktree add`
- `internal/config/config.go` — `Load()` is the single construction point for `Config`; defaults now, file-reading later. Has a `log()` method called by `Load()` to emit debug state.

**Key constraint:** `main.go` does not import `internal/config/`. Config is a CLI concern — only `internal/cli/` depends on it.

Planned commands (stubbed in comments in `internal/cli/root.go`): `add`, `switch`, `purge`, `update`, `tui`, `status`.
