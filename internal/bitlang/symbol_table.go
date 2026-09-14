package bitlang

// SymbolTable stores symbols by canonical Bitlang identifier.
type SymbolTable struct {
	symbols map[string]Symbol
}

// NewSymbolTable creates an empty symbol table.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]Symbol)}
}

// Define inserts one symbol and rejects canonical-name collisions.
func (t *SymbolTable) Define(name string, value any) (Symbol, error) {
	parsed, err := NewCanonicalName(name)
	if err != nil {
		return Symbol{}, err
	}
	if previous, ok := t.symbols[parsed.Canonical]; ok {
		return Symbol{}, DuplicateSymbolError{Existing: previous.Name.Spelling, Incoming: parsed.Spelling}
	}
	symbol := Symbol{Name: parsed, Value: value}
	t.symbols[parsed.Canonical] = symbol
	return symbol, nil
}

// Find resolves one symbol using case-insensitive identifier semantics.
func (t *SymbolTable) Find(name string) (Symbol, bool) {
	canonical, err := CanonicalizeIdentifier(name)
	if err != nil {
		return Symbol{}, false
	}
	symbol, ok := t.symbols[canonical]
	return symbol, ok
}

// Get returns the stored value for one identifier.
func (t *SymbolTable) Get(name string) (any, bool) {
	symbol, ok := t.Find(name)
	if !ok {
		return nil, false
	}
	return symbol.Value, true
}

// Len returns the number of stored symbols.
func (t *SymbolTable) Len() int {
	return len(t.symbols)
}
