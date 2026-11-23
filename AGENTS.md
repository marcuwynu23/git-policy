# AGENTS.md — git-policy

> Generated: 2026-11-24 · by project-aware-agent-maker · evidence-based, re-run before trusting · project-manager updates 2026-11-24

## Project status

git-policy is a released, CI-covered Go CLI (v1.1.0 tag history) that installs global Git hooks and runs six built-in policies plus YAML-defined custom rules; `main.go:5` calls `cmd.Execute` (`cmd/root.go:25`), and `cmd/run.go:18` delegates to `internal/runner/runner.go:19` → `internal/engine/engine.go:45`. All local quality gates executed this session exited 0 (`make lint`, `make vet`, `make test`, `make build`, CI-form coverage test). Repository is `marcuwynu23/git-policy` (`gh repo view`), with no root `PLAN.md` and several documentation drift items recorded below.

- CI: green — `gh run list --limit 10 --json databaseId,name,headBranch,status,conclusion,workflowName` returned 10/10 `conclusion: success`, including release run `30564856094` on branch `v1.1.0` and test run `30564832509` on `main`; no failed run required `--log-failed`.
- Board: `git-policy Roadmap` project 25 (https://github.com/users/marcuwynu23/projects/25) was created and synced 2026-11-24 after explicit user approval; 12 draft items are all `Todo`. Repository-linked `projectsV2` nodes were `[]` before creation.

| Area | Finding | Evidence |
| --- | --- | --- |
| Languages / framework / versions | Go module `github.com/marcuwynu23/git-policy`, `go 1.26.1`; Cobra `v1.10.2`, yaml.v3 `v3.0.1`; installed Go `go1.26.1 windows/amd64` | `go.mod:1-8`; `go version` output |
| Architecture / entry points | Cobra CLI in `cmd/`; execution pipeline `cmd/run.go:18-24` → `internal/runner/runner.go:19-82` → `internal/engine/engine.go:45-58`; six built-ins registered at `internal/runner/runner.go:40-45` | `main.go:5-7`, `cmd/root.go:12-35`, cited files |
| Quality gates | `make lint`, `make vet`, `make test`, `make build` all exit 0; CI coverage test also exit 0 | `Makefile:23-41`; command output this session |
| CI | `test.yml`: 3-OS matrix (ubuntu/macos/windows) with coverage + lint + vet; `release.yml`: `v*` tags and `workflow_dispatch`, cross-build + `softprops/action-gh-release@v2.2.2` | `.github/workflows/test.yml:10-60`, `.github/workflows/release.yml:4-50` |
| Docs state | Root docs present; no `PLAN.md`; README/CONTRIBUTING/CHANGELOG/SECURITY contain drift (see Open questions) | `Get-ChildItem -Force`, `README.md`, `CONTRIBUTING.md`, `CHANGELOG.md`, `SECURITY.md` |
| Coverage policy | `CONTRIBUTING.md:269` requires 90%+ for `internal/`; not met today and not enforced in CI | `go test -count=1 -covermode=atomic -coverprofile=... ./...` output (e.g. `internal/runner` 19.0%, `internal/config` 46.5%, `internal/policy` 92.2%) |
| Planning / board | `ROADMAP.md` rewritten 2026-11-24 with current v1.1.0 status and ranked queue; project 25 has 12 draft items with Status/Priority/Effort fields | `ROADMAP.md`, `gh project item-list 25 --owner marcuwynu23 --limit 30 --format json` |

## Environment & toolchain

| Tool | Version | Evidence |
| --- | --- | --- |
| Go | go1.26.1 windows/amd64 | `go version` |
| Module Go directive | go 1.26.1 | `go.mod:3` |
| golangci-lint | v1.64.8 (built with go1.26.1) | `golangci-lint --version` |
| Git | 2.36.1.windows.1 | `git --version` |
| GNU Make | 3.81 | `make --version` |
| GitHub CLI | 2.73.0 (2025-05-19), authenticated as `marcuwynu23` | `gh --version`, `gh auth status` |

## Commands

| Purpose | Command | Verified |
| --- | --- | --- |
| install dependencies | `go mod download` | yes (exit 0 this session) |
| test (full suite) | `make test` | yes (exit 0; all packages PASS) |
| test with CI coverage flags (Windows PowerShell form) | `go test -count=1 -covermode=atomic -coverprofile="$env:TEMP\git-policy-coverage.out" ./...` | yes (exit 0 this session) |
| lint (local) | `make lint` | yes (exit 0 this session) |
| lint (exact CI form) | `golangci-lint run ./... --out-format=colored-line-number` | yes (exit 0 this session) |
| vet | `make vet` | yes (exit 0 this session) |
| build | `make build` | yes (exit 0 this session) |
| format check | `gofmt -l .` | yes (exit 0; lists 21 files — see Open questions) |
| working tree status | `git status --short` | yes (exit 0 this session) |
| diff whitespace check | `git diff --check` | yes (exit 0; no whitespace errors) |
| diff summary | `git diff --stat` | yes (ROADMAP.md: 123 insertions, 211 deletions) |
| recent commits | `git log --oneline -15` | yes (history through `3b763ec`) |
| open issues | `gh issue list --state open --limit 50 --json number,title,labels,milestone,updatedAt` | yes (returned `[]`) |
| closed issues | `gh issue list --state closed --limit 50 --json number,title,labels,milestone,updatedAt` | yes (returned `[]`) |
| labels | `gh label list` | yes (default labels only) |
| releases | `gh release list --limit 5` | yes (`v1.1.0` latest, 2026-07-30) |
| CI history | `gh run list --limit 10 --json databaseId,name,headBranch,status,conclusion,workflowName` | yes (10/10 success) |
| project command help | `gh project create --help`, `gh project field-create --help`, `gh project field-list --help`, `gh project item-create --help`, `gh project item-edit --help`, `gh project item-list --help` | yes (flags captured) |
| linked projects (GraphQL) | `gh api graphql -f query='query { repository(owner: "marcuwynu23", name: "git-policy") { projectsV2(first: 100) { nodes { number title url } } } }'` | yes (nodes `[]` before creation) |
| create project | `gh project create --owner marcuwynu23 --title "git-policy Roadmap" --format json` | yes (project 25, id `PVT_kwHOA7B5Z84Bs58E`) |
| list project fields | `gh project field-list 25 --owner marcuwynu23 --limit 30 --format json` | yes (13 default fields) |
| create project fields | `gh project field-create 25 --owner marcuwynu23 --name "Priority" --data-type SINGLE_SELECT --single-select-options "P0,P1,P2,P3" --format json` (and `Effort` with `S,M,L,XL`) | yes (option IDs captured) |
| create project item | `gh project item-create 25 --owner marcuwynu23 --title "..." --body "..." --format json` | yes (12 draft items created) |
| set project item fields | `gh project item-edit --id <item-id> --project-id PVT_kwHOA7B5Z84Bs58E --field-id <field-id> --single-select-option-id <option-id> --format json` | yes (all 12 items) |
| list project items | `gh project item-list 25 --owner marcuwynu23 --limit 30 --format json` | yes (`totalCount: 12`) |

## Testing & loop engineering (mandatory)

- **Verified commands only:** use only the commands above unless a new command is executed first and added to the table with real output; never mark an unexecuted command `yes`.
- **Test gates:** run `make lint`, `make vet`, `make test`, and `make build` before the first change and after every change. A failing test, a skipped test without a recorded reason, or a failed lint/vet/build blocks completion.
- **Regression:** the final pass must re-run the full suite (`make test`) plus lint/vet/build, not only the package containing the edit. `Makefile:23-24` uses `-count=1`; preserve that to avoid cached results.
- **CI parity:** `Makefile:23-24` runs `go test -v -count=1 ./...`, while `.github/workflows/test.yml:27` runs `go test -v -count=1 -coverprofile=${{ runner.temp }}/coverage.out -covermode=atomic ./...` on three OSes. Both local forms are in the Commands table; do not claim parity from `make test` alone. `Makefile:38` runs `golangci-lint run ./...`, CI adds `--out-format=colored-line-number` (`.github/workflows/test.yml:57`); the CI form is verified locally.
- **Coverage/flaky policy:** keep the 90%+ `internal/` target from `CONTRIBUTING.md:269`; never lower or delete it. Current measurements are below it in several packages, so report the miss with numbers instead of weakening the gate. CI uploads coverage to Codecov with `fail_ci_if_error: false` (`.github/workflows/test.yml:29-33`). No flaky-test policy was found in the repository; record intermittent failures with evidence rather than silently retrying until green.
- **Change discipline:** plan → execute → verify → correct → repeat; never declare done on the first pass. Tests exercise the real `git` CLI and temporary repositories (`internal/git/git_test.go`, `cmd/cmd_test.go`), so preserve test isolation and avoid mutating global Git config during tests.

## Conventions

- Follow Go idioms and the enabled linters in `.golangci.yml:4-11` (`errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`); keep focused functions, wrap errors with `%w`, avoid panics in production code, add GoDoc to exported symbols, and use `slog` via `internal/logger/` (`CONTRIBUTING.md:260-271`).
- Cobra commands live in `cmd/*.go` and register themselves in `init()` (`cmd/root.go:32-35`, `cmd/run.go:27-29`); `rule` is the primary subcommand and `policy` is an alias (`README.md:193`, `CONTRIBUTING.md:242`).
- CLI rule names are kebab-case and map to internal policy names through `config.PolicyNames` (`internal/runner/runner.go:103-120`, `CONTRIBUTING.md:242-249`).
- New built-in policies follow the documented sequence: implement in `internal/policy/`, add config fields/defaults, register in `internal/runner/runner.go:40-45`, add the CLI name, and add pass/fail/edge-case tests next to the code (`CONTRIBUTING.md:151-257`, `:130-147`).
- Tests use Go's standard `testing` package, table-driven patterns, and files alongside the code (`CONTRIBUTING.md:130-147`).
- Commit messages must satisfy the project's own conventional-commit policy (`feat:`, `fix:`, `docs:`, etc.) enforced by `CommitMessage` (`README.md:133`, `USER-GUIDE.md:510-527`); use one Conventional Commit per logical change.
- Builds are static (`CGO_ENABLED=0`, `Makefile:9,17`); artifacts go to `dist/` and are ignored (`.gitignore:6`).
- Project-manager workflow explicitly owns `ROADMAP.md` and the linked GitHub Project; the user approved board creation on 2026-11-24.

## Planning documents (ROADMAP.md maintained by project-manager workflow; others reference only — do not rewrite)

- `ROADMAP.md` — rewritten 2026-11-24 with current v1.1.0 status, implemented/verdict tables, fetched research, ranked Now/Next/Later queue, issue triage, and project 25 link.
- `README.md` — product context, rule list, command overview.
- `CONTRIBUTING.md` — development setup, code standards, test and release process.
- `USER-GUIDE.md` — configuration reference, rule management, troubleshooting.
- `CHANGELOG.md`, `CODE_OF_CONDUCT.md`, `CONTRIBUTORS.md`, `MAINTAINERS.md`, `SECURITY.md`, `project-aware-agent-maker.md` — read-only references.
- `PLAN.md` — not present at the repository root at generation time (`Get-ChildItem -Force`); do not create or invent one without asking.

## Known priorities (from gh issues and roadmap audit)

- `gh issue list --state open --limit 50` returned `[]` (no open issues as of 2026-11-24); no issue numbers, labels, or milestones to report. All 12 board items are `needs-issue` drafts.
- Roadmap work (`ROADMAP.md` Prioritized execution queue): P1 docs drift, internal coverage, gofmt; P2 RequiredFiles policy, pre-push, sync, SECURITY/CHANGELOG, per-repo overrides, Windows native hooks; P3 plugin scope, observability, config/GUI spike.
- CI is green (see Project status); watch pushed runs to success rather than assuming the last result persists.

## Open questions

- `README.md:7` and `CONTRIBUTING.md:18` say Go 1.25+, while `go.mod:3` and `go version` show 1.26.1. Which is authoritative for docs?
- `README.md:225-226` says custom rules/plugins are not implemented; `internal/plugins/plugins.go`, `internal/policy/custom.go`, `cmd/policy.go`, and `USER-GUIDE.md:384-489` show they are. Confirm which statement to keep.
- `CONTRIBUTING.md:286,303-304` describe goreleaser; `.github/workflows/release.yml:31-50` uses a manual cross-build loop plus `softprops/action-gh-release@v2.2.2`. Reconcile the release documentation; the 2026-11-24 roadmap rewrite already removed this claim from `ROADMAP.md`.
- `docs/releases/v1.1.0/RELEASE-NOTES.md`, `git log --oneline -10`, and release run `30564856094` reference v1.1.0; the 2026-11-24 roadmap rewrite removed the stale "Current Release" v1.0 label. Other docs should be checked for the same drift.
- `CONTRIBUTING.md:269` requires 90%+ coverage for `internal/`; measured values include 0.0% (`internal/plugins`, `internal/sync`), 19.0% (`internal/runner`), and 46.5% (`internal/config`). Close the gap; never silently lower it.
- `gofmt -l .` lists 21 Go files; sampled diffs (e.g. `main.go`, `internal/branch/branch.go`) show missing final newline, and `.golangci.yml` does not enable `gofmt`, so lint passes. Ask before bulk reformatting.
- `SECURITY.md:5-17` is an unfilled template listing versions 5.1.x/5.0.x/4.0.x, unrelated to the v1.x releases in `docs/releases/`.
- No root `PLAN.md` exists; confirm whether planning belongs only in `ROADMAP.md`.

## Verification loop

1. **Plan** — state the intended change, the docs to leave untouched, and the gates that will prove success.
2. **Execute** — derive every claim from real file, command, `gh`, or manifest output.
3. **Verify** — run `make lint`, `make vet`, `make test`, `make build`, the CI-form coverage test, and `git status --short`; read actual output.
4. **Correct** — fix each failed check and re-run the full loop, including the full-suite regression.
5. **Stop condition** — all gates exit 0 with no skipped-without-reason tests, CI parity is explicit, only explicitly requested planning docs and no workflows changed, and the working tree contains only intended changes. Otherwise report exactly what is pending or failed.

## Never do

- Do not edit `PLAN.md`, `ROADMAP.md`, `README.md`, `CONTRIBUTING.md`, `USER-GUIDE.md`, other root `.md` files, `.github/workflows/**`, GitHub issues, or release tags without an explicit request in the same session. The project-manager workflow owns `ROADMAP.md`; user approval for the 2026-11-24 rewrite and project 25 sync is recorded above.
- Do not create, close, relabel, or edit issues; `gh` issue/tag writes are prohibited. GitHub Project creation/sync is allowed only with explicit user approval (recorded 2026-11-24 for project 25).
- Do not commit or push until a per-file Conventional Commit plan has been presented and explicitly confirmed; after pushing, watch the triggered runs with `gh run watch --exit-status` / `gh run view` (or match `git ls-remote origin refs/heads/<branch>` to local `HEAD` if no workflow triggers).
- Do not run `make dev`, `make install`, `make uninstall`, `make clean`, `git-policy uninstall --all`, or `git-policy history --clear` without explicit approval — they change global hooks/config or delete artifacts/history (`.github/workflows` not involved; `Makefile:53-57,91`, `USER-GUIDE.md:238-246,814-820`).
- Do not commit build artifacts, coverage output, `.env`, keys, tokens, or secrets; `.gitignore:1-24` covers binaries, `dist/`, and coverage files, and `SecretScan`/`BlockFiles` enforce the rest (`README.md:128-135`).
