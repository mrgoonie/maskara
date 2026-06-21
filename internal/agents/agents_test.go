package agents

import (
	"path/filepath"
	"testing"
)

func TestNormalizeRequestedAgentAliases(t *testing.T) {
	tests := map[string]string{
		"Cursor":             "cursor",
		"Antigravity CLI":    "antigravity",
		"Kimi Code CLI":      "kimi",
		"Droid":              "droid",
		"Gemini CLI":         "gemini",
		"GitHub Copilot":     "github-copilot",
		"Hermes Agent":       "hermes",
		"OpenClaw":           "openclaw",
		"Kilo Code":          "kilo",
		"Kiro CLI":           "kiro",
		"Pi CLI":             "pi",
		"Qoder":              "qoder",
		"Qwen Code":          "qwen",
		"Trae":               "trae",
		"claude-code":        "claude",
		"open-code":          "opencode",
		"github-copilot-cli": "github-copilot",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := Normalize(input); got != want {
				t.Fatalf("Normalize(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestSupportedIncludesRequestedAgents(t *testing.T) {
	for _, agent := range []string{
		"cursor",
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
	} {
		t.Run(agent, func(t *testing.T) {
			if !IsSupported(agent) {
				t.Fatalf("expected %q to be supported", agent)
			}
			if len(defaultTargetsFor(t.TempDir(), agent)) == 0 {
				t.Fatalf("expected default targets for %q", agent)
			}
			if PrimaryConfigRoot(agent, t.TempDir()) == "" {
				t.Fatalf("expected primary config root for %q", agent)
			}
		})
	}
}

func TestResolveAllWithCustomRootScansOnce(t *testing.T) {
	root := t.TempDir()
	targets, err := Resolve("all", root)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected one target, got %d", len(targets))
	}
	if targets[0].Root != root && targets[0].Root != filepath.Clean(root) {
		t.Fatalf("expected target root %q, got %q", root, targets[0].Root)
	}
	if targets[0].Agent != "all" {
		t.Fatalf("expected agent all, got %q", targets[0].Agent)
	}
}

func TestResolveCustomRootPreservesExplicitAgentLabels(t *testing.T) {
	root := t.TempDir()
	targets, err := Resolve("gemini-cli", root)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected one target, got %d", len(targets))
	}
	if targets[0].Agent != "gemini" {
		t.Fatalf("expected normalized agent gemini, got %q", targets[0].Agent)
	}
	if targets[0].Root == "" {
		t.Fatal("expected absolute target root")
	}
}

func TestResolveUnsupportedAgentWithoutRootErrors(t *testing.T) {
	_, err := Resolve("unknown-agent", "")
	if err == nil {
		t.Fatal("expected unsupported agent error")
	}
}
