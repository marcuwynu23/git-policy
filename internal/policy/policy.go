package policy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
)

type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
)

type Result struct {
	PolicyName string
	Status     Status
	Message    string
	Fix        string
}

type Context struct {
	RepoPath    string
	StagedFiles []string
	CommitMsg   string
	BranchName  string
}

type Policy interface {
	Name() string
	Validate(ctx Context) Result
}

type BlockFilesPolicy struct {
	cfg *config.Config
}

func NewBlockFilesPolicy(cfg *config.Config) *BlockFilesPolicy {
	return &BlockFilesPolicy{cfg: cfg}
}

func (p *BlockFilesPolicy) Name() string {
	return "BlockFiles"
}

func (p *BlockFilesPolicy) Validate(ctx Context) Result {
	for _, file := range ctx.StagedFiles {
		for _, pattern := range p.cfg.Policies.BlockFiles {
			if matched, _ := filepath.Match(pattern, filepath.Base(file)); matched {
				return Result{
					PolicyName: p.Name(),
					Status:     StatusFail,
					Message:    fmt.Sprintf("%s detected: %s", pattern, file),
					Fix:        fmt.Sprintf("Remove %s from staging or add it to .gitignore", file),
				}
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

type CommitMessagePolicy struct {
	cfg *config.Config
}

func NewCommitMessagePolicy(cfg *config.Config) *CommitMessagePolicy {
	return &CommitMessagePolicy{cfg: cfg}
}

func (p *CommitMessagePolicy) Name() string {
	return "CommitMessage"
}

var validPrefixes = []string{
	"feat:", "fix:", "refactor:", "docs:", "test:",
	"build:", "ci:", "style:", "perf:", "chore:", "revert:",
}

func (p *CommitMessagePolicy) Validate(ctx Context) Result {
	if !p.cfg.Policies.ConventionalCommits {
		return Result{PolicyName: p.Name(), Status: StatusPass}
	}
	msg := ctx.CommitMsg
	if msg == "" {
		return Result{PolicyName: p.Name(), Status: StatusPass}
	}
	firstLine := strings.SplitN(msg, "\n", 2)[0]
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(firstLine, prefix) || strings.HasPrefix(firstLine, strings.ToUpper(prefix[:1])+prefix[1:]) {
			return Result{PolicyName: p.Name(), Status: StatusPass}
		}
	}
	return Result{
		PolicyName: p.Name(),
		Status:     StatusFail,
		Message:    fmt.Sprintf("Commit message does not follow conventional commits: %s", firstLine),
		Fix:        fmt.Sprintf("Use one of: %s", strings.Join(validPrefixes, ", ")),
	}
}

type FileSizePolicy struct {
	cfg *config.Config
}

func NewFileSizePolicy(cfg *config.Config) *FileSizePolicy {
	return &FileSizePolicy{cfg: cfg}
}

func (p *FileSizePolicy) Name() string {
	return "FileSize"
}

func parseMaxSize(sizeStr string) (int64, error) {
	if sizeStr == "" {
		return 10 * 1024 * 1024, nil
	}
	var multiplier int64 = 1
	s := strings.ToUpper(strings.TrimSpace(sizeStr))
	switch {
	case strings.HasSuffix(s, "GB"):
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GB")
	case strings.HasSuffix(s, "MB"):
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "MB")
	case strings.HasSuffix(s, "KB"):
		multiplier = 1024
		s = strings.TrimSuffix(s, "KB")
	case strings.HasSuffix(s, "B"):
		multiplier = 1
		s = strings.TrimSuffix(s, "B")
	}
	var size int64
	if _, err := fmt.Sscanf(s, "%d", &size); err != nil {
		return 10 * 1024 * 1024, nil
	}
	return size * multiplier, nil
}

func (p *FileSizePolicy) Validate(ctx Context) Result {
	maxSize, err := parseMaxSize(p.cfg.Policies.MaxFileSize)
	if err != nil {
		return Result{PolicyName: p.Name(), Status: StatusFail, Message: fmt.Sprintf("Invalid max file size: %v", err)}
	}
	for _, file := range ctx.StagedFiles {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.Size() > maxSize {
			return Result{
				PolicyName: p.Name(),
				Status:     StatusFail,
				Message:    fmt.Sprintf("File %s is %d bytes (max %d bytes)", file, info.Size(), maxSize),
				Fix:        fmt.Sprintf("Reduce file size or increase maxFileSize in config"),
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

type BinaryFilePolicy struct {
	cfg *config.Config
}

func NewBinaryFilePolicy(cfg *config.Config) *BinaryFilePolicy {
	return &BinaryFilePolicy{cfg: cfg}
}

func (p *BinaryFilePolicy) Name() string {
	return "BinaryFile"
}

func (p *BinaryFilePolicy) Validate(ctx Context) Result {
	blocked := p.cfg.Policies.BlockBinaries
	if len(blocked) == 0 {
		blocked = []string{".exe", ".dll", ".so", ".iso", ".zip"}
	}
	for _, file := range ctx.StagedFiles {
		ext := strings.ToLower(filepath.Ext(file))
		for _, blockedExt := range blocked {
			if ext == blockedExt {
				return Result{
					PolicyName: p.Name(),
					Status:     StatusFail,
					Message:    fmt.Sprintf("Binary file detected: %s (type: %s)", file, ext),
					Fix:        fmt.Sprintf("Remove %s from the commit or add it to allowed binaries", file),
				}
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

type SecretScanPolicy struct {
	cfg *config.Config
}

func NewSecretScanPolicy(cfg *config.Config) *SecretScanPolicy {
	return &SecretScanPolicy{cfg: cfg}
}

func (p *SecretScanPolicy) Name() string {
	return "SecretScan"
}

var secretPatterns = []struct {
	Name string
	Prefix string
}{
	{"AWS Access Key ID", "AKIA"},
	{"GitHub PAT", "ghp_"},
	{"GitLab PAT", "glpat-"},
	{"Google API Key", "AIza"},
	{"Stripe Key", "sk_live_"},
	{"OpenAI Key", "sk-"},
	{"Slack Token", "xoxb-"},
}

func (p *SecretScanPolicy) Validate(ctx Context) Result {
	if !p.cfg.Policies.SecretScan {
		return Result{PolicyName: p.Name(), Status: StatusPass}
	}
	for _, file := range ctx.StagedFiles {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			for _, pattern := range secretPatterns {
				if strings.Contains(line, pattern.Prefix) {
					return Result{
						PolicyName: p.Name(),
						Status:     StatusFail,
						Message:    fmt.Sprintf("Potential %s found in %s on line %d", pattern.Name, file, lineNum),
						Fix:        fmt.Sprintf("Remove the secret from %s or use environment variables", file),
					}
				}
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

type BranchPolicy struct {
	cfg *config.Config
}

func NewBranchPolicy(cfg *config.Config) *BranchPolicy {
	return &BranchPolicy{cfg: cfg}
}

func (p *BranchPolicy) Name() string {
	return "BranchProtection"
}

func (p *BranchPolicy) Validate(ctx Context) Result {
	protected := p.cfg.Policies.ProtectedBranches
	if len(protected) == 0 {
		return Result{PolicyName: p.Name(), Status: StatusPass}
	}
	for _, branch := range protected {
		if ctx.BranchName == branch {
			return Result{
				PolicyName: p.Name(),
				Status:     StatusFail,
				Message:    fmt.Sprintf("Commits to protected branch '%s' are not allowed", branch),
				Fix:        fmt.Sprintf("Create a feature branch and use a pull request to merge into %s", branch),
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}
