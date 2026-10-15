package bitlang

import (
	"fmt"
	"strings"
)

//	{
//	  責務: [
//	    ParseMixedPreprocessedDeclarations: 宣言種別を混在させたトークン列を解析する
//	  ]
//	  処理: [
//	    1: 各位置で仕様に合う宣言種別を特定する
//	    2: 宣言列を入力順に返す
//	  ]
//	  引数: [
//	    specs: 宣言形式とプロパティの検証済み仕様
//	    tokens: 種別が混在するプリプロセストークン列
//	  ]
//	  戻り値: [
//	    declarations: 解析済み宣言列
//	    error: 宣言形式が曖昧または不正な場合の失敗
//	  ]
//	}
func ParseMixedPreprocessedDeclarations(specs SpecificationSet, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	return ParseMixedPreprocessedDeclarationsInScope(specs, nil, tokens)
}

//	{
//	  責務: [
//	    ParseMixedPreprocessedDeclarationsInContext: 文字列指定の互換入口から文脈付き宣言解析を行う
//	  ]
//	  処理: [
//	    1: 文脈名を型付きスコープへ変換する
//	    2: スコープを使って宣言種別を絞り込み解析する
//	  ]
//	  引数: [
//	    specs: 宣言形式とプロパティの検証済み仕様
//	    context: 現在の宣言文脈名
//	    tokens: 解析対象のプリプロセストークン列
//	  ]
//	  戻り値: [
//	    declarations: 文脈に適合した解析済み宣言列
//	    error: 文脈または宣言解析に失敗した場合の内容
//	  ]
//	}
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

//	{
//	  責務: [
//	    ParseMixedPreprocessedDeclarationsInScope: 型付きスコープに適合する宣言列を解析する
//	  ]
//	  処理: [
//	    1: スコープ文脈に適用できる宣言種別だけを候補にする
//	    2: 候補種別に基づいて宣言列を解析する
//	  ]
//	  引数: [
//	    specs: 宣言形式とプロパティの検証済み仕様
//	    scope: 宣言文脈を表す型付きスコープ
//	    tokens: 解析対象のプリプロセストークン列
//	  ]
//	  戻り値: [
//	    declarations: 文脈に適合した解析済み宣言列
//	    error: 対応する宣言種別がない場合などの失敗
//	  ]
//	}
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
