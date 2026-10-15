package bitlang

import "fmt"

//	{
//	  責務: [
//	    Lexer: ソーステキストを位置情報付きの字句トークンへ変換する
//	  ]
//	  フィールド: [
//	    source: 元のソースパスと変更していない入力文字列
//	    index: 現在の入力バイト位置
//	    line: 現在位置の 1 始まり行番号
//	    column: UTF-8 バイト単位で数える現在位置の 1 始まり列番号
//	  ]
//	}
type Lexer struct {
	source SourceText
	index  int
	line   int
	column int
}

//	{
//	  責務: [
//	    NewLexer: 入力先頭から走査する Lexer を作る
//	  ]
//	  引数: [
//	    source: 字句解析する元ソース
//	  ]
//	  戻り値: [
//	    lexer: 行・列を 1 から数える新しい走査器
//	  ]
//	}
func NewLexer(source SourceText) *Lexer {
	return &Lexer{source: source, line: 1, column: 1}
}

//	{
//	  責務: [
//	    Lex: ソース全体を走査し、元表記を保ったトークン列を返す
//	  ]
//	  処理: [
//	    1: 空白とコメントを読み飛ばす
//	    2: 識別子・数値・リテラル・記号をトークン化する
//	    3: ASCII 外の通常トークンを誤変換せず拒否する
//	    4: 終端トークンを追加する
//	  ]
//	  引数: [
//	    l: 元ソースと現在の走査位置を持つ Lexer
//	  ]
//	  戻り値: [
//	    tokens: ソース順の字句トークン列
//	    error: 不正バイトや閉じていないコメント等の失敗
//	  ]
//	}
func (l *Lexer) Lex() ([]Token, error) {
	tokens := make([]Token, 0)
	for !l.atEnd() {
		if l.skipWhitespace() {
			continue
		}
		skipped, err := l.skipComment()
		if err != nil {
			return nil, err
		}
		if skipped {
			continue
		}
		token, err := l.scanToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, NewToken(TokenEOF, "", l.line, l.column))
	return tokens, nil
}

func (l *Lexer) scanToken() (Token, error) {
	current := l.peek()
	if isIdentifierStart(current) {
		return l.scanIdentifier(), nil
	}
	if isDecimalDigit(current) {
		return l.scanNumber(), nil
	}
	if current == '"' {
		return l.scanQuotedLiteral('"', TokenString)
	}
	if current == '\'' {
		return l.scanQuotedLiteral('\'', TokenCharacter)
	}
	if current >= 0x80 {
		return Token{}, fmt.Errorf("unsupported non-ASCII source byte 0x%02X at %d:%d", current, l.line, l.column)
	}
	return l.scanSymbol(), nil
}

func (l *Lexer) scanIdentifier() Token {
	startIndex := l.index
	startLine := l.line
	startColumn := l.column
	for !l.atEnd() && isIdentifierContinue(l.peek()) {
		l.advance()
	}
	lexeme := l.source.Text[startIndex:l.index]
	kind := TokenIdentifier
	if keyword, found := keywordKind(lexeme); found {
		kind = keyword
	}
	return NewToken(kind, lexeme, startLine, startColumn)
}
func (l *Lexer) scanNumber() Token {
	startIndex := l.index
	startLine := l.line
	startColumn := l.column
	for !l.atEnd() && isDecimalDigit(l.peek()) {
		l.advance()
	}
	return NewToken(TokenNumber, l.source.Text[startIndex:l.index], startLine, startColumn)
}

func (l *Lexer) scanQuotedLiteral(quote byte, kind TokenKind) (Token, error) {
	startIndex := l.index
	startLine := l.line
	startColumn := l.column
	l.advance()
	for !l.atEnd() {
		current := l.advance()
		if current == '\\' && !l.atEnd() {
			l.advance()
			continue
		}
		if current == quote {
			return NewToken(kind, l.source.Text[startIndex:l.index], startLine, startColumn), nil
		}
	}
	return Token{}, fmt.Errorf("unterminated %s literal at %d:%d", kind, startLine, startColumn)
}

func (l *Lexer) scanSymbol() Token {
	startLine := l.line
	startColumn := l.column
	value := l.advance()
	return NewToken(TokenSymbol, string(value), startLine, startColumn)
}

func (l *Lexer) skipWhitespace() bool {
	consumed := false
	for !l.atEnd() {
		switch l.peek() {
		case ' ', '\t', '\r', '\n':
			l.advance()
			consumed = true
		default:
			return consumed
		}
	}
	return consumed
}

func (l *Lexer) advance() byte {
	value := l.source.Text[l.index]
	l.index++
	if value == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return value
}

func (l *Lexer) peek() byte {
	return l.source.Text[l.index]
}

func (l *Lexer) atEnd() bool {
	return l.index >= len(l.source.Text)
}

func isIdentifierStart(value byte) bool {
	return value == '_' || value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z'
}

func isIdentifierContinue(value byte) bool {
	return isIdentifierStart(value) || isDecimalDigit(value)
}

func isDecimalDigit(value byte) bool {
	return value >= '0' && value <= '9'
}
