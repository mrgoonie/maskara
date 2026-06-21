package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteVersionFlags(t *testing.T) {
	for _, args := range [][]string{
		{"version"},
		{"--version"},
		{"-v"},
		{"-V"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute(args, &stdout, &stderr)
			if code != exitClean {
				t.Fatalf("expected exit %d, got %d; stderr=%q", exitClean, code, stderr.String())
			}
			if !strings.HasPrefix(stdout.String(), "maskara ") {
				t.Fatalf("unexpected version output: %q", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestExecuteHelpListsExpandedAgents(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"--help"}, &stdout, &stderr)
	if code != exitClean {
		t.Fatalf("expected exit %d, got %d; stderr=%q", exitClean, code, stderr.String())
	}
	output := stdout.String()
	for _, agent := range []string{"cursor", "gemini", "github-copilot", "qwen", "trae"} {
		if !strings.Contains(output, agent) {
			t.Fatalf("expected help to contain %q; output=%q", agent, output)
		}
	}
}

func TestExecuteScanAcceptsExpandedAgentAliasWithRoot(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"scan", "--agent", "gemini-cli", "--root", t.TempDir()}, &stdout, &stderr)
	if code != exitClean {
		t.Fatalf("expected exit %d, got %d; stderr=%q", exitClean, code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "found 0 findings") {
		t.Fatalf("unexpected scan output: %q", stdout.String())
	}
}

func TestExecuteGuardrailsDryRunAcceptsExpandedAgent(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"guardrails", "-a", "qwen-code", "--root", t.TempDir(), "--dry-run"}, &stdout, &stderr)
	if code != exitClean {
		t.Fatalf("expected exit %d, got %d; stderr=%q", exitClean, code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "would-create") {
		t.Fatalf("unexpected guardrails output: %q", stdout.String())
	}
}
