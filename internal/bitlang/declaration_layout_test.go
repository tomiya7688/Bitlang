package bitlang

import "testing"

func TestParseDeclarationLayoutSupportsPropertiesLast(t *testing.T) {
	tokens := []PreprocessedToken{
		{Kind: TokenIdentifier, Lexeme: "Int"},
		{Kind: TokenIdentifier, Lexeme: "PlayerHP"},
		{Kind: TokenIdentifier, Lexeme: "Private"},
		{Kind: TokenIdentifier, Lexeme: "unnullable"},
	}
	parts, err := parseDeclarationLayout([]string{"type", "name", "properties"}, tokens)
	if err != nil {
		t.Fatal(err)
	}
	if parts.typeToken.Lexeme != "Int" || parts.nameToken.Lexeme != "PlayerHP" {
		t.Fatalf("unexpected scalar parts: %#v", parts)
	}
	if len(parts.properties) != 2 || parts.properties[0].Lexeme != "Private" {
		t.Fatalf("unexpected property parts: %#v", parts.properties)
	}
}

func TestParseDeclarationLayoutSupportsPropertiesMiddle(t *testing.T) {
	tokens := []PreprocessedToken{
		{Kind: TokenIdentifier, Lexeme: "Int"},
		{Kind: TokenIdentifier, Lexeme: "Private"},
		{Kind: TokenIdentifier, Lexeme: "unnullable"},
		{Kind: TokenIdentifier, Lexeme: "PlayerHP"},
	}
	parts, err := parseDeclarationLayout([]string{"type", "properties", "name"}, tokens)
	if err != nil {
		t.Fatal(err)
	}
	if parts.typeToken.Lexeme != "Int" || parts.nameToken.Lexeme != "PlayerHP" {
		t.Fatalf("unexpected scalar parts: %#v", parts)
	}
	if len(parts.properties) != 2 {
		t.Fatalf("property count = %d, want 2", len(parts.properties))
	}
}

func TestValidDeclarationLayoutRejectsDuplicateComponent(t *testing.T) {
	if validDeclarationLayout([]string{"properties", "name", "name"}) {
		t.Fatal("duplicate declaration layout component accepted")
	}
}
