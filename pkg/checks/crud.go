package checks

import (
	"fmt"

	"github.com/PET-dev-projects/ArchLint/pkg/model"
	"github.com/PET-dev-projects/ArchLint/pkg/types"
)

const crudRuleID = "ARCH-CRUD"

type crudRule struct{}

type crudConfig struct {
	AllowedTags   []string `json:"allowedTags"`
	ExclusiveTags []string `json:"exclusiveTags"`
}

var defaultCrudConfig = crudConfig{
	AllowedTags:   []string{"crud", "repo", "relay"},
	ExclusiveTags: []string{"repo"},
}

// NewCRUDRule creates CRUD rule implementation.
func NewCRUDRule() Rule { return &crudRule{} }

func (r *crudRule) ID() string { return crudRuleID }

func (r *crudRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultCrudConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(crudRuleID, err)}
	}

	allowedTag := toStringSet(conf.AllowedTags)
	exclusiveTag := toStringSet(conf.ExclusiveTags)

	containerIndex := m.ContainerMap()
	outgoing := buildOutgoing(m)

	findings := make([]types.Finding, 0)

	for _, relRef := range m.Relations() {
		rel := relRef.Relation
		from, okFrom := containerIndex[rel.From]
		to, okTo := containerIndex[rel.To]
		if !okFrom || !okTo {
			continue
		}
		if to.Container.Type == model.ContainerDatabase {
			if rel.Kind != model.RelationKindDB {
				findings = append(findings, types.Finding{
					RuleID:   crudRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("relation to database %s must use kind 'db'", to.Container.Name),
					Path:     relRef.Path,
				})
			}
			if !hasTag(from.Container.Tags, allowedTag) {
				findings = append(findings, types.Finding{
					RuleID:   crudRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("container %s must declare one of %v to access databases", from.Container.Name, conf.AllowedTags),
					Path:     relRef.Path,
				})
			}
		}
	}

	for _, ref := range m.Containers() {
		if !hasTag(ref.Container.Tags, exclusiveTag) {
			continue
		}
		rels := outgoing[ref.Container.Name]
		for _, relRef := range rels {
			to := containerIndex[relRef.Relation.To]
			if to.Container.Type != model.ContainerDatabase || relRef.Relation.Kind != model.RelationKindDB {
				findings = append(findings, types.Finding{
					RuleID:   crudRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("container %s is restricted to database relations", ref.Container.Name),
					Path:     relRef.Path,
				})
			}
		}
	}

	return findings
}
