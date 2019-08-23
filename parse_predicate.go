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
