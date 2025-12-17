package engine

import (
	"sort"

	"github.com/NovokshanovE/archlint/pkg/checks"
	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

// Options tune engine execution behaviour.
type Options struct {
	EnabledRules []string
	RuleConfig   map[string]map[string]any
}

// RunAll executes all enabled rules against the provided model.
func RunAll(m *model.Architecture, opts Options) []types.Finding {
	registry := checks.DefaultRegistry()
	var rules []checks.Rule
	if len(opts.EnabledRules) > 0 {
		for _, id := range opts.EnabledRules {
			if rule, ok := registry.Find(id); ok {
				rules = append(rules, rule)
			}
		}
	} else {
		rules = registry.Rules()
	}

	findings := make([]types.Finding, 0)
	for _, rule := range rules {
		cfg := map[string]any(nil)
		if opts.RuleConfig != nil {
			cfg = opts.RuleConfig[rule.ID()]
		}
		findings = append(findings, rule.Run(m, cfg)...)
	}

	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].RuleID == findings[j].RuleID {
			if findings[i].Path == findings[j].Path {
				return findings[i].Message < findings[j].Message
			}
			return findings[i].Path < findings[j].Path
		}
		return findings[i].RuleID < findings[j].RuleID
	})

	return findings
}
