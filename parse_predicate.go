package predicate

import (
	"fmt"
	"io"
)

func Parse(rd io.Reader) (p *Predicate, e error) {
	e = fmt.Errorf("Not implemented")
	return
}

func parse(tf func() (*token, error)) (p *Predicate, e error) {
	return
}

func predicateSym(curr func() *Predicate) (r func() *Predicate,
	branch func(string)) {
	// curr() is defined
	p := curr()
	c := p
	branch = func(s string) {
		n := curr()
		old := &Predicate{
			Operator: c.Operator,
			A:        c.A,
			B:        c.B,
			String:   c.String,
		}
		c.Operator, c.A, c.B, c.String = s, old, n, ""
		c = n
	}
	r = func() *Predicate { return p }
	return
}
