package docpairchecker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateBaseCommitRejectsOptionLikeInput(t *testing.T) {
	if err := validateBaseCommit("--output=/tmp/escape"); err == nil {
		t.Fatal("validateBaseCommit accepted option-like input")
	}
}

func TestValidateBaseCommitAcceptsFullHexObjectID(t *testing.T) {
	if err := validateBaseCommit("0123456789abcdef0123456789abcdef01234567"); err != nil {
		t.Fatalf("validateBaseCommit rejected valid object ID: %v", err)
	}
}

func TestReadRegisteredDocumentationRejectsSymlink(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.md")
	if err := os.WriteFile(target, []byte("target"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	link := filepath.Join(root, "README.md")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	if _, err := readRegisteredDocumentation(root, "README.md"); err == nil {
		t.Fatal("readRegisteredDocumentation accepted a symlink")
	}
}

func TestReadRegisteredDocumentationRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	if _, err := readRegisteredDocumentation(root, "../outside.md"); err == nil {
		t.Fatal("readRegisteredDocumentation accepted traversal")
	}
}
