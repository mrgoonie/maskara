package cli

import (
	"fmt"
	"io"

	"github.com/mrgoonie/maskara/internal/agents"
)

func printHelp(writer io.Writer) {
	fmt.Fprintf(writer, `maskara scans coding-agent session logs for sensitive information.

Usage:
  maskara                         scan, redact with backups, and write a report
  maskara --version | -v           print version information
  maskara scan [flags]             scan only
  maskara report [flags]           scan and write a report
  maskara guardrails [flags]       install privacy guardrails for coding agents

Common flags:
  -a, --agent <name>      scan target agent
      --root <path>       override scan or config root
  -o, --output <path>     report file or directory
      --json              write JSON report

Agent names:
  %s

Examples:
  maskara scan -a codex
  maskara report --json -o ./maskara-report.json
  maskara guardrails -a codex --dry-run
`, agents.HelpList())
}

func printScanHelp(writer io.Writer) {
	fmt.Fprintf(writer, `Usage:
  maskara scan [flags]

Scan coding-agent logs and print a summary. Does not redact.

Flags:
  -a, --agent <name>      scan target agent
      --root <path>       scan root override
  -o, --output <path>     optional report file or directory
      --json              write JSON when --output is set

Agent names:
  %s
`, agents.HelpList())
}

func printReportHelp(writer io.Writer) {
	fmt.Fprintf(writer, `Usage:
  maskara report [flags]

Scan coding-agent logs and write a Markdown or JSON report. Does not redact.

Flags:
  -a, --agent <name>      scan target agent
      --root <path>       scan root override
  -o, --output <path>     report file or directory
      --json              write JSON report

Agent names:
  %s
`, agents.HelpList())
}

func printGuardrailsHelp(writer io.Writer) {
	fmt.Fprintf(writer, `Usage:
  maskara guardrails [flags]

Install local instructions, skills, and hook scripts for coding agents.

Flags:
  -a, --agent <name>      target agent
      --root <path>       home/config root override
      --dry-run           show planned files without writing

Agent names:
  %s
`, agents.HelpList())
}
