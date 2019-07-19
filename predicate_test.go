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
	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
	"strings"
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
					B:        &Predicate{Operator: NotOp, A: NewTerm("A")},
				},
			},
			res: &Predicate{
				Operator: NotOp,
				A:        &Predicate{Operator: NotOp, A: NewTerm("A")},
			},
		},
		{
			tov: &Predicate{
				Operator: AndOp,
				A:        NewTerm("A"),
				B:        NewTerm("A"),
			},
			res: NewTerm("A"),
		},
		{
			tov: &Predicate{
				Operator: ImpliesOp,
				A:        True(),
				B:        False(),
			},
			res: False(),
		},
		{
			tov: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("A"),
				B:        True(),
			},
			res: NewTerm("A"),
		},
		{
			tov: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("A"),
				B:        False(),
			},
			res: negate(NewTerm("A")),
		},
		{
			tov: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("A"),
				B:        NewTerm("A"),
			},
			res: True(),
		},
		{
			// A ≡ ¬A ≡ false
			tov: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("A"),
				B:        negate(NewTerm("A")),
			},
			res: False(),
		},
		{
			// A ≢ A ≡ false
			tov: &Predicate{
				Operator: NotEquivalesOp,
				A:        NewTerm("A"),
				B:        NewTerm("A"),
			},
			res: False(),
		},
		{
			// A ⇐ true ≡ A
			tov: &Predicate{
				Operator: FollowsOp,
				A:        NewTerm("A"),
				B:        True(),
			},
			res: NewTerm("A"),
		},
		{
			// B ≡ A ≡ C ∧ ¬true → B ≡ ¬A
			tov: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("B"),
				B: &Predicate{
					Operator: EquivalesOp,
					A:        NewTerm("A"),
					B: &Predicate{
						Operator: AndOp,
						A:        NewTerm("C"),
						B:        negate(True()),
					},
				},
			},
			res: &Predicate{
				Operator: EquivalesOp,
				A:        NewTerm("B"),
				B:        negate(NewTerm("A")),
			},
		},
		{
			// with X = true
			tov: NewTerm("X"),
			res: True(),
		},
		{
			tov: &Predicate{
				Operator: AndOp,
				A:        False(),
				B:        NewTerm("Y"),
			},
			res: False(),
		},
		{
			tov: &Predicate{
				Operator: OrOp,
				A:        True(),
				B:        NewTerm("Y"),
			},
			res: True(),
		},
	}
	itp := func(n string) (v, ok bool) {
		v, ok = n == TrueStr || n == "X",
			n == TrueStr || n == FalseStr || n == "X"
		if n == "Y" {
			t.Fatalf("%s cannot be evaluated when zero found", n)
		}
		return
	}
	inf := func(i int) {
		require.True(t, ps[i].tov.Valid())
		require.True(t, ps[i].res.Valid())
		r := Reduce(ps[i].tov, itp)
		stov, sr := String(ps[i].tov), String(r)
		require.Equal(t, String(ps[i].res), sr)
		t.Logf("%s → %s", stov, sr)
	}
	alg.Forall(inf, len(ps))
}

func TestNot(t *testing.T) {
	itp := func(n string) (v, ok bool) {
		v, ok = n == "true", n != "A"
		return
	}
	p := &Predicate{
		Operator: NotOp,
		A:        &Predicate{Operator: NotOp, A: NewTerm("A")},
	}
	nr := new(Predicate)
	reduceNot(p, nr, itp)
	require.Equal(t, String(p), String(nr))
}

type predStr struct {
	p *Predicate
	s string
}

