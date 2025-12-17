package model

import "fmt"

// ContainerRef exposes discovery metadata for a container.
type ContainerRef struct {
	Container    *Container
	Path         string
	Boundary     *Boundary
	BoundaryPath string
}

// RelationRef exposes discovery metadata for a relation.
type RelationRef struct {
	Relation     *Relation
	Path         string
	Boundary     *Boundary
	BoundaryPath string
}

// Containers returns every container declared in the architecture, including externals.
func (a *Architecture) Containers() []ContainerRef {
	refs := make([]ContainerRef, 0)
	for i := range a.Boundaries {
		gatherBoundaryContainers(&refs, &a.Boundaries[i], fmt.Sprintf("boundaries[%d]", i))
	}
	for i := range a.Externals {
		refs = append(refs, ContainerRef{
			Container: &a.Externals[i],
			Path:      fmt.Sprintf("externals[%d]", i),
		})
	}
	return refs
}

// ContainerMap returns a map keyed by container name.
func (a *Architecture) ContainerMap() map[string]ContainerRef {
	indexed := make(map[string]ContainerRef)
	for _, ref := range a.Containers() {
		if ref.Container.Name == "" {
			continue
		}
		indexed[ref.Container.Name] = ref
	}
	return indexed
}

// Relations returns all relations declared within boundaries.
func (a *Architecture) Relations() []RelationRef {
	refs := make([]RelationRef, 0)
	for i := range a.Boundaries {
		gatherBoundaryRelations(&refs, &a.Boundaries[i], fmt.Sprintf("boundaries[%d]", i))
	}
	return refs
}

func gatherBoundaryContainers(dst *[]ContainerRef, b *Boundary, path string) {
	for i := range b.Containers {
		*dst = append(*dst, ContainerRef{
			Container:    &b.Containers[i],
			Path:         fmt.Sprintf("%s.containers[%d]", path, i),
			Boundary:     b,
			BoundaryPath: path,
		})
	}
	for i := range b.Boundaries {
		gatherBoundaryContainers(dst, &b.Boundaries[i], fmt.Sprintf("%s.boundaries[%d]", path, i))
	}
}

func gatherBoundaryRelations(dst *[]RelationRef, b *Boundary, path string) {
	for i := range b.Relations {
		*dst = append(*dst, RelationRef{
			Relation:     &b.Relations[i],
			Path:         fmt.Sprintf("%s.relations[%d]", path, i),
			Boundary:     b,
			BoundaryPath: path,
		})
	}
	for i := range b.Boundaries {
		gatherBoundaryRelations(dst, &b.Boundaries[i], fmt.Sprintf("%s.boundaries[%d]", path, i))
	}
}
