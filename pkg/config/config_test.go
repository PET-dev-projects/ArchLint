package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PET-dev-projects/ArchLint/pkg/config"
)

func TestLoadOptionsFromFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "rules.yaml")
	content := `
rules:
  - id: ARCH-ACL
    enabled: true
    config:
      allowedTags: ["acl", "edge"]
  - id: ARCH-CRUD
    enabled: false
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	opts, err := config.LoadOptionsFromFile(file)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(opts.EnabledRules) != 1 || opts.EnabledRules[0] != "ARCH-ACL" {
		t.Fatalf("unexpected enabled rules: %+v", opts.EnabledRules)
	}
	if opts.RuleConfig["ARCH-ACL"]["allowedTags"].([]any)[1].(string) != "edge" {
		t.Fatalf("failed to parse config overrides: %+v", opts.RuleConfig)
	}
}
