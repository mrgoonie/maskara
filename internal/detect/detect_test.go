package detect

import (
	"strings"
	"testing"
)

func TestFindMasksSecretPreview(t *testing.T) {
	secret := "sk-" + strings.Repeat("a", 24)
	content := "OPENAI_API_KEY=" + secret + "\n"

	findings := Find(content, "session.jsonl", "codex", DefaultRules())
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}

	first := findings[0]
	if first.Preview == secret {
		t.Fatal("preview leaked full secret")
	}
	if first.SHA256 == "" {
		t.Fatal("expected hash")
	}
	if first.Line != 1 {
		t.Fatalf("expected line 1, got %d", first.Line)
	}
}

func TestFindUsesEnvAssignmentValueRange(t *testing.T) {
	value := "maskara-test-secret-value"
	content := "client_secret=\"" + value + "\"\n"

	findings := Find(content, "history.log", "claude", DefaultRules())
	if len(findings) != 1 {
		t.Fatalf("expected one env finding, got %d", len(findings))
	}

	got := content[findings[0].Start:findings[0].End]
	if got != value {
		t.Fatalf("expected value-only range %q, got %q", value, got)
	}
}
