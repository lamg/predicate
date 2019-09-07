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
	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	txt := "true¬∧∨≡≢⇒⇐()bla9   x3  (Abla)true"
	tks := []string{"true",
		NotOp, AndOp, OrOp, EquivalesOp, NotEquivalesOp,
		ImpliesOp, FollowsOp, OPar, CPar, "bla9", "", "x3", "",
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
	testScan(t, ss, txt, tks)
}

func testScan(t *testing.T, ss []scanner, txt string,
	tks []string) {
	scan := tokens(strings.NewReader(txt), ss)

	inf := func(i int) {
		tk, e := scan()
		require.NoError(t, e)
		require.Equal(t, tks[i], tk.value)
		t.Log(tk)
	}
	alg.Forall(inf, len(tks))
}

func TestFuncPointer(t *testing.T) {
	var f func()
	g := func() {
		f()
	}
	a := false
	f = func() { a = true }
	g()
	require.True(t, a)
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
	var tk *token
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

func TestParseOp(t *testing.T) {
	ps := []string{"a", "¬a", "a ∧ b", "a ∧ b ∧ ¬c ∧ d"}
	ss := []scanner{spaceScan, identScan, strScan(AndOp),
		strScan(NotOp)}
	inf := func(i int) {
		rd := strings.NewReader(ps[i])
		st := &predState{tkf: tokens(rd, ss)}
		p, e := st.parseOp("factor", st.factor(), false, AndOp)()
		require.NoError(t, e, "At %d", i)
		require.Equal(t, ps[i], String(p), "At %d", i)
	}
	alg.Forall(inf, len(ps))
}

func TestParse(t *testing.T) {
	ps := []struct {
		pred string
		e    error
	}{
		{"", errIdentOrOpening()},
		{"true ∧ false", nil},
		{"true ∧", errIdentOrOpening()},
		{"¬A", nil},
		{"¬A ∧ (B ∨ C)", nil},
		{"A ∨ ¬(B ∧ C)", nil},
		{"A ≡ B ≢ ¬C ⇒ D", nil},
		{"A ≡ B ≡ ¬C ⇐ D", nil},
		{"A ≡ B ≡ ¬(C ⇐ D)", nil},
		{"A ∨ B ∨ C", nil},
		{"A ∨ B ∧ C", errUnkChar("∧")},
		{"A ⇒ B ⇐ C", errUnkChar("⇐")},
		{"A ∨ (B ∧ C)", nil},
		{"A ⇒ (B ⇐ C)", nil},
		{"a ≡ b ≢ c ≡ ¬x ∧ (¬z ≡ y) ≢ true", nil},
	}
	inf := func(i int) {
		np, e := Parse(strings.NewReader(ps[i].pred))
		require.Equal(t, e == nil, ps[i].e == nil,
			"At %d: %s %v", i, ps[i].pred, e)
		if e == nil {
			s := String(np)
			t.Logf("'%s'", s)
			require.Equal(t, ps[i].pred, s, "At %d", i)
		} else {
			t.Logf("'%s' → %s", ps[i].pred, e.Error())
			require.Equal(t, ps[i].e.Error(), e.Error(), "At '%s'",
				ps[i].pred)
		}
	}
	alg.Forall(inf, len(ps))
}
