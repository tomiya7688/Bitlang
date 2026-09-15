package bitlang

import "testing"

func FuzzLexerNeverPanics(f *testing.F) {
	seeds := []string{
		"",
		"identifier 123 symbol",
		"\"quoted string\"",
		"'c'",
		"\"unterminated",
		"'\\\\'",
		"line1\nline2\r\nline3",
		string([]byte{0x00, 0xff, 0x80, '\n'}),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		tokens, err := NewLexer(NewSourceText("<fuzz>", text)).Lex()
		if err != nil {
			return
		}
		if len(tokens) == 0 {
			t.Fatal("successful lex returned no EOF token")
		}
		if tokens[len(tokens)-1].Kind != TokenEOF {
			t.Fatal("successful lex did not end with EOF")
		}
		for _, token := range tokens {
			if token.Line < 1 || token.Column < 1 {
				t.Fatalf("invalid token position: %d:%d", token.Line, token.Column)
			}
		}
	})
}
