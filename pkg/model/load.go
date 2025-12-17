package model

import (
	"io"

	"gopkg.in/yaml.v3"
)

// LoadModelFromYAML parses an architecture definition from YAML.
func LoadModelFromYAML(r io.Reader) (*Architecture, error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	var arch Architecture
	if err := dec.Decode(&arch); err != nil {
		return nil, err
	}
	return &arch, nil
}
