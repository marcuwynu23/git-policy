// Package history handles storing, logging, and querying git‑policy run history
package history

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

// RuleResult represents the result of a single policy rule
type RuleResult struct {
	Rule    string `json:"rule"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Record represents a single git-policy run history entry
type Record struct {
	Timestamp string       `json:"timestamp"`
	Repo      string       `json:"repo"`
	Branch    string       `json:"branch"`
	Commit    string       `json:"commit,omitempty"`
	Results   []RuleResult `json:"results"`
	Overall   string       `json:"overall"`
}

// Log records a history entry for a git-policy run
func Log(cfg *config.Config, configPath string, results []policy.Result) error {
	if !cfg.History.Enabled {
		return nil
	}
	repoPath, err := getRepoPath()
	if err != nil {
		return err
	}
	branch, err := git.GetBranchName()
	if err != nil {
		return err
	}
	commit, err := getShortCommitHash()
	if err != nil {
		commit = ""
	}
	ruleResults := make([]RuleResult, len(results))
	overall := "pass"
	for i, r := range results {
		status := "pass"
		if r.Status == policy.StatusFail {
			status = "fail"
			overall = "fail"
		}
		ruleResults[i] = RuleResult{
			Rule:    r.PolicyName,
			Status:  status,
			Message: r.Message,
		}
	}
	record := Record{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Repo:      repoPath,
		Branch:    branch,
		Commit:    commit,
		Results:   ruleResults,
		Overall:   overall,
	}
	return appendRecord(cfg, configPath, repoPath, record)
}

// QueryOptions defines filtering/querying options
type QueryOptions struct {
	Limit     int
	RepoPath  string
	Status    string
}

// Query retrieves history records
func Query(cfg *config.Config, configPath string, opts QueryOptions) ([]Record, error) {
	if !cfg.History.Enabled {
		return nil, nil
	}
	targetRepo := opts.RepoPath
	if targetRepo == "" {
		var err error
		targetRepo, err = getRepoPath()
		if err != nil {
			return nil, err
		}
	}
	lines, err := readRecords(cfg, configPath, targetRepo)
	if err != nil {
		return nil, err
	}
	var records []Record
	for _, line := range lines {
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			continue
		}
		if opts.Status != "" && rec.Overall != opts.Status {
			continue
		}
		records = append(records, rec)
	}
	// reverse to get most recent first
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
	if opts.Limit > 0 && len(records) > opts.Limit {
		records = records[:opts.Limit]
	}
	return records, nil
}

// Clear clears history records
func Clear(cfg *config.Config, configPath string, repoPath string) error {
	if repoPath == "" {
		var err error
		repoPath, err = getRepoPath()
		if err != nil {
			return err
		}
	}
	historyFile := getHistoryFilePath(configPath, repoPath)
	return os.Remove(historyFile)
}

func getRepoPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func getShortCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getHistoryFilePath(configPath string, repoPath string) string {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		absRepoPath = repoPath
	}
	hash := sha256.Sum256([]byte(absRepoPath))
	hashStr := hex.EncodeToString(hash[:])
	historyDir := config.HistoryDir(configPath)
	return filepath.Join(historyDir, hashStr+".jsonl")
}

func appendRecord(cfg *config.Config, configPath string, repoPath string, record Record) error {
	historyFile := getHistoryFilePath(configPath, repoPath)
	historyDir := filepath.Dir(historyFile)
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	lines, err := readAllRecords(historyFile)
	if err != nil {
		return err
	}
	lines = append(lines, data)
	if len(lines) > cfg.History.MaxRecords {
		lines = lines[len(lines)-cfg.History.MaxRecords:]
	}
	var strLines []string
	for _, l := range lines {
		strLines = append(strLines, string(l))
	}
	return os.WriteFile(historyFile, []byte(strings.Join(strLines, "\n")+"\n"), 0644)
}

func readRecords(cfg *config.Config, configPath string, repoPath string) ([][]byte, error) {
	historyFile := getHistoryFilePath(configPath, repoPath)
	return readAllRecords(historyFile)
}

func readAllRecords(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()
	lines := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 {
			lines = append(lines, append([]byte{}, line...))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
