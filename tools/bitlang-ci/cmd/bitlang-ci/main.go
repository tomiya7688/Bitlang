package main

import (
	"os"

	bitlangci "github.com/tomiya7688/Bitlang/tools/bitlang-ci"
)

// main runs the shared Bitlang CI gate and maps its result to the process exit status.
func main() {
	os.Exit(bitlangci.Run(".", os.Stdout, os.Stderr))
}
