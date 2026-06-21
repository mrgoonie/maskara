package redact

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrgoonie/maskara/internal/agents"
	"github.com/mrgoonie/maskara/internal/scanner"
)

func TestApplyRedactsJSONLAndCreatesBackup(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "session.jsonl")
	secret := "sk-" + strings.Repeat("b", 24)
	original := []byte(`{"message":"` + secret + `"}` + "\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := scanner.Scan(scanner.Options{Targets: []agents.Target{{Agent: "codex", Root: root}}})
	if err != nil {
		t.Fatal(err)
	}
	summary, err := Apply(result)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Replaced == 0 || len(summary.Files) != 1 {
		t.Fatalf("expected one redacted file, got %+v", summary)
	}
	if _, err := os.Stat(summary.Files[0].BackupPath); err != nil {
		t.Fatalf("expected backup: %v", err)
	}

	redacted, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(redacted), secret) {
		t.Fatal("secret still present after redaction")
	}
	assertJSONLValid(t, redacted)
}

func TestApplyRedactsJSONLWithEscapedDatabaseURLDelimiter(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "session.jsonl")
	databaseURL := "postgres" + "://user:redacted@example.com/db"
	line, err := json.Marshal(map[string]string{
		"message": `dsn="` + databaseURL + `"`,
	})
	if err != nil {
		t.Fatal(err)
	}
	original := append(line, '\n')
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := scanner.Scan(scanner.Options{Targets: []agents.Target{{Agent: "codex", Root: root}}})
	if err != nil {
		t.Fatal(err)
	}
	summary, err := Apply(result)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Replaced != 1 || len(summary.Files) != 1 {
		t.Fatalf("expected one database URL replacement, got %+v", summary)
	}

	redacted, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(redacted), databaseURL) {
		t.Fatal("database URL still present after redaction")
	}
	assertJSONLValid(t, redacted)
}

func assertJSONLValid(t *testing.T, content []byte) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		if !json.Valid([]byte(line)) {
			t.Fatalf("redacted JSONL is invalid: %s", line)
		}
	}
}
