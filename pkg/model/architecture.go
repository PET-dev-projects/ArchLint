package model

// Architecture represents the full architecture YAML document.
type Architecture struct {
	Version    int         `yaml:"version"`
	Boundaries []Boundary  `yaml:"boundaries"`
	Externals  []Container `yaml:"externals,omitempty"`
	Meta       Metadata    `yaml:"meta,omitempty"`
}

// Metadata allows attaching arbitrary key/value pairs.
type Metadata map[string]string

// Boundary is a logical grouping of containers and relations.
type Boundary struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
	Owner       string      `yaml:"owner,omitempty"`
	Containers  []Container `yaml:"containers"`
	Boundaries  []Boundary  `yaml:"boundaries,omitempty"`
	Relations   []Relation  `yaml:"relations,omitempty"`
	Meta        Metadata    `yaml:"meta,omitempty"`
}

// ContainerType enumerates supported container kinds.
type ContainerType string

const (
	ContainerService  ContainerType = "service"
	ContainerDatabase ContainerType = "database"
	ContainerExternal ContainerType = "external"
)

// Container models a service, database, or external dependency.
type Container struct {
	Name        string        `yaml:"name"`
	Type        ContainerType `yaml:"type"`
	Description string        `yaml:"description,omitempty"`
	Owner       string        `yaml:"owner,omitempty"`
	Technology  string        `yaml:"technology,omitempty"`
	Protocol    string        `yaml:"protocol,omitempty"`
	Tags        []string      `yaml:"tags,omitempty"`
	Meta        Metadata      `yaml:"meta,omitempty"`
}

// RelationKind enumerates supported relation kinds.
type RelationKind string

const (
	RelationKindSync  RelationKind = "sync"
	RelationKindAsync RelationKind = "async"
	RelationKindDB    RelationKind = "db"
)

// Relation describes a dependency between two containers.
type Relation struct {
	From        string       `yaml:"from"`
	To          string       `yaml:"to"`
	Kind        RelationKind `yaml:"kind"`
	Description string       `yaml:"description,omitempty"`
	Protocol    string       `yaml:"protocol,omitempty"`
	Tags        []string     `yaml:"tags,omitempty"`
	Meta        Metadata     `yaml:"meta,omitempty"`
}
