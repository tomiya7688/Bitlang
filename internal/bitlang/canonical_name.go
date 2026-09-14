package bitlang

import (
	"fmt"
	"strings"
)

// CanonicalizeIdentifier returns the comparison form of a Bitlang identifier.
// String and character literal contents must never be passed through this function.
func CanonicalizeIdentifier(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("identifier must not be empty")
	}
	// TODO: Replace Go lowercase mapping when Bitlang Unicode identifier semantics are fixed.
	return strings.ToLower(name), nil
}

// CanonicalName keeps both diagnostic spelling and comparison form.
type CanonicalName struct {
	Spelling  string
	Canonical string
}

// NewCanonicalName creates the two forms of one Bitlang identifier.
func NewCanonicalName(spelling string) (CanonicalName, error) {
	canonical, err := CanonicalizeIdentifier(spelling)
	if err != nil {
		return CanonicalName{}, err
	}
	return CanonicalName{Spelling: spelling, Canonical: canonical}, nil
}
