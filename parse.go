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

predicate = term {('≡'|'≢') term}.
term = junction ({'⇒' junction} | {'⇐' junction}).
junction = factor ({'∨' factor} | {'∧' factor}).
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
*/

func Parse(rd io.Reader) (p *Predicate, e error) {
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

	st := &predState{
		tkf: tokens(rd, ss),
	}

	p, e = st.predicate()
	if e == nil && st.token.value != string(eof) {
		e = errUnkChar(st.token.value)
	}
	return
}

type predState struct {
	tkf   func() (*token, error)
	token *token
}

func (s *predState) next() (e error) {
	s.token, e = s.tkf()
	if e == nil && s.token.value == "" {
		s.token, e = s.tkf()
	}
	return
}

func (s *predState) predicate() (p *Predicate, e error) {
	factor := s.factor()
	junction := s.parseOp("junction", factor, false, OrOp, AndOp)
	term := s.parseOp("term", junction, false, ImpliesOp, FollowsOp)
	p, e = s.parseOp("predicate", term, true, EquivalesOp, NotEquivalesOp)()
	return
}

func (s *predState) parseOp(
	name string,
	sym func() (*Predicate, error),
	mixAlt bool, ops ...string) func() (*Predicate, error) {
	return func() (p *Predicate, e error) {
		p, e = sym()
		if e == nil {
			var o string
			o, e = s.moreOps(ops)
			if e == nil && o != "" && !mixAlt {
				// restrict the set of operators to the detected
				ops = []string{o}
			}
			curr := p
			for e == nil && o != "" {
				var b *Predicate
				b, e = sym()
				if e == nil {
					old := &Predicate{
						Operator: curr.Operator,
						A:        curr.A,
						B:        curr.B,
						String:   curr.String,
					}
					curr.Operator, curr.A, curr.B = o, old, b
					curr = curr.B
					o, e = s.moreOps(ops)
				}
			}
		}
		return
	}
}

func (s *predState) moreOps(ops []string) (
	op string, e error,
) {
	if s.token.value != CPar {
		ib := func(i int) bool { return ops[i] == s.token.value }
		ok, n := alg.BLnSrch(ib, len(ops))
		if ok {
			op = ops[n]
		}
	}
	return
}

func errUnkChar(value string) error {
	return fmt.Errorf("Unknown char %s", value)
}

func (s *predState) factor() func() (*Predicate, error) {
	return func() (p *Predicate, e error) {
		e = s.next()
		var nt *Predicate
		if e == nil {
			if s.token.value == NotOp {
				nt = &Predicate{Operator: NotOp}
				e = s.next()
			}
			if e == nil {
				if s.token.isIdent {
					p = &Predicate{Operator: Term, String: s.token.value}
				} else if s.token.value == OPar {
					p, e = s.predicate()
					if e == nil && s.token.value != CPar {
						e = errClosingPar()
					}
				} else {
					e = errIdentOrOpening()
				}
			}
		}
		if e == nil && nt != nil {
			nt.B = p
			p = nt
		}
		if e == nil {
			e = s.next()
		}
		return
	}
}

func errClosingPar() error {
	return fmt.Errorf("Expecting closing parenthesis")
}

func errIdentOrOpening() error {
	return fmt.Errorf("Expecting identifier or opening parenthesis")
}

const (
	OPar = "("
	CPar = ")"
)

type token struct {
	value    string
	isIdent  bool
	isNumber bool
}

const (
	// 0x3 is the end of file character
	eof = 0x3
)

func tokens(source io.Reader, ss []scanner) (
	tf func() (*token, error)) {
	rd := bufio.NewReader(source)
	ss = append(ss, eofScan)
	var rn rune
	var sc func(rune) (*token, bool, bool)
	n, end, read, search, scan := 0, false, true, false, false
	tf = func() (t *token, e error) {
		if end {
			t = &token{value: string(eof)}
		}
		for !end {
			if read {
				rn, _, e = rd.ReadRune()
				if e != nil {
					rn = eof
					if e == io.EOF {
						e = nil
					}
				}
				read, search = false, !scan
			} else if search {
				if n == len(ss) {
					e, end = notRec(string(rn)), true
				} else {
					sc, n, search = ss[n](), n+1, false
				}
			} else if !search {
				t, read, end = sc(rn)
				search = !read || end
				scan = !search
			}
		}
		n, end = 0, e != nil || t.value == string(eof)
		return
	}
	return
}

func notRec(s string) (e error) {
	return fmt.Errorf("Not recognized '%s'", s)
}

type scanner func() func(rune) (*token, bool, bool)

func identScan() func(rune) (*token, bool, bool) {
	var ident string
	return func(rn rune) (t *token, cont, prod bool) {
		cont = unicode.IsLetter(rn) ||
			(ident != "" && unicode.IsDigit(rn))
		if cont {
			ident = ident + string(rn)
		} else if ident != "" {
			t, prod = &token{value: ident, isIdent: true}, true
		}
		return
	}
}

func strScan(strScan string) (s scanner) {
	s = func() func(rune) (*token, bool, bool) {
		str := strScan
		return func(rn rune) (t *token, cont, prod bool) {
			sr, size := utf8.DecodeRuneInString(str)
			cont = sr != utf8.RuneError && sr == rn
			if cont {
				str = str[size:]
			}
			prod = len(str) == 0
			if prod {
				t, cont = &token{value: strScan}, true
			}
			return
		}
	}
	return
}

func spaceScan() func(rune) (*token, bool, bool) {
	start := false
	return func(rn rune) (t *token, cont, prod bool) {
		cont = unicode.IsSpace(rn)
		if cont {
			start = true
		}
		prod = start && !cont
		if prod {
			t, start = new(token), false
		}
		return
	}
}

func eofScan() func(rune) (*token, bool, bool) {
	return func(r rune) (t *token, cont, prod bool) {
		if r == eof {
			t, prod = &token{value: string(r)}, true
		}
		return
	}
}
