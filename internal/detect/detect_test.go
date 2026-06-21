package detect

import (
	"encoding/json"
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

func TestFindPreservesBackslashesInsideEnvAssignmentValue(t *testing.T) {
	value := `maskara\test-secret-value`
	content := "client_secret=\"" + value + "\"\n"

	findings := Find(content, "history.log", "claude", DefaultRules())
	if len(findings) != 1 {
		t.Fatalf("expected one env finding, got %d", len(findings))
	}

	got := content[findings[0].Start:findings[0].End]
	if got != value {
		t.Fatalf("expected backslash-preserving range %q, got %q", value, got)
	}
}

func TestFindStopsBoundedRulesBeforeJSONEscapedQuote(t *testing.T) {
	databaseURL := "postgres" + "://user:redacted@example.com/db"
	envName := "api" + "_key"
	envValue := "maskara-test-redactable-value"
	tests := []struct {
		name    string
		message string
		ruleID  string
		want    string
	}{
		{
			name:    "database URL",
			message: `dsn="` + databaseURL + `"`,
			ruleID:  "database-url",
			want:    databaseURL,
		},
		{
			name:    "env value",
			message: envName + `=` + envValue + `"`,
			ruleID:  "env-secret",
			want:    envValue,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := json.Marshal(map[string]string{"message": test.message})
			if err != nil {
				t.Fatal(err)
			}

			findings := Find(string(encoded), "session.jsonl", "codex", DefaultRules())
			var got string
			for _, finding := range findings {
				if finding.RuleID == test.ruleID {
					got = string(encoded)[finding.Start:finding.End]
					break
				}
			}
			if got != test.want {
				t.Fatalf("expected %q range %q, got %q", test.ruleID, test.want, got)
			}
		})
	}
}
