package checks

import (
	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

// Rule describes a reusable architecture rule implementation.
type Rule interface {
	ID() string
	Run(*model.Architecture, map[string]any) []types.Finding
}
