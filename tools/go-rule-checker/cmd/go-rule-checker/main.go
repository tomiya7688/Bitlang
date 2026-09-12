package main

import (
	"fmt"
	"os"

	gorulechecker "github.com/tomiya7688/Bitlang/tools/go-rule-checker"
)

// main starts the Go implementation rule checker and maps its result to an exit status.
func main() {
	code, err := gorulechecker.Run(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(code)
}
