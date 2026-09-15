package main

import (
	"flag"
	"os"

	docpairchecker "github.com/tomiya7688/Bitlang/tools/doc-pair-checker"
)

// main validates bilingual documentation pairs and optionally checks PR changes against a base ref.
func main() {
	baseRef := flag.String("base", "", "git base ref used to require paired document updates")
	flag.Parse()
	os.Exit(docpairchecker.Run(".", *baseRef, os.Stdout, os.Stderr))
}
