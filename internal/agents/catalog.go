package agents

import (
	"os"
	"path/filepath"
	"runtime"
)

type agentSpec struct {
	targets func(home string) []Target
	markers func(root string) []string
}

var agentSpecs = map[string]agentSpec{
	"claude": {
		targets: func(home string) []Target {
			return []Target{{Agent: "claude", Root: filepath.Join(home, ".claude", "projects")}}
		},
		markers: func(root string) []string {
			return []string{filepath.Join(root, ".claude")}
		},
	},
	"codex": {
		targets: func(home string) []Target {
			return []Target{{Agent: "codex", Root: filepath.Join(home, ".codex", "sessions")}}
		},
		markers: func(root string) []string {
			return []string{filepath.Join(root, ".codex")}
		},
	},
	"cursor":         dotAndPlatformSpec("cursor", ".cursor", "Cursor", "cursor"),
	"opencode":       dotAndPlatformSpec("opencode", ".opencode", "opencode", "opencode"),
	"antigravity":    dotAndPlatformSpec("antigravity", ".antigravity", "Antigravity", "antigravity"),
	"kimi":           dotAndPlatformSpec("kimi", ".kimi", "Kimi", "kimi"),
	"droid":          dotAndPlatformSpec("droid", ".droid", "Droid", "droid"),
	"gemini":         dotAndPlatformSpec("gemini", ".gemini", "Gemini", "gemini"),
	"github-copilot": dotAndPlatformSpec("github-copilot", ".github-copilot", "GitHub Copilot", "github-copilot"),
	"hermes":         dotAndPlatformSpec("hermes", ".hermes", "Hermes Agent", "hermes"),
	"openclaw":       dotAndPlatformSpec("openclaw", ".openclaw", "OpenClaw", "openclaw"),
	"kilo":           dotAndPlatformSpec("kilo", ".kilo-code", "Kilo Code", "kilo-code"),
	"kiro":           dotAndPlatformSpec("kiro", ".kiro", "Kiro", "kiro"),
	"pi":             dotAndPlatformSpec("pi", ".pi", "Pi", "pi"),
	"qoder":          dotAndPlatformSpec("qoder", ".qoder", "Qoder", "qoder"),
	"qwen":           dotAndPlatformSpec("qwen", ".qwen", "Qwen Code", "qwen-code"),
	"trae":           dotAndPlatformSpec("trae", ".trae", "Trae", "trae"),
}

func dotAndPlatformSpec(agent, dotDir, appName, xdgName string) agentSpec {
	return agentSpec{
		targets: func(home string) []Target {
			return targetsFromRoots(agent, defaultRoots(home, dotDir, appName, xdgName))
		},
		markers: func(root string) []string {
			return markerRoots(root, dotDir, appName, xdgName)
		},
	}
}

func targetsFromRoots(agent string, roots []string) []Target {
	targets := make([]Target, 0, len(roots))
	for _, root := range uniqueNonEmpty(roots) {
		targets = append(targets, Target{Agent: agent, Root: root})
	}
	return targets
}

func defaultRoots(home, dotDir, appName, xdgName string) []string {
	roots := []string{filepath.Join(home, dotDir)}
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			roots = append(roots, filepath.Join(appData, appName))
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			roots = append(roots, filepath.Join(localAppData, appName))
		}
	case "darwin":
		roots = append(roots, filepath.Join(home, "Library", "Application Support", appName))
	default:
		roots = append(roots, filepath.Join(home, ".config", xdgName), filepath.Join(home, ".local", "share", xdgName))
	}
	return roots
}

func markerRoots(root, dotDir, appName, xdgName string) []string {
	return uniqueNonEmpty([]string{
		filepath.Join(root, dotDir),
		filepath.Join(root, ".config", xdgName),
		filepath.Join(root, ".local", "share", xdgName),
		filepath.Join(root, "Library", "Application Support", appName),
		filepath.Join(root, "AppData", "Roaming", appName),
		filepath.Join(root, "AppData", "Local", appName),
	})
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
