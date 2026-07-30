# Roadmap

## v1.0 — Current Release

### Built-in Policies (6 rules)

- [x] **BlockFiles** — prevents committing `.env`, `*.pem`, `*.key`, `*.p12`, `*.crt` (configurable patterns)
- [x] **SecretScan** — detects AWS keys, GitHub tokens, OpenAI keys, Stripe keys, Slack tokens in staged files
- [x] **BranchProtection** — blocks direct commits to `main`, `master`, `production` (configurable)
- [x] **CommitMessage** — enforces conventional commits: `feat:`, `fix:`, `docs:`, `test:`, etc.
- [x] **FileSize** — rejects files larger than configured limit (default 10MB)
- [x] **BinaryFile** — blocks `.exe`, `.dll`, `.so`, `.iso`, `.zip` from being committed

### CLI Commands

- [x] `install` — install global hooks
- [x] `uninstall` — remove hooks (`--all` also removes config)
- [x] `run` — run policies against current repository
- [x] `history` — show git-policy run history
- [x] `doctor` — system health check
- [x] `validate` — validate YAML config
- [x] `version` — print version
- [x] `rule list` — list rules and enabled/disabled status
- [x] `rule enable` / `rule disable` — toggle rules via CLI
- [x] `rule skip [name...]` — temporarily skip rules for current commit
- [x] `rule skip --list` — show currently skipped rules
- [x] `rule skip --clear` — clear all skipped rules
- [x] `rule add <name>` — add a custom rule directly to config (--type, --pattern, --message, --fix)
- [x] `rule remove <name>` — remove a custom rule from config
- [x] `rule export <name>` — export a custom rule as a YAML file
- [x] `rule import <file>` — import a custom rule from a YAML file
- [x] `plugins install <file>` — install a plugin from YAML descriptor
- [x] `plugins install --disabled <file>` — install with rules disabled by default
- [x] `plugins uninstall <name>` — remove a plugin and its descriptor file from disk
- [x] `plugins list` — list installed plugins with file paths
- [ ] `sync` — sync policies from remote source (Git, HTTP, S3, GCS)

### History (`history` subcommand)

Track git-policy activity across all repositories — what rules ran, what passed, what was blocked.

- [x] `git-policy history` — show recent activity (pass/fail per rule, timestamps, repository, branch)
- [x] `git-policy history --limit <n>` — restrict to the last N records (default: 20)
- [x] `git-policy history --repo` — filter by current or specific repository
- [x] `git-policy history --status pass|fail` — filter results by outcome
- [x] Enable/disable via config: `history.enabled: true` in `config.yaml`
- [x] Configurable retention: `history.maxRecords: 1000` in `config.yaml`
- [x] Storage: append-only JSONL for all repositories at `<config-dir>/history/history.jsonl`
- [x] Auto-log on every `git-policy run` — records date, repo path, branch, each rule name + status + message
- [x] `git-policy history --clear` — wipe all history, or `--clear --repo <path>` to wipe for specific repository

**Log format (JSONL — one JSON object per line):**

```json
{
  "timestamp": "2026-07-22T14:30:00Z",
  "repo":    "/home/user/projects/my-app",
  "branch":  "feature/login",
  "commit":  "abc1234",
  "results": [
    {"rule": "BlockFiles",     "status": "pass"},
    {"rule": "SecretScan",    "status": "pass"},
    {"rule": "BranchProtection", "status": "pass"},
    {"rule": "CommitMessage", "status": "fail", "message": "Commit message does not follow conventional commits: fixed bug"}
  ],
  "overall": "fail"
}
```

Fields:
- `timestamp` — ISO 8601 UTC
- `repo` — absolute path to the repository
- `branch` — branch name at time of run
- `commit` — short commit hash (if available)
- `results` — array of per-rule outcomes; each entry has `rule`, `status`, and optional `message` on failure
- `overall` — `"pass"` if all rules passed, `"fail"` otherwise

### Plugin System (YAML-Driven Custom Rules)

Custom rules defined in plain YAML — no Go compilation needed. Works on all platforms.

- [x] Plugin descriptor YAML format with `name` and `rules` array
- [x] 4 custom rule types:

| Type             | Description                              |
|------------------|------------------------------------------|
| `file-block`     | Block files matching a glob pattern      |
| `file-content`   | Scan file contents for a string pattern  |
| `branch-name`    | Block commits to branches matching glob  |
| `commit-message` | Block commits with messages matching glob|

- [x] `CustomPolicy` in `internal/policy/custom.go` — validates each rule type
- [x] `plugins.Loader.PoliciesFromPlugins` — loads descriptor files from disk
- [x] `plugins install <file>` — copies descriptor to `<config-dir>/plugins/<name>.yaml`
- [x] `plugins install --disabled` — installs with all rules disabled
- [x] `plugins uninstall <name>` — removes config entry + descriptor file from disk
- [x] `plugins list` — shows name, enabled/disabled, and file path
- [x] Relative plugin path resolution against config file directory
- [x] Cross-platform: pure YAML, no `.so` files, works on Windows/Linux/macOS

