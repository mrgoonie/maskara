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
