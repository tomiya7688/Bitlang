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

type Symbol[T any] struct {
	Name  CanonicalName
	Value T
}

type DuplicateSymbolError struct {
	Existing string
	Incoming string
}

func (e DuplicateSymbolError) Error() string {
	return fmt.Sprintf("duplicate symbol: %q conflicts with existing %q", e.Incoming, e.Existing)
}

type SymbolTable[T any] struct {
	symbols map[string]Symbol[T]
}

func NewSymbolTable[T any]() *SymbolTable[T] {
	return &SymbolTable[T]{symbols: make(map[string]Symbol[T])}
}

func (t *SymbolTable[T]) Define(name string, value T) (Symbol[T], error) {
	parsed, err := NewCanonicalName(name)
	if err != nil {
		return Symbol[T]{}, err
	}
	if previous, ok := t.symbols[parsed.Canonical]; ok {
		return Symbol[T]{}, DuplicateSymbolError{
			Existing: previous.Name.Spelling,
			Incoming: parsed.Spelling,
		}
	}

	symbol := Symbol[T]{Name: parsed, Value: value}
	t.symbols[parsed.Canonical] = symbol
	return symbol, nil
}

func (t *SymbolTable[T]) Find(name string) (Symbol[T], bool) {
	canonical, err := CanonicalizeIdentifier(name)
	if err != nil {
		return Symbol[T]{}, false
	}
	symbol, ok := t.symbols[canonical]
	return symbol, ok
}

func (t *SymbolTable[T]) Get(name string) (T, bool) {
	symbol, ok := t.Find(name)
	if !ok {
		var zero T
		return zero, false
	}
	return symbol.Value, true
}

func (t *SymbolTable[T]) Len() int {
	return len(t.symbols)
}
