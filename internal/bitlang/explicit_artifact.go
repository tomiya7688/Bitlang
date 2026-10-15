package bitlang

import "fmt"

//	{
//	  責務: [
//	    ExplicitArtifact: 完全明示プロパティを検証済みの宣言成果物を保持する
//	  ]
//	  フィールド: [
//	    source: 元入力全文と字句情報
//	    declarations: 適用可能な全propertyを検証済みの宣言列
//	    scope: 宣言解析に使った意味スコープ
//	  ]
//	  補足: [
//	    フィールドを非公開にし、生成後のproperty欠落や外部変更を防ぐ
//	  ]
//	}
type ExplicitArtifact struct {
	source       LexicalSource
	declarations []PreprocessedDeclaration
	scope        *PreprocessedScope
}

//	{
//	  責務: [
//	    NewExplicitArtifact: 現行の完全明示宣言だけからExplicit成果物を作る
//	  ]
//	  処理: [
//	    1: 原文を保持したまま字句解析する
//	    2: 現行の宣言仕様で宣言を解析する
//	    3: 適用可能な全propertyとスコープを検証する
//	    4: 元入力と解析結果を独立した成果物として保持する
//	  ]
//	  引数: [
//	    specs: propertyおよび宣言の検証済み仕様
//	    source: 完全明示形式として解析する入力
//	    scope: 宣言解析に使う意味スコープ
//	  ]
//	  戻り値: [
//	    artifact: 原文と完全明示宣言を保持する成果物
//	    error: 字句・宣言・property・scopeのいずれかが不正な場合の失敗
//	  ]
//	  補足: [
//	    現時点の宣言文法のみを受け入れ、class/member糖衣は後続の正規化まで拒否する
//	    省略propertyの推測と旧形式の移行は行わず、不完全な入力を拒否する
//	  ]
//	}
func NewExplicitArtifact(specs SpecificationSet, source SourceText, scope *PreprocessedScope) (ExplicitArtifact, error) {
	lexicalSource, err := LexicalPreprocess(source)
	if err != nil {
		return ExplicitArtifact{}, fmt.Errorf("lex source for Explicit artifact: %w", err)
	}
	clonedScope, err := cloneExplicitScope(scope, make(map[*PreprocessedScope]bool))
	if err != nil {
		return ExplicitArtifact{}, err
	}
	declarations, err := parseExplicitDeclarations(specs, clonedScope, lexicalSource.Tokens)
	if err != nil {
		return ExplicitArtifact{}, fmt.Errorf("parse Explicit declarations: %w", err)
	}
	return ExplicitArtifact{
		source:       cloneLexicalSource(lexicalSource),
		declarations: cloneExplicitDeclarations(declarations),
		scope:        clonedScope,
	}, nil
}

func parseExplicitDeclarations(specs SpecificationSet, scope *PreprocessedScope, tokens []PreprocessedToken) ([]PreprocessedDeclaration, error) {
	groups, err := splitMixedDeclarationTokens(specs.Declarations, tokens)
	if err != nil {
		return nil, err
	}
	declarations := make([]PreprocessedDeclaration, 0, len(groups))
	for _, group := range groups {
		var matches []PreprocessedDeclaration
		for _, kind := range specs.Declarations.Kinds {
			if !declarationKindAppliesToScope(kind, scope) {
				continue
			}
			declaration, err := parseExplicitDeclaration(specs, kind, group)
			if err == nil {
				matches = append(matches, declaration)
			}
		}
		if len(matches) != 1 {
			if len(matches) == 0 {
				return nil, fmt.Errorf("%d:%d: declaration matches no configured kind", group[0].Line, group[0].Column)
			}
			return nil, fmt.Errorf("%d:%d: declaration matches multiple configured kinds", group[0].Line, group[0].Column)
		}
		declarations = append(declarations, matches[0])
	}
	return declarations, nil
}

