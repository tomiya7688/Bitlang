package bitlang

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSourceClassesPreservesGenericSyntax(t *testing.T) {
	source := NewSourceText("generic.bit", `PUBLIC CLASS Box<T> {
	public Dictionary<Str, List<T>> Values { get; }
	public T Convert<U>(U input) { return input; }
}`)
	lexical, err := LexicalPreprocess(source)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseSourceClassesInSourceResult(NewSourceParseResult(lexical))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Classes) != 1 {
		t.Fatalf("class count = %d, want 1", len(parsed.Classes))
	}
	class := parsed.Classes[0]
	if len(class.Modifiers) != 1 || class.Modifiers[0] != "PUBLIC" {
		t.Fatalf("class modifiers = %#v", class.Modifiers)
	}
	if got := tokenLexemes(class.TypeParameters); !reflect.DeepEqual(got, []string{"<", "T", ">"}) {
		t.Fatalf("type parameters = %#v", got)
	}
	if len(class.Members) != 2 {
		t.Fatalf("member count = %d, want 2", len(class.Members))
	}
	property := class.Members[0]
	if got := tokenLexemes(property.TypeArguments); !reflect.DeepEqual(got, []string{"<", "Str", ",", "List", "<", "T", ">", ">"}) {
		t.Fatalf("property type arguments = %#v", got)
	}
	if len(property.Accessors) != 1 || property.Accessors[0].Spelling != "get" {
		t.Fatalf("property accessors = %#v", property.Accessors)
	}
	method := class.Members[1]
	if got := tokenLexemes(method.MethodTypeParameters); !reflect.DeepEqual(got, []string{"<", "U", ">"}) {
		t.Fatalf("method type parameters = %#v", got)
	}
	if len(method.Parameters) != 1 || method.Parameters[0].Name.Spelling != "input" {
		t.Fatalf("method parameters = %#v", method.Parameters)
	}
	if got := tokenLexemes(method.Body); !reflect.DeepEqual(got, []string{"return", "input", ";"}) {
		t.Fatalf("method body = %#v", got)
	}
}

func TestParseSourceClassesRejectsMalformedGenericLists(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{"empty parameters", "class Box<> {}"},
		{"missing parameter", "class Box<,T> {}"},
		{"trailing parameter comma", "class Box<T,> {}"},
		{"missing parameter comma", "class Box<T U> {}"},
		{"qualified parameter", "class Box<A.T> {}"},
		{"nested parameter", "class Box<List<T>> {}"},
		{"keyword parameter", "class Box<public> {}"},
		{"method parameter", "class Box { T Convert<U,>() {} }"},
		{"empty arguments", "class Box { List<> values; }"},
		{"missing argument", "class Box { List<,T> values; }"},
		{"trailing argument comma", "class Box { List<T,> values; }"},
		{"missing argument comma", "class Box { List<T U> values; }"},
		{"empty nested arguments", "class Box { List<List<>> values; }"},
		{"literal argument", "class Box { List<1> values; }"},
		{"unfinished qualified argument", "class Box { List<A.> values; }"},
		{"generic qualified suffix", "class Box { List<A<T>.B> values; }"},
		{"unterminated list", "class Box { List<T values; }"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			source, err := LexicalPreprocess(NewSourceText("invalid.bit", test.text))
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ParseSourceClassesInSourceResult(NewSourceParseResult(source))
			if err == nil {
				t.Fatalf("accepted malformed generic syntax: %s", test.text)
			}
			if len(parsed.Classes) != 0 || !strings.Contains(err.Error(), " at 1:") {
				t.Fatalf("expected no partial classes and a source-position error, got %#v, %v", parsed.Classes, err)
			}
		})
	}
}

func TestParseSourceClassesRejectsUnbalancedBodies(t *testing.T) {
	for _, text := range []string{
		"class Box { T Run() { call(]; } }",
		"class Box { T Run() { call(; } }",
		"class Box { T Run() { values[0); } }",
		"class Box { Box() { call(; } }",
	} {
		source, err := LexicalPreprocess(NewSourceText("invalid.bit", text))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseSourceClassesInSourceResult(NewSourceParseResult(source)); err == nil {
			t.Errorf("accepted unbalanced body: %s", text)
		}
	}
}

func TestParseSourceClassesPreservesBalancedBodiesAndMetadata(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText("balanced.bit", `class Box<T, U> {
	Ns.Map<Ns.Key, List<U>> Values = create(values[0], { "}" });
	Box() { call(values[0], { ")" }); }
	T Convert<V, W>(Ns.List<V> input) { { call(input[0], "}"); } }
}`))
	if err != nil {
		t.Fatal(err)
	}
	scope, err := NewPreprocessedScope("member", nil)
	if err != nil {
		t.Fatal(err)
	}
	parseResult := NewSourceParseResult(source)
	parseResult.Scope = &scope
	parseResult.Declarations = []PreprocessedDeclaration{{Kind: "existing"}}
	parsed, err := ParseSourceClassesInSourceResult(parseResult)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Source.Path != source.Path || parsed.Scope != parseResult.Scope ||
		!reflect.DeepEqual(parsed.Source.Tokens, source.Tokens) || !reflect.DeepEqual(parsed.Declarations, parseResult.Declarations) {
		t.Fatal("parsing replaced source metadata")
	}
	if len(parsed.Classes) != 1 || len(parsed.Classes[0].Members) != 3 {
		t.Fatalf("unexpected classes: %#v", parsed.Classes)
	}
	members := parsed.Classes[0].Members
	if members[0].Type.Spelling != "Ns.Map" || members[1].Kind != SourceConstructorMember || members[2].Kind != SourceMethodMember {
		t.Fatalf("unexpected members: %#v", members)
	}
	wantBodies := [][]string{
		{"create", "(", "values", "[", "0", "]", ",", "{", `"}"`, "}", ")"},
		{"call", "(", "values", "[", "0", "]", ",", "{", `")"`, "}", ")", ";"},
		{"{", "call", "(", "input", "[", "0", "]", ",", `"}"`, ")", ";", "}"},
	}
	sequences := [][]PreprocessedToken{members[0].Initializer, members[1].Body, members[2].Body}
	for index, sequence := range sequences {
		if got := tokenLexemes(sequence); !reflect.DeepEqual(got, wantBodies[index]) {
			t.Fatalf("body %d = %#v", index, got)
		}
	}
	sequences = append(sequences, parsed.Classes[0].TypeParameters, members[0].TypeArguments,
		members[2].MethodTypeParameters, members[2].Parameters[0].TypeArguments)
	originals := make(map[[2]int]PreprocessedToken)
	for _, token := range source.Tokens {
		originals[[2]int{token.Line, token.Column}] = token
	}
	for _, sequence := range sequences {
		for _, token := range sequence {
			if original, found := originals[[2]int{token.Line, token.Column}]; !found || token != original {
				t.Fatalf("token or source position changed: %#v", token)
			}
		}
	}
}

func tokenLexemes(tokens []PreprocessedToken) []string {
	lexemes := make([]string, len(tokens))
	for index, token := range tokens {
		lexemes[index] = token.Lexeme
	}
	return lexemes
}
