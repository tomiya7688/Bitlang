package bitlang

import "fmt"

// ParsePreprocessedDeclaration parses one already-expanded declaration.
// Bootstrap grammar:
//   <properties...> <type> <name> ;
func ParsePreprocessedDeclaration(spec PropertySpecification, target string, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	body := declarationBody(tokens)
	if len(body) < 3 {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration requires properties, type, and name")
	}
	if body[len(body)-1].Lexeme != ";" {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration must end with semicolon")
	}
	typeToken := body[len(body)-3]
	nameToken := body[len(body)-2]
	if typeToken.Kind != TokenIdentifier || nameToken.Kind != TokenIdentifier {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration type and name must be identifiers")
	}

	properties := make([]PreprocessedProperty, 0, len(body)-3)
	for _, token := range body[:len(body)-3] {
		if token.Kind != TokenIdentifier {
			return PreprocessedDeclaration{}, fmt.Errorf("property %q must be an identifier", token.Lexeme)
		}
		properties = append(properties, PreprocessedProperty(token.Lexeme))
	}
	if err := ValidateProperties(spec, target, properties); err != nil {
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
		Name: name, Type: typeName, Properties: properties,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func declarationBody(tokens []PreprocessedToken) []PreprocessedToken {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind == TokenEOF {
		return tokens[:len(tokens)-1]
	}
	return tokens
}
