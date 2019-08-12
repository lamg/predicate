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
	testScan(t, ss, txt, tks)
}

func testScan(t *testing.T, ss []scanner, txt string,
	tks []string) {
	scan := tokens(strings.NewReader(txt), ss)

	inf := func(i int) {
		tk, e := scan()
		require.NoError(t, e)
		println("value:", tk.value)
		require.Equal(t, tks[i], tk.value)
		t.Log(tk)
	}
	alg.Forall(inf, len(tks))
}

func TestPredicateSym(t *testing.T) {
	ps, i := []Predicate{
		{Operator: Term, String: TrueStr},
		{Operator: Term, String: FalseStr},
	}, 0
	curr := func() (p *Predicate) {
		p, i = &ps[i], i+1
		return
	}
	r, b := predicateSym(curr)
	b(AndOp)
	p := r()
	expected := &Predicate{Operator: AndOp, A: True(), B: False()}
	require.Equal(t, expected, p)
}
