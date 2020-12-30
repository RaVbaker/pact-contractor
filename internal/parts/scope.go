package parts

import (
	"strconv"
)

type Scope struct {
	num int
	total int
	merged bool
}

func (s Scope) Merged() bool {
	return s.merged
}

func NewScope(num, total int) Scope {
	return Scope{num: num, total: total}
}

func (s Scope) Defined() bool {
	return s.total != 0
}

func (s *Scope) MarkAsMerged() {
	s.merged = true
}

func (s Scope) Total() int {
	return s.total
}

func (s Scope) Current() int {
	return s.num
}

func (s Scope) Name() string {
	return strconv.Itoa(s.num)
}