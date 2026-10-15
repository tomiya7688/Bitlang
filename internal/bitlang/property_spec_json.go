package bitlang

import (
	"encoding/json"
	"fmt"
)

func validatePropertySpecificationJSONFields(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if err := ValidateJSONFieldNames(fields, "version", "axes"); err != nil {
		return err
	}
	var axes []json.RawMessage
	if err := json.Unmarshal(fields["axes"], &axes); err != nil {
		return err
	}
	for _, data := range axes {
		var axisFields map[string]json.RawMessage
		if err := json.Unmarshal(data, &axisFields); err != nil {
			return err
		}
		if err := ValidateJSONFieldNames(axisFields, "name", "states", "exclusive", "required", "applies_to", "source_default", "source_default_when"); err != nil {
			return err
		}
		if err := validatePropertySourceDefaultJSONFields(axisFields["source_default_when"]); err != nil {
			return err
		}
	}
	return nil
}

func validatePropertySourceDefaultJSONFields(data json.RawMessage) error {
	if len(data) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if fields == nil {
		return fmt.Errorf("source_default_when must be an object")
	}
	return ValidateJSONFieldNames(fields, "axis", "states")
}
