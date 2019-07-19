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
	"bufio"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"unicode"
	"unicode/utf8"
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
factor =	[unaryOp] (identifier | '(' predicate ')'| extension).
unaryOp = '¬'.
*/

const (
	OPar = "("
	CPar = ")"
)

func Parse(source io.Reader) (p *Predicate, e error) {
	ss := []scanner{
		identScan,
		spaceScan,
		strScan(NotOp),
		strScan(AndOp),
		strScan(OrOp),
		strScan(EquivalesOp),
		strScan(NotEquivalesOp),
		strScan(ImpliesOp),
		strScan(FollowsOp),
		strScan(OPar),
		strScan(CPar),
	}
	ts, e := tokens(source, ss)
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
			p, n, e = parseTerm(ts, i, nil)
		}
	}
	return
}

func parseTerm(ts []token, i int, ext parser) (p *Predicate, n int,
	e error) {
	ops := []string{ImpliesOp, FollowsOp}
	ib := func(k int) (b bool) {
		p, n, e = parseOp(ts, i, []string{ops[k]}, parseJunction)
		b = e == nil
		return
	}
	ok, _ := alg.BLnSrch(ib, len(ops))
	if !ok {
		p, n, e = parseJunction(ts, i, ext)
	}
	return
}

func parseJunction(ts []token, i int, ext parser) (p *Predicate, n int,
	e error) {
	ops := []string{OrOp, AndOp}
	ib := func(k int) (b bool) {
		p, n, e = parseOp(ts, i, []string{ops[k]}, parseFactor)
		b = e == nil
		return
	}
	ok, _ := alg.BLnSrch(ib, len(ops))
	if !ok {
		p, n, e = parseFactor(ts, i, ext)
	}
	return
}

type parser func([]token, int, parser) (*Predicate, int, error)

func parseOp(ts []token, i int, ops []string,
	s parser) (p *Predicate, n int, e error) {
	p, n, e = s(ts, i, nil)
	op := moreOps(ts, ops, n)
	if op == "" {
		e = notFound(ops[0])
	} else if n == len(ts)-1 {
		e = malformedOp(op)
	}
	curr := p
	for e == nil && op != "" {
		var b *Predicate
		b, n, e = s(ts, n+1, nil)
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

func parseFactor(ts []token, i int,
	ext parser) (p *Predicate, n int, e error) {
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
		} else if t.value == OPar {
			np, n, e = parse(ts, n+1)
			if n != len(ts) && ts[n].value == CPar {
				e, n = nil, n+1
			} else if e == nil {
				e = noMatchPar()
			}
		} else if ext != nil {
			p, n, e = ext(ts, i, nil)
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
	e = fmt.Errorf("Not recognized symbol \"%s\"", t)
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

func malformedOp(t string) (e error) {
	e = fmt.Errorf("Malformed operation '%s'", t)
	return
}

type token struct {
	value   string
	isIdent bool
}

func tokens(source io.Reader, ss []scanner) (ts []token,
	e error) {
	var rn rune
	var sc func(rune) (token, bool, bool)
	var t token
	rd := bufio.NewReader(source)
	n, prod, read, search, scan := 0, false, true, false, false
	for prod || scan || e == nil {
		if read {
			rn, _, e = rd.ReadRune()
			if e != nil {
				rn = ' '
			}
			read, search = false, !scan
		} else if search {
			if n == len(ss) {
				e = fmt.Errorf("Not recognized '%s'", string(rn))
			} else {
				sc, n, scan, search = ss[n](), n+1, true, false
			}
		} else if scan && !prod {
			t, scan, prod = sc(rn)
			search, read = !scan && !prod, scan
		} else if prod {
			if t.value != "" {
				ts = append(ts, t)
			}
			n, prod, search = 0, false, true
		}
	}
	if e == io.EOF {
		e = nil
	}
	return
}

type scanner func() func(rune) (token, bool, bool)

func identScan() func(rune) (token, bool, bool) {
	var ident string
	return func(rn rune) (t token, cont, prod bool) {
		cont = unicode.IsLetter(rn) ||
			(ident != "" && unicode.IsDigit(rn))
		if cont {
			ident = ident + string(rn)
		} else if ident != "" {
			t, prod = token{value: ident, isIdent: true}, true
		}
		return
	}
}

func strScan(strScan string) (s scanner) {
	s = func() func(rune) (token, bool, bool) {
		str := strScan
		return func(rn rune) (t token, cont, prod bool) {
			sr, size := utf8.DecodeRuneInString(str)
			cont = sr != utf8.RuneError && sr == rn
			if cont {
				str = str[size:]
			}
			prod = len(str) == 0
			if prod {
				t, cont = token{value: strScan}, true
			}
			return
		}
	}
	return
}

func spaceScan() func(rune) (token, bool, bool) {
	start := false
	return func(rn rune) (t token, cont, prod bool) {
		cont = unicode.IsSpace(rn)
		if cont {
			start = true
		}
		prod = start && !cont
		return
	}
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
