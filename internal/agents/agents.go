package agents

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Target struct {
	Agent string `json:"agent"`
	Root  string `json:"root"`
}

var supportedAgents = []string{"claude", "codex", "opencode", "antigravity"}

func Supported() []string {
	out := make([]string, len(supportedAgents))
	copy(out, supportedAgents)
	return out
}

func Resolve(agentName, root string) ([]Target, error) {
	agent := Normalize(agentName)
	if root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			return nil, err
		}
		if agent == "auto" {
			agent = "custom"
		}
		return []Target{{Agent: agent, Root: abs}}, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if agent == "auto" || agent == "all" {
		return existingOrDefaultTargets(home), nil
	}
	if !isSupported(agent) {
		return nil, errors.New("unsupported agent: " + agentName)
	}
	return defaultTargetsFor(home, agent), nil
}

func Normalize(agent string) string {
	value := strings.ToLower(strings.TrimSpace(agent))
	switch value {
	case "", "auto":
		return "auto"
	case "claude-code", "claudecode":
		return "claude"
	case "open-code":
		return "opencode"
	default:
		return value
	}
}

func existingOrDefaultTargets(home string) []Target {
	var existing []Target
	for _, agent := range supportedAgents {
		for _, target := range defaultTargetsFor(home, agent) {
			if info, err := os.Stat(target.Root); err == nil && info.IsDir() {
				existing = append(existing, target)
			}
		}
	}
	if len(existing) > 0 {
		return existing
	}

	var defaults []Target
	for _, agent := range []string{"claude", "codex"} {
		defaults = append(defaults, defaultTargetsFor(home, agent)...)
	}
	return defaults
}

func defaultTargetsFor(home, agent string) []Target {
	switch agent {
	case "claude":
		return []Target{{Agent: "claude", Root: filepath.Join(home, ".claude", "projects")}}
	case "codex":
		return []Target{{Agent: "codex", Root: filepath.Join(home, ".codex", "sessions")}}
	case "opencode":
		return opencodeTargets(home)
	case "antigravity":
		return antigravityTargets(home)
	default:
		return nil
	}
}

func opencodeTargets(home string) []Target {
	targets := []Target{{Agent: "opencode", Root: filepath.Join(home, ".opencode")}}
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			targets = append(targets, Target{Agent: "opencode", Root: filepath.Join(appData, "opencode")})
		}
	case "darwin":
		targets = append(targets, Target{Agent: "opencode", Root: filepath.Join(home, "Library", "Application Support", "opencode")})
	default:
		targets = append(targets, Target{Agent: "opencode", Root: filepath.Join(home, ".local", "share", "opencode")})
	}
	return targets
}

func antigravityTargets(home string) []Target {
	targets := []Target{{Agent: "antigravity", Root: filepath.Join(home, ".antigravity")}}
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			targets = append(targets, Target{Agent: "antigravity", Root: filepath.Join(appData, "Antigravity")})
		}
	case "darwin":
		targets = append(targets, Target{Agent: "antigravity", Root: filepath.Join(home, "Library", "Application Support", "Antigravity")})
	default:
		targets = append(targets, Target{Agent: "antigravity", Root: filepath.Join(home, ".config", "antigravity")})
	}
	return targets
}

func isSupported(agent string) bool {
	for _, supported := range supportedAgents {
		if agent == supported {
			return true
		}
	}
	return false
}
