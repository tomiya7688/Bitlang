package gorulechecker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectGoFilesRejectsSymlinkedGoFile(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	if err := os.WriteFile(target, []byte("not Go source"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	link := filepath.Join(root, "escape.go")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	if _, err := collectGoFiles([]string{root}); err == nil {
		t.Fatal("collectGoFiles accepted a symlinked Go file")
	}
}

func TestCollectGoFilesRejectsSymlinkRoot(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.go")
	if err := os.WriteFile(target, []byte("package test\n"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	link := filepath.Join(root, "link.go")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	if _, err := collectGoFiles([]string{link}); err == nil {
		t.Fatal("collectGoFiles accepted a symlink root")
	}
}
