package checks

import (
	"fmt"

	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

const aclRuleID = "ARCH-ACL"

type aclRule struct{}

type aclConfig struct {
	AllowedTags []string `json:"allowedTags"`
}

var defaultACLConfig = aclConfig{
	AllowedTags: []string{"acl"},
}

// NewACLRule creates ACL rule implementation.
func NewACLRule() Rule { return &aclRule{} }

func (r *aclRule) ID() string { return aclRuleID }

func (r *aclRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultACLConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(aclRuleID, err)}
	}

	allowed := toStringSet(conf.AllowedTags)
	containerIndex := m.ContainerMap()

	findings := make([]types.Finding, 0)
	for _, relRef := range m.Relations() {
		rel := relRef.Relation
		from, okFrom := containerIndex[rel.From]
		to, okTo := containerIndex[rel.To]
		if !okFrom || !okTo {
			continue
		}
		if to.Container.Type == model.ContainerExternal {
			if !hasTag(from.Container.Tags, allowed) {
				findings = append(findings, types.Finding{
					RuleID:   aclRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("container %s must declare one of %v to talk to external %s", from.Container.Name, conf.AllowedTags, to.Container.Name),
					Path:     relRef.Path,
				})
			}
		}
	}

	return findings
}
