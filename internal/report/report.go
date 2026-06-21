package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mrgoonie/maskara/internal/detect"
	"github.com/mrgoonie/maskara/internal/redact"
	"github.com/mrgoonie/maskara/internal/scanner"
)

type Document struct {
	Result    scanner.Result `json:"result"`
	Redaction redact.Summary `json:"redaction"`
}

func Write(writer io.Writer, result scanner.Result, redaction redact.Summary, jsonMode bool) error {
	if jsonMode {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(Document{Result: result, Redaction: redaction})
	}
	_, err := writer.Write([]byte(Markdown(result, redaction)))
	return err
}

func WriteOutput(output string, result scanner.Result, redaction redact.Summary, jsonMode bool) (string, error) {
	path, err := resolveOutputPath(output, jsonMode)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := Write(file, result, redaction, jsonMode); err != nil {
		return "", err
	}
	return path, nil
}

func Markdown(result scanner.Result, redaction redact.Summary) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# Maskara Secret Exposure Report\n\n")
	fmt.Fprintf(&builder, "- Generated: `%s`\n", result.GeneratedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(&builder, "- Files scanned: `%d`\n", result.FilesScanned)
	fmt.Fprintf(&builder, "- Files skipped: `%d`\n", result.FilesSkipped)
	fmt.Fprintf(&builder, "- Findings: `%d`\n", len(result.Findings))
	fmt.Fprintf(&builder, "- Redacted: `%d`\n\n", redaction.Replaced)

	if len(result.Targets) > 0 {
		builder.WriteString("## Scan Targets\n\n")
		builder.WriteString("| Agent | Root |\n|---|---|\n")
		for _, target := range result.Targets {
			fmt.Fprintf(&builder, "| `%s` | `%s` |\n", target.Agent, escapePipes(target.Root))
		}
		builder.WriteString("\n")
	}

	if len(result.Warnings) > 0 {
		builder.WriteString("## Warnings\n\n")
		for _, warning := range result.Warnings {
			fmt.Fprintf(&builder, "- %s\n", warning)
		}
		builder.WriteString("\n")
	}

	if len(result.Findings) == 0 {
		builder.WriteString("## Findings\n\nNo sensitive values detected.\n")
		return builder.String()
	}

	builder.WriteString("## Summary By Rule\n\n")
	builder.WriteString("| Rule | Count |\n|---|---:|\n")
	for _, entry := range summarizeRules(result.Findings) {
		fmt.Fprintf(&builder, "| %s | %d |\n", entry.Name, entry.Count)
	}
	builder.WriteString("\n")

	builder.WriteString("## Findings\n\n")
	builder.WriteString("| Agent | File | Line | Rule | Severity | Masked Preview | SHA-256 |\n")
	builder.WriteString("|---|---|---:|---|---|---|---|\n")
	for _, finding := range result.Findings {
		fmt.Fprintf(
			&builder,
			"| `%s` | `%s` | %d | %s | `%s` | `%s` | `%s` |\n",
			finding.Agent,
			escapePipes(finding.File),
			finding.Line,
			escapePipes(finding.RuleName),
			finding.Severity,
			finding.Preview,
			finding.SHA256[:12],
		)
	}
	builder.WriteString("\n## Rotation Guidance\n\n")
	builder.WriteString("Rotate every credential listed above. Redaction removes local copies from agent logs, but it cannot revoke credentials already shared with a provider or remote service.\n")
	if redaction.Replaced > 0 {
		builder.WriteString("\n## Redaction Backups\n\n")
		for _, file := range redaction.Files {
			fmt.Fprintf(&builder, "- `%s` -> backup `%s`\n", file.Path, file.BackupPath)
		}
	}
	return builder.String()
}

type ruleSummary struct {
	Name  string
	Count int
}

func summarizeRules(findings []detect.Finding) []ruleSummary {
	counts := map[string]int{}
	for _, finding := range findings {
		counts[finding.RuleName]++
	}
	summaries := make([]ruleSummary, 0, len(counts))
	for name, count := range counts {
		summaries = append(summaries, ruleSummary{Name: name, Count: count})
	}
	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].Count == summaries[j].Count {
			return summaries[i].Name < summaries[j].Name
		}
		return summaries[i].Count > summaries[j].Count
	})
	return summaries
}

func resolveOutputPath(output string, jsonMode bool) (string, error) {
	extension := ".md"
	if jsonMode {
		extension = ".json"
	}
	if output == "" {
		output = "."
	}
	info, err := os.Stat(output)
	if err == nil && info.IsDir() {
		return filepath.Join(output, "maskara-report"+extension), nil
	}
	if err != nil && strings.HasSuffix(output, string(os.PathSeparator)) {
		return filepath.Join(output, "maskara-report"+extension), nil
	}
	if filepath.Ext(output) == "" {
		return filepath.Join(output, "maskara-report"+extension), nil
	}
	return output, nil
}

func escapePipes(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}
