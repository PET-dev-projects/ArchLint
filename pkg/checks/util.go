package checks

import "github.com/NovokshanovE/archlint/pkg/model"

type relationRef = model.RelationRef

type relationList map[string][]relationRef

func toStringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

func hasTag(tags []string, allowed map[string]struct{}) bool {
	if len(allowed) == 0 {
		return false
	}
	for _, tag := range tags {
		if _, ok := allowed[tag]; ok {
			return true
		}
	}
	return false
}

func buildOutgoing(m *model.Architecture) relationList {
	res := make(relationList)
	for _, rel := range m.Relations() {
		res[rel.Relation.From] = append(res[rel.Relation.From], rel)
	}
	return res
}
