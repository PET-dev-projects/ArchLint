package checks

// Registry exposes default rules used by the engine.
type Registry struct {
	rules []Rule
}

// DefaultRegistry returns registry with built-in rules.
func DefaultRegistry() Registry {
	return Registry{
		rules: []Rule{
			NewAcyclicRule(),
			NewCRUDRule(),
			NewACLRule(),
			NewBoundariesRule(),
			NewExternalProtocolRule(),
			NewDatabaseIsolationRule(),
		},
	}
}

// Rules returns all registered rules.
func (r Registry) Rules() []Rule {
	return r.rules
}

// Find looks up rule by ID.
func (r Registry) Find(id string) (Rule, bool) {
	for _, rule := range r.rules {
		if rule.ID() == id {
			return rule, true
		}
	}
	return nil, false
}
