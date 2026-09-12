package gorulechecker

import "fmt"

type finding struct {
	level   string
	rule    string
	path    string
	line    int
	message string
}

func (f finding) String() string {
	if f.line > 0 {
		return fmt.Sprintf("%s %s:%d %s %s", f.level, f.path, f.line, f.rule, f.message)
	}
	return fmt.Sprintf("%s %s %s %s", f.level, f.path, f.rule, f.message)
}
