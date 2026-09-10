package main

import (
	"fmt"
	"os"

	core "github.com/tomiya7688/Bitlang/internal/bitlang"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  bitlang canonicalize <identifier> [identifier ...]")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "canonicalize":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		for _, identifier := range os.Args[2:] {
			canonical, err := core.CanonicalizeIdentifier(identifier)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("%s\t%s\n", identifier, canonical)
		}
	default:
		usage()
		os.Exit(2)
	}
}
