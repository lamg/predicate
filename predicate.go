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
)

type Predicate struct {
	Operator string              `json:"operator"`
	A        *Predicate          `json:"a"`
	B        *Predicate          `json:"b"`
	Val      func() (bool, bool) `json:"-"`
	String   string              `json:"string"`
}

const (
	NotOp          = "¬"
	AndOp          = "∧"
	OrOp           = "∨"
	EquivalesOp    = "≡"
	NotEquivalesOp = "≢"
	ImpliesOp      = "⇒"
	Term           = "term"
)

func Reduce(p *Predicate) (r *Predicate) {
	r = new(Predicate)
	id := func(p, r *Predicate) { *r = *p }
	fps := []func(*Predicate, *Predicate){
		reduceNot,
		reduceAnd,
		reduceOr,
		reduceImplies,
		reduceEquivales,
		reduceNotEquivales,
		id,
	}
	ops := []string{
		NotOp,
		AndOp,
		OrOp,
		ImpliesOp,
		EquivalesOp,
		NotEquivalesOp,
		Term,
	}
	fs := make([]kFunc, len(fps))
	inf := func(i int) {
		fs[i] = kFunc{ops[i], func() { fps[i](p, r) }}
	}
	forall(inf, len(fs))
	execKF(fs, p.Operator)
	return
}

func NoVal() (v, ok bool) {
	v, ok = false, false
	return
}

func reduceNot(p, r *Predicate) {
	nr := Reduce(p.A)
	if nr.Operator == Term {
		r.Val = func() (v, ok bool) {
			v, ok = nr.Val()
			v = !v
			return
		}
		v, ok := r.Val()
		if ok {
			r.String = fmt.Sprint(v)
		} else {
			r.String = nr.String
		}
		r.Operator = Term
	} else {
		r.A = nr
		r.Operator = NotOp
	}
}

func reduceAnd(p, r *Predicate) {
	reduceUnit(p, r, true)
}

func reduceOr(p, r *Predicate) {
	reduceUnit(p, r, false)
}

func reduceUnit(p, r *Predicate, unit bool) {
	ps := []*Predicate{Reduce(p.A), Reduce(p.B)}
	unitF, un := false, 0
	ib := func(i int) (b bool) {
		v, ok := ps[i].Val()
		b = ps[i].Operator == Term && ok && v != unit
		if ps[i].Operator == Term && ok && v == unit {
			unitF, un = true, i
		}
		return
	}
	zeroF, _ := bLnSrch(ib, len(ps))
	if zeroF {
		r.Operator = Term
		r.Val = func() (v bool, ok bool) {
			v, ok = !unit, true
			return
		}
		r.String = fmt.Sprint(!unit)
	} else if unitF {
		*r = *ps[len(ps)-1-un]
	} else {
		r.Operator = p.Operator
		r.A, r.B = ps[0], ps[1]
	}
}

func reduceEquivales(p, r *Predicate) {
	notImplemented()
}

func reduceImplies(p, r *Predicate) {
	notImplemented()
}

func reduceNotEquivales(p, r *Predicate) {
	notImplemented()
}

func True() (r *Predicate) {
	r = NewTerm("", true)
	return
}

func False() (r *Predicate) {
	r = NewTerm("", false)
	return
}

func StrF(s string) (f func() string) {
	f = func() string {
		return s
	}
	return
}

func NewVar(s string) (p *Predicate) {
	p = NewTerm(s, false)
	return
}

func NewTerm(s string, v bool) (p *Predicate) {
	p = &Predicate{
		Operator: Term,
		Val: func() (u, ok bool) {
			u, ok = v, s == ""
			return
		},
	}
	if s != "" {
		p.String = s
	} else {
		p.String = fmt.Sprint(v)
	}
	return
}

func String(p *Predicate) (r string) {
	if p.Operator == Term {
		r = p.String
	} else if p.Operator == NotOp {
		var sfm string
		if p.A.Operator == Term {
			sfm = "%s"
		} else {
			sfm = "(%s)"
		}
		r = fmt.Sprintf("%s"+sfm, NotOp, String(p.A))
	} else {
		r = fmt.Sprintf(
			format(p.Operator, p.A.Operator)+" %s "+
				format(p.Operator, p.B.Operator),
			String(p.A), p.Operator, String(p.B))
	}
	return
}

func format(oa, ob string) (r string) {
	assocOps := []string{
		AndOp, AndOp,
		OrOp, OrOp,
		EquivalesOp, EquivalesOp,
		EquivalesOp, NotEquivalesOp,
		AndOp, Term,
		OrOp, Term,
		EquivalesOp, Term,
		NotEquivalesOp, Term,
		AndOp, NotOp,
		OrOp, NotOp,
		EquivalesOp, NotOp,
		NotEquivalesOp, NotOp,
	}
	ib := func(i int) (b bool) {
		oi, oi1 := assocOps[2*i], assocOps[2*i+1]
		b = (oi == oa && oi1 == ob) || (oi == ob && oi1 == oa)
		return
	}
	ok, _ := bLnSrch(ib, len(assocOps)/2)
	if ok {
		r = "%s"
	} else {
		r = "(%s)"
	}
	return
}

const (
	OperandAK = "a"
	OperandBK = "b"
	StrK      = "str"
)

func (p *Predicate) To() (m map[string]interface{}) {
	notImplemented()
	return
}

func (p *Predicate) From(m map[string]interface{},
	f func(string) (func() (bool, bool), func() string)) (e error) {
	notImplemented()
	return
}

func (p *Predicate) ReplacePostorderAt(r *Predicate,
	n int) (ok bool) {
	notImplemented()
	return
}

func (p *Predicate) Valid() (ok bool) {
	if p.Operator == NotOp {
		ok = p.A != nil && p.B == nil
		ok = ok && p.A.Valid()
	} else if p.Operator == Term {
		ok = p.String != ""
	} else {
		ops := []string{AndOp, OrOp, ImpliesOp, EquivalesOp,
			NotEquivalesOp}
		ib := func(i int) bool { return p.Operator == ops[i] }
		ok, _ = bLnSrch(ib, len(ops))
		ok = ok && p.A != nil && p.B != nil
		ok = ok && p.A.Valid() && p.B.Valid()
	}
	return
}

func notImplemented() {
	panic("Not implemented")
}

type kFunc struct {
	k string
	f func()
}

func execKF(kf []kFunc, key string) (ok bool) {
	ib := func(i int) (b bool) {
		b = kf[i].k == key
		return
	}
	ok, n := bLnSrch(ib, len(kf))
	if ok {
		kf[n].f()
	}
	return
}

type intBool func(int) bool

// bLnSrch is the bounded lineal search algorithm
// { n ≥ 0 ∧ forall.n.(def.ib) }
// { i =⟨↑j: 0 ≤ j ≤ n ∧ ⟨∀k: 0 ≤ k < j: ¬ib.k⟩: j⟩
//   ∧ b ≡ i ≠ n }
func bLnSrch(ib intBool, n int) (b bool, i int) {
	b, i, udb := false, 0, true
	// udb: undefined b for i
	for !b && i != n {
		if udb {
			// udb ∧ i ≠ n
			b, udb = ib(i), false
		} else {
			// ¬udb ∧ ¬b
			i, udb = i+1, true
		}
	}
	return
}

type intF func(int)

func forall(inf intF, n int) {
	for i := 0; i != n; i++ {
		inf(i)
	}
}
