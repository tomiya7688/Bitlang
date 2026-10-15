package bitlang

//	{
//	  責務: [
//	    ParseSourceClasses: トークン列からクラスとメンバーの宣言境界を解析する
//	  ]
//	  処理: [
//	    1: クラス宣言とメンバー宣言の境界を保持する
//	    2: 式・文の本文は解析せず元トークンとして残す
//	  ]
//	  引数: [
//	    tokens: プリプロセス済みの元表記付きトークン列
//	  ]
//	  戻り値: [
//	    classes: 入力順のクラスとメンバー情報
//	    error: 宣言境界や区切りが不正な場合の失敗
//	  ]
//	}
func ParseSourceClasses(tokens []PreprocessedToken) ([]SourceClass, error) {
	parser := sourceClassParser{tokens: tokens}
	return parser.parseFile()
}

type sourceClassParser struct {
	tokens []PreprocessedToken
	index  int
}

func (p *sourceClassParser) parseFile() ([]SourceClass, error) {
	classes := make([]SourceClass, 0)
	for p.current().Kind != TokenEOF {
		class, err := p.parseClass()
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (p *sourceClassParser) parseClass() (SourceClass, error) {
	modifiers := p.parseModifiers()
	if _, err := p.expectWord("class"); err != nil {
		return SourceClass{}, err
	}
	nameToken, err := p.expectIdentifier("class name")
	if err != nil {
		return SourceClass{}, err
	}
	name, err := sourceName(nameToken)
	if err != nil {
		return SourceClass{}, err
	}
	typeParameters := make([]PreprocessedToken, 0)
	if p.isSymbol("<") {
		typeParameters, err = p.parseGenericTokenSequence(true)
		if err != nil {
			return SourceClass{}, err
		}
	}
	if _, err := p.expectSymbol("{"); err != nil {
		return SourceClass{}, err
	}

	class := SourceClass{
		Name: name, TypeParameters: typeParameters,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers,
		Members: make([]SourceMember, 0),
		Line:    nameToken.Line, Column: nameToken.Column,
	}
	for !p.isSymbol("}") {
		if p.current().Kind == TokenEOF {
			return SourceClass{}, p.errorAtCurrent("unterminated class body")
		}
		member, err := p.parseMember(class.Name)
		if err != nil {
			return SourceClass{}, err
		}
		class.Members = append(class.Members, member)
	}
	p.advance()
	return class, nil
}

//	{
//	  責務: [
//	    ParseSourceClassesInSourceResult: 既存の解析結果へクラス情報を追加する
//	  ]
//	  処理: [
//	    1: 関連付けられた字句トークンからクラスを解析する
//	    2: ソース情報や他の解析結果を保ったまま Classes を設定する
//	  ]
//	  引数: [
//	    source: 字句ソースと既存の解析情報を持つ値
//	  ]
//	  戻り値: [
//	    result: クラス解析結果を追加した値
//	    error: クラス宣言の解析に失敗した場合の内容
//	  ]
//	}
func ParseSourceClassesInSourceResult(source SourceParseResult) (SourceParseResult, error) {
	classes, err := ParseSourceClasses(source.Source.Tokens)
	if err != nil {
		return SourceParseResult{}, err
	}
	source.Classes = classes
	return source, nil
}
