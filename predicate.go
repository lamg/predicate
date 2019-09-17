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
	alg "github.com/lamg/algorithms"
)

type Predicate struct {
	Operator string     `json:"operator"`
	A        *Predicate `json:"a"`
	B        *Predicate `json:"b"`
	String   string     `json:"string"`
	AltRef   int        `json:"-"`
}

const (
	// the associated to each constant indicate how to type the
	// character using Vim with the digraph feature
	NotOp          = "¬" // C-k NO
	AndOp          = "∧" // C-k AN
	OrOp           = "∨" // C-k OR
	EquivalesOp    = "≡" // C-k 3=
	NotEquivalesOp = "≢" // C-k ne (custom def. `:digraph ne 8802`)
	ImpliesOp      = "⇒" // C-k =>
	FollowsOp      = "⇐" // C-k <=
	Term           = "term"
)

type NameBool func(string) (bool, bool)

func Reduce(p *Predicate, interp NameBool) (r *Predicate) {
	r = new(Predicate)
	fps := []func(*Predicate, *Predicate, NameBool) bool{
		reduceNot,
		reduceAnd,
		reduceOr,
		reduceImplies,
		reduceFollows,
		reduceEquivales,
		reduceNotEquivales,
		reduceTerm,
	}
	ops := []string{
		NotOp,
		AndOp,
		OrOp,
		ImpliesOp,
		FollowsOp,
		EquivalesOp,
		NotEquivalesOp,
		Term,
	}
	fs := make([]alg.KFunc, len(fps))
	inf := func(i int) {
		fs[i] = alg.KFunc{ops[i], func() { fps[i](p, r, interp) }}
	}
	alg.Forall(inf, len(fs))
	alg.ExecF(fs, p.Operator)
	return
}

func reduceTerm(p, r *Predicate, itp NameBool) (ok bool) {
	v, ok := itp(p.String)
	if ok {
		if v {
			tr := True()
			*r = *tr
		} else {
			tr := False()
			*r = *tr
		}
	} else {
		*r = *p
	}
	return
}

func reduceNot(p, r *Predicate, itp NameBool) (ok bool) {
	nr := Reduce(p.B, itp)
	v, ok := false, nr.Operator == Term
	if ok {
		v, ok = itp(nr.String)
	}
	if ok {
		r.String = fmt.Sprint(!v)
		r.Operator = Term
	} else {
		r.B = nr
		r.Operator = NotOp
	}
	return
}

func reduceAnd(p, r *Predicate, itp NameBool) (ok bool) {
	ok = reduceUnit(p, r, true, itp)
	return
}

func reduceOr(p, r *Predicate, itp NameBool) (ok bool) {
	ok = reduceUnit(p, r, false, itp)
	return
}

func reduceUnit(p, r *Predicate, unit bool,
	itp NameBool) (ok bool) {
	ps0 := make([]*Predicate, 2)
	var pr *Predicate
	ps := []func(){
		func() { pr = Reduce(p.A, itp); ps0[0] = pr },
		func() { pr = Reduce(p.B, itp); ps0[1] = pr },
	}
	unitF, un := false, 0
	ib := func(i int) (b bool) {
		ps[i]() // this avoids superflous
		// evaluation if zero found
		v, ok := itp(pr.String)
		b = pr.Operator == Term && ok && v != unit
		if pr.Operator == Term && ok && v == unit {
			unitF, un = true, i
		}
		return
	}
	zeroF, _ := alg.BLnSrch(ib, len(ps))
	if zeroF {
		r.Operator = Term
		r.String = fmt.Sprint(!unit)
		ok = true
	} else if unitF {
		*r = *ps0[len(ps)-1-un]
		ok = true
	} else {
		if String(ps0[0]) == String(ps0[1]) {
			*r = *ps0[0]
		} else {
			r.Operator = p.Operator
			r.A, r.B = ps0[0], ps0[1]
		}
	}
	return
}

func reduceEquivales(p, r *Predicate, itp NameBool) (ok bool) {
	ps := []*Predicate{Reduce(p.A, itp), Reduce(p.B, itp)}
	// A ≡ true ≡ A
	// A ≡ false ≡ ¬A
	ib := func(i int) (b bool) {
		b = ps[i].String == TrueStr || ps[i].String == FalseStr
		return
	}
	ok, n := alg.BLnSrch(ib, len(ps))
	if ok {
		if ps[n].String == TrueStr {
			*r = *ps[len(ps)-1-n]
		} else {
			*r = *negate(ps[len(ps)-1-n])
		}
	} else if String(ps[0]) == String(ps[1]) {
		*r = *True()
		ok = true
	} else if String(ps[0]) == String(negate(ps[1])) ||
		String(negate(ps[0])) == String(ps[1]) {
		*r = *False()
		ok = true
	} else {
		r.Operator = EquivalesOp
		r.A = ps[0]
		r.B = ps[1]
	}
	return
}

