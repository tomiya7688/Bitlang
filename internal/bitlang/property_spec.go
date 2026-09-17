package bitlang

// PropertySpecification is the machine-readable definition of Preprocessed
// property axes. It intentionally contains no Bitlang-specific state names in
// Go code.
type PropertySpecification struct {
	Version int                `json:"version"`
	Axes    []PropertyAxisSpec `json:"axes"`
}

// PropertyAxisSpec defines one property axis and its allowed final states.
type PropertyAxisSpec struct {
	Name      string   `json:"name"`
	States    []string `json:"states"`
	Exclusive bool     `json:"exclusive"`
}
