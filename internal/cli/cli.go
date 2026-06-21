package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/mrgoonie/maskara/internal/agents"
	"github.com/mrgoonie/maskara/internal/guardrails"
	"github.com/mrgoonie/maskara/internal/redact"
	"github.com/mrgoonie/maskara/internal/report"
	"github.com/mrgoonie/maskara/internal/scanner"
	"github.com/mrgoonie/maskara/internal/version"
)

const (
	exitClean    = 0
	exitFindings = 1
	exitError    = 2
)

type commonOptions struct {
	agent  string
	root   string
	output string
	json   bool
}

func Execute(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "help", "-h", "--help":
			printHelp(stdout)
			return exitClean
		case "version", "--version", "-V":
			fmt.Fprintf(stdout, "maskara %s (%s, %s)\n", version.Version, version.Commit, version.Date)
			return exitClean
		case "scan":
			if wantsHelp(args[1:]) {
				printScanHelp(stdout)
				return exitClean
			}
			return runScan(args[1:], stdout, stderr)
		case "report":
			if wantsHelp(args[1:]) {
				printReportHelp(stdout)
				return exitClean
			}
			return runReport(args[1:], stdout, stderr)
		case "guardrails":
			if wantsHelp(args[1:]) {
				printGuardrailsHelp(stdout)
				return exitClean
			}
			return runGuardrails(args[1:], stdout, stderr)
		}
	}
	return runFull(args, stdout, stderr)
}

func runFull(args []string, stdout, stderr io.Writer) int {
	options, err := parseCommon("maskara", args, true)
	if err != nil {
		return fail(stderr, err)
	}
	result, err := scan(options)
	if err != nil {
		return fail(stderr, err)
	}
	var redaction redact.Summary
	if len(result.Findings) > 0 {
		redaction, err = redact.Apply(result)
		if err != nil {
			return fail(stderr, err)
		}
	}
	path, err := report.WriteOutput(options.output, result, redaction, options.json)
	if err != nil {
		return fail(stderr, err)
	}
	fmt.Fprintf(stdout, "maskara: scanned %d files, found %d findings, redacted %d values\n", result.FilesScanned, len(result.Findings), redaction.Replaced)
	fmt.Fprintf(stdout, "maskara: report written to %s\n", path)
	if len(result.Findings) > 0 {
		return exitFindings
	}
	return exitClean
}

func runScan(args []string, stdout, stderr io.Writer) int {
	options, err := parseCommon("maskara scan", args, false)
	if err != nil {
		return fail(stderr, err)
	}
	result, err := scan(options)
	if err != nil {
		return fail(stderr, err)
	}
	if options.output != "" {
		path, err := report.WriteOutput(options.output, result, redact.Summary{}, options.json)
		if err != nil {
			return fail(stderr, err)
		}
		fmt.Fprintf(stdout, "maskara: scan report written to %s\n", path)
	} else {
		fmt.Fprintf(stdout, "maskara: scanned %d files, found %d findings\n", result.FilesScanned, len(result.Findings))
		for _, warning := range result.Warnings {
			fmt.Fprintf(stderr, "warning: %s\n", warning)
		}
	}
	if len(result.Findings) > 0 {
		return exitFindings
	}
	return exitClean
}

func runReport(args []string, stdout, stderr io.Writer) int {
	options, err := parseCommon("maskara report", args, true)
	if err != nil {
		return fail(stderr, err)
	}
	result, err := scan(options)
	if err != nil {
		return fail(stderr, err)
	}
	path, err := report.WriteOutput(options.output, result, redact.Summary{}, options.json)
	if err != nil {
		return fail(stderr, err)
	}
	fmt.Fprintf(stdout, "maskara: report written to %s\n", path)
	if len(result.Findings) > 0 {
		return exitFindings
	}
	return exitClean
}

func runGuardrails(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("maskara guardrails", flag.ContinueOnError)
	fs.SetOutput(stderr)
	agentName := fs.String("agent", "auto", "agent to install guardrails for: auto, claude, codex, opencode, antigravity")
	fs.StringVar(agentName, "a", "auto", "shorthand for --agent")
	root := fs.String("root", "", "home/config root override, useful for tests")
	dryRun := fs.Bool("dry-run", false, "show planned files without writing")
	if err := fs.Parse(args); err != nil {
		return exitError
	}
	changes, err := guardrails.Install(guardrails.Options{Agent: *agentName, Root: *root, DryRun: *dryRun})
	if err != nil {
		return fail(stderr, err)
	}
	if len(changes) == 0 {
		fmt.Fprintln(stdout, "maskara: no guardrail changes")
		return exitClean
	}
	for _, change := range changes {
		if change.BackupPath != "" {
			fmt.Fprintf(stdout, "%s: %s (backup: %s)\n", change.Action, change.Path, change.BackupPath)
			continue
		}
		fmt.Fprintf(stdout, "%s: %s\n", change.Action, change.Path)
	}
	return exitClean
}

func parseCommon(name string, args []string, defaultOutput bool) (commonOptions, error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	options := commonOptions{agent: "auto"}
	fs.StringVar(&options.agent, "agent", "auto", "agent to scan: auto, claude, codex, opencode, antigravity")
	fs.StringVar(&options.agent, "a", "auto", "shorthand for --agent")
	fs.StringVar(&options.root, "root", "", "scan root override")
	fs.StringVar(&options.output, "output", "", "report output file or directory")
	fs.StringVar(&options.output, "o", "", "shorthand for --output")
	fs.BoolVar(&options.json, "json", false, "write JSON report")
	if err := fs.Parse(args); err != nil {
		return options, err
	}
	if fs.NArg() > 0 {
		return options, errors.New("unknown command or argument: " + strings.Join(fs.Args(), " "))
	}
	if defaultOutput && options.output == "" {
		options.output = "."
	}
	return options, nil
}

func scan(options commonOptions) (scanner.Result, error) {
	targets, err := agents.Resolve(options.agent, options.root)
	if err != nil {
		return scanner.Result{}, err
	}
	return scanner.Scan(scanner.Options{Targets: targets})
}

func fail(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "maskara: %v\n", err)
	return exitError
}

func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}
