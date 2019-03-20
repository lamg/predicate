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
	"io"
	"io/ioutil"
	"unicode"
)

/*
Grammar in EBNF syntax

predicate = predicate binaryOp predicate | unaryOp predicate |
	'(' predicate ')'| identifier.
binaryOp = '∧'|'∨'|'⇒'|'⇐'|'≡'|'≢'.
unaryOp = '¬'.

simplifying first production:
predicate =
	(identifier | unaryOp (identifier|'(' predicate ')') |
	'(' predicate ')') [binaryOp predicate].
*/

func Parse(source io.Reader) (p *Predicate, e error) {
	ts, e := tokens(source)
	if e == nil {
		p, _, e = parse(ts, 0)
	}
	return
}

func parse(ts []token, i int) (p *Predicate, n int, e error) {
	p, n = new(Predicate), i
	if n == len(ts) {
		e = unexpEnd()
	} else {
		t := ts[n]
		if t.isIdent {
			p, n = &Predicate{Operator: Term, String: t.value}, n+1
		} else if isUnary(t.value) {
			p = &Predicate{
				Operator: NotOp,
			}
			if n+1 != len(ts) && ts[n+1].isIdent {
				p.A, n = &Predicate{
					Operator: Term,
					String:   ts[n+1].value,
				},
					n+2
			} else {
				p.A, n, e = parse(ts, n+1)
			}
		} else if t.value == string(opar) {
			p, n, e = parse(ts, n+1)
			if n != len(ts) && ts[n].value == string(cpar) {
				e, n = nil, n+1
			} else if e == nil {
				e = noMatchPar()
			}
		} else {
			e = notRec(t.value)
		}
	}
	if e == nil && n != len(ts) {
		if isBinary(ts[n].value) {
			p = &Predicate{
				Operator: ts[n].value,
				A:        p,
			}
			p.B, n, e = parse(ts, n+1)
		} else {
			e = notRec(ts[n].value)
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
				ok, _ := bLnSrch(ib, len(rns))
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
	ok, _ = bLnSrch(ib, len(bs))
	return
}
