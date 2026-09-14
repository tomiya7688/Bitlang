package bitlang

import "fmt"

// DuplicateSymbolError reports two spellings that resolve to the same Bitlang identifier.
type DuplicateSymbolError struct {
	Existing string
	Incoming string
}

// Error formats the duplicate-symbol diagnostic.
func (e DuplicateSymbolError) Error() string {
	return fmt.Sprintf("duplicate symbol: %q conflicts with existing %q", e.Incoming, e.Existing)
}
