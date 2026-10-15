package bitlang

//	{
//	  責務: [
//	    PreprocessedToken: 元表記・位置・必要な比較表記を持つトークンを表す
//	  ]
//	  フィールド: [
//	    Kind: 字句カテゴリ
//	    Lexeme: 常に保持する元のソース表記
//	    Canonical: 識別子とキーワードだけに設定する比較表記
//	    Line: 元ソースの行番号
//	    Column: 元ソースの列番号
//	  ]
//	  補足: [
//	    リテラル内容は識別子の大文字小文字正規化から除外する
//	  ]
//	}
type PreprocessedToken struct {
	Kind      TokenKind
	Lexeme    string
	Canonical string
	Line      int
	Column    int
}

func newPreprocessedToken(token Token) (PreprocessedToken, error) {
	result := PreprocessedToken{
		Kind: token.Kind, Lexeme: token.Lexeme,
		Line: token.Line, Column: token.Column,
	}
	if token.Kind != TokenIdentifier && token.Kind != TokenKeyword {
		return result, nil
	}
	canonical, err := CanonicalizeIdentifier(token.Lexeme)
	if err != nil {
		return PreprocessedToken{}, err
	}
	result.Canonical = canonical
	return result, nil
}
