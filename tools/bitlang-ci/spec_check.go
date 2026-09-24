package bitlangci

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tomiya7688/Bitlang/internal/bitlang"
)

func runSpecificationCheck(root string, out io.Writer, errOut io.Writer) bool {
	fmt.Fprintln(out, "==> machine-readable specifications")

	properties, err := loadPropertySpecificationFile(root)
	if err != nil {
		fmt.Fprintf(errOut, "FAIL machine-readable specifications: %v\n", err)
		return false
	}
	declarations, err := loadDeclarationSpecificationFile(root)
	if err != nil {
		fmt.Fprintf(errOut, "FAIL machine-readable specifications: %v\n", err)
		return false
	}
	if err := bitlang.ValidateSpecificationConsistency(properties, declarations); err != nil {
		fmt.Fprintf(errOut, "FAIL machine-readable specifications: %v\n", err)
		return false
	}

	fmt.Fprintln(out, "OK machine-readable specifications")
	return true
}

func loadPropertySpecificationFile(root string) (bitlang.PropertySpecification, error) {
	path := filepath.Join(root, "spec", "preprocessed", "properties.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return bitlang.PropertySpecification{}, fmt.Errorf("read properties specification: %w", err)
	}
	return bitlang.LoadPropertySpecification(data)
}

func loadDeclarationSpecificationFile(root string) (bitlang.DeclarationSpecification, error) {
	path := filepath.Join(root, "spec", "preprocessed", "declarations.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return bitlang.DeclarationSpecification{}, fmt.Errorf("read declarations specification: %w", err)
	}
	return bitlang.LoadDeclarationSpecification(data)
}
