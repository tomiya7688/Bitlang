package bitlang

import "fmt"

type declarationParts struct {
	properties []PreprocessedToken
	typeToken  PreprocessedToken
	nameToken  PreprocessedToken
}

func parseDeclarationLayout(layout []string, tokens []PreprocessedToken) (declarationParts, error) {
	propertiesIndex := declarationLayoutIndex(layout, "properties")
	if propertiesIndex < 0 {
		return declarationParts{}, fmt.Errorf("declaration layout has no properties component")
	}

	fixedBefore := propertiesIndex
	fixedAfter := len(layout) - propertiesIndex - 1
	if len(tokens) < fixedBefore+fixedAfter {
		return declarationParts{}, fmt.Errorf("declaration does not satisfy configured layout")
	}

	propertyEnd := len(tokens) - fixedAfter
	parts := declarationParts{
		properties: append([]PreprocessedToken(nil), tokens[fixedBefore:propertyEnd]...),
	}
	for layoutIndex, component := range layout {
		if component == "properties" {
			continue
		}
		tokenIndex := layoutIndex
		if layoutIndex > propertiesIndex {
			tokenIndex = propertyEnd + layoutIndex - propertiesIndex - 1
		}
		if err := assignDeclarationPart(&parts, component, tokens[tokenIndex]); err != nil {
			return declarationParts{}, err
		}
	}
	return parts, nil
}

func assignDeclarationPart(parts *declarationParts, component string, token PreprocessedToken) error {
	switch component {
	case "type":
		parts.typeToken = token
	case "name":
		parts.nameToken = token
	default:
		return fmt.Errorf("unsupported declaration layout component %q", component)
	}
	return nil
}

func validDeclarationLayout(layout []string) bool {
	if len(layout) != 3 {
		return false
	}
	return declarationLayoutIndex(layout, "properties") >= 0 &&
		declarationLayoutIndex(layout, "type") >= 0 &&
		declarationLayoutIndex(layout, "name") >= 0 &&
		declarationLayoutComponentsUnique(layout)
}

func declarationLayoutIndex(layout []string, component string) int {
	for index, candidate := range layout {
		if candidate == component {
			return index
		}
	}
	return -1
}

func declarationLayoutComponentsUnique(layout []string) bool {
	for index, component := range layout {
		for other := index + 1; other < len(layout); other++ {
			if layout[other] == component {
				return false
			}
		}
	}
	return true
}
