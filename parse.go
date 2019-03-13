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
	"fmt"
	"io"
	"text/scanner"
	"unicode"
)

/*
Grammar in EBNF syntax

predicate = predicate binaryOp predicate | unaryOp predicate |
	'(' predicate ')'| identifier.
binaryOp = '∧'|'∨'|'⇒'|'≡'|'≢'.
unaryOp = '¬'.

*/

func Parse(source io.Reader,
	fv func(string) (func() (bool, bool), bool)) (p *Predicate,
	e error) {
	notImplemented()
	return
}

const (
	not  = '¬'
	and  = '∧'
	or   = '∨'
	eq   = '≡'
	neq  = '≢'
	imp  = '⇒'
	opar = '('
	cpar = ')'
)

func tokens(source io.Reader) (ts []string, e error) {
	s := new(scanner.Scanner)
	s.Mode = scanner.ScanIdents
	s.IsIdentRune = func(ch rune, i int) (b bool) {
		chs := []rune{not, and, or, eq, neq, imp, opar, cpar}
		ib := func(i int) bool { return chs[i] == ch }
		ok, _ := bLnSrch(ib, len(chs))
		b = (i == 0 && ok) || unicode.IsLetter(ch) ||
			(i > 0 && unicode.IsDigit(ch))
		return
	}
	s.Error = func(n *scanner.Scanner, msg string) {
		if n.Position.IsValid() {
			e = fmt.Errorf("%d:%d %s", n.Position.Line,
				n.Position.Column, msg)
		} else {
			e = fmt.Errorf(msg)
		}
	}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		ts = append(ts, s.TokenText())
	}
	return
}
