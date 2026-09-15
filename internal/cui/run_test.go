package cui

import (
	"bytes"
	"testing"
)

func TestRunCanonicalize(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run([]string{"canonicalize", "PlayerHP"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%q", code, errOut.String())
	}
	if out.String() != "PlayerHP\tplayerhp\n" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestRunRejectsMissingCommand(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run(nil, &out, &errOut)
	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if errOut.Len() == 0 {
		t.Fatal("expected usage on stderr")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := Run([]string{"unknown"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
}
