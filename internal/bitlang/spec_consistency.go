package bitlang

import "fmt"

// ValidateSpecificationConsistency verifies invariants that span multiple
// machine-readable Bitlang specification files.
func ValidateSpecificationConsistency(properties PropertySpecification, declarations DeclarationSpecification) error {
	targets := propertyTargets(properties)
	for _, kind := range declarations.Kinds {
		if !targets[kind.PropertyTarget] {
			return fmt.Errorf("declaration kind %q references unknown property target %q", kind.Name, kind.PropertyTarget)
		}
	}
	return validateUniquePropertyStates(properties)
}

func propertyTargets(spec PropertySpecification) map[string]bool {
	targets := map[string]bool{}
	for _, axis := range spec.Axes {
		for _, target := range axis.AppliesTo {
			targets[target] = true
		}
	}
	return targets
}

func validateUniquePropertyStates(spec PropertySpecification) error {
	owners := map[string]string{}
	for _, axis := range spec.Axes {
		for _, state := range axis.States {
			if state == "" {
				return fmt.Errorf("property axis %q contains an empty state", axis.Name)
			}
			if owner, exists := owners[state]; exists {
				return fmt.Errorf("property state %q is shared by axes %q and %q", state, owner, axis.Name)
			}
			owners[state] = axis.Name
		}
	}
	return nil
}
