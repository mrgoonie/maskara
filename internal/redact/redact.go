package redact

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mrgoonie/maskara/internal/agents"
	"github.com/mrgoonie/maskara/internal/detect"
	"github.com/mrgoonie/maskara/internal/scanner"
)

type Summary struct {
	Files    []FileSummary `json:"files,omitempty"`
	Replaced int           `json:"replaced"`
	Skipped  int           `json:"skipped"`
}

type FileSummary struct {
	Path       string `json:"path"`
	BackupPath string `json:"backup_path"`
	Replaced   int    `json:"replaced"`
}

func Apply(result scanner.Result) (Summary, error) {
	grouped := map[string][]detect.Finding{}
	for _, finding := range result.Findings {
		grouped[finding.File] = append(grouped[finding.File], finding)
	}

	var summary Summary
	for path, findings := range grouped {
		fileSummary, err := redactFile(path, findings, result.Targets)
		if err != nil {
			return summary, err
		}
		if fileSummary.Replaced == 0 {
			summary.Skipped++
			continue
		}
		summary.Replaced += fileSummary.Replaced
		summary.Files = append(summary.Files, fileSummary)
	}
	sort.Slice(summary.Files, func(i, j int) bool {
		return summary.Files[i].Path < summary.Files[j].Path
	})
	return summary, nil
}

func redactFile(path string, findings []detect.Finding, targets []agents.Target) (FileSummary, error) {
	if !insideTargets(path, targets) {
		return FileSummary{}, errors.New("refusing to redact outside scan targets: " + path)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return FileSummary{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return FileSummary{}, errors.New("refusing to redact symlink: " + path)
	}

	original, err := os.ReadFile(path)
	if err != nil {
		return FileSummary{}, err
	}
	rewritten := make([]byte, len(original))
	copy(rewritten, original)

	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Start > findings[j].Start
	})
	lastStart := len(original) + 1
	replaced := 0
	for _, finding := range findings {
		if finding.Start < 0 || finding.End > len(original) || finding.Start >= finding.End {
			continue
		}
		if finding.End > lastStart {
			continue
		}
		replacement := []byte(finding.Redaction)
		rewritten = append(rewritten[:finding.Start], append(replacement, rewritten[finding.End:]...)...)
		lastStart = finding.Start
		replaced++
	}
	if replaced == 0 || bytes.Equal(original, rewritten) {
		return FileSummary{Path: path}, nil
	}
	if err := validateStructured(path, original, rewritten); err != nil {
		return FileSummary{}, err
	}

	backupPath := backupName(path)
	if err := os.WriteFile(backupPath, original, 0o600); err != nil {
		return FileSummary{}, err
	}
	if err := writeReplace(path, rewritten, info.Mode().Perm()); err != nil {
		return FileSummary{}, err
	}
	return FileSummary{Path: path, BackupPath: backupPath, Replaced: replaced}, nil
}

func insideTargets(path string, targets []agents.Target) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	for _, target := range targets {
		absRoot, err := filepath.Abs(target.Root)
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(absRoot, absPath)
		if err != nil {
			continue
		}
		if rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))) {
			return true
		}
	}
	return false
}
