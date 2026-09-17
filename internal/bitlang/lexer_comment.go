package bitlang

import "fmt"

func (l *Lexer) skipComment() (bool, error) {
	if !l.hasNext() || l.peek() != '/' {
		return false, nil
	}
	switch l.peekNext() {
	case '/':
		l.skipLineComment()
		return true, nil
	case '*':
		return true, l.skipBlockComment()
	default:
		return false, nil
	}
}

func (l *Lexer) skipLineComment() {
	l.advance()
	l.advance()
	for !l.atEnd() && l.peek() != '\n' {
		l.advance()
	}
}

func (l *Lexer) skipBlockComment() error {
	startLine := l.line
	startColumn := l.column
	l.advance()
	l.advance()
	for !l.atEnd() {
		if l.hasNext() && l.peek() == '*' && l.peekNext() == '/' {
			l.advance()
			l.advance()
			return nil
		}
		l.advance()
	}
	return fmt.Errorf("unterminated block comment at %d:%d", startLine, startColumn)
}

func (l *Lexer) hasNext() bool {
	return l.index+1 < len(l.source.Text)
}

func (l *Lexer) peekNext() byte {
	return l.source.Text[l.index+1]
}
