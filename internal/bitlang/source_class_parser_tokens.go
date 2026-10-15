package bitlang

import "fmt"

func (p *sourceClassParser) parseBlock() ([]PreprocessedToken, error) {
	if _, err := p.expectSymbol("{"); err != nil {
		return nil, err
	}
	stack := []string{"}"}
	body := make([]PreprocessedToken, 0)
	for len(stack) > 0 {
		token := p.current()
		if token.Kind == TokenEOF {
			return nil, p.errorAtCurrent("unterminated method body")
		}
		if token.Kind == TokenSymbol {
			if isOpeningDelimiter(token.Lexeme) {
				stack = append(stack, matchingDelimiter(token.Lexeme))
			} else if isClosingDelimiter(token.Lexeme) {
				if stack[len(stack)-1] != token.Lexeme {
					return nil, p.errorAt(token, fmt.Sprintf("unmatched delimiter %q in method body", token.Lexeme))
				}
				stack = stack[:len(stack)-1]
			}
		}
		p.advance()
		if len(stack) > 0 {
			body = append(body, token)
		}
	}
	return body, nil
}

func (p *sourceClassParser) parseUntilSemicolon() ([]PreprocessedToken, error) {
	body := make([]PreprocessedToken, 0)
	stack := make([]string, 0)
	for {
		token := p.current()
		if token.Kind == TokenEOF {
			return nil, p.errorAtCurrent("expected semicolon after field initializer")
		}
		if token.Kind == TokenSymbol {
			if token.Lexeme == ";" && len(stack) == 0 {
				return body, nil
			}
			if isOpeningDelimiter(token.Lexeme) {
				stack = append(stack, matchingDelimiter(token.Lexeme))
			} else if isClosingDelimiter(token.Lexeme) {
				if len(stack) == 0 || stack[len(stack)-1] != token.Lexeme {
					return nil, p.errorAt(token, fmt.Sprintf("unmatched delimiter %q in field initializer", token.Lexeme))
				}
				stack = stack[:len(stack)-1]
			}
		}
		body = append(body, token)
		p.advance()
	}
}

func (p *sourceClassParser) parseModifiers() []PreprocessedToken {
	modifiers := make([]PreprocessedToken, 0)
	for isModifier(p.current()) {
		modifiers = append(modifiers, p.current())
		p.advance()
	}
	return modifiers
}

//	{
//	  責務: [
//	    sourceModifierSpellings: 修飾子トークンをアダプター向け元表記列へ写す
//	  ]
//	  引数: [
//	    tokens: 元表記と位置を持つ修飾子トークン列
//	  ]
//	  戻り値: [
//	    spellings: 入力順・重複を保ったソース上の文字列列
//	  ]
//	  補足: [
//	    正規化・並べ替え・重複除去を行わない
//	  ]
//	}
func sourceModifierSpellings(tokens []PreprocessedToken) []string {
	spellings := make([]string, len(tokens))
	for index, token := range tokens {
		spellings[index] = token.Lexeme
	}
	return spellings
}

func isModifier(token PreprocessedToken) bool {
	if token.Kind != TokenIdentifier && token.Kind != TokenKeyword {
		return false
	}
	switch token.Canonical {
	case "public", "private", "protected", "internal", "static", "readonly", "const", "abstract", "sealed", "virtual", "override", "direct":
		return true
	default:
		return false
	}
}

func sourceName(token PreprocessedToken) (CanonicalName, error) {
	if token.Canonical == "" {
		return CanonicalName{}, fmt.Errorf("identifier %q has no canonical form at %d:%d", token.Lexeme, token.Line, token.Column)
	}
	return CanonicalName{Spelling: token.Lexeme, Canonical: token.Canonical}, nil
}

func (p *sourceClassParser) expectIdentifier(description string) (PreprocessedToken, error) {
	token := p.current()
	if token.Kind != TokenIdentifier {
		return PreprocessedToken{}, p.errorAt(token, fmt.Sprintf("expected %s", description))
	}
	p.advance()
	return token, nil
}

func (p *sourceClassParser) expectIdentifierOrKeyword(description string) (PreprocessedToken, error) {
	token := p.current()
	if token.Kind != TokenIdentifier && token.Kind != TokenKeyword {
		return PreprocessedToken{}, p.errorAt(token, fmt.Sprintf("expected %s", description))
	}
	p.advance()
	return token, nil
}

func (p *sourceClassParser) expectWord(word string) (PreprocessedToken, error) {
	token := p.current()
	if (token.Kind != TokenIdentifier && token.Kind != TokenKeyword) || token.Canonical != word {
		return PreprocessedToken{}, p.errorAt(token, fmt.Sprintf("expected %q", word))
	}
	p.advance()
	return token, nil
}

func (p *sourceClassParser) expectSymbol(symbol string) (PreprocessedToken, error) {
	token := p.current()
	if token.Kind != TokenSymbol || token.Lexeme != symbol {
		return PreprocessedToken{}, p.errorAt(token, fmt.Sprintf("expected %q", symbol))
	}
	p.advance()
	return token, nil
}

func (p *sourceClassParser) isSymbol(symbol string) bool {
	return p.isSymbolAt(p.index, symbol)
}

func (p *sourceClassParser) isSymbolAt(index int, symbol string) bool {
	if index >= len(p.tokens) {
		return false
	}
	token := p.tokens[index]
	return token.Kind == TokenSymbol && token.Lexeme == symbol
}

func (p *sourceClassParser) current() PreprocessedToken {
	if p.index >= len(p.tokens) {
		return PreprocessedToken{Kind: TokenEOF}
	}
	return p.tokens[p.index]
}

func (p *sourceClassParser) advance() {
	if p.index < len(p.tokens) {
		p.index++
	}
}

func (p *sourceClassParser) errorAtCurrent(message string) error {
	return p.errorAt(p.current(), message)
}

func (p *sourceClassParser) errorAt(token PreprocessedToken, message string) error {
	return fmt.Errorf("%s at %d:%d", message, token.Line, token.Column)
}

func isOpeningDelimiter(value string) bool {
	return value == "(" || value == "[" || value == "{"
}

func isClosingDelimiter(value string) bool {
	return value == ")" || value == "]" || value == "}"
}

func matchingDelimiter(value string) string {
	switch value {
	case "(":
		return ")"
	case "[":
		return "]"
	default:
		return "}"
	}
}
