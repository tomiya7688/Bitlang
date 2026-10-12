package bitlang

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSourceClassesPreservesModifierTokens(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText("modifiers.bit", `PUBLIC
sealed class Box {
    private readonly Int value;
    protected Box(CONST Int input) {}
    internal static Direct Int Run(readonly Int arg) {}
    public Int Count { get; private protected set; }
    PUBLIC public Int Duplicate;
    Int Plain;
}`))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseSourceClassesInSourceResult(NewSourceParseResult(source))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Classes) != 1 || len(parsed.Classes[0].Members) != 6 {
		t.Fatalf("unexpected declarations: %#v", parsed.Classes)
	}
	class := parsed.Classes[0]
	if len(class.Members[1].Parameters) != 1 || len(class.Members[2].Parameters) != 1 || len(class.Members[3].Accessors) != 2 {
		t.Fatalf("unexpected parameters or accessors: %#v", class.Members)
	}
	constructor := class.Members[1].Parameters[0]
	method := class.Members[2].Parameters[0]
	get, set := class.Members[3].Accessors[0], class.Members[3].Accessors[1]
	cases := []struct {
		name      string
		spellings []string
		tokens    []PreprocessedToken
		want      []string
		positions [][2]int
	}{
		{"class", class.Modifiers, class.ModifierTokens, []string{"PUBLIC", "sealed"}, [][2]int{{1, 1}, {2, 1}}},
		{"field", class.Members[0].Modifiers, class.Members[0].ModifierTokens, []string{"private", "readonly"}, [][2]int{{3, 5}, {3, 13}}},
		{"constructor", class.Members[1].Modifiers, class.Members[1].ModifierTokens, []string{"protected"}, [][2]int{{4, 5}}},
		{"constructor parameter", constructor.Modifiers, constructor.ModifierTokens, []string{"CONST"}, [][2]int{{4, 19}}},
		{"method", class.Members[2].Modifiers, class.Members[2].ModifierTokens, []string{"internal", "static", "Direct"}, [][2]int{{5, 5}, {5, 14}, {5, 21}}},
		{"method parameter", method.Modifiers, method.ModifierTokens, []string{"readonly"}, [][2]int{{5, 36}}},
		{"property", class.Members[3].Modifiers, class.Members[3].ModifierTokens, []string{"public"}, [][2]int{{6, 5}}},
		{"get", get.Modifiers, get.ModifierTokens, []string{}, [][2]int{}},
		{"set", set.Modifiers, set.ModifierTokens, []string{"private", "protected"}, [][2]int{{6, 29}, {6, 37}}},
		{"duplicates", class.Members[4].Modifiers, class.Members[4].ModifierTokens, []string{"PUBLIC", "public"}, [][2]int{{7, 5}, {7, 12}}},
		{"unmodified", class.Members[5].Modifiers, class.Members[5].ModifierTokens, []string{}, [][2]int{}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.spellings, test.want) || !reflect.DeepEqual(tokenLexemes(test.tokens), test.want) {
				t.Fatalf("modifier spelling/order changed: %#v, %#v", test.spellings, test.tokens)
			}
			for index, token := range test.tokens {
				want := PreprocessedToken{
					Kind: TokenKeyword, Lexeme: test.want[index], Canonical: strings.ToLower(test.want[index]),
					Line: test.positions[index][0], Column: test.positions[index][1],
				}
				if token != want {
					t.Fatalf("modifier token = %#v, want %#v", token, want)
				}
			}
		})
	}
	if !reflect.DeepEqual(parsed.Source.Tokens, source.Tokens) {
		t.Fatal("source tokens changed during modifier parsing")
	}
}
