package checks

import (
	"fmt"

	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

const acyclicRuleID = "ARCH-ACYCLIC"

type acyclicRule struct{}

type acyclicConfig struct {
	AllowedKinds     []model.RelationKind `json:"allowedKinds"`
	IgnoreContainers []string             `json:"ignoreContainers"`
}

var defaultAcyclicConfig = acyclicConfig{
	AllowedKinds: []model.RelationKind{
		model.RelationKindSync,
		model.RelationKindAsync,
		model.RelationKindDB,
	},
}

// NewAcyclicRule returns the default implementation of the Acyclic rule.
func NewAcyclicRule() Rule {
	return &acyclicRule{}
}

func (r *acyclicRule) ID() string { return acyclicRuleID }

func (r *acyclicRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultAcyclicConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(acyclicRuleID, err)}
	}

	allowedKind := map[model.RelationKind]struct{}{}
	for _, kind := range conf.AllowedKinds {
		allowedKind[kind] = struct{}{}
	}
	ignored := map[string]struct{}{}
	for _, name := range conf.IgnoreContainers {
		ignored[name] = struct{}{}
	}

	type edge struct {
		to   string
		path string
	}

	graph := map[string][]edge{}
	for _, relRef := range m.Relations() {
		rel := relRef.Relation
		if _, skip := ignored[rel.From]; skip {
			continue
		}
		if len(allowedKind) > 0 {
			if _, ok := allowedKind[rel.Kind]; !ok {
				continue
			}
		}
		graph[rel.From] = append(graph[rel.From], edge{to: rel.To, path: relRef.Path})
	}

	findings := make([]types.Finding, 0)
	visited := map[string]bool{}
	stack := make([]string, 0)
	stackIdx := map[string]int{}

	var visit func(node string)
	visit = func(node string) {
		visited[node] = true
		stackIdx[node] = len(stack)
		stack = append(stack, node)

		for _, edge := range graph[node] {
			if _, ignore := ignored[edge.to]; ignore {
				continue
			}
			if idx, onStack := stackIdx[edge.to]; onStack {
				cycle := append([]string{}, stack[idx:]...)
				cycle = append(cycle, edge.to)
				findings = append(findings, types.Finding{
					RuleID:   acyclicRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("cycle detected: %v", cycle),
					Path:     edge.path,
					Meta: map[string]any{
						"cycle": cycle,
					},
				})
				continue
			}
			if !visited[edge.to] {
				visit(edge.to)
			}
		}

		stack = stack[:len(stack)-1]
		delete(stackIdx, node)
	}

	for node := range graph {
		if !visited[node] {
			visit(node)
		}
	}

	return findings
}
