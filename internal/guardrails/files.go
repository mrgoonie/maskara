package guardrails

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type filePlan struct {
	path    string
	content string
	mode    os.FileMode
	append  bool
}

func appendPlan(path, content string) filePlan {
	return filePlan{path: path, content: content, mode: 0o644, append: true}
}

func writePlan(path, content string, mode os.FileMode) filePlan {
	return filePlan{path: path, content: content, mode: mode}
}

func installFiles(root string, dryRun bool, plans []filePlan) ([]Change, error) {
	var changes []Change
	for _, plan := range plans {
		if !insideRoot(plan.path, root) {
			return changes, errors.New("refusing to write outside root: " + plan.path)
		}
		change, err := installFile(plan, dryRun)
		if err != nil {
			return changes, err
		}
		changes = append(changes, change)
	}
	return changes, nil
}

func installFile(plan filePlan, dryRun bool) (Change, error) {
	action := "create"
	finalContent := plan.content
	var backupPath string

	existing, err := os.ReadFile(plan.path)
	exists := err == nil
	if exists {
		if plan.append {
			if strings.Contains(string(existing), marker) {
				return Change{Path: plan.path, Action: "unchanged"}, nil
			}
			finalContent = string(existing) + "\n\n" + wrap(plan.content)
			action = "append"
		} else if string(existing) == plan.content {
			return Change{Path: plan.path, Action: "unchanged"}, nil
		} else {
			action = "replace"
		}
		backupPath = plan.path + ".maskara.bak"
		if _, statErr := os.Stat(backupPath); statErr == nil {
			backupPath += "." + time.Now().UTC().Format("20060102150405")
		}
	} else if plan.append {
		finalContent = wrap(plan.content)
	}

	if dryRun {
		return Change{Path: plan.path, Action: "would-" + action, BackupPath: backupPath}, nil
	}
	if err := os.MkdirAll(filepath.Dir(plan.path), 0o755); err != nil {
		return Change{}, err
	}
	if exists {
		if err := os.WriteFile(backupPath, existing, 0o600); err != nil {
			return Change{}, err
		}
	}
	if err := os.WriteFile(plan.path, []byte(finalContent), plan.mode); err != nil {
		return Change{}, err
	}
	return Change{Path: plan.path, Action: action, BackupPath: backupPath}, nil
}

func hookAsset() string {
	if runtime.GOOS == "windows" {
		return "assets/maskara-privacy-hook.ps1"
	}
	return "assets/maskara-privacy-hook.sh"
}

func hookName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".ps1"
	}
	return base + ".sh"
}

func hookMode() os.FileMode {
	if runtime.GOOS == "windows" {
		return 0o644
	}
	return 0o755
}

func mustAsset(name string) string {
	content, err := assetFS.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("missing embedded asset %s: %v", name, err))
	}
	return strings.TrimRight(string(content), "\r\n") + "\n"
}

func wrap(content string) string {
	return "<!-- " + marker + ":START -->\n" + strings.TrimSpace(content) + "\n<!-- " + marker + ":END -->\n"
}

func insideRoot(path, root string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)))
}
