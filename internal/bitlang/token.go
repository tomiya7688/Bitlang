package bitlang

//	{
//	  責務: [
//	    Token: Bitlang ソース中の一字句と元位置を表す
//	  ]
//	  フィールド: [
//	    Kind: 字句カテゴリ
//	    Lexeme: 変更しないソース上の表記
//	    Line: 1 始まりの行番号
//	    Column: UTF-8 バイト単位で数える 1 始まり列番号
//	  ]
//	  補足: [
//	    識別子の比較用正規化を別段階に分け、文字・文字列リテラルを変更しない
//	  ]
//	}
type Token struct {
	Kind   TokenKind
	Lexeme string
	Line   int
	Column int
}

//	{
//	  責務: [
//	    NewToken: 字句カテゴリ・元表記・位置からトークン値を作る
//	  ]
//	  引数: [
//	    kind: 字句カテゴリ
//	    lexeme: ソース上の元表記
//	    line: 1 始まりの行番号
//	    column: 1 始まりの列番号
//	  ]
//	  戻り値: [
//	    token: 指定された情報をそのまま保持する値
//	  ]
//	}
func NewToken(kind TokenKind, lexeme string, line int, column int) Token {
	return Token{Kind: kind, Lexeme: lexeme, Line: line, Column: column}
}
