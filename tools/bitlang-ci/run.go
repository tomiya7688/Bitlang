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
	name    string
	command *exec.Cmd
}

// Run executes the strict local CI gate used by developers and GitHub Actions.
func Run(root string, out io.Writer, errOut io.Writer) int {
	failures := 0

	if !runFormatCheck(root, out, errOut) {
		failures++
	}
	if !runSpecificationCheck(root, out, errOut) {
		failures++
	}

	checks := []commandCheck{
		{name: "go-rule-checker", command: exec.Command("go", "run", "./tools/go-rule-checker/cmd/go-rule-checker", ".")},
		{name: "documentation pairs", command: exec.Command("go", "run", "./tools/doc-pair-checker/cmd/doc-pair-checker")},
		{name: "go mod tidy", command: exec.Command("go", "mod", "tidy", "-diff")},
		{name: "go vet", command: exec.Command("go", "vet", "./...")},
		{name: "go test", command: exec.Command("go", "test", "./...")},
		{name: "go test shuffled/repeated", command: exec.Command("go", "test", "-shuffle=on", "-count=3", "./...")},
		{name: "go build packages", command: exec.Command("go", "build", "./...")},
		{name: "git diff check", command: exec.Command("git", "diff", "--check")},
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
	check.command.Dir = root
	check.command.Stdout = out
	check.command.Stderr = errOut
	if err := check.command.Run(); err != nil {
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
	// #nosec G204 -- output is generated inside a process-owned temporary directory; executable and other arguments are constants.
	command := exec.Command("go", "build", "-o", output, "./cmd/bitlang")
	return runCommand(root, commandCheck{name: "go build", command: command}, out, errOut)
}
