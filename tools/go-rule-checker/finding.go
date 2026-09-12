package gorulechecker

import "fmt"

type finding struct {
	level   string
	path    string
	line    int
	message string
}

func (f finding) String() string {
	if f.line > 0 {
		return fmt.Sprintf("%s %s:%d: %s", f.level, f.path, f.line, f.message)
	}
	return fmt.Sprintf("%s %s: %s", f.level, f.path, f.message)
}
