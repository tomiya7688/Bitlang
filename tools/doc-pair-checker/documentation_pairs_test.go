package docpairchecker

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunAcceptsCompleteDocumentationPairs(t *testing.T) {
	root := t.TempDir()
	writeDocumentationPairFixtures(t, root)

	var output bytes.Buffer
	var errorOutput bytes.Buffer
	if exitCode := Run(root, "", &output, &errorOutput); exitCode != 0 {
		t.Fatalf("Run returned %d: %s", exitCode, errorOutput.String())
	}
}

func TestRunRejectsMissingJapaneseDocument(t *testing.T) {
	root := t.TempDir()
	writeDocumentationPairFixtures(t, root)

	missingPath := filepath.Join(root, filepath.FromSlash(requiredDocumentationPairs[0].japanesePath))
	if err := os.Remove(missingPath); err != nil {
		t.Fatalf("remove fixture: %v", err)
	}

	var output bytes.Buffer
	var errorOutput bytes.Buffer
	if exitCode := Run(root, "", &output, &errorOutput); exitCode == 0 {
		t.Fatalf("Run unexpectedly succeeded")
	}
}

func writeDocumentationPairFixtures(t *testing.T, root string) {
	t.Helper()
	for _, pair := range requiredDocumentationPairs {
		englishPath := filepath.Join(root, filepath.FromSlash(pair.englishPath))
		japanesePath := filepath.Join(root, filepath.FromSlash(pair.japanesePath))
		if err := os.MkdirAll(filepath.Dir(englishPath), 0o755); err != nil {
			t.Fatalf("create English fixture directory: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(japanesePath), 0o755); err != nil {
			t.Fatalf("create Japanese fixture directory: %v", err)
		}
		if err := os.WriteFile(englishPath, []byte("[Japanese]("+path.Base(pair.japanesePath)+")\n"), 0o644); err != nil {
			t.Fatalf("write English fixture: %v", err)
		}
		if err := os.WriteFile(japanesePath, []byte("[English]("+path.Base(pair.englishPath)+")\n"), 0o644); err != nil {
			t.Fatalf("write Japanese fixture: %v", err)
		}
	}
}
