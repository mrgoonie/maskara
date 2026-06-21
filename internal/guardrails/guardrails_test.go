package guardrails

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	changes, err := Install(Options{Agent: "codex", Root: root, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected planned changes")
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatal("dry run wrote files")
	}
}

func TestInstallCodexWritesSkillAndBacksUpExistingAgentsFile(t *testing.T) {
	root := t.TempDir()
	agentsPath := filepath.Join(root, ".codex", "AGENTS.md")
	if err := os.MkdirAll(filepath.Dir(agentsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agentsPath, []byte("# Existing\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	changes, err := Install(Options{Agent: "codex", Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected changes")
	}
	if _, err := os.Stat(agentsPath + ".maskara.bak"); err != nil {
		t.Fatalf("expected backup: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "skills", "maskara-privacy", "SKILL.md")); err != nil {
		t.Fatalf("expected skill: %v", err)
	}
}
