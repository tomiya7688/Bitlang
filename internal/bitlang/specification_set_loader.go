package bitlang

// LoadSpecificationSet decodes all machine-readable specifications needed by
// the bootstrap semantic pipeline and validates their cross-file invariants.
func LoadSpecificationSet(propertiesData []byte, declarationsData []byte) (SpecificationSet, error) {
	properties, err := LoadPropertySpecification(propertiesData)
	if err != nil {
		return SpecificationSet{}, err
	}
	declarations, err := LoadDeclarationSpecification(declarationsData)
	if err != nil {
		return SpecificationSet{}, err
	}
	if err := ValidateSpecificationConsistency(properties, declarations); err != nil {
		return SpecificationSet{}, err
	}
	return SpecificationSet{
		Properties:   properties,
		Declarations: declarations,
	}, nil
}
