package checks

import (
	"fmt"

	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

const boundariesRuleID = "ARCH-BOUNDARIES"

type boundariesRule struct{}

type boundariesConfig struct {
	MinInternalToCrossRatio float64 `json:"minInternalToCrossRatio"`
	MaxCrossRelations       int     `json:"maxCrossRelations"`
}

var defaultBoundariesConfig = boundariesConfig{
	MinInternalToCrossRatio: 1,
}

// NewBoundariesRule returns the cohesion/coupling rule implementation.
func NewBoundariesRule() Rule { return &boundariesRule{} }

func (r *boundariesRule) ID() string { return boundariesRuleID }

func (r *boundariesRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultBoundariesConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(boundariesRuleID, err)}
	}

	outgoing := buildOutgoing(m)
	metrics := collectBoundaryMetrics(m, outgoing)

	findings := make([]types.Finding, 0)
	for _, metric := range metrics {
		if metric.cross == 0 {
			continue
		}
		ratio := 0.0
		if metric.internal > 0 {
			ratio = float64(metric.internal) / float64(metric.cross)
		}

		if ratio < conf.MinInternalToCrossRatio {
			findings = append(findings, types.Finding{
				RuleID:   boundariesRuleID,
				Severity: types.SeverityWarn,
				Message:  fmt.Sprintf("boundary %s cohesion/coupling ratio %.2f below minimum %.2f", metric.name, ratio, conf.MinInternalToCrossRatio),
				Path:     metric.path,
				Meta: map[string]any{
					"internal": metric.internal,
					"cross":    metric.cross,
					"ratio":    ratio,
				},
			})
		}

		if conf.MaxCrossRelations > 0 && metric.cross > conf.MaxCrossRelations {
			findings = append(findings, types.Finding{
				RuleID:   boundariesRuleID,
				Severity: types.SeverityWarn,
				Message:  fmt.Sprintf("boundary %s has %d cross-boundary relations (max %d)", metric.name, metric.cross, conf.MaxCrossRelations),
				Path:     metric.path,
				Meta: map[string]any{
					"internal": metric.internal,
					"cross":    metric.cross,
				},
			})
		}
	}

	return findings
}

type boundaryMetric struct {
	name       string
	path       string
	containers map[string]struct{}
	internal   int
	cross      int
}

func collectBoundaryMetrics(m *model.Architecture, outgoing relationList) []boundaryMetric {
	metrics := make([]boundaryMetric, 0)
	for idx := range m.Boundaries {
		path := fmt.Sprintf("boundaries[%d]", idx)
		metrics = append(metrics, computeMetrics(&m.Boundaries[idx], path, outgoing)...)
	}
	return metrics
}

func computeMetrics(b *model.Boundary, path string, outgoing relationList) []boundaryMetric {
	current := boundaryMetric{
		name:       b.Name,
		path:       path,
		containers: make(map[string]struct{}),
	}
	for _, c := range b.Containers {
		if c.Name == "" {
			continue
		}
		current.containers[c.Name] = struct{}{}
	}

	metrics := []boundaryMetric{}
	for idx := range b.Boundaries {
		nestedPath := fmt.Sprintf("%s.boundaries[%d]", path, idx)
		nestedMetrics := computeMetrics(&b.Boundaries[idx], nestedPath, outgoing)
		for _, nm := range nestedMetrics {
			for name := range nm.containers {
				current.containers[name] = struct{}{}
			}
		}
		metrics = append(metrics, nestedMetrics...)
	}

	for name := range current.containers {
		for _, rel := range outgoing[name] {
			if _, ok := current.containers[rel.Relation.To]; ok {
				current.internal++
			} else {
				current.cross++
			}
		}
	}

	metrics = append([]boundaryMetric{current}, metrics...)
	return metrics
}
