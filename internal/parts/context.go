package parts

import (
	"fmt"
)

type Context struct {
	num    int
	total  int
	merged bool
}

// NewScope creates a fresh context from two numbers
func NewScope(num, total int) Context {
	return Context{num: num, total: total}
}

func (s Context) Current() int {
	return s.num
}

func (s Context) Total() int {
	return s.total
}

func (s Context) Merged() bool {
	return s.merged
}

// Name returns string that contains part number and total
func (s Context) Name() string {
	return fmt.Sprintf("%d-%d", s.num, s.total)
}

// Defined is true when total is bigger than 0
func (s Context) Defined() bool {
	return s.total != 0
}

// MarkAsMerged records the fact that part was merged
func (s *Context) MarkAsMerged() {
	s.merged = true
}
