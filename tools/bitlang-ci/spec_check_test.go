package bitlangci

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSpecificationCheckFailsClosedForIncompleteFamilyPipeline(t *testing.T) {
	root := writeSpecificationFixture(t, "variable")
	var out bytes.Buffer
	var errOut bytes.Buffer
	if runSpecificationCheck(root, &out, &errOut) {
		t.Fatal("expected an incomplete family pipeline to fail the specification gate")
	}
	if !strings.Contains(errOut.String(), "cross-repository acceptance fixtures are not implemented") {
		t.Fatalf("expected the family fixture readiness error, got: %s", errOut.String())
	}
}

func TestRunSpecificationCheckRejectsUnknownPropertyTarget(t *testing.T) {
	root := writeSpecificationFixture(t, "unknown")
	var out bytes.Buffer
	var errOut bytes.Buffer
	if runSpecificationCheck(root, &out, &errOut) {
		t.Fatal("expected specification check failure")
	}
}

func writeSpecificationFixture(t *testing.T, propertyTarget string) string {
	t.Helper()
	root := t.TempDir()
	specDir := filepath.Join(root, "spec", "preprocessed")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	properties := []byte(`{"version":1,"axes":[{"name":"nullability","states":["nullable","unnullable"],"exclusive":true,"required":true,"applies_to":["variable"]}]}`)
	declarations := []byte(`{"version":1,"kinds":[{"name":"variable","property_target":"` + propertyTarget + `","terminator":";","layout":["properties","type","name"]}]}`)
	if err := os.WriteFile(filepath.Join(specDir, "properties.json"), properties, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "declarations.json"), declarations, 0o600); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "spec", "family-pipeline.json")
	manifest, err := os.ReadFile(filepath.Join("..", "..", "spec", "family-pipeline.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}
