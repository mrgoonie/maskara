package detect

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
)

type Finding struct {
	RuleID    string   `json:"rule_id"`
	RuleName  string   `json:"rule_name"`
	Severity  Severity `json:"severity"`
	Agent     string   `json:"agent"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Column    int      `json:"column"`
	Start     int      `json:"start"`
	End       int      `json:"end"`
	Preview   string   `json:"preview"`
	SHA256    string   `json:"sha256"`
	Redaction string   `json:"redaction"`
}

type Rule struct {
	ID         string
	Name       string
	Severity   Severity
	Expression string
	ValueGroup int

	re *regexp.Regexp
}

func DefaultRules() []Rule {
	specs := []Rule{
		{ID: "openai-api-key", Name: "OpenAI API key", Severity: SeverityCritical, Expression: `\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b`},
		{ID: "anthropic-api-key", Name: "Anthropic API key", Severity: SeverityCritical, Expression: `\bsk-ant-[A-Za-z0-9_-]{20,}\b`},
		{ID: "github-token", Name: "GitHub token", Severity: SeverityCritical, Expression: `\b(?:gh[pousr]_[A-Za-z0-9_]{36,}|github_pat_[A-Za-z0-9_]{20,})\b`},
		{ID: "aws-access-key", Name: "AWS access key ID", Severity: SeverityHigh, Expression: `\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`},
		{ID: "google-api-key", Name: "Google API key", Severity: SeverityHigh, Expression: `\bAIza[0-9A-Za-z_-]{35}\b`},
		{ID: "slack-token", Name: "Slack token", Severity: SeverityHigh, Expression: `\bxox[baprs]-[A-Za-z0-9-]{20,}\b`},
		{ID: "stripe-live-key", Name: "Stripe live key", Severity: SeverityCritical, Expression: `\b(?:sk|rk)_live_[A-Za-z0-9]{16,}\b`},
		{ID: "jwt", Name: "JSON Web Token", Severity: SeverityHigh, Expression: `\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`},
		{ID: "database-url", Name: "Database URL", Severity: SeverityCritical, Expression: `(?i)\b(?:postgres(?:ql)?|mysql|mongodb(?:\+srv)?|redis)://[^\s"'<>` + "`" + `]+`},
		{ID: "private-key", Name: "Private key block", Severity: SeverityCritical, Expression: `(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`},
		{ID: "env-secret", Name: "Secret-like env assignment", Severity: SeverityMedium, Expression: `(?i)\b(?:api[_-]?key|secret|token|password|passwd|pwd|private[_-]?key|client[_-]?secret)\b\s*[:=]\s*["']?([^\s"',` + "`" + `]{8,})`, ValueGroup: 1},
	}
	for index := range specs {
		specs[index].re = regexp.MustCompile(specs[index].Expression)
	}
	return specs
}

func Find(content, file, agent string, rules []Rule) []Finding {
	var findings []Finding
	for _, rule := range rules {
		matches := rule.re.FindAllStringSubmatchIndex(content, -1)
		for _, match := range matches {
			start, end := valueRange(match, rule.ValueGroup)
			if start < 0 || end <= start {
				continue
			}
			value := content[start:end]
			if strings.Contains(value, "MASKARA_REDACTED") {
				continue
			}
			line, column := lineColumn(content, start)
			findings = append(findings, Finding{
				RuleID:    rule.ID,
				RuleName:  rule.Name,
				Severity:  rule.Severity,
				Agent:     agent,
				File:      file,
				Line:      line,
				Column:    column,
				Start:     start,
				End:       end,
				Preview:   Mask(value),
				SHA256:    hash(value),
				Redaction: "[MASKARA_REDACTED:" + rule.ID + "]",
			})
		}
	}
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].File == findings[j].File {
			return findings[i].Start < findings[j].Start
		}
		return findings[i].File < findings[j].File
	})
	return dedupeOverlaps(findings)
}

func Mask(value string) string {
	clean := strings.ReplaceAll(value, "\r", "")
	clean = strings.ReplaceAll(clean, "\n", "\\n")
	if len(clean) <= 8 {
		return "***"
	}
	return clean[:4] + "..." + clean[len(clean)-4:]
}

func valueRange(match []int, group int) (int, int) {
	if group > 0 {
		offset := group * 2
		if len(match) > offset+1 && match[offset] >= 0 && match[offset+1] >= 0 {
			return match[offset], match[offset+1]
		}
	}
	return match[0], match[1]
}

func lineColumn(content string, offset int) (int, int) {
	line, column := 1, 1
	for index := 0; index < offset && index < len(content); index++ {
		if content[index] == '\n' {
			line++
			column = 1
			continue
		}
		column++
	}
	return line, column
}

func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func dedupeOverlaps(findings []Finding) []Finding {
	var out []Finding
	for _, finding := range findings {
		overlaps := false
		for _, kept := range out {
			if finding.File == kept.File && rangesOverlap(finding.Start, finding.End, kept.Start, kept.End) {
				overlaps = true
				break
			}
		}
		if !overlaps {
			out = append(out, finding)
		}
	}
	return out
}

func rangesOverlap(aStart, aEnd, bStart, bEnd int) bool {
	return aStart < bEnd && bStart < aEnd
}