func negate(p *Predicate) (r *Predicate) {
	if p.String == TrueStr {
		r = False()
	} else if p.String == FalseStr {
		r = True()
	} else {
		r = &Predicate{Operator: NotOp, B: p}
	}
	return
}

func reduceImplies(p, r *Predicate, itp NameBool) (ok bool) {
	ps := []*Predicate{Reduce(p.A, itp), Reduce(p.B, itp)}
	ib := func(i int) (b bool) {
		b = ps[i].String == TrueStr || ps[i].String == FalseStr
		return
	}
	ok, _ = alg.BLnSrch(ib, len(ps))
	if ok {
		// a ⇒ b ≡ ¬a ∨ b
		np := &Predicate{
			Operator: OrOp,
			A:        negate(p.A),
			B:        p.B,
		}
		reduceOr(np, r, itp)
	} else {
		r.Operator = ImpliesOp
		r.A = ps[0]
		r.B = ps[1]
	}
	return
}

func reduceFollows(p, r *Predicate, itp NameBool) (ok bool) {
	// b ⇐ a ≡ a ⇒ b
	np := &Predicate{Operator: ImpliesOp, A: p.B, B: p.A}
	ok = reduceImplies(np, r, itp)
	if !ok {
		r.Operator = FollowsOp
		r.A, r.B = r.B, r.A
	}
	return
}

func reduceNotEquivales(p, r *Predicate, itp NameBool) (ok bool) {
	// a ≢ b ≡ a ≡ ¬b"
	np := &Predicate{Operator: EquivalesOp, A: p.A, B: negate(p.B)}
	ok = reduceEquivales(np, r, itp)
	if !ok {
		r.Operator = NotEquivalesOp
		r.B = r.B.A
	}
	return
}

const (
	TrueStr  = "true"
	FalseStr = "false"
)

func True() (r *Predicate) {
	r = NewTerm(TrueStr)
	return
}

func False() (r *Predicate) {
	r = NewTerm(FalseStr)
	return
}

func NewTerm(s string) (p *Predicate) {
	p = &Predicate{
		Operator: Term,
		String:   s,
	}
	return
}

func String(p *Predicate) (r string) {
	if p.Operator == Term {
		r = p.String
	} else if p.Operator == NotOp {
		if p.B == nil {
			panic("Malformed ¬ predicate:" + p.String)
		}
		var sfm string
		if p.B.Operator == Term {
			sfm = "%s"
		} else {
			sfm = "(%s)"
		}
		r = fmt.Sprintf("%s"+sfm, NotOp, String(p.B))
	} else {
		r = fmt.Sprintf(
			format(p.Operator, p.A.Operator)+" %s "+
				format(p.Operator, p.B.Operator),
			String(p.A), p.Operator, String(p.B))
	}
	return
}

func format(oa, ob string) (r string) {
	priority := map[string]int{
		Term:           3,
		NotOp:          3,
		AndOp:          2,
		OrOp:           2,
		ImpliesOp:      1,
		FollowsOp:      1,
		EquivalesOp:    0,
		NotEquivalesOp: 0,
	}
	pa, pb := priority[oa], priority[ob]
	if pa <= pb && !(pa == pb && (pa == 1 || pa == 2) && oa != ob) {
		// the second conjunct is for excluding the case
		// when the pair contains ∧,∨ or ⇒,⇐ which need
		// parenthesis if appear in sequence since they aren't
		// associative, like ≡,≢
		r = "%s"
	} else {
		r = "(%s)"
	}
	return
}

func (p *Predicate) Valid() (ok bool) {
	if p.Operator == NotOp {
		ok = p.A == nil && p.B != nil
		ok = ok && p.B.Valid()
	} else if p.Operator == Term {
		ok = p.String != ""
	} else {
		ops := []string{AndOp, OrOp, ImpliesOp, EquivalesOp,
			NotEquivalesOp, FollowsOp}
		ib := func(i int) bool { return p.Operator == ops[i] }
		ok, _ = alg.BLnSrch(ib, len(ops))
		ok = ok && p.A != nil && p.B != nil
		ok = ok && p.A.Valid() && p.B.Valid()
		ok = ok && p.String == ""
	}
	return
}
