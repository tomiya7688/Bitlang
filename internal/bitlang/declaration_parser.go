package bitlang

import "fmt"

// ParsePreprocessedDeclaration parses one declaration using its data-defined
// strict grammar and property target.
func ParsePreprocessedDeclaration(properties PropertySpecification, kind DeclarationKindSpec, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	body := declarationBody(tokens)
	if len(body) < 3 {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration requires properties, type, and name")
	}
	if body[len(body)-1].Lexeme != kind.Terminator {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration must end with %q", kind.Terminator)
	}
	typeToken := body[len(body)-3]
	nameToken := body[len(body)-2]
	if typeToken.Kind != TokenIdentifier || nameToken.Kind != TokenIdentifier {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration type and name must be identifiers")
	}

	explicit := make([]PreprocessedProperty, 0, len(body)-3)
	for _, token := range body[:len(body)-3] {
		if token.Kind != TokenIdentifier {
			return PreprocessedDeclaration{}, fmt.Errorf("property %q must be an identifier", token.Lexeme)
		}
		explicit = append(explicit, PreprocessedProperty(token.Lexeme))
	}
	if err := ValidateProperties(properties, kind.PropertyTarget, explicit); err != nil {
		return PreprocessedDeclaration{}, err
	}

	name, err := NewCanonicalName(nameToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	typeName, err := NewCanonicalName(typeToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	return PreprocessedDeclaration{
		Name: name, Type: typeName, Properties: explicit,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func declarationBody(tokens []PreprocessedToken) []PreprocessedToken {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind == TokenEOF {
		return tokens[:len(tokens)-1]
	}
	return tokens
}
