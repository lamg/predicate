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
	"github.com/stretchr/testify/require"
	"testing"
)

type toEvalResult struct {
	tov, res *Predicate
}

func TestReduce(t *testing.T) {

	ps := []*toEvalResult{
		{tov: True(), res: True()},
		{tov: &Predicate{Operator: NotOp, A: False()}, res: True()},
		{tov: &Predicate{Operator: NotOp, A: True()}, res: False()},
		{
			tov: &Predicate{Operator: AndOp, A: True(), B: False()},
			res: False(),
		},
		{
			tov: &Predicate{Operator: AndOp, A: False(), B: False()},
			res: False(),
		},
		{
			tov: &Predicate{Operator: OrOp, A: False(), B: False()},
			res: False(),
		},
		{
			tov: &Predicate{Operator: OrOp, A: False(), B: True()},
			res: True(),
		},
		{
			tov: &Predicate{
				Operator: NotOp,
				A: &Predicate{
					Operator: AndOp, A: True(), B: True(),
				},
			},
			res: False(),
		},
		{
			tov: &Predicate{
				Operator: NotOp,
				A: &Predicate{
					Operator: AndOp,
					A:        True(),
					B:        &Predicate{Operator: Term, Val: NoVal},
				},
			},
			res: &Predicate{Operator: Term, Val: NoVal},
		},
	}
	inf := func(i int) {
		r := Reduce(ps[i].tov)
		require.Equal(t, ps[i].res.Operator, r.Operator, "At %d", i)
		if ps[i].tov.Operator == Term {
			v, ok := ps[i].res.Val()
			av, aok := r.Val()
			require.Equal(t, ok, aok, "At %d", i)
			require.Equal(t, v, av, "At %d", i)
		}
	}
	forall(inf, len(ps))
}

type predStr struct {
	p *Predicate
	s string
}

func TestString(t *testing.T) {
	ts := []*predStr{
		{p: True(), s: "true"},
		{
			p: &Predicate{
				Operator: Term,
				String:   func() string { return "X" },
			},
			s: "X",
		},
		{
			p: &Predicate{
				Operator: NotOp,
				A: &Predicate{
					Operator: Term,
					String:   func() string { return "Y" },
				},
			},
			s: "¬(Y)",
		},
		{
			p: &Predicate{
				Operator: AndOp,
				A:        True(),
				B:        False(),
			},
			s: "(true) ∧ (false)",
		},
	}
	inf := func(i int) {
		rs := String(ts[i].p)
		require.Equal(t, ts[i].s, rs, "At %d", i)
	}
	forall(inf, len(ts))
}
