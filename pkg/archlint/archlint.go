package archlint

import (
	"io"

	"github.com/NovokshanovE/archlint/pkg/engine"
	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

// Options alias to engine Options for convenience.
type Options = engine.Options

// LoadModelFromYAML reads architecture YAML.
func LoadModelFromYAML(r io.Reader) (*model.Architecture, error) {
	return model.LoadModelFromYAML(r)
}

// ValidateModel runs structural validation.
func ValidateModel(m *model.Architecture) []types.Finding {
	return model.ValidateModel(m)
}

// RunAll executes all enabled checks.
func RunAll(m *model.Architecture, opts Options) []types.Finding {
	return engine.RunAll(m, opts)
}
