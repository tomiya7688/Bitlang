package bitlangci

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type commandCheck struct {
	name string
	cmd  string
	args []string
}

// Run executes the strict local CI gate used by developers and GitHub Actions.
func Run(root string, out io.Writer, errOut io.Writer) int {
	failures := 0

	if !runFormatCheck(root, out, errOut) {
		failures++
	}

	checks := []commandCheck{
		{name: "go-rule-checker", cmd: "go", args: []string{"run", "./tools/go-rule-checker/cmd/go-rule-checker", "."}},
		{name: "documentation pairs", cmd: "go", args: []string{"run", "./tools/doc-pair-checker/cmd/doc-pair-checker"}},
		{name: "go mod tidy", cmd: "go", args: []string{"mod", "tidy", "-diff"}},
		{name: "go vet", cmd: "go", args: []string{"vet", "./..."}},
		{name: "go test", cmd: "go", args: []string{"test", "./..."}},
		{name: "go test shuffled/repeated", cmd: "go", args: []string{"test", "-shuffle=on", "-count=3", "./..."}},
		{name: "go build packages", cmd: "go", args: []string{"build", "./..."}},
		{name: "git diff check", cmd: "git", args: []string{"diff", "--check"}},
	}

	for _, check := range checks {
		if !runCommand(root, check, out, errOut) {
			failures++
		}
	}

	if !runBuildCheck(root, out, errOut) {
		failures++
	}

	if failures > 0 {
		fmt.Fprintf(errOut, "NG bitlang-ci: %d check(s) failed\n", failures)
		return 1
	}

	fmt.Fprintln(out, "OK bitlang-ci")
	return 0
}

func runFormatCheck(root string, out io.Writer, errOut io.Writer) bool {
	fmt.Fprintln(out, "==> gofmt")
	command := exec.Command("gofmt", "-l", ".")
	command.Dir = root
	var formatted bytes.Buffer
	command.Stdout = &formatted
	command.Stderr = errOut
	if err := command.Run(); err != nil {
		fmt.Fprintf(errOut, "FAIL gofmt: %v\n", err)
		return false
	}
	if strings.TrimSpace(formatted.String()) != "" {
		fmt.Fprintln(errOut, "FAIL gofmt: files require formatting")
		fmt.Fprint(errOut, formatted.String())
		return false
	}
	fmt.Fprintln(out, "OK gofmt")
	return true
}

func runCommand(root string, check commandCheck, out io.Writer, errOut io.Writer) bool {
	fmt.Fprintf(out, "==> %s\n", check.name)
	command := exec.Command(check.cmd, check.args...)
	command.Dir = root
	command.Stdout = out
	command.Stderr = errOut
	if err := command.Run(); err != nil {
		fmt.Fprintf(errOut, "FAIL %s: %v\n", check.name, err)
		return false
	}
	fmt.Fprintf(out, "OK %s\n", check.name)
	return true
}

func runBuildCheck(root string, out io.Writer, errOut io.Writer) bool {
	tempDir, err := os.MkdirTemp("", "bitlang-ci-build-*")
	if err != nil {
		fmt.Fprintf(errOut, "FAIL go build: create temp dir: %v\n", err)
		return false
	}
	defer os.RemoveAll(tempDir)

	output := filepath.Join(tempDir, "bitlang")
	if runtime.GOOS == "windows" {
		output += ".exe"
	}
	return runCommand(root, commandCheck{
		name: "go build",
		cmd:  "go",
		args: []string{"build", "-o", output, "./cmd/bitlang"},
	}, out, errOut)
}
