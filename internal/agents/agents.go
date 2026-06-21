package agents

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Target struct {
	Agent string `json:"agent"`
	Root  string `json:"root"`
}

var supportedAgents = []string{
	"claude",
	"codex",
	"cursor",
	"opencode",
	"antigravity",
	"kimi",
	"droid",
	"gemini",
	"github-copilot",
	"hermes",
	"openclaw",
	"kilo",
	"kiro",
	"pi",
	"qoder",
	"qwen",
	"trae",
}

func Supported() []string {
	out := make([]string, len(supportedAgents))
	copy(out, supportedAgents)
	return out
}

func HelpList() string {
	names := append([]string{"auto", "all"}, supportedAgents...)
	return strings.Join(names, ", ")
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
	if agent == "auto" {
		return existingOrDefaultTargets(home), nil
	}
	if agent == "all" {
		return allDefaultTargets(home), nil
	}
	if !IsSupported(agent) {
		return nil, errors.New("unsupported agent: " + agentName)
	}
	return defaultTargetsFor(home, agent), nil
}

func Normalize(agent string) string {
	value := strings.ToLower(strings.TrimSpace(agent))
	value = strings.NewReplacer("_", "-", " ", "-").Replace(value)
	switch value {
	case "", "auto":
		return "auto"
	case "claude-code", "claudecode":
		return "claude"
	case "open-code":
		return "opencode"
	case "antigravity-cli", "antigravity-code":
		return "antigravity"
	case "kimi-code", "kimi-code-cli", "kimi-cli":
		return "kimi"
	case "gemini-cli":
		return "gemini"
	case "github-copilot-cli", "copilot":
		return "github-copilot"
	case "hermes-agent":
		return "hermes"
	case "open-claw", "openclaw-cli":
		return "openclaw"
	case "kilo-code":
		return "kilo"
	case "kiro-cli":
		return "kiro"
	case "pi-cli":
		return "pi"
	case "qwen-code":
		return "qwen"
	case "qoder-cli":
		return "qoder"
	case "trae-cli":
		return "trae"
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

func allDefaultTargets(home string) []Target {
	var targets []Target
	for _, agent := range supportedAgents {
		targets = append(targets, defaultTargetsFor(home, agent)...)
	}
	return targets
}

func defaultTargetsFor(home, agent string) []Target {
	spec, ok := agentSpecs[agent]
	if !ok {
		return nil
	}
	return spec.targets(home)
}

func RootMarkers(agent, root string) []string {
	spec, ok := agentSpecs[agent]
	if !ok {
		return nil
	}
	return spec.markers(root)
}

func PrimaryConfigRoot(agent, root string) string {
	markers := RootMarkers(agent, root)
	if len(markers) == 0 {
		return ""
	}
	return markers[0]
}

func IsSupported(agent string) bool {
	for _, supported := range supportedAgents {
		if agent == supported {
			return true
		}
	}
	return false
}
