package checks

import (
	"github.com/PET-dev-projects/ArchLint/pkg/model"
	"github.com/PET-dev-projects/ArchLint/pkg/types"
)

// Rule describes a reusable architecture rule implementation.
type Rule interface {
	ID() string
	Run(*model.Architecture, map[string]any) []types.Finding
}
