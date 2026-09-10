package bitlang

import (
	"fmt"
	"strings"
)

// CanonicalizeIdentifier returns the comparison form of a Bitlang identifier.
//
// Bitlang identifiers are case-insensitive. String and character literal
// contents must never be passed through this function; canonicalization belongs
// to identifier/symbol handling after lexing.
//
// Unicode identifier semantics are not fixed yet. strings.ToLower gives a
// deterministic standard-library-only baseline without committing the language
// specification to a particular Unicode normalization profile.
func CanonicalizeIdentifier(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("identifier must not be empty")
	}
	return strings.ToLower(name), nil
}

type CanonicalName struct {
	Spelling  string
	Canonical string
}

func NewCanonicalName(spelling string) (CanonicalName, error) {
	canonical, err := CanonicalizeIdentifier(spelling)
	if err != nil {
		return CanonicalName{}, err
	}
	return CanonicalName{Spelling: spelling, Canonical: canonical}, nil
}

// Symbol intentionally avoids Go generics. Semantic compiler data structures
// should be straightforward to reproduce in other implementation languages.
type Symbol struct {
	Name  CanonicalName
	Value any
}

type DuplicateSymbolError struct {
	Existing string
	Incoming string
}

func (e DuplicateSymbolError) Error() string {
	return fmt.Sprintf("duplicate symbol: %q conflicts with existing %q", e.Incoming, e.Existing)
}

type SymbolTable struct {
	symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]Symbol)}
}

func (t *SymbolTable) Define(name string, value any) (Symbol, error) {
	parsed, err := NewCanonicalName(name)
	if err != nil {
		return Symbol{}, err
	}
	if previous, ok := t.symbols[parsed.Canonical]; ok {
		return Symbol{}, DuplicateSymbolError{
			Existing: previous.Name.Spelling,
			Incoming: parsed.Spelling,
		}
	}

	symbol := Symbol{Name: parsed, Value: value}
	t.symbols[parsed.Canonical] = symbol
	return symbol, nil
}

func (t *SymbolTable) Find(name string) (Symbol, bool) {
	canonical, err := CanonicalizeIdentifier(name)
	if err != nil {
		return Symbol{}, false
	}
	symbol, ok := t.symbols[canonical]
	return symbol, ok
}

func (t *SymbolTable) Get(name string) (any, bool) {
	symbol, ok := t.Find(name)
	if !ok {
		return nil, false
	}
	return symbol.Value, true
}

func (t *SymbolTable) Len() int {
	return len(t.symbols)
}
