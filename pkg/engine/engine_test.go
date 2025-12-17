package engine_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PET-dev-projects/ArchLint/pkg/engine"
	"github.com/PET-dev-projects/ArchLint/pkg/model"
)

func TestRunAllDefaultRules(t *testing.T) {
	arch := loadArch(t, "arch_valid.yaml")
	findings := engine.RunAll(arch, engine.Options{})
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestRunAllEnabledRules(t *testing.T) {
	arch := loadArch(t, "arch_acl_violation.yaml")
	findings := engine.RunAll(arch, engine.Options{EnabledRules: []string{"ARCH-ACL"}})
	if len(findings) == 0 {
		t.Fatal("expected acl finding")
	}
	for _, f := range findings {
		if f.RuleID != "ARCH-ACL" {
			t.Fatalf("unexpected rule id: %s", f.RuleID)
		}
	}
}

func TestRunAllRuleConfig(t *testing.T) {
	arch := loadArch(t, "arch_boundary_weak.yaml")
	opts := engine.Options{
		RuleConfig: map[string]map[string]any{
			"ARCH-BOUNDARIES": {
				"minInternalToCrossRatio": 0.1,
				"maxCrossRelations":       0,
			},
		},
	}
	findings := engine.RunAll(arch, opts)
	if len(findings) == 0 {
		t.Fatal("expected findings because ratio still below threshold")
	}
}

func TestRunAllMusicStreamingExample(t *testing.T) {
	arch := loadExample(t, "music_streaming.yaml")
	findings := engine.RunAll(arch, engine.Options{})
	if len(findings) != 0 {
		t.Fatalf("expected no findings for music streaming example, got %v", findings)
	}
}

func loadArch(t *testing.T, name string) *model.Architecture {
	t.Helper()
	return loadArchFromDir(t, filepath.Join("..", "..", "testdata"), name)
}

func loadExample(t *testing.T, name string) *model.Architecture {
	t.Helper()
	return loadArchFromDir(t, filepath.Join("..", "..", "examples"), name)
}

func loadArchFromDir(t *testing.T, dir, name string) *model.Architecture {
	t.Helper()
	path := filepath.Join(dir, name)
	fh, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer fh.Close()
	arch, err := model.LoadModelFromYAML(fh)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}
	return arch
}
