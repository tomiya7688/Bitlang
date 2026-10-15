package bitlang

import "strings"

//	{
//	  責務: [
//	    keywordKinds: 字句解析で予約語として扱う語を定義する
//	  ]
//	  補足: [
//	    走査処理から分離し、別実装でも同じ予約語集合を再現できるようにする
//	  ]
//	}
var keywordKinds = map[string]TokenKind{
	"namespace": TokenKeyword,
	"using":     TokenKeyword,
	"class":     TokenKeyword,
	"struct":    TokenKeyword,
	"interface": TokenKeyword,
	"enum":      TokenKeyword,
	"public":    TokenKeyword,
	"private":   TokenKeyword,
	"protected": TokenKeyword,
	"internal":  TokenKeyword,
	"static":    TokenKeyword,
	"readonly":  TokenKeyword,
	"const":     TokenKeyword,
	"get":       TokenKeyword,
	"set":       TokenKeyword,
	"init":      TokenKeyword,
	"direct":    TokenKeyword,
	"virtual":   TokenKeyword,
	"override":  TokenKeyword,
	"abstract":  TokenKeyword,
	"sealed":    TokenKeyword,
	"new":       TokenKeyword,
	"return":    TokenKeyword,
	"if":        TokenKeyword,
	"else":      TokenKeyword,
	"while":     TokenKeyword,
	"for":       TokenKeyword,
	"true":      TokenKeyword,
	"false":     TokenKeyword,
	"null":      TokenKeyword,
	"this":      TokenKeyword,
}

func keywordKind(lexeme string) (TokenKind, bool) {
	kind, found := keywordKinds[strings.ToLower(lexeme)]
	return kind, found
}
