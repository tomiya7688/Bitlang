package bitlang

// Symbol represents one named semantic value.
// It intentionally avoids Go generics so the representation is easy to reproduce elsewhere.
type Symbol struct {
	Name  CanonicalName
	Value any
}
