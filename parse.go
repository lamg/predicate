// Copyright © 2019 Luis Ángel Méndez Gort

// This file is part of Predicate.

// Predicate is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General
// Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your
// option) any later version.

// Predicate is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the
// implied warranty of MERCHANTABILITY or FITNESS FOR A
// PARTICULAR PURPOSE. See the GNU Lesser General Public
// License for more details.

// You should have received a copy of the GNU Lesser General
// Public License along with Predicate.  If not, see
// <https://www.gnu.org/licenses/>.

package predicate

import (
	"bytes"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"io/ioutil"
	"unicode"
)

/*
Grammar in EBNF syntax

predicate = term ('≡'|'≢') term {('≡'|'≢') term}| term.
term = implication | consequence | junction.
implication = junction '⇒' junction {'⇒' junction}.
consequence = junction '⇐' junction {'⇐' junction}.
junction = disjunction | conjunction | factor.
disjunction = factor '∨' factor {'∨' factor}.
conjunction = factor '∧' factor {'∧' factor}.
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
*/

func Parse(source io.Reader) (p *Predicate, e error) {
	ts, e := tokens(source)
	if e == nil {
		var n int
		p, n, e = parse(ts, 0)
		if e == nil && n != len(ts) {
			e = notRec(ts[n].value)
		}
	}
	return
}

func parse(ts []token, i int) (p *Predicate, n int, e error) {
	if i == len(ts) {
		e = unexpEnd()
	} else {
		ops := []string{EquivalesOp, NotEquivalesOp}
		p, n, e = parseOp(ts, i, ops, parseTerm)
		if e != nil && e.Error() == notFound(EquivalesOp).Error() {
			p, n, e = parseTerm(ts, i)
		}
	}
	return
}

func parseTerm(ts []token, i int) (p *Predicate, n int,
	e error) {
	ops := []string{ImpliesOp, FollowsOp}
	ib := func(k int) (b bool) {
		p, n, e = parseOp(ts, i, []string{ops[k]}, parseJunction)
		b = e == nil
		return
	}
	ok, _ := alg.BLnSrch(ib, len(ops))
	if !ok {
		p, n, e = parseJunction(ts, i)
	}
	return
}

func parseJunction(ts []token, i int) (p *Predicate, n int,
	e error) {
	ops := []string{OrOp, AndOp}
	ib := func(k int) (b bool) {
		p, n, e = parseOp(ts, i, []string{ops[k]}, parseFactor)
		b = e == nil
		return
	}
	ok, _ := alg.BLnSrch(ib, len(ops))
	if !ok {
		p, n, e = parseFactor(ts, i)
	}
	return
}

type parser func([]token, int) (*Predicate, int, error)

func parseOp(ts []token, i int, ops []string,
	s parser) (p *Predicate, n int, e error) {
	p, n, e = s(ts, i)
	op := moreOps(ts, ops, n)
	if op == "" {
		e = notFound(ops[0])
	}
	curr := p
	for e == nil && op != "" {
		var b *Predicate
		b, n, e = s(ts, n+1)
		if e == nil {
			old := &Predicate{
				Operator: curr.Operator,
				A:        curr.A,
				B:        curr.B,
				String:   curr.String,
			}
			curr.Operator = op
			curr.A = old
			curr.B = b
			curr = curr.B
		}
		op = moreOps(ts, ops, n)
	}
	return
}

func moreOps(ts []token, ops []string, n int) (op string) {
	if n != len(ts) {
		op = ts[n].value
		ib := func(i int) bool { return ops[i] == op }
		ok, _ := alg.BLnSrch(ib, len(ops))
		if !ok {
			op = ""
		}
	}
	return
}

func notFound(op string) (e error) {
	e = fmt.Errorf("Not found %s", op)
	return
}

func parseFactor(ts []token, i int) (p *Predicate, n int, e error) {
	n = i
	t := ts[n]
	var neg *Predicate
	if isUnary(t.value) {
		neg, n = &Predicate{Operator: t.value}, n+1
	}
	if n != len(ts) {
		t = ts[n]
	} else {
		e = unexpEnd()
	}
	var np *Predicate
	if e == nil {
		if t.isIdent {
			np, n = &Predicate{Operator: Term, String: t.value}, n+1
		} else if t.value == string(opar) {
			np, n, e = parse(ts, n+1)
			if n != len(ts) && ts[n].value == string(cpar) {
				e, n = nil, n+1
			} else if e == nil {
				e = noMatchPar()
			}
		} else {
			e = notRec(t.value)
		}
	}
	if e == nil {
		if neg != nil {
			neg.A = np
			p = neg
		} else {
			p = np
		}
	}
	return
}

func notRec(t string) (e error) {
	e = fmt.Errorf("Not recognized token \"%s\"", t)
	return
}

func unexpEnd() (e error) {
	e = fmt.Errorf("Unexpected end of input")
	return
}

func noMatchPar() (e error) {
	e = fmt.Errorf("No matching )")
	return
}

const (
	not     = '¬'
	and     = '∧'
	or      = '∨'
	eq      = '≡'
	neq     = '≢'
	implies = '⇒'
	follows = '⇐'
	opar    = '('
	cpar    = ')'
)

type token struct {
	value   string
	isIdent bool
}

func tokens(source io.Reader) (ts []token, e error) {
	bs, e := ioutil.ReadAll(source)
	if e == nil {
		rd := bytes.NewReader(bs)
		var ident string
		for e == nil {
			var rn rune
			rn, _, e = rd.ReadRune()
			if e == nil {
				rns := []rune{not, and, or, eq, neq, implies, follows,
					opar, cpar}
				ib := func(i int) (b bool) {
					b = rns[i] == rn
					return
				}
				ok, _ := alg.BLnSrch(ib, len(rns))
				if ok {
					if ident != "" {
						ts, ident = append(ts, token{ident, true}), ""
					}
					ts = append(ts, token{string(rn), false})
				} else if unicode.IsLetter(rn) ||
					(unicode.IsDigit(rn) && len(ident) != 0) {
					ident = ident + string(rn)
				} else if unicode.IsSpace(rn) {
					if ident != "" {
						ts, ident = append(ts, token{ident, true}), ""
					}
				} else {
					e = notRec(string(rn))
				}
			}
		}
		if e == io.EOF {
			if ident != "" {
				ts = append(ts, token{ident, true})
			}
			e = nil
		}
	}
	return
}

func isUnary(t string) (ok bool) {
	return t == NotOp
}

func isBinary(t string) (ok bool) {
	bs := []string{AndOp, OrOp, EquivalesOp, NotEquivalesOp,
		ImpliesOp, FollowsOp}
	ib := func(i int) (b bool) {
		b = string(bs[i]) == t
		return
	}
	ok, _ = alg.BLnSrch(ib, len(bs))
	return
}
