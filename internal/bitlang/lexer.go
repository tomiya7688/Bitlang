package bitlang

import "fmt"

// Lexer converts SourceText into a sequence of lexical tokens.
type Lexer struct {
	source SourceText
	index  int
	line   int
	column int
}

// NewLexer creates a lexer positioned at the first byte of the source.
func NewLexer(source SourceText) *Lexer {
	return &Lexer{source: source, line: 1, column: 1}
}

// Lex scans the complete source and returns its token sequence.
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
	return l.scanSymbol(), nil
}

func (l *Lexer) scanIdentifier() Token {
	startIndex := l.index
	startLine := l.line
	startColumn := l.column
	for !l.atEnd() && isIdentifierContinue(l.peek()) {
		l.advance()
	}
	return NewToken(TokenIdentifier, l.source.Text[startIndex:l.index], startLine, startColumn)
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
