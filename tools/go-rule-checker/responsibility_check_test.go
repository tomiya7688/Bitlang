package gorulechecker

import "testing"

func TestFullRepositoryCheck(t *testing.T) {
	cases := []struct {
		paths []string
		want  bool
	}{
		{paths: []string{"."}, want: true},
		{paths: []string{"./"}, want: true},
		{paths: []string{"internal/bitlang"}, want: false},
		{paths: []string{".", "tools"}, want: false},
	}
	for _, item := range cases {
		if got := isFullRepositoryCheck(item.paths); got != item.want {
			t.Fatalf("isFullRepositoryCheck(%v) = %v, want %v", item.paths, got, item.want)
		}
	}
}

func TestNormalizeResponsibilityPath(t *testing.T) {
	if got := normalizeResponsibilityPath("./internal/bitlang/lexer.go"); got != "internal/bitlang/lexer.go" {
		t.Fatalf("normalizeResponsibilityPath() = %q", got)
	}
}
