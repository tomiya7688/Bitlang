package bitlang

import "fmt"

// ResolveProperties maps source property spelling to the canonical spelling
// declared by the machine-readable property specification.
func ResolveProperties(spec PropertySpecification, target string, properties []PreprocessedProperty) ([]PreprocessedProperty, error) {
	resolved := make([]PreprocessedProperty, 0, len(properties))
	for _, property := range properties {
		state, err := resolveProperty(spec, target, property)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, state)
	}
	return resolved, nil
}

func resolveProperty(spec PropertySpecification, target string, property PreprocessedProperty) (PreprocessedProperty, error) {
	propertyName, err := CanonicalizeIdentifier(string(property))
	if err != nil {
		return "", err
	}
	for _, axis := range spec.Axes {
		if !propertyAxisApplies(axis, target) {
			continue
		}
		for _, state := range axis.States {
			stateName, err := CanonicalizeIdentifier(state)
			if err != nil {
				return "", err
			}
			if stateName == propertyName {
				return PreprocessedProperty(state), nil
			}
		}
	}
	return "", fmt.Errorf("property %q is not applicable to %s", property, target)
}