func parseExplicitDeclaration(specs SpecificationSet, kind DeclarationKindSpec, tokens []PreprocessedToken) (PreprocessedDeclaration, error) {
	body := declarationBody(tokens)
	if len(body) < 2 || body[len(body)-1].Lexeme != kind.Terminator {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration does not satisfy configured layout")
	}
	parts, err := parseDeclarationLayout(kind.Layout, body[:len(body)-1])
	if err != nil || parts.typeToken.Kind != TokenIdentifier || parts.nameToken.Kind != TokenIdentifier {
		return PreprocessedDeclaration{}, fmt.Errorf("declaration does not satisfy configured layout")
	}
	properties := make([]PreprocessedProperty, 0, len(parts.properties))
	for _, token := range parts.properties {
		if token.Kind != TokenIdentifier && token.Kind != TokenKeyword {
			return PreprocessedDeclaration{}, fmt.Errorf("property %q must be an identifier or keyword", token.Lexeme)
		}
		properties = append(properties, PreprocessedProperty(token.Lexeme))
	}
	properties, err = ApplySourcePropertyDefaults(specs.Properties, kind.PropertyTarget, properties)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	if err := ValidateProperties(specs.Properties, kind.PropertyTarget, properties); err != nil {
		return PreprocessedDeclaration{}, err
	}
	name, err := NewCanonicalName(parts.nameToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	typeName, err := NewCanonicalName(parts.typeToken.Lexeme)
	if err != nil {
		return PreprocessedDeclaration{}, err
	}
	return PreprocessedDeclaration{Kind: kind.Name, Name: name, Type: typeName, Properties: properties, Line: parts.nameToken.Line, Column: parts.nameToken.Column}, nil
}

//	{
//	  責務: [
//	    Source: 外部変更を防ぐ複製として元ソースと字句情報を返す
//	  ]
//	  引数: [
//	    a: 読み取り対象のExplicit成果物
//	  ]
//	  戻り値: [
//	    source: 原文と字句情報の独立した複製
//	  ]
//	}
func (a ExplicitArtifact) Source() LexicalSource {
	return cloneLexicalSource(a.source)
}

//	{
//	  責務: [
//	    Declarations: 完全明示プロパティを持つ宣言の独立した複製を返す
//	  ]
//	  引数: [
//	    a: 読み取り対象のExplicit成果物
//	  ]
//	  戻り値: [
//	    declarations: 検証済み宣言列の独立した複製
//	  ]
//	}
func (a ExplicitArtifact) Declarations() []PreprocessedDeclaration {
	return cloneExplicitDeclarations(a.declarations)
}

//	{
//	  責務: [
//	    Scope: 宣言解析に使ったスコープの独立した複製を返す
//	  ]
//	  引数: [
//	    a: 読み取り対象のExplicit成果物
//	  ]
//	  戻り値: [
//	    scope: スコープの複製またはスコープ未指定を表すnil
//	  ]
//	}
func (a ExplicitArtifact) Scope() *PreprocessedScope {
	scope, _ := cloneExplicitScope(a.scope, make(map[*PreprocessedScope]bool))
	return scope
}

//	{
//	  責務: [
//	    AsArtifact: 型付きExplicit成果物を正規stage identity付き成果物へ包む
//	  ]
//	  引数: [
//	    a: 完全明示プロパティを検証済みの成果物
//	  ]
//	  戻り値: [
//	    artifact: ArtifactExplicit identityと型付きpayloadを持つ成果物
//	  ]
//	}
func (a ExplicitArtifact) AsArtifact() Artifact {
	return Artifact{Kind: ArtifactExplicit, Payload: a}
}

func cloneLexicalSource(source LexicalSource) LexicalSource {
	cloned := source
	cloned.Tokens = append([]PreprocessedToken(nil), source.Tokens...)
	return cloned
}

func cloneExplicitDeclarations(declarations []PreprocessedDeclaration) []PreprocessedDeclaration {
	cloned := make([]PreprocessedDeclaration, len(declarations))
	for index, declaration := range declarations {
		cloned[index] = declaration
		cloned[index].Properties = append([]PreprocessedProperty(nil), declaration.Properties...)
	}
	return cloned
}

func cloneExplicitScope(scope *PreprocessedScope, seen map[*PreprocessedScope]bool) (*PreprocessedScope, error) {
	if scope == nil {
		return nil, nil
	}
	if seen[scope] {
		return nil, fmt.Errorf("Explicit scope contains a parent cycle")
	}
	seen[scope] = true
	parent, err := cloneExplicitScope(scope.Parent, seen)
	delete(seen, scope)
	if err != nil {
		return nil, err
	}
	return &PreprocessedScope{Context: scope.Context, Parent: parent}, nil
}
