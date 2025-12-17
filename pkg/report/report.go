package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/PET-dev-projects/ArchLint/pkg/types"
)

// WriteText renders findings as a plain-text list.
func WriteText(w io.Writer, findings []types.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "No findings")
		return err
	}
	for _, f := range findings {
		var meta string
		if len(f.Meta) > 0 {
			data, err := json.Marshal(f.Meta)
			if err == nil {
				meta = string(data)
			}
		}
		line := fmt.Sprintf("%s\t%s\t%s\t%s", f.RuleID, f.Severity, f.Path, f.Message)
		if meta != "" {
			line += "\t" + meta
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON serializes findings to JSON array.
func WriteJSON(w io.Writer, findings []types.Finding) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}
