package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/PET-dev-projects/ArchLint/pkg/archlint"
	"github.com/PET-dev-projects/ArchLint/pkg/config"
	"github.com/PET-dev-projects/ArchLint/pkg/engine"
	"github.com/PET-dev-projects/ArchLint/pkg/report"
	"github.com/PET-dev-projects/ArchLint/pkg/types"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "check":
		if err := runCheck(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	file := fs.String("f", "", "path to architecture YAML file")
	format := fs.String("format", "text", "output format: text or json")
	failOn := fs.String("fail-on", "error", "fail on severity: error|warn|info|none")
	configPath := fs.String("config", "", "YAML file describing enabled rules and their configs")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *file == "" {
		return errors.New("-f is required")
	}

	var opts engine.Options
	if *configPath != "" {
		loaded, err := config.LoadOptionsFromFile(*configPath)
		if err != nil {
			return err
		}
		opts = loaded
	}

	fh, err := os.Open(*file)
	if err != nil {
		return err
	}
	defer fh.Close()

	model, err := archlint.LoadModelFromYAML(fh)
	if err != nil {
		return err
	}

	findings := make([]types.Finding, 0)
	findings = append(findings, archlint.ValidateModel(model)...)
	findings = append(findings, archlint.RunAll(model, opts)...)

	switch *format {
	case "text":
		if err := report.WriteText(os.Stdout, findings); err != nil {
			return err
		}
	case "json":
		if err := report.WriteJSON(os.Stdout, findings); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format %s", *format)
	}

	if shouldFail(findings, *failOn) {
		return errors.New("fail-on threshold reached")
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: archlint <command> [options]

Commands:
  check   Run architecture checks

Examples:
  archlint check -f examples/payments.yaml --config configs/rules.yaml
`)
}

func shouldFail(findings []types.Finding, failOn string) bool {
	failOn = strings.ToLower(failOn)
	severityOrder := map[string]int{
		string(types.SeverityInfo):  1,
		string(types.SeverityWarn):  2,
		string(types.SeverityError): 3,
	}

	threshold := severityOrder[failOn]
	if failOn == "none" || threshold == 0 {
		return false
	}

	for _, f := range findings {
		if severityOrder[strings.ToLower(string(f.Severity))] >= threshold {
			return true
		}
	}
	return false
}
