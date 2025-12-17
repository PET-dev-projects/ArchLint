package checks

import (
	"fmt"
	"strings"

	"github.com/NovokshanovE/archlint/pkg/model"
	"github.com/NovokshanovE/archlint/pkg/types"
)

const externalProtocolRuleID = "ARCH-EXTERNAL-PROTOCOL"

type externalProtocolRule struct{}

type externalProtocolConfig struct {
	AllowedPrefixes []string `json:"allowedPrefixes"`
	RequireProtocol bool     `json:"requireProtocol"`
}

var defaultExternalProtocolConfig = externalProtocolConfig{
	AllowedPrefixes: []string{
		"https://gateway.",
		"kafka://",
	},
	RequireProtocol: true,
}

// NewExternalProtocolRule enforces allowed protocols when hitting externals.
func NewExternalProtocolRule() Rule { return &externalProtocolRule{} }

func (r *externalProtocolRule) ID() string { return externalProtocolRuleID }

func (r *externalProtocolRule) Run(m *model.Architecture, cfg map[string]any) []types.Finding {
	conf := defaultExternalProtocolConfig
	if err := decodeConfig(cfg, &conf); err != nil {
		return []types.Finding{configFinding(externalProtocolRuleID, err)}
	}

	containerIndex := m.ContainerMap()

	findings := make([]types.Finding, 0)
	for _, relRef := range m.Relations() {
		rel := relRef.Relation
		to, ok := containerIndex[rel.To]
		if !ok || to.Container.Type != model.ContainerExternal {
			continue
		}
		protocol := rel.Protocol
		if strings.TrimSpace(protocol) == "" {
			if conf.RequireProtocol {
				findings = append(findings, types.Finding{
					RuleID:   externalProtocolRuleID,
					Severity: types.SeverityError,
					Message:  fmt.Sprintf("relation from %s to external %s must define protocol", rel.From, rel.To),
					Path:     relRef.Path,
				})
			}
			continue
		}
		if len(conf.AllowedPrefixes) == 0 {
			continue
		}
		valid := false
		lower := strings.ToLower(protocol)
		for _, prefix := range conf.AllowedPrefixes {
			if strings.HasPrefix(lower, strings.ToLower(prefix)) {
				valid = true
				break
			}
		}
		if !valid {
			findings = append(findings, types.Finding{
				RuleID:   externalProtocolRuleID,
				Severity: types.SeverityError,
				Message:  fmt.Sprintf("protocol %q for external %s is not allowed", protocol, rel.To),
				Path:     relRef.Path + ".protocol",
				Meta: map[string]any{
					"allowedPrefixes": conf.AllowedPrefixes,
				},
			})
		}
	}

	return findings
}
