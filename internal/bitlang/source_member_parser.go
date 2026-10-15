package bitlang

import "fmt"

func (p *sourceClassParser) parseMember(className CanonicalName) (SourceMember, error) {
	modifiers := p.parseModifiers()
	first, err := p.expectIdentifier("member type or constructor name")
	if err != nil {
		return SourceMember{}, err
	}
	if first.Canonical == className.Canonical && p.isSymbolAt(p.index, "(") {
		return p.parseConstructor(first, modifiers)
	}
	return p.parseTypedMember(first, modifiers)
}

func (p *sourceClassParser) parseConstructor(nameToken PreprocessedToken, modifiers []PreprocessedToken) (SourceMember, error) {
	name, err := sourceName(nameToken)
	if err != nil {
		return SourceMember{}, err
	}
	parameters, err := p.parseParameters()
	if err != nil {
		return SourceMember{}, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return SourceMember{}, err
	}
	return SourceMember{
		Kind: SourceConstructorMember, Name: name,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers,
		Parameters: parameters, Body: body, Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func (p *sourceClassParser) parseTypedMember(first PreprocessedToken, modifiers []PreprocessedToken) (SourceMember, error) {
	typeName, typeArguments, err := p.parseTypeAfterFirst(first)
	if err != nil {
		return SourceMember{}, err
	}
	nameToken, err := p.expectIdentifier("member name")
	if err != nil {
		return SourceMember{}, err
	}
	name, err := sourceName(nameToken)
	if err != nil {
		return SourceMember{}, err
	}
	methodTypeParameters := make([]PreprocessedToken, 0)
	if p.isSymbol("<") {
		methodTypeParameters, err = p.parseGenericTokenSequence(true)
		if err != nil {
			return SourceMember{}, err
		}
	}
	if p.isSymbol("{") {
		if len(methodTypeParameters) != 0 {
			return SourceMember{}, p.errorAt(nameToken, "property cannot declare method type parameters")
		}
		return p.parseProperty(name, typeName, typeArguments, nameToken, modifiers)
	}
	if p.isSymbol("(") {
		return p.parseMethod(name, typeName, typeArguments, methodTypeParameters, nameToken, modifiers)
	}
	if len(methodTypeParameters) != 0 {
		return SourceMember{}, p.errorAt(nameToken, "generic member must be a method")
	}
	return p.parseField(name, typeName, typeArguments, nameToken, modifiers)
}
func (p *sourceClassParser) parseProperty(name, typeName CanonicalName, typeArguments []PreprocessedToken, nameToken PreprocessedToken, modifiers []PreprocessedToken) (SourceMember, error) {
	accessors, err := p.parseAccessors()
	if err != nil {
		return SourceMember{}, err
	}
	return SourceMember{
		Kind: SourcePropertyMember, Name: name, Type: typeName, TypeArguments: typeArguments,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers, Accessors: accessors,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func (p *sourceClassParser) parseMethod(name, typeName CanonicalName, typeArguments, methodTypeParameters []PreprocessedToken, nameToken PreprocessedToken, modifiers []PreprocessedToken) (SourceMember, error) {
	parameters, err := p.parseParameters()
	if err != nil {
		return SourceMember{}, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return SourceMember{}, err
	}
	return SourceMember{
		Kind: SourceMethodMember, Name: name, Type: typeName, TypeArguments: typeArguments, MethodTypeParameters: methodTypeParameters,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers,
		Parameters: parameters, Body: body,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func (p *sourceClassParser) parseField(name, typeName CanonicalName, typeArguments []PreprocessedToken, nameToken PreprocessedToken, modifiers []PreprocessedToken) (SourceMember, error) {
	initializer := make([]PreprocessedToken, 0)
	if p.isSymbol("=") {
		p.advance()
		var err error
		initializer, err = p.parseUntilSemicolon()
		if err != nil {
			return SourceMember{}, err
		}
		if len(initializer) == 0 {
			return SourceMember{}, p.errorAtCurrent("field initializer must not be empty")
		}
	}
	if _, err := p.expectSymbol(";"); err != nil {
		return SourceMember{}, err
	}
	return SourceMember{
		Kind: SourceFieldMember, Name: name, Type: typeName, TypeArguments: typeArguments,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers, Initializer: initializer,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func (p *sourceClassParser) parseTypeAfterFirst(first PreprocessedToken) (CanonicalName, []PreprocessedToken, error) {
	name, err := sourceName(first)
	if err != nil {
		return CanonicalName{}, nil, err
	}
	for p.isSymbol(".") {
		p.advance()
		part, err := p.expectIdentifier("qualified type name")
		if err != nil {
			return CanonicalName{}, nil, err
		}
		if part.Canonical == "" {
			return CanonicalName{}, nil, p.errorAt(part, "qualified type name has no canonical form")
		}
		name.Spelling += "." + part.Lexeme
		name.Canonical += "." + part.Canonical
	}
	arguments := make([]PreprocessedToken, 0)
	if p.isSymbol("<") {
		arguments, err = p.parseGenericTokenSequence(false)
		if err != nil {
			return CanonicalName{}, nil, err
		}
	}
	return name, arguments, nil
}
func (p *sourceClassParser) parseParameters() ([]SourceParameter, error) {
	if _, err := p.expectSymbol("("); err != nil {
		return nil, err
	}
	parameters := make([]SourceParameter, 0)
	if p.isSymbol(")") {
		p.advance()
		return parameters, nil
	}
	for {
		parameter, err := p.parseParameter()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, parameter)
		if p.isSymbol(")") {
			p.advance()
			return parameters, nil
		}
		if _, err := p.expectSymbol(","); err != nil {
			return nil, err
		}
	}
}

func (p *sourceClassParser) parseParameter() (SourceParameter, error) {
	modifiers := p.parseModifiers()
	typeToken, err := p.expectIdentifier("parameter type")
	if err != nil {
		return SourceParameter{}, err
	}
	typeName, typeArguments, err := p.parseTypeAfterFirst(typeToken)
	if err != nil {
		return SourceParameter{}, err
	}
	nameToken, err := p.expectIdentifier("parameter name")
	if err != nil {
		return SourceParameter{}, err
	}
	name, err := sourceName(nameToken)
	if err != nil {
		return SourceParameter{}, err
	}
	return SourceParameter{
		Name: name, Type: typeName, TypeArguments: typeArguments,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers,
		Line: nameToken.Line, Column: nameToken.Column,
	}, nil
}

func (p *sourceClassParser) parseAccessors() ([]SourceAccessor, error) {
	if _, err := p.expectSymbol("{"); err != nil {
		return nil, err
	}
	accessors := make([]SourceAccessor, 0, 2)
	seen := make(map[string]bool)
	for !p.isSymbol("}") {
		if p.current().Kind == TokenEOF {
			return nil, p.errorAtCurrent("unterminated property body")
		}
		accessor, err := p.parseAccessor()
		if err != nil {
			return nil, err
		}
		if seen[accessor.Kind] {
			return nil, p.errorAt(PreprocessedToken{Lexeme: accessor.Spelling, Line: accessor.Line, Column: accessor.Column}, fmt.Sprintf("duplicate property accessor %q", accessor.Spelling))
		}
		seen[accessor.Kind] = true
		accessors = append(accessors, accessor)
	}
	if len(accessors) == 0 {
		return nil, p.errorAtCurrent("property requires at least one accessor")
	}
	p.advance()
	return accessors, nil
}

func (p *sourceClassParser) parseAccessor() (SourceAccessor, error) {
	modifiers := p.parseModifiers()
	kindToken, err := p.expectIdentifierOrKeyword("property accessor")
	if err != nil {
		return SourceAccessor{}, err
	}
	kind := kindToken.Canonical
	if kind != "get" && kind != "set" && kind != "init" {
		return SourceAccessor{}, p.errorAt(kindToken, fmt.Sprintf("unknown property accessor %q", kindToken.Lexeme))
	}
	if _, err := p.expectSymbol(";"); err != nil {
		return SourceAccessor{}, err
	}
	return SourceAccessor{
		Kind: kind, Spelling: kindToken.Lexeme,
		Modifiers: sourceModifierSpellings(modifiers), ModifierTokens: modifiers,
		Line: kindToken.Line, Column: kindToken.Column,
	}, nil
}
