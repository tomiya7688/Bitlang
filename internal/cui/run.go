package cui

import (
	"fmt"
	"io"

	core "github.com/tomiya7688/Bitlang/internal/bitlang"
)

// Run handles one CUI invocation and returns the process exit status.
func Run(args []string, out io.Writer, errOut io.Writer) int {
	if len(args) < 1 {
		writeUsage(errOut)
		return 2
	}

	switch args[0] {
	case "canonicalize":
		if len(args) < 2 {
			writeUsage(errOut)
			return 2
		}
		for _, identifier := range args[1:] {
			canonical, err := core.CanonicalizeIdentifier(identifier)
			if err != nil {
				fmt.Fprintln(errOut, err)
				return 1
			}
			fmt.Fprintf(out, "%s\t%s\n", identifier, canonical)
		}
		return 0
	default:
		writeUsage(errOut)
		return 2
	}
}

func writeUsage(out io.Writer) {
	fmt.Fprintln(out, "usage:")
	fmt.Fprintln(out, "  bitlang canonicalize <identifier> [identifier ...]")
}
