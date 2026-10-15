package bitlang

import "fmt"

//	{
//	  責務: [
//	    ParsePreprocessedDeclaration: 仕様で定義された形式に従って宣言を解析する
//	  ]
//	  処理: [
//	    1: 共有仕様から宣言種別を解決する
//	    2: トークン列を宣言の構成要素へ分解する
//	    3: 宣言名と型名を正規表記で保持する
//	  ]
//	  引数: [
//	    specs: 検証済みの共有仕様
//	    kindName: 仕様で定義された宣言種別名
//	    tokens: 解析対象のプリプロセストークン列
//	  ]
//	  戻り値: [
//	    declaration: 解析した宣言
//	    error: 種別または宣言形式が無効な場合の失敗
//	  ]
//	}
func ParsePreprocessedDeclaration(specs SpecificationSet, kindName string, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	kind, err := specs.Declarations.DeclarationKind(kindName)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}

	body := declarationBody(tokens)
	if len(body) < 2 {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration does not satisfy configured layout")
	}
	if body[len(body)-1].Lexeme != kind.Terminator {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration must end with %q", kind.Terminator)
	}

	parts, err := parseDeclarationLayout(kind.Layout, body[:len(body)-1])
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	if parts.typeToken.Kind != TokenIdentifier || parts.nameToken.Kind != TokenIdentifier {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration type and name must be identifiers")
	}

	explicit := make([]PreprocessedProperty, 0, len(parts.properties))
	for _, token := range parts.properties {
		if token.Kind != TokenIdentifier && token.Kind != TokenKeyword {
			return PreprocessedDeclaration{}, fmt.Errorf("property %q must be an identifier or keyword", token.Lexeme)
		}
		explicit = append(explicit, PreprocessedProperty(token.Lexeme))
	}
	resolved, err := ResolveProperties(specs.Properties, kind.PropertyTarget, explicit)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	if err := ValidateProperties(specs.Properties, kind.PropertyTarget, resolved); err != nil {
		return PreprocessedDeclaration{}, err
	}

	name, err := NewCanonicalName(parts.nameToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	typeName, err := NewCanonicalName(parts.typeToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	return PreprocessedDeclaration{
		Kind: kind.Name, Name: name, Type: typeName, Properties: resolved,
		Line: parts.nameToken.Line, Column: parts.nameToken.Column,
	}, nil
}

func declarationBody(tokens []PreprocessedToken) []PreprocessedToken {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind == TokenEOF {
		return tokens[:len(tokens)-1]
	}
	return tokens
}