func TestString(t *testing.T) {
	ts := []*predStr{
		{p: True(), s: "true"},
		{
			p: NewTerm("X"),
			s: "X",
		},
		{
			p: &Predicate{
				Operator: NotOp,
				A:        NewTerm("Y"),
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
				A:        NewTerm("A"),
				B: &Predicate{
					Operator: OrOp,
					A:        NewTerm("B"),
					B: &Predicate{
						Operator: OrOp,
						A:        NewTerm("C"),
						B: &Predicate{
							Operator: AndOp,
							A:        NewTerm("R"),
							B: &Predicate{
								Operator: NotOp,
								A:        NewTerm("T"),
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
					A:        NewTerm("A"),
					B:        NewTerm("B"),
				},
			},
			s: "¬(A ∨ B)",
		},
	}
	inf := func(i int) {
		rs := String(ts[i].p)
		require.Equal(t, ts[i].s, rs, "At %d", i)
	}
	alg.Forall(inf, len(ts))
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
	alg.Forall(inf, len(ps))
}

func TestStrScan(t *testing.T) {
	ns := strScan(NotOp)()
	tk, cont, prod := ns('¬')
	require.True(t, prod)
	require.True(t, cont)
	require.Equal(t, "¬", tk.value)
}

func TestIdentScan(t *testing.T) {
	ids := identScan()
	rs := []rune{'a', 'b', 'c', '0'}
	var tk token
	var cont, prod bool
	inf := func(i int) {
		_, cont, prod = ids(rs[i])
		require.True(t, cont)
		require.False(t, prod)
	}
	alg.Forall(inf, len(rs))
	tk, cont, prod = ids(' ')
	require.False(t, cont)
	require.True(t, prod)
	require.Equal(t, "abc0", tk.value)

	ids0 := identScan()
	_, cont, prod = ids0(' ')
	require.False(t, cont)
	require.False(t, prod)
}

func TestScan(t *testing.T) {
	txt := "true¬∧∨≡≢⇒⇐()bla9   x3  (Abla)true"
	tks := []string{"true",
		NotOp, AndOp, OrOp, EquivalesOp, NotEquivalesOp,
		ImpliesOp, FollowsOp, OPar, CPar, "bla9", "x3",
		"(", "Abla", ")", "true"}
	ss := []scanner{
		strScan(NotOp),
		strScan(AndOp),
		strScan(OrOp),
		strScan(EquivalesOp),
		strScan(NotEquivalesOp),
		strScan(ImpliesOp),
		strScan(FollowsOp),
		strScan(OPar),
		strScan(CPar),
		identScan,
		spaceScan,
	}
	scanned, e := tokens(strings.NewReader(txt), ss)

	require.NoError(t, e)
	require.Equal(t, len(tks), len(scanned))
	inf := func(i int) {
		require.Equal(t, tks[i], scanned[i].value)
		t.Log(scanned[i])
	}
	alg.Forall(inf, len(tks))
}

func TestParse(t *testing.T) {
	ps := []struct {
		pred string
		e    error
	}{
		{"true ∧ false", nil},
		{"true ∧", notRec("∧")},
		{"¬A", nil},
		{"¬A ∧ (B ∨ C)", nil},
		{"A ∨ ¬(B ∧ C)", nil},
		{"A ≡ B ≡ ¬C ⇒ D", nil},
		{"A ≡ B ≡ ¬C ⇐ D", nil},
		{"A ≡ B ≡ ¬(C ⇐ D)", nil},
		{"A ∨ B ∨ C", nil},
		{"A ∨ B ∧ C", notRec("∧")},
		{"A ⇒ B ⇐ C", notRec("⇐")},
		{"A ∨ (B ∧ C)", nil},
		{"A ⇒ (B ⇐ C)", nil},
	}
	inf := func(i int) {
		np, e := Parse(strings.NewReader(ps[i].pred))
		require.Equal(t, e == nil, ps[i].e == nil,
			"At %d: %s %v", i, ps[i].pred, e)
		if e == nil {
			s := String(np)
			t.Logf("'%s'", s)
			require.Equal(t, ps[i].pred, s)
		} else {
			t.Logf("'%s' → %s", ps[i].pred, e.Error())
			require.Equal(t, ps[i].e.Error(), e.Error(), "At '%s'",
				ps[i].pred)
		}
	}
	alg.Forall(inf, len(ps))
}
