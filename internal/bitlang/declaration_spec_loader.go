package bitlang

import (
	"encoding/json"
	"fmt"
)

// LoadDeclarationSpecification decodes and validates declaration grammar data.
func LoadDeclarationSpecification(data []byte) (DeclarationSpecification, error) {
	var spec DeclarationSpecification
	if err := json.Unmarshal(data, &spec); err != nil {
		return DeclarationSpecification{}, fmt.Errorf("decode declaration specification: %w", err)
	}
	if spec.Version <= 0 || len(spec.Kinds) == 0 {
		return DeclarationSpecification{}, fmt.Errorf("declaration specification requires version and kinds")
	}
	seen := map[string]bool{}
	for _, kind := range spec.Kinds {
		if kind.Name == "" || kind.PropertyTarget == "" || kind.Terminator == "" {
			return DeclarationSpecification{}, fmt.Errorf("declaration kind requires name, property target, and terminator")
		}
		if seen[kind.Name] {
			return DeclarationSpecification{}, fmt.Errorf("duplicate declaration kind %q", kind.Name)
		}
		seen[kind.Name] = true
		if !validDeclarationLayout(kind.Layout) {
			return DeclarationSpecification{}, fmt.Errorf("unsupported declaration layout for %q", kind.Name)
		}
	}
	return spec, nil
}

// DeclarationKind resolves one declaration kind by its data-defined name.
func (s DeclarationSpecification) DeclarationKind(name string) (DeclarationKindSpec, error) {
	for _, kind := range s.Kinds {
		if kind.Name == name {
			return kind, nil
		}
	}
	return DeclarationKindSpec{}, fmt.Errorf("unknown declaration kind %q", name)
}