### Rule Skip (`rule skip`)

Temporarily bypass one or more rules for the current commit without globally disabling them.

```
git policy skip block-files secret-scan    # skip multiple rules
git policy skip --list                     # show currently skipped rules
git policy skip --clear                    # clear all skipped rules
```

**Storage:** Per-repository via local git config (`git config --local`) — stored in `.git/config` under `git-policy.skip`. This keeps skips isolated to a single repo and doesn't pollute global state.

**Manual clear required:** The skip list persists across all hooks and commits until explicitly cleared. Use `git-policy rule skip --clear` to remove all skipped rules. This ensures the skip works through all Git hooks (pre-commit, commit-msg) that fire during a single commit.

**Behavior:**

| Scenario | Result |
|----------|--------|
| `skip` with valid rule names | Rule names written to local git config |
| `skip` with invalid rule name | Error with available rules list |
| `skip` (no args) | Show usage / current skips |
| Already-skipped rule in `skip` | Idempotent — no duplicate |
| Rule passes after skip | Skip persists — clear manually with `--clear` |
| Non-skipped rule blocks | Skip preserved, commit blocked |

### Architecture

- [x] Go static binary, single dependency on `git` CLI
- [x] Cobra CLI framework
- [x] YAML config at `~/.config/git-policy/config.yaml`
- [x] Global hooks via `core.hooksPath`
- [x] 14 internal packages: `config`, `policy`, `engine`, `runner`, `hook`, `git`, `plugins`, `sync`, `secrets`, `commitmsg`, `files`, `branch`, `logger`, `utils`
- [x] Cross-platform: Windows, macOS, Linux
- [x] CI/CD: GitHub Actions (test matrix + goreleaser)

---

## v2 — Planned

### Remote Sync (Team Policy Distribution)

Share a single policy config across an entire team via remote sources.

- [ ] Wire `internal/sync` into the `sync` command
- [ ] `sync` command with `--source` flag (git, http, s3, gcs)
- [ ] **GitProvider** — clone/fetch policy from a Git repository
- [ ] **HTTPProvider** — fetch policy from an HTTP/HTTPS endpoint
- [ ] S3 / GCS providers
- [ ] Auto-sync on `git-policy run` (with TTL/cache)
- [ ] Sync validation — dry-run, diff local vs remote
- [ ] Signing / verification of remote policy sources

Provider interface stubs exist in `internal/sync/`.

### Plugin Enhancements

- [ ] Plugin-level enable/disable (currently all-or-nothing per plugin)
- [ ] Per-rule enable/disable within a plugin
- [ ] Plugin update command (`plugins update <name>`)
- [ ] Plugin SDK documentation for contributing custom rule types
- [ ] Plugin versioning and dependency tracking

### Per-Repository Overrides

Allow repositories to customize or extend the global policy without disabling it entirely.

- [ ] `git-policy init` — create `.git-policy.yaml` in a repository
- [ ] Merge strategy — repo-level config merges on top of global config
- [ ] Per-repo rule enable/disable
- [ ] Per-repo policy config (block patterns, branch lists, size limits)
- [ ] `git-policy rule list --local` view

### Pre-Push Enhancements

Push-specific checks beyond what `pre-commit` covers.

- [ ] Block force-push to protected branches
- [ ] Block push of large refs (exceeds size limit)
- [ ] Block push of sensitive commit history (commit message scan)
- [ ] `--max-commits` flag to limit push size

### Required Files Policy

The `requiredFiles` config field exists but no policy validates it yet.

- [ ] `RequiredFiles` policy — ensure staged commits include required files (e.g., `README.md`, `LICENSE`)
- [ ] Config-driven list of required files per commit type

### Rule Export / Import

Share individual rule configurations across machines or teams.

- [ ] `rule export <name> --format yaml|json` — export a rule's config
- [ ] `rule import <file>` — import a rule from a file
- [ ] Rule registry / marketplace (future)

### Windows Native Hooks

Replace Git Bash dependency with native PowerShell hooks.

- [ ] PowerShell `pre-commit`, `pre-push`, `commit-msg` scripts
- [ ] Detect Windows and install `.ps1` hooks instead of `.sh`
- [ ] Handle PowerShell execution policy

### Config Version Migration

- [ ] `version: 2` config schema support
- [ ] Migration path from v1 to v2 config
- [ ] Backward-compatible loading

### Observability

- [ ] `git-policy stats` — hit counters per rule, block rate, execution time
- [ ] `git-policy log` — recent policy runs with results
- [ ] JSON output mode (`--json`) for tooling integration

---

## Backlog / Ideas

These are under consideration but not yet scoped for a specific release:

- [ ] **GUI** — system tray app with policy status, block notifications, quick enable/disable
- [ ] **Pre-receive hooks** — server-side enforcement for self-hosted Git
- [ ] **VS Code extension** — inline policy warnings in the editor
- [ ] **Policy templates** — pre-built config profiles (solo dev, team, OSS, enterprise)
- [ ] **git-policy server** — central policy management dashboard for teams
