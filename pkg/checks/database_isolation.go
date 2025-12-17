package checks

import (
	"fmt"

	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

const databaseIsolationRuleID = "ARCH-DB-ISOLATION"

type databaseIsolationRule struct{}

type databaseIsolationConfig struct {
	RequireInbound bool `json:"requireInbound"`
}

var defaultDatabaseIsolationConfig = databaseIsolationConfig{
	RequireInbound: true,
}

// NewDatabaseIsolationRule ensures databases remain passive dependencies.
func NewDatabaseIsolationRule() Rule { return &databaseIsolationRule{} }

func (r *databaseIsolationRule) ID() string { return databaseIsolationRuleID }

func (r *databaseIsolationRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultDatabaseIsolationConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(databaseIsolationRuleID, err)}
	}

	containerIndex := m.ContainerMap()
	inbound := map[string]int{}

	findings := make([]types.Finding, 0)
	for _, relRef := range m.Relations() {
		rel := relRef.Relation
		fromRef, okFrom := containerIndex[rel.From]
		if okFrom && fromRef.Container.Type == model.ContainerDatabase {
			findings = append(findings, types.Finding{
				RuleID:   databaseIsolationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("database %s must not initiate relations", rel.From),
				Path:     relRef.Path,
			})
		}
		toRef, okTo := containerIndex[rel.To]
		if okTo && toRef.Container.Type == model.ContainerDatabase {
			inbound[rel.To]++
		}
	}

	if conf.RequireInbound {
		for _, ref := range m.Containers() {
			if ref.Container.Type != model.ContainerDatabase {
				continue
			}
			if inbound[ref.Container.Name] == 0 {
				findings = append(findings, types.Finding{
					RuleID:   databaseIsolationRuleID,
					Severity: types.SeverityWarn,
					Message:  fmt.Sprintf("database %s has no inbound relations", ref.Container.Name),
					Path:     ref.Path,
				})
			}
		}
	}

	return findings
}
