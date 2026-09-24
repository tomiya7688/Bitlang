package bitlang

import "fmt"

// ParsePreprocessedDeclarations parses a sequence of strict declarations of one
// target kind. Target-kind discovery remains separate until declaration syntax
// is defined by machine-readable grammar data.
func ParsePreprocessedDeclarations(spec PropertySpecification, target string, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	groups, err := splitDeclarationTokens(tokens)
	if err != nil {
		return nil, err
	}
	declarations := make([]PreprocessedDeclaration, 0, len(groups))
	for _, group := range groups {
		declaration, err := ParsePreprocessedDeclaration(spec, target, group)
		if err != nil {
			return nil, fmt.Errorf("%d:%d: %w", group[0].Line, group[0].Column, err)
		}
		declarations = append(declarations, declaration)
	}
	return declarations, nil
}

func splitDeclarationTokens(tokens []PreprocessedToken) ([][]PreprocessedToken, error) {
	var groups [][]PreprocessedToken
	start := 0
	for index, token := range tokens {
		if token.Kind == TokenEOF {
			if index != start {
				return nil, fmt.Errorf("%d:%d: declaration missing semicolon", tokens[start].Line, tokens[start].Column)
			}
			break
		}
		if token.Lexeme != ";" {
			continue
		}
		group := append([]PreprocessedToken(nil), tokens[start:index+1]...)
		groups = append(groups, group)
		start = index + 1
	}
	if start < len(tokens) && tokens[start].Kind != TokenEOF {
		return nil, fmt.Errorf("%d:%d: declaration missing semicolon", tokens[start].Line, tokens[start].Column)
	}
	return groups, nil
}
