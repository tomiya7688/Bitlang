package bitlang

// DeclarationSpecification defines strict declaration forms without embedding
// declaration-kind names in the Go implementation.
type DeclarationSpecification struct {
	Version int                   `json:"version"`
	Kinds   []DeclarationKindSpec `json:"kinds"`
}

// DeclarationKindSpec defines one strict declaration kind.
type DeclarationKindSpec struct {
	Name           string   `json:"name"`
	PropertyTarget string   `json:"property_target"`
	Terminator     string   `json:"terminator"`
	Layout         []string `json:"layout"`
}
