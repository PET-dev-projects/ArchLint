package model_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/NovokshanovE/archlint/pkg/model"
)

func TestValidateModel(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		m := mustLoadModel(t, "../testdata/arch_valid.yaml")
		findings := model.ValidateModel(m)
		if len(findings) != 0 {
			t.Fatalf("expected no findings, got %v", findings)
		}
	})

	t.Run("invalid version and relation", func(t *testing.T) {
		m := &model.Architecture{
			Version: 99,
			Boundaries: []model.Boundary{{
				Name: "bad",
				Containers: []model.Container{{
					Name: "svc",
					Type: model.ContainerService,
				}},
				Relations: []model.Relation{{
					From: "svc",
					To:   "unknown",
					Kind: model.RelationKindSync,
				}},
			}},
		}
		findings := model.ValidateModel(m)
		if len(findings) < 2 {
			t.Fatalf("expected at least 2 findings, got %v", findings)
		}
	})
}

func mustLoadModel(t *testing.T, rel string) *model.Architecture {
	t.Helper()
	path := filepath.Join("..", rel)
	fh, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer fh.Close()
	m, err := model.LoadModelFromYAML(fh)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}
	return m
}
