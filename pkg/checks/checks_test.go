package checks_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PET-dev-projects/ArchLint/pkg/checks"
	"github.com/PET-dev-projects/ArchLint/pkg/model"
)

func TestAcyclicRule(t *testing.T) {
	t.Run("valid graph", func(t *testing.T) {
		arch := loadArch(t, "arch_valid.yaml")
		findings := checks.NewAcyclicRule().Run(arch, nil)
		if len(findings) != 0 {
			t.Fatalf("expected no findings, got %v", findings)
		}
	})

	t.Run("cycle detected", func(t *testing.T) {
		arch := loadArch(t, "arch_cycle.yaml")
		findings := checks.NewAcyclicRule().Run(arch, nil)
		if len(findings) == 0 {
			t.Fatal("expected cycle finding")
		}
	})
}

func TestCRUDRule(t *testing.T) {
	arch := loadArch(t, "arch_crud_violation.yaml")
	findings := checks.NewCRUDRule().Run(arch, nil)
	if len(findings) < 2 {
		t.Fatalf("expected multiple findings, got %v", findings)
	}
}

func TestACLRule(t *testing.T) {
	arch := loadArch(t, "arch_acl_violation.yaml")
	findings := checks.NewACLRule().Run(arch, nil)
	if len(findings) == 0 {
		t.Fatalf("expected acl finding, got %v", findings)
	}
}

func TestBoundariesRule(t *testing.T) {
	arch := loadArch(t, "arch_boundary_weak.yaml")
	findings := checks.NewBoundariesRule().Run(arch, nil)
	if len(findings) == 0 {
		t.Fatalf("expected boundary finding")
	}
}

func TestExternalProtocolRule(t *testing.T) {
	arch := loadArch(t, "arch_external_protocol.yaml")
	findings := checks.NewExternalProtocolRule().Run(arch, nil)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %v", findings)
	}
}

func TestDatabaseIsolationRule(t *testing.T) {
	arch := loadArch(t, "arch_db_isolation.yaml")
	findings := checks.NewDatabaseIsolationRule().Run(arch, nil)
	if len(findings) != 2 {
		t.Fatalf("expected outbound + orphan findings, got %v", findings)
	}
}

func loadArch(t *testing.T, name string) *model.Architecture {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", name)
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
