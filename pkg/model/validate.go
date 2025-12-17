package model

import (
	"fmt"
	"strings"

	"github.com/NovokshanovE/archlint/pkg/types"
)

const validationRuleID = "MODEL-0001"

// ValidateModel performs structural validation and returns findings.
func ValidateModel(m *Architecture) []types.Finding {
	if m == nil {
		return []types.Finding{{
			RuleID:   validationRuleID,
			Severity: types.SeverityError,
			Message:  "model is nil",
			Path:     "$",
		}}
	}

	findings := make([]types.Finding, 0)

	if m.Version != 1 {
		findings = append(findings, types.Finding{
			RuleID:   validationRuleID,
			Severity: types.SeverityError,
			Message:  fmt.Sprintf("unsupported version %d (only version 1 is supported)", m.Version),
			Path:     "version",
		})
	}

	if len(m.Boundaries) == 0 {
		findings = append(findings, types.Finding{
			RuleID:   validationRuleID,
			Severity: types.SeverityError,
			Message:  "at least one boundary is required",
			Path:     "boundaries",
		})
	}

	containers := m.Containers()
	nameIndex := map[string]string{}
	for _, ref := range containers {
		c := ref.Container
		if strings.TrimSpace(c.Name) == "" {
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  "container name is required",
				Path:     ref.Path + ".name",
			})
		}
		if _, ok := nameIndex[c.Name]; ok {
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("duplicate container name %q", c.Name),
				Path:     ref.Path + ".name",
				Meta: map[string]any{
					"container": c.Name,
				},
			})
		} else if c.Name != "" {
			nameIndex[c.Name] = ref.Path
		}

		switch c.Type {
		case ContainerService, ContainerDatabase, ContainerExternal:
		default:
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("invalid container type %q", c.Type),
				Path:     ref.Path + ".type",
			})
		}
	}

	relations := m.Relations()
	for _, ref := range relations {
		rel := ref.Relation
		if rel.From == "" || rel.To == "" {
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  "relation must define both from and to",
				Path:     ref.Path,
			})
			continue
		}
		if _, ok := nameIndex[rel.From]; !ok {
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("relation references unknown container %q", rel.From),
				Path:     ref.Path + ".from",
			})
		}
		if _, ok := nameIndex[rel.To]; !ok {
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("relation references unknown container %q", rel.To),
				Path:     ref.Path + ".to",
			})
		}
		switch rel.Kind {
		case RelationKindSync, RelationKindAsync, RelationKindDB:
		default:
			findings = append(findings, types.Finding{
				RuleID:   validationRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("invalid relation kind %q", rel.Kind),
				Path:     ref.Path + ".kind",
			})
		}
	}

	return findings
}
