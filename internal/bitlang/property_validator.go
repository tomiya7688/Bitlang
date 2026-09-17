package bitlang

import "fmt"

// ValidateProperties verifies that a declaration contains exactly the property
// states required by the machine-readable specification for its target kind.
func ValidateProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) error {
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) {
			continue
		}
		matches := countPropertyAxisStates(axis, properties)
		if axis.Required && matches == 0 {
			return fmt.Errorf("missing required property axis %q for %s", axis.Name, target)
		}
		if axis.Exclusive && matches > 1 {
			return fmt.Errorf("conflicting property states for axis %q", axis.Name)
		}
	}
	return validateKnownProperties(spec, target, properties)
}

func propertyAxisApplies(axis PropertyAxisSpec, target string) bool {
	for _, candidate := range axis.AppliesTo {
		if candidate == target {
			return true
		}
	}
	return false
}

func countPropertyAxisStates(axis PropertyAxisSpec, properties []PreprocessedProperty) int {
	count := 0
	for _, property := range properties {
		for _, state := range axis.States {
			if string(property) == state {
				count++
			}
		}
	}
	return count
}

func validateKnownProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) error {
	for _, property := range properties {
		known := false
		for _, axis := range spec.Axes {
			if propertyAxisApplies(axis, target) && propertyInAxis(axis, property) {
				known = true
				break
			}
		}
		if !known {
			return fmt.Errorf("property %q is not applicable to %s", property, target)
		}
	}
	return nil
}

func propertyInAxis(axis PropertyAxisSpec, property PreprocessedProperty) bool {
	for _, state := range axis.States {
		if state == string(property) {
			return true
		}
	}
	return false
}
