package bitlang

import (
	"encoding/json"
	"fmt"
)

// LoadPropertySpecification decodes and minimally validates the machine-readable
// Preprocessed property definition supplied by the caller.
func LoadPropertySpecification(data []byte) (PropertySpecification, error) {
	var spec PropertySpecification
	if err := json.Unmarshal(data, &spec); err != nil {
		return PropertySpecification{}, fmt.Errorf("decode property specification: %w", err)
	}
	if spec.Version <= 0 {
		return PropertySpecification{}, fmt.Errorf("property specification version must be positive")
	}
	if len(spec.Axes) == 0 {
		return PropertySpecification{}, fmt.Errorf("property specification must define at least one axis")
	}
	for _, axis := range spec.Axes {
		if axis.Name == "" || len(axis.States) == 0 {
			return PropertySpecification{}, fmt.Errorf("property axis requires name and states")
		}
	}
	return spec, nil
}
