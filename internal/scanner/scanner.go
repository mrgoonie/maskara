package scanner

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mrgoonie/maskara/internal/agents"
	"github.com/mrgoonie/maskara/internal/detect"
)

const defaultMaxFileBytes int64 = 25 * 1024 * 1024

type Options struct {
	Targets      []agents.Target
	MaxFileBytes int64
}

type Result struct {
	GeneratedAt  time.Time        `json:"generated_at"`
	Targets      []agents.Target  `json:"targets"`
	Findings     []detect.Finding `json:"findings"`
	Warnings     []string         `json:"warnings,omitempty"`
	FilesScanned int              `json:"files_scanned"`
	FilesSkipped int              `json:"files_skipped"`
}

func Scan(options Options) (Result, error) {
	if options.MaxFileBytes <= 0 {
		options.MaxFileBytes = defaultMaxFileBytes
	}
	result := Result{
		GeneratedAt: time.Now().UTC(),
		Targets:     options.Targets,
	}
	rules := detect.DefaultRules()

	for _, target := range options.Targets {
		info, err := os.Stat(target.Root)
		if err != nil {
			result.Warnings = append(result.Warnings, "missing root: "+target.Root)
			continue
		}
		if !info.IsDir() {
			findings, scanned, skipped := scanFile(target.Agent, target.Root, options.MaxFileBytes, rules)
			result.Findings = append(result.Findings, findings...)
			result.FilesScanned += scanned
			result.FilesSkipped += skipped
			continue
		}
		err = filepath.WalkDir(target.Root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				result.Warnings = append(result.Warnings, walkErr.Error())
				return nil
			}
			if entry.IsDir() {
				if shouldSkipDir(entry.Name()) && path != target.Root {
					return filepath.SkipDir
				}
				return nil
			}
			findings, scanned, skipped := scanFile(target.Agent, path, options.MaxFileBytes, rules)
			result.Findings = append(result.Findings, findings...)
			result.FilesScanned += scanned
			result.FilesSkipped += skipped
			return nil
		})
		if err != nil {
			return result, err
		}
	}

	sort.SliceStable(result.Findings, func(i, j int) bool {
		if result.Findings[i].File == result.Findings[j].File {
			return result.Findings[i].Start < result.Findings[j].Start
		}
		return result.Findings[i].File < result.Findings[j].File
	})
	return result, nil
}

func scanFile(agent, path string, maxBytes int64, rules []detect.Rule) ([]detect.Finding, int, int) {
	info, err := os.Lstat(path)
	if err != nil || info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, 0, 1
	}
	if info.Size() > maxBytes || !looksLikeSessionText(path) {
		return nil, 0, 1
	}
	content, err := os.ReadFile(path)
	if err != nil || isBinary(content) {
		return nil, 0, 1
	}
	return detect.Find(string(content), path, agent, rules), 1, 0
}

func shouldSkipDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", "node_modules", ".venv", "venv", "target", "dist", "build", ".next", "__pycache__":
		return true
	default:
		return false
	}
}

func looksLikeSessionText(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".jsonl", ".md", ".txt", ".log", ".yaml", ".yml", ".toml", ".env":
		return true
	}
	if strings.HasPrefix(name, ".env") {
		return true
	}
	return strings.Contains(name, "session") ||
		strings.Contains(name, "conversation") ||
		strings.Contains(name, "transcript") ||
		strings.Contains(name, "history")
}

func isBinary(content []byte) bool {
	sample := content
	if len(sample) > 8192 {
		sample = sample[:8192]
	}
	return bytes.IndexByte(sample, 0) >= 0
}
