package route

import (
	"simple-open-library/helper"
	"strings"
)

var exists = struct{}{}

type Prefixes struct {
	Set map[string]struct{}
}

func NewPrefixes() *Prefixes {
	p := &Prefixes{}
	p.Set = make(map[string]struct{})
	return p
}

func (p *Prefixes) Add(value string) {
	p.Set[value] = exists
}

func (p *Prefixes) Remove(value string) {
	delete(p.Set, value)
}

func (p *Prefixes) ContainsPrefix(prefixes string) bool {
	for _, prefix := range helper.SplitSubpaths(prefixes) {
		_, hasValue := p.Set[strings.TrimPrefix(prefix,"//:")]
		if (hasValue) {
			return hasValue
		}
	}
	return false
}