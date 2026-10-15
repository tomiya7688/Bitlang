package bitlangci

import "testing"

func TestValidateFamilyPipelineManifestJSONFieldsRejectsIncorrectCase(t *testing.T) {
	cases := []struct {
		name string
		data string
	}{
		{"root field", `{"Schema_version":1}`},
		{"repository field", `{"repositories":[{"ID":"source"}]}`},
		{"transition field", `{"transitions":[{"From":"source"}]}`},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if err := validateFamilyPipelineManifestJSONFields([]byte(test.data)); err == nil {
				t.Fatal("expected incorrectly cased manifest field to fail")
			}
		})
	}
}
