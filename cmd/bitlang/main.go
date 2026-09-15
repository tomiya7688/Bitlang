package main

import (
	"os"

	"github.com/tomiya7688/Bitlang/internal/cui"
)

func main() {
	os.Exit(cui.Run(os.Args[1:], os.Stdout, os.Stderr))
}
