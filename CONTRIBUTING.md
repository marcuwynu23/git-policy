# Contributing to git-policy

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Structure](#project-structure)
3. [Development Setup](#development-setup)
4. [Building](#building)
5. [Testing](#testing)
6. [Adding a New Policy](#adding-a-new-policy)
7. [Code Standards](#code-standards)
8. [Release Process](#release-process)

---

## Prerequisites

- **Go 1.25+** — download from [go.dev](https://go.dev/dl)
- **Make** — optional, for Makefile convenience targets
- **Git** — for version control

## Project Structure

```
├── cmd/                  # CLI commands (Cobra)
│   ├── root.go           # Root command & config init
│   ├── install.go        # Hook installation
│   ├── uninstall.go      # Hook removal
│   ├── run.go            # Policy execution
│   ├── doctor.go         # System health check
│   ├── validate.go       # Config validation
│   ├── version.go        # Version info
│   ├── sync.go           # Remote sync (stub)
│   └── policy.go         # Rule enable/disable/list (`policy` alias)
├── internal/
│   ├── config/           # YAML config load/save
│   ├── hook/             # Global hook installer
│   ├── policy/           # Policy interface & built-in policies
│   ├── engine/           # Policy execution engine
│   ├── runner/           # Hook runner wiring
│   ├── git/              # Git CLI wrapper
│   ├── secrets/          # Secret pattern scanning
│   ├── files/            # File blocking & size checks
│   ├── commitmsg/        # Conventional commit validation
│   ├── branch/           # Branch protection
│   ├── plugins/          # YAML-driven custom rules loader
│   ├── sync/             # Sync providers (stub)
│   ├── logger/           # slog setup
│   └── utils/            # Shared utilities
├── hooks/                # Sample hook scripts
├── examples/             # Example config
├── testdata/             # Test fixtures
├── Makefile              # Build & test targets
└── main.go               # Entry point
```

### Package Responsibilities

| Package | Responsibility |
|---------|---------------|
| `cmd/` | Cobra command definitions, CLI entry points |
| `config/` | Load, parse, save YAML configuration |
| `hook/` | Install/uninstall global Git hooks |
| `policy/` | Policy interface + 6 built-in policy implementations |
| `engine/` | Orchestrates policy registration & execution |
| `runner/` | Wires Git context (files, branch, commit msg) into the engine |
| `git/` | Git CLI wrappers (branch, staged files, commit message) |
| `secrets/` | Secret pattern matching engine |
| `files/` | File blocking, binary detection, size checks |
| `commitmsg/` | Conventional commit format validation |
| `branch/` | Branch protection checks |
| `plugins/` | YAML-driven custom rules loader (4 rule types) |
| `logger/` | Structured logging via slog |

---

## Development Setup

```bash
# Clone the repository
git clone https://github.com/marcuwynu23/git-policy
cd git-policy

# Full dev setup: build + link to PATH + install hooks
make dev

# Or step by step:
make build          # Build the binary
make install-binary # Copy binary to C:\Bin\tools (Windows) or /usr/local/bin
make install        # Install global hooks
```

---

## Building

```bash
make build          # Build for current platform
make dist           # Build for all platforms (Windows, Linux, macOS)
```

The binary is output as `git-policy` (or `git-policy.exe` on Windows).

Cross-platform builds output to `dist/`:

- `git-policy-windows-amd64.exe`
- `git-policy-linux-amd64`
- `git-policy-linux-arm64`
- `git-policy-darwin-amd64`
- `git-policy-darwin-arm64`

---

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run a specific package
go test -v ./internal/policy/...

# Run tests without cache
go test -count=1 ./internal/...
```

### Writing Tests

Tests use Go's standard `testing` package and table-driven patterns:

```go
func TestBlockFilesPolicy_Pass(t *testing.T) {
    cfg := config.DefaultConfig()
    p := NewBlockFilesPolicy(cfg)
    result := p.Validate(Context{
        StagedFiles: []string{"main.go"},
    })
    if result.Status != StatusPass {
        t.Errorf("expected pass, got %s", result.Status)
    }
}
```

Test files live alongside the code they test (e.g., `policy_test.go` next to `policy.go`).

---

## Adding a New Rule

### Step 1: Create the Rule File

Add a new file in `internal/policy/` and implement the `Policy` interface:

```go
package policy

type MyPolicy struct {
    cfg *config.Config
}

func NewMyPolicy(cfg *config.Config) *MyPolicy {
    return &MyPolicy{cfg: cfg}
}

func (p *MyPolicy) Name() string {
    return "MyPolicy"
}

func (p *MyPolicy) Validate(ctx Context) Result {
    // Your validation logic here
    return Result{
        PolicyName: p.Name(),
        Status:     StatusPass,
    }
}
```

### Policy Interface

```go
type Policy interface {
    Name() string
    Validate(ctx Context) Result
}
```

### Context

The `Context` struct provides data from the current Git operation:

```go
type Context struct {
    RepoPath    string   // Path to the Git repo
    StagedFiles []string // Files staged for commit
    CommitMsg   string   // Commit message (from COMMIT_EDITMSG)
    BranchName  string   // Current branch name
}
```

### Result

Return a `Result` to pass or fail:

```go
type Result struct {
    PolicyName string // Your policy name
    Status     Status // StatusPass or StatusFail
    Message    string // Human-readable explanation
    Fix        string // Suggestion on how to fix
}
```

### Step 2: Add Config Fields

If your rule needs configuration, add fields to `PoliciesConfig` in `internal/config/config.go`:

```go
type PoliciesConfig struct {
    // ... existing fields ...
    MySetting  bool     `yaml:"mySetting"`
}
```

Also update `DefaultConfig()` with sensible defaults.

### Step 3: Register the Rule

In `internal/runner/runner.go`, register your policy:

```go
eng.Register(policy.NewBlockFilesPolicy(cfg))
eng.Register(policy.NewCommitMessagePolicy(cfg))
eng.Register(policy.NewMyPolicy(cfg))          // ← add this line
// ... other policies ...
```

### Step 4: Add Enable/Disable Support

Add the CLI name in `internal/config/config.go`'s `PolicyNames` map (the command uses `rule`; `policy` is accepted as an alias):

```go
var PolicyNames = map[string]string{
    // ... existing ...
    "my-policy": "MyPolicy",
}
```

### Step 5: Write Tests

Cover at least:
- **Pass case** — valid input, policy should pass
- **Fail case** — invalid input, policy should fail
- **Edge cases** — empty input, nil values, boundary conditions

---

## Code Standards

- Follow standard **Go idioms** and `golangci-lint`.
- Keep **functions focused** and small — one responsibility per function.
- Favor **dependency injection** over globals or package-level state.
- Use **interfaces** only where they improve testability.
- **Wrap errors** with context using `fmt.Errorf("doing x: %w", err)`.
- **Never panic** in production code. Return errors.
- Keep **packages cohesive** — if a package grows too broad, split it.
- Maintain **90%+ test coverage** for `internal/` packages.
- All **exported symbols** must have GoDoc comments.
- Use `slog` for logging — see `internal/logger/` for setup.

### Common Tasks Quick Reference

```bash
make fmt         # Format all Go code
make vet         # Static analysis
make lint        # golangci-lint
make tidy        # go mod tidy + verify
```

---

## Release Process

Releases are built with [goreleaser](https://goreleaser.com):

```bash
make dist
```

This produces static binaries for:
- Windows (amd64)
- Linux (amd64, arm64)
- macOS (amd64, arm64)

Before releasing:

1. Update the `version` variable in `cmd/version.go`
2. Run `make test` and verify all tests pass
3. Run `make lint` and `make vet`
4. Tag the release: `git tag v0.1.0 && git push --tags`
5. Build and upload: `goreleaser release`

---

## Getting Help

- Open an issue on GitHub
- See [USER-GUIDE.md](USER-GUIDE.md) for usage questions
