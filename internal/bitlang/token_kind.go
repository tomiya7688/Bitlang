package bitlang

//	{
//	  責務: [
//	    TokenKind: Bitlang トークンの言語中立な字句カテゴリを識別する
//	  ]
//	  補足: [
//	    他のコンパイラー実装が Go 固有の内部構造に依存せず同じ列を再現できる
//	  ]
//	}
type TokenKind string

//	{
//	  責務: [
//	    TokenKind の定義値: 字句解析器が出力するカテゴリを列挙する
//	  ]
//	  項目: [
//	    TokenIdentifier: 識別子
//	    TokenKeyword: キーワード
//	    TokenNumber: 数値リテラル
//	    TokenString: 文字列リテラル
//	    TokenCharacter: 文字リテラル
//	    TokenSymbol: 記号
//	    TokenEOF: 入力終端
//	  ]
//	}
const (
	TokenIdentifier TokenKind = "identifier"
	TokenKeyword    TokenKind = "keyword"
	TokenNumber     TokenKind = "number"
	TokenString     TokenKind = "string"
	TokenCharacter  TokenKind = "character"
	TokenSymbol     TokenKind = "symbol"
	TokenEOF        TokenKind = "eof"
)
