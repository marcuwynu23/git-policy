# Git Policy User Guide

## Table of Contents

1. [Installation](#installation)
2. [Setup](#setup)
3. [Configuration](#configuration)
4. [Rule Management](#rule-management)
5. [Testing Rules](#testing-rules)
6. [Workflows](#workflows)
7. [Troubleshooting](#troubleshooting)
8. [FAQ](#faq)

---

## Installation

### From Source (Development)

```bash
git clone https://github.com/marcuwynu23/git-policy
cd git-policy
make dev
```

This runs three steps automatically:
1. **Build** the binary
2. **Copy** it to a PATH directory (`C:\Bin\tools` on Windows)
3. **Install** global Git hooks

### Pre-built Binary

Download from the [releases page](https://github.com/marcuwynu23/git-policy/releases),
place the binary somewhere in your PATH, then:

```bash
git-policy install
```

### Verify Installation

```bash
git-policy doctor
```

Expected output:
```
Running git-policy doctor...

PASS  Git found: git version 2.36.1.windows.1
PASS  Global hooks installed

All checks passed.
```

---

## Setup

### 1. Install Hooks

```bash
git-policy install
```

This:
- Creates `~/.config/git-policy/hooks/` (or `%APPDATA%\git-policy\hooks\` on Windows)
- Writes `pre-commit`, `pre-push`, `commit-msg`, `post-merge` hook scripts
- Sets Git global `core.hooksPath` to point at that directory

After this, **every repository** on your machine uses these hooks.

### 2. Configure Rules

```bash
# View current configuration
git-policy validate

# See active rules
git-policy rule list
```

Edit the config file manually:

- **Linux/macOS:** `~/.config/git-policy/config.yaml`
- **Windows:** `%APPDATA%\git-policy\config.yaml`

```yaml
policies:
  blockFiles:
    - ".env"
    - "*.pem"
    - "*.key"
  maxFileSize: "10MB"
  secretScan: true
  protectedBranches:
    - main
    - production
  conventionalCommits: true
```

### 3. Disable a Rule (No Config Editing)

No need to edit YAML — use the CLI:

```bash
git-policy rule disable secret-scan
git-policy rule enable secret-scan
```

---

## Configuration

### File Location

| OS      | Path                                    |
|---------|-----------------------------------------|
| Windows | `%APPDATA%\git-policy\config.yaml`      |
| macOS   | `~/.config/git-policy/config.yaml`      |
| Linux   | `~/.config/git-policy/config.yaml`      |

### Custom Config Path

```bash
git-policy --config /path/to/config.yaml run
```

### Full Reference

```yaml
version: 1

hooks:
  pre-commit:
    enabled: true       # Run on git commit
  commit-msg:
    enabled: true       # Validate commit message
  pre-push:
    enabled: true       # Run on git push
  post-merge:
    enabled: false      # Run on git merge

policies:
  # Files that should never be committed
  blockFiles:
    - ".env"
    - "*.pem"
    - "*.key"
    - "*.p12"
    - "*.crt"

  # Maximum file size before rejection
  maxFileSize: "10MB"    # Supports KB, MB, GB

  # Scan for API keys and tokens
  secretScan: true

  # Branches that block direct commits
  protectedBranches:
    - main
    - master
    - production

  # Require conventional commit format
  conventionalCommits: true

  # Binary extensions to block
  blockBinaries:
    - ".exe"
    - ".dll"
    - ".so"
    - ".iso"
    - ".zip"
```

---

## Rule Management

### List Rules

```bash
git-policy rule list
```

Output:
```
Policies:
  binary-file          enabled
  block-files          enabled
  branch-protection    enabled
  commit-message       enabled
  file-size            enabled
  secret-scan          enabled
```

### Enable / Disable

```bash
git-policy rule disable secret-scan
git-policy rule enable secret-scan
```

Available rule names: `block-files`, `commit-message`, `file-size`, `binary-file`, `secret-scan`, `branch-protection`

### Run Rules Manually

```bash
cd /path/to/repo
git-policy run
```

This is what the hooks execute. Useful for testing without committing.

---

## Testing Rules

### Test BlockFiles

```bash
cd /path/to/repo
echo "DB_PASSWORD=secret" > .env
git add .env
git commit -m "test"
```

Expected:
```
BLOCKED: BlockFiles - .env detected
  Fix: Remove .env from staging or add it to .gitignore
Commit blocked by git-policy.
```

### Test Conventional Commits

```bash
git add main.go
git commit -m "fixed bug"
```

Expected:
```
BLOCKED: CommitMessage - Commit message does not follow conventional commits: fixed bug
  Fix: Use one of: feat:, fix:, refactor:, docs:, ...
```

Passing:
```bash
git commit -m "fix: resolve login timeout"
# All policies passed.
```

### Test Secret Scan

```bash
echo "AWS_KEY=AKIA1234567890" > config.txt
git add config.txt
git commit -m "feat: add config"
```

Expected:
```
BLOCKED: SecretScan - Potential AWS Access Key ID found in config.txt on line 1
```

### Test Branch Protection

```bash
git checkout -b main
echo "test" > test.txt
git add test.txt
git commit -m "feat: test"
```

Expected:
```
BLOCKED: BranchProtection - Commits to protected branch 'main' are not allowed
```

### Test File Size

Create a file larger than 10MB:

```bash
# Windows (creates a 15MB file)
fsutil file createnew large.bin 15728640
git add large.bin
git commit -m "feat: add large file"
```

### Test Binary Files

```bash
echo "dummy" > program.exe
git add program.exe
git commit -m "feat: add binary"
```

---

## Workflows

### Solo Developer

```bash
# One-time setup
git-policy install

# Day-to-day — git-policy enforces rules automatically
git add .
git commit -m "feat: implement login"
# If you try to commit a .env or bad message, it blocks
```

### Team Setup

1. **One person** creates the config:
   ```bash
   # Write the team's rules
   # ~/.config/git-policy/config.yaml
   ```

2. **Share** the config file in a Git repo or internal wiki.

3. **Each teammate** runs:
   ```bash
   git-policy install
   ```

4. **Sync** config updates (future):
   ```bash
   git-policy sync
   ```

### CI/CD Integration

While git-policy is designed for local development, you can also run it in CI:

```yaml
# .github/workflows/ci.yml
jobs:
  policy-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install git-policy
        run: |
          wget https://github.com/marcuwynu23/git-policy/releases/download/v0.1.0/git-policy-linux-amd64
          chmod +x git-policy-linux-amd64
          sudo mv git-policy-linux-amd64 /usr/local/bin/git-policy
      - name: Check rules
        run: git-policy run
```

---

## Troubleshooting

### "git-policy: not found"

The binary is not in your PATH.

```bash
# Development: rebuild and copy
cd git-policy
make dev

# Or manually copy the binary to a PATH directory
cp git-policy /usr/local/bin/          # Linux/macOS
copy git-policy C:\Bin\tools\          # Windows
```

### Hook not running

```bash
# Check if hooks are installed
git-policy doctor

# Reinstall
git-policy uninstall
git-policy install
```

### Rule blocked something incorrectly

```bash
# Check what was blocked
git-policy run

# Temporarily disable the rule
git-policy rule disable block-files

# After commit, re-enable
git-policy rule enable block-files
```

### "Not a git repository"

Run `git-policy run` from inside a Git repository.

```bash
cd /path/to/your/repo
git-policy run
```

### Config changes not taking effect

git-policy reads the config file on every run. Verify:

```bash
# Validate your config
git-policy validate

# Your config file
cat ~/.config/git-policy/config.yaml
```

---

## FAQ

### Does it work on Windows?

Yes. git-policy is tested on Windows, macOS, and Linux.
On Windows, hooks run via Git Bash (shipped with Git for Windows).

### Does it slow down git commits?

No. Rules run in-process with zero external dependencies.
A typical `git commit` takes < 50ms with all 6 rules enabled.

### Can I use it alongside Husky / Lefthook?

Yes. git-policy sets `core.hooksPath` globally. If a repo has its own
hooks directory (e.g., `.husky/`), you can override it per-repo with:

```bash
git config core.hooksPath .husky
```

### Can I write custom rules?

Not yet. Custom rules via the plugin system are planned for v2.
Currently you can extend the built-in rules via config
(e.g., add more blocked file patterns, secret patterns, etc.).

### How do I uninstall completely?

```bash
git-policy uninstall --all
```

This removes hooks, unsets `core.hooksPath`, and deletes `~/.config/git-policy/`.

### Is there a remote sync feature?

Planned for v2. You'll be able to sync rules from a Git repo,
HTTP endpoint, S3, or GCS — enabling team-wide rule distribution.

### Can I use different configs for different projects?

Not yet. git-policy uses a single global config. Per-repo overrides
are a planned feature. For now, use `--config` to point at different files:

```bash
git-policy --config ~/team-project-policy.yaml run
```
