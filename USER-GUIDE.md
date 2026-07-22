# Git Policy User Guide

## Table of Contents

1. [Installation](#installation)
2. [Setup](#setup)
3. [Configuration](#configuration)
4. [Commands at a Glance](#commands-at-a-glance)
5. [Rule Management](#rule-management)
6. [Plugin System](#plugin-system)
7. [Testing Rules](#testing-rules)
8. [Workflows](#workflows)
9. [Troubleshooting](#troubleshooting)
10. [FAQ](#faq)

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

### Plugin Reference in Config

Plugins are referenced by their descriptor file path (managed via `plugins install`):

```yaml
plugins:
  - name: my-custom-rules
    path: /home/user/.config/git-policy/plugins/my-custom-rules.yaml
    enabled: true
```

The path is stored as an absolute path and is managed by the `plugins` subcommands —
you rarely edit this section manually.

### Custom Rules in Config

Rules defined via `rule add` are stored directly in `config.yaml`:

```yaml
customRules:
  - name: no-todo
    type: file-content
    pattern: "TODO:"
    message: "Commits containing TODO are not allowed"
    fix: "Resolve the TODO before committing"
  - name: no-zip-files
    type: file-block
    pattern: "*.zip"
    message: "Zip files are not allowed"
```

This section is managed by the `rule add`, `rule remove`, `rule import`, and `rule export` commands.

---

## Commands at a Glance

| Command                                   | Description                                 |
|-------------------------------------------|---------------------------------------------|
| `git-policy install`                      | Install global Git hooks                    |
| `git-policy uninstall`                    | Remove global Git hooks                     |
| `git-policy run`                          | Execute all policies against staged files   |
| `git-policy doctor`                       | Verify installation and hooks               |
| `git-policy validate`                     | Show current config and validate it         |
| `git-policy version`                      | Print version information                   |
| `git-policy rule list`                    | List all rules with enabled/disabled status |
| `git-policy rule enable <name>`           | Enable a rule                               |
| `git-policy rule disable <name>`          | Disable a rule                              |
| `git-policy rule skip [name...]`          | Skip rules for the current commit           |
| `git-policy rule skip --list`             | Show currently skipped rules                |
| `git-policy rule skip --clear`            | Clear all skipped rules                     |
| `git-policy rule add <name>`              | Add a custom rule (--type, --pattern, --message, --fix) |
| `git-policy rule remove <name>`           | Remove a custom rule from config            |
| `git-policy rule export <name>`           | Export a custom rule as a YAML file         |
| `git-policy rule import <file>`           | Import a custom rule from a YAML file       |
| `git-policy plugins install <file>`       | Install a plugin from a YAML descriptor     |
| `git-policy plugins install --disabled`   | Install with rules disabled by default      |
| `git-policy plugins uninstall <name>`     | Remove a plugin and its descriptor file     |
| `git-policy plugins list`                 | List installed plugins with paths           |

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

### Skip Rules Temporarily

When a rule blocks a legitimate commit, skip it for that one commit:

```bash
# Skip one or more rules
git-policy rule skip block-files
git-policy rule skip block-files secret-scan

# List currently skipped rules
git-policy rule skip --list

# Clear all skipped rules
git-policy rule skip --clear
```

Skipped rules are stored in the repository's local Git config and are
**automatically cleared** after a successful commit. This means you never
accidentally leave rules disabled.

### Run Rules Manually

```bash
cd /path/to/repo
git-policy run
```

This is what the hooks execute. Useful for testing without committing.

### Custom Rules (rule add / remove)

Add a rule directly to `config.yaml` without creating a plugin file:

```bash
git-policy rule add no-todo \
  --type file-content \
  --pattern "TODO:" \
  --message "Commits containing TODO are not allowed" \
  --fix "Resolve the TODO before committing"
```

Remove a rule:

```bash
git-policy rule remove no-todo
```

List all custom rules by viewing the config or the full rule list:

```bash
git-policy validate
```

Custom rules are stored in `config.yaml` under `customRules:` and run alongside
built-in rules on every commit. They can be skipped with `custom:<name>`:

```bash
git-policy rule skip custom:no-todo
```

### Export / Import Rules

Share a custom rule as a standalone YAML file:

```bash
# Export a rule
git-policy rule export no-todo -o ./my-rule.yaml

# Import on another machine
git-policy rule import ./my-rule.yaml
```

The exported file is a single rule definition that can be version-controlled
and shared with your team.

---

## Plugin System

Plugins let you define **custom rules** in YAML without writing Go code.
They are loaded from descriptor files stored alongside your config.

### Plugin Descriptor Format

Create a `.yaml` file with your custom rules:

```yaml
name: my-custom-rules
rules:
  - name: no-todo
    type: file-content
    pattern: "TODO:"
    message: "Commits containing TODO are not allowed"
    fix: "Resolve the TODO before committing"
  - name: no-zip-files
    type: file-block
    pattern: "*.zip"
    message: "Zip files are not allowed"
    fix: "Remove zip files from the commit"
  - name: no-draft-branches
    type: branch-name
    pattern: "draft-*"
    message: "Commits to draft branches are blocked"
    fix: "Switch to a non-draft branch"
  - name: no-wip-commits
    type: commit-message
    pattern: "WIP:*"
    message: "WIP commits are not allowed"
    fix: "Write a proper commit message"
```

### Supported Rule Types

| Type             | Description                              | Pattern matching     |
|------------------|------------------------------------------|----------------------|
| `file-block`     | Block files matching a glob pattern      | `filepath.Match`     |
| `file-content`   | Scan file contents for a string pattern  | `strings.Contains`   |
| `branch-name`    | Block commits to branches matching glob  | `filepath.Match`     |
| `commit-message` | Block commits with messages matching glob| `filepath.Match`     |

Each rule requires:
- **name** — unique identifier within the plugin
- **type** — one of the four types above
- **pattern** — the glob or text pattern to match
- **message** — error message shown when the rule blocks a commit
- **fix** (optional) — suggested fix shown to the user

### Install a Plugin

```bash
git-policy plugins install ./my-custom-rules.yaml
```

This:
1. Validates the descriptor file and all rules
2. Copies the file to `~/.config/git-policy/plugins/my-custom-rules.yaml`
3. Adds an entry to your config with the absolute path
4. Enables all rules by default

Install with rules disabled:

```bash
git-policy plugins install --disabled ./my-custom-rules.yaml
```

### List Installed Plugins

```bash
git-policy plugins list
```

Output:
```
Installed plugins:
  my-custom-rules      enabled   /home/user/.config/git-policy/plugins/my-custom-rules.yaml
```

### Uninstall a Plugin

```bash
git-policy plugins uninstall my-custom-rules
```

This:
1. Removes the plugin entry from the config
2. Deletes the descriptor file from the plugins directory

### Custom Rules at Runtime

Once installed, custom rules are executed alongside built-in rules on every
`git commit` or `git-policy run`. They appear with a `Custom:` prefix:

```
BLOCKED: Custom:no-todo - Commits containing TODO are not allowed
  Fix: Resolve the TODO before committing
```

You can skip a custom rule just like any built-in rule:

```bash
git-policy rule skip custom:no-todo
```

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

### Test a Custom Plugin Rule

```bash
# Create a plugin with a file-block rule
cat > test-plugin.yaml << 'EOF'
name: test-plugin
rules:
  - name: no-zips
    type: file-block
    pattern: "*.zip"
    message: "Zip files are blocked by custom plugin"
    fix: "Remove the zip file"
EOF

# Install it
git-policy plugins install test-plugin.yaml

# Test it
echo "test" > archive.zip
git add archive.zip
git commit -m "feat: add archive"
```

Expected:
```
BLOCKED: Custom:no-zips - Zip files are blocked by custom plugin
  Fix: Remove the zip file
Commit blocked by git-policy.
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

2. **Share** the config file and plugin descriptor files in a Git repo
   or internal wiki.

3. **Each teammate** runs:
   ```bash
   git-policy install
   git-policy plugins install ./team-rules.yaml
   ```

4. **Sync** config updates (future):
   ```bash
   git-policy sync
   ```

### Using Plugins Across a Team

Create a shared plugin descriptor and distribute it via version control:

1. Maintain a `git-policy-plugins/` directory in your team's repository:
   ```
   team-repo/
   ├── plugins/
   │   └── team-rules.yaml
   ├── docs/
   └── src/
   ```

2. Each developer installs it:
   ```bash
   git clone https://github.com/team/team-repo
   git-policy plugins install ./plugins/team-rules.yaml
   ```

3. To update, replace the file and re-run install (it overwrites).

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

# Temporarily skip the rule for this commit
git-policy rule skip block-files

# It auto-clears after a successful commit
```

### Plugin installation fails

```bash
# Ensure the descriptor file is valid YAML
git-policy plugins install ./my-plugin.yaml

# Common errors:
#   "name is required"     — add `name: my-plugin` to the descriptor
#   "at least one rule..." — add at least one rule under `rules:`
#   "type is required"     — each rule needs a type (file-block, file-content, etc.)
#   "pattern is required"  — each rule needs a pattern to match
#   "message is required"  — each rule needs an error message
```

### Plugin not loaded at runtime

```bash
# Check if the plugin is installed and enabled
git-policy plugins list

# Verify the descriptor file exists at the stored path
cat ~/.config/git-policy/plugins/my-plugin.yaml

# Reinstall if needed
git-policy plugins uninstall my-plugin
git-policy plugins install ./my-plugin.yaml
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
A typical `git commit` takes < 50ms with all built-in rules enabled.

### Can I use it alongside Husky / Lefthook?

Yes. git-policy sets `core.hooksPath` globally. If a repo has its own
hooks directory (e.g., `.husky/`), you can override it per-repo with:

```bash
git config core.hooksPath .husky
```

### Can I write custom rules?

Yes. Use the plugin system to define custom rules in YAML:

```bash
git-policy plugins install ./my-rules.yaml
```

See the [Plugin System](#plugin-system) section for details on all
supported rule types and the descriptor format.

### Are plugins cross-platform?

Yes. Plugin descriptors are plain YAML files. The plugin system works
identically on Windows, macOS, and Linux. Plugin files are stored
alongside your config and loaded at runtime.

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
