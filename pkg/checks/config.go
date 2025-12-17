package checks

import "encoding/json"

func decodeConfig(cfg map[string]any, target any) error {
	if cfg == nil {
		cfg = map[string]any{}
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
