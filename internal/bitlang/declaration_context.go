package bitlang

import "fmt"

func declarationKindAppliesToContext(kind DeclarationKindSpec, context string) bool {
	if len(kind.Contexts) == 0 {
		return true
	}
	if context == "" {
		return false
	}
	canonicalContext, err := CanonicalizeIdentifier(context)
	if err != nil {
		return false
	}
	for _, candidate := range kind.Contexts {
		canonicalCandidate, err := CanonicalizeIdentifier(candidate)
		if err == nil && canonicalCandidate == canonicalContext {
			return true
		}
	}
	return false
}

func validateDeclarationContexts(kind DeclarationKindSpec) error {
	seen := map[string]bool{}
	for _, context := range kind.Contexts {
		canonical, err := CanonicalizeIdentifier(context)
		if err != nil {
			return fmt.Errorf("declaration kind %q has invalid context: %w", kind.Name, err)
		}
		if seen[canonical] {
			return fmt.Errorf("declaration kind %q has duplicate context %q", kind.Name, context)
		}
		seen[canonical] = true
	}
	return nil
}
