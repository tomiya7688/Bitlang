package bitlangci

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tomiya7688/Bitlang/internal/bitlang"
)

func loadSpecificationSetFromRoot(root string) (bitlang.SpecificationSet, error) {
	propertiesPath := filepath.Join(root, "spec", "preprocessed", "properties.json")
	declarationsPath := filepath.Join(root, "spec", "preprocessed", "declarations.json")

	propertiesData, err := os.ReadFile(propertiesPath)
	if err != nil {
		return bitlang.SpecificationSet{}, fmt.Errorf("read properties specification: %w", err)
	}
	declarationsData, err := os.ReadFile(declarationsPath)
	if err != nil {
		return bitlang.SpecificationSet{}, fmt.Errorf("read declarations specification: %w", err)
	}
	return bitlang.LoadSpecificationSet(propertiesData, declarationsData)
}
