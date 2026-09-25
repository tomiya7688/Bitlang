package bitlang

import (
	"fmt"
	"strings"
)

// ParseMixedPreprocessedDeclarations parses a token stream containing multiple
// declaration kinds by testing each declaration against data-defined kinds.
func ParseMixedPreprocessedDeclarations(specs SpecificationSet, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	groups, err := splitMixedDeclarationTokens(specs.Declarations, tokens)
	if err != nil {
		return nil, err
	}

	declarations := make([]PreprocessedDeclaration, 0, len(groups))
	for _, group := range groups {
		declaration, err := detectDeclarationKind(specs, group)
		if err != nil {
			return nil, fmt.Errorf("%d:%d: %w", group[0].Line, group[0].Column, err)
		}
		declarations = append(declarations, declaration)
	}
	return declarations, nil
}

func detectDeclarationKind(specs SpecificationSet, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	var matches []PreprocessedDeclaration
	for _, kind := range specs.Declarations.Kinds {
		declaration, err := ParsePreprocessedDeclaration(specs, kind.Name, tokens)
		if err == nil {
			matches = append(matches, declaration)
		}
	}
	if len(matches) == 0 {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration matches no configured kind")
	}
	if len(matches) > 1 {
		names := make([]string, 0, len(matches))
		for _, match := range matches {
			names = append(names, match.Kind)
		}
		return PreprocessedDeclaration{}, fmt.Errorf("ambiguous declaration matches kinds: %s", strings.Join(names, ", "))
	}
	return matches[0], nil
}

func splitMixedDeclarationTokens(spec DeclarationSpecification, tokens []PreprocessedToken) ([][]PreprocessedToken, error) {
	terminators := declarationTerminators(spec)
	var groups [][]PreprocessedToken
	start := 0
	for index, token := range tokens {
		if token.Kind == TokenEOF {
			if index != start {
				return nil, fmt.Errorf("%d:%d: declaration missing terminator", tokens[start].Line, tokens[start].Column)
			}
			break
		}
		if !terminators[token.Lexeme] {
			continue
		}
		groups = append(groups, append([]PreprocessedToken(nil), tokens[start:index+1]...))
		start = index + 1
	}
	if start < len(tokens) && tokens[start].Kind != TokenEOF {
		return nil, fmt.Errorf("%d:%d: declaration missing terminator", tokens[start].Line, tokens[start].Column)
	}
	return groups, nil
}

func declarationTerminators(spec DeclarationSpecification) map[string]bool {
	terminators := make(map[string]bool)
	for _, kind := range spec.Kinds {
		terminators[kind.Terminator] = true
	}
	return terminators
}
