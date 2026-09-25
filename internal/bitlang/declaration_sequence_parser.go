package bitlang

import "fmt"

// ParsePreprocessedDeclarations parses a sequence of declarations using a
// validated shared specification set and one declaration kind name.
func ParsePreprocessedDeclarations(specs SpecificationSet, kindName string, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	kind, err := specs.Declarations.DeclarationKind(kindName)
	if err != nil {
		return nil, err
	}
	groups, err := splitDeclarationTokens(tokens, kind.Terminator)
	if err != nil {
		return nil, err
	}
	declarations := make([]PreprocessedDeclaration, 0, len(groups))
	for _, group := range groups {
		declaration, err := ParsePreprocessedDeclaration(specs, kindName, group)
		if err != nil {
			return nil, fmt.Errorf("%d:%d: %w", group[0].Line, group[0].Column, err)
		}
		declarations = append(declarations, declaration)
	}
	return declarations, nil
}

func splitDeclarationTokens(tokens []PreprocessedToken, terminator string) ([][]PreprocessedToken, error) {
	var groups [][]PreprocessedToken
	start := 0
	for index, token := range tokens {
		if token.Kind == TokenEOF {
			if index != start {
				return nil, fmt.Errorf("%d:%d: declaration missing %q", tokens[start].Line, tokens[start].Column, terminator)
			}
			break
		}
		if token.Lexeme != terminator {
			continue
		}
		group := append([]PreprocessedToken(nil), tokens[start:index+1]...)
		groups = append(groups, group)
		start = index + 1
	}
	if start < len(tokens) && tokens[start].Kind != TokenEOF {
		return nil, fmt.Errorf("%d:%d: declaration missing %q", tokens[start].Line, tokens[start].Column, terminator)
	}
	return groups, nil
}
