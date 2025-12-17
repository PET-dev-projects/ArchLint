package checks

import "github.com/PET-dev-projects/ArchLint/pkg/types"

func configFinding(ruleID string, err error) types.Finding {
	return types.Finding{
		RuleID:   ruleID,
		Severity: types.SeverityError,
		Message:  "invalid rule configuration: " + err.Error(),
		Path:     "options.ruleConfig[" + ruleID + "]",
	}
}
