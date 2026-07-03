# Release Notes — v1.0.0

**Release date:** 2026-07-03

git-policy is a cross-platform CLI that provides global Git rule management.
Install once and protect every repository on your machine — block secrets,
enforce commit conventions, protect branches, reject large files, and more.

---

## What's New

### Global Git Hook Management

No per-repo setup, no `npm install`, no CI dependency. Run `git-policy install`
once and every repository on your machine is protected.

### 6 Built-in Rules

| Rule | What it does |
|------|-------------|
| **BlockFiles** | Prevents committing `.env`, `*.pem`, `*.key`, `*.p12`, `*.crt` (configurable) |
| **SecretScan** | Detects AWS keys, GitHub PATs, OpenAI keys, Stripe keys, Slack tokens, JWTs, and more |
| **BranchProtection** | Blocks direct commits to `main`, `master`, `production` (configurable) |
| **CommitMessage** | Enforces conventional commits: `feat:`, `fix:`, `docs:`, `test:`, etc. |
| **FileSize** | Rejects files over the configured limit (default 10MB) |
| **BinaryFile** | Blocks `.exe`, `.dll`, `.so`, `.iso`, `.zip` from being committed |

Each rule can be enabled/disabled via CLI or config file.

### CLI Commands

```bash
git-policy install          # Install global hooks
git-policy uninstall        # Remove hooks
git-policy uninstall --all  # Remove hooks + config
git-policy run              # Run rules against current repo
git-policy doctor           # Check system health
git-policy validate         # Validate config
git-policy rule list        # Show rules with status
git-policy rule enable      # Enable a rule
git-policy rule disable     # Disable a rule
git-policy version          # Print version
```

### Configuration

Single YAML file at `~/.config/git-policy/config.yaml` (Linux/macOS) or
`%APPDATA%\git-policy\config.yaml` (Windows). Sensible defaults if no
config file exists.

---

## Breaking Changes

None — this is the initial release.

---

## Known Issues

- **No per-repo overrides** — rules are global only
- **No custom rules** — plugin system not yet available
- **No remote sync** — team policy distribution not yet supported
- **No pre-push specific rules** — push-only checks not yet implemented

---

## Installation

```bash
# From source
git clone https://github.com/marcuwynu23/git-policy
cd git-policy
make dev

# Or download a pre-built binary from the release
```

## Contributors

- Mark Wayne Menorca

---

## Links

- [GitHub Repository](https://github.com/marcuwynu23/git-policy)
- [Changelog](../../../CHANGELOG.md)
- [Documentation](../../../GUIDE.md)
- [Contributing](../../../CONTRIBUTING.md)
