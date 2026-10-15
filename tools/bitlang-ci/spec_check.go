package bitlangci

import (
	"fmt"
	"io"
)

func runSpecificationCheck(root string, out io.Writer, errOut io.Writer) bool {
	fmt.Fprintln(out, "==> machine-readable specifications")

	if _, err := loadSpecificationSetFromRoot(root); err != nil {
		fmt.Fprintf(errOut, "FAIL machine-readable specifications: %v\n", err)
		return false
	}
	if err := loadFamilyPipelineManifest(root); err != nil {
		fmt.Fprintf(errOut, "FAIL family pipeline contract: %v\n", err)
		return false
	}

	fmt.Fprintln(out, "OK family pipeline contract")
	fmt.Fprintln(out, "OK machine-readable specifications")
	return true
}
