package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrgoonie/maskara/internal/agents"
)

func TestScanFindsSecretsInSessionFiles(t *testing.T) {
	root := t.TempDir()
	secret := "ghp_" + strings.Repeat("a", 36)
	if err := os.WriteFile(filepath.Join(root, "session.jsonl"), []byte(`{"content":"`+secret+`"}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := Scan(Options{Targets: []agents.Target{{Agent: "codex", Root: root}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesScanned != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.FilesScanned)
	}
	if len(result.Findings) == 0 {
		t.Fatal("expected findings")
	}
}

func TestScanSkipsBinaryFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "session.log"), []byte{0, 1, 2}, 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := Scan(Options{Targets: []agents.Target{{Agent: "codex", Root: root}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesScanned != 0 {
		t.Fatalf("expected no scanned files, got %d", result.FilesScanned)
	}
}
