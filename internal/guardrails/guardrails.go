package guardrails

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/mrgoonie/maskara/internal/agents"
)

//go:embed assets/*
var assetFS embed.FS

const marker = "MASKARA-GUARDRAILS"

type Options struct {
	Agent  string
	Root   string
	DryRun bool
}

type Change struct {
	Path       string `json:"path"`
	Action     string `json:"action"`
	BackupPath string `json:"backup_path,omitempty"`
}

func Install(options Options) ([]Change, error) {
	root := options.Root
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		root = home
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	targetAgents := resolveAgents(options.Agent, absRoot)
	var changes []Change
	for _, agent := range targetAgents {
		planned, err := installAgent(agent, absRoot, options.DryRun)
		if err != nil {
			return changes, err
		}
		changes = append(changes, planned...)
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})
	return changes, nil
}

func resolveAgents(agentName, root string) []string {
	agent := agents.Normalize(agentName)
	if agent != "auto" && agent != "all" {
		return []string{agent}
	}

	candidates := agents.Supported()
	var existing []string
	for _, candidate := range candidates {
		for _, path := range agents.RootMarkers(candidate, root) {
			if _, err := os.Stat(path); err == nil {
				existing = append(existing, candidate)
				break
			}
		}
	}
	if len(existing) > 0 {
		return existing
	}
	return []string{"claude", "codex"}
}

func installAgent(agent, root string, dryRun bool) ([]Change, error) {
	switch agent {
	case "claude":
		return installFiles(root, dryRun, []filePlan{
			appendPlan(filepath.Join(root, ".claude", "CLAUDE.md"), mustAsset("assets/guardrail-instructions.md")),
			writePlan(filepath.Join(root, ".claude", "skills", "maskara-privacy", "SKILL.md"), mustAsset("assets/skill.md"), 0o644),
			writePlan(filepath.Join(root, ".claude", "hooks", hookName("maskara-privacy-hook")), mustAsset(hookAsset()), hookMode()),
		})
	case "codex":
		return installFiles(root, dryRun, []filePlan{
			appendPlan(filepath.Join(root, ".codex", "AGENTS.md"), mustAsset("assets/guardrail-instructions.md")),
			writePlan(filepath.Join(root, ".codex", "skills", "maskara-privacy", "SKILL.md"), mustAsset("assets/skill.md"), 0o644),
			writePlan(filepath.Join(root, ".codex", "hooks", hookName("maskara-privacy-hook")), mustAsset(hookAsset()), hookMode()),
		})
	case "opencode":
		return installFiles(root, dryRun, []filePlan{
			appendPlan(filepath.Join(root, ".config", "opencode", "maskara-guardrails.md"), mustAsset("assets/guardrail-instructions.md")),
			writePlan(filepath.Join(root, ".config", "opencode", "hooks", hookName("maskara-privacy-hook")), mustAsset(hookAsset()), hookMode()),
		})
	case "antigravity":
		return installFiles(root, dryRun, []filePlan{
			appendPlan(filepath.Join(root, ".antigravity", "maskara-guardrails.md"), mustAsset("assets/guardrail-instructions.md")),
			writePlan(filepath.Join(root, ".antigravity", "hooks", hookName("maskara-privacy-hook")), mustAsset(hookAsset()), hookMode()),
		})
	default:
		if !agents.IsSupported(agent) {
			return nil, errors.New("unsupported guardrails agent: " + agent)
		}
		configRoot := agents.PrimaryConfigRoot(agent, root)
		if configRoot == "" {
			return nil, errors.New("unsupported guardrails agent: " + agent)
		}
		return installFiles(root, dryRun, []filePlan{
			appendPlan(filepath.Join(configRoot, "maskara-guardrails.md"), mustAsset("assets/guardrail-instructions.md")),
			writePlan(filepath.Join(configRoot, "hooks", hookName("maskara-privacy-hook")), mustAsset(hookAsset()), hookMode()),
		})
	}
}
