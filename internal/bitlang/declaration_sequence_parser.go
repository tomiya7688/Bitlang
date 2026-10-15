package bitlang

import "fmt"

//	{
//	  責務: [
//	    ParsePreprocessedDeclarations: 同じ種別の宣言列を個別の宣言へ分割して解析する
//	  ]
//	  処理: [
//	    1: 終端記号でトークン列を宣言単位に分割する
//	    2: 各宣言を共有仕様に従って解析する
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    kindName: 宣言列に適用する宣言種別名
//	    tokens: 解析対象のプリプロセストークン列
//	  ]
//	  戻り値: [
//	    declarations: 入力順の解析済み宣言
//	    error: 分割または宣言解析に失敗した場合の内容
//	  ]
//	}
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
