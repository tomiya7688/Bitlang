package bitlang

import (
	"fmt"
	"strings"
)

// ParseMixedPreprocessedDeclarations parses a token stream containing multiple
// declaration kinds without applying a declaration context filter.
func ParseMixedPreprocessedDeclarations(specs SpecificationSet, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	return ParseMixedPreprocessedDeclarationsInScope(specs, nil, tokens)
}

// ParseMixedPreprocessedDeclarationsInContext preserves the string-based
// bootstrap entrypoint while routing through the typed scope representation.
func ParseMixedPreprocessedDeclarationsInContext(specs SpecificationSet, context string, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	if context == "" {
		return ParseMixedPreprocessedDeclarationsInScope(specs, nil, tokens)
	}
	scope, err := NewPreprocessedScope(context, nil)
	if err != nil {
		return nil, err
	}
	return ParseMixedPreprocessedDeclarationsInScope(specs, &scope, tokens)
}

// ParseMixedPreprocessedDeclarationsInScope parses mixed declarations after
// filtering candidate kinds by the current Preprocessed scope.
func ParseMixedPreprocessedDeclarationsInScope(specs SpecificationSet, scope *PreprocessedScope, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	groups, err := splitMixedDeclarationTokens(specs.Declarations, tokens)
	if err != nil {
		return nil, err
	}

	declarations := make([]PreprocessedDeclaration, 0, len(groups))
	for _, group := range groups {
		declaration, err := detectDeclarationKind(specs, scope, group)
		if err != nil {
			return nil, fmt.Errorf("%d:%d: %w", group[0].Line, group[0].Column, err)
		}
		declarations = append(declarations, declaration)
	}
	return declarations, nil
}

func detectDeclarationKind(specs SpecificationSet, scope *PreprocessedScope, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	var matches []PreprocessedDeclaration
	for _, kind := range specs.Declarations.Kinds {
		if !declarationKindAppliesToScope(kind, scope) {
			continue
		}
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
