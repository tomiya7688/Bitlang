package bitlang

// SpecificationSet groups machine-readable Bitlang specifications after
// individual decoding and cross-specification validation.
type SpecificationSet struct {
	Properties   PropertySpecification
	Declarations DeclarationSpecification
}
