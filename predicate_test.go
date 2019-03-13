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
	"encoding/json"
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
					B:        NewVar("A"),
				},
			},
			res: NewVar("A"),
		},
	}
	inf := func(i int) {
		require.True(t, ps[i].tov.Valid())
		require.True(t, ps[i].res.Valid())
		r := Reduce(ps[i].tov)
		require.Equal(t, ps[i].res.String, r.String)
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
			p: NewVar("X"),
			s: "X",
		},
		{
			p: &Predicate{
				Operator: NotOp,
				A:        NewVar("Y"),
			},
			s: "¬Y",
		},
		{
			p: &Predicate{
				Operator: AndOp,
				A:        True(),
				B:        False(),
			},
			s: "true ∧ false",
		},
		{
			p: &Predicate{
				Operator: OrOp,
				A:        NewVar("A"),
				B: &Predicate{
					Operator: OrOp,
					A:        NewVar("B"),
					B: &Predicate{
						Operator: OrOp,
						A:        NewVar("C"),
						B: &Predicate{
							Operator: AndOp,
							A:        NewVar("R"),
							B: &Predicate{
								Operator: NotOp,
								A:        NewVar("T"),
							},
						},
					},
				},
			},
			s: "A ∨ B ∨ C ∨ (R ∧ ¬T)",
		},
		{
			p: &Predicate{
				Operator: NotOp,
				A: &Predicate{
					Operator: OrOp,
					A:        NewVar("A"),
					B:        NewVar("B"),
				},
			},
			s: "¬(A ∨ B)",
		},
	}
	inf := func(i int) {
		rs := String(ts[i].p)
		require.Equal(t, ts[i].s, rs, "At %d", i)
	}
	forall(inf, len(ts))
}

func TestMarshal(t *testing.T) {
	ps := []*predStr{
		{
			p: &Predicate{
				Operator: AndOp,
				A:        True(),
				B:        False(),
			},
			s: `{"operator":"∧",` +
				`"a":{"operator":"term",` +
				`"a":null,"b":null,"string":"true"},` +
				`"b":{"operator":"term",` +
				`"a":null,"b":null,"string":"false"},` +
				`"string":""}`,
		},
	}
	inf := func(i int) {
		bs, e := json.Marshal(ps[i].p)
		require.NoError(t, e)
		require.Equal(t, ps[i].s, string(bs))
	}
	forall(inf, len(ps))
}
