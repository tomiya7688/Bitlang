package gorulechecker

import "fmt"

type finding struct {
	level   string
	path    string
	line    int
	message string
}

func (f finding) String() string {
	level := "W"
	if f.level == "ERROR" {
		level = "E"
	}
	if f.line > 0 {
		return fmt.Sprintf("%s %s:%d %s", level, f.path, f.line, f.message)
	}
	return fmt.Sprintf("%s %s %s", level, f.path, f.message)
}
