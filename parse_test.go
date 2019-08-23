package predicate

import (
	"fmt"
	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
	"io"
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
		if e != io.EOF {
			require.NoError(t, e)
			require.Equal(t, tks[i], tk.value)
			t.Log(tk)
		}
	}
	alg.Forall(inf, len(tks))
}

func TestParseOp(t *testing.T) {
	txt := "pepe ∧ false ∧ coco"
	ss := []scanner{
		spaceScan,
		identScan,
		strScan(AndOp),
	}
	stp := &scanStatePreserver{
		tkf: tokens(strings.NewReader(txt), ss),
	}
	root := new(Predicate)
	curr := root
	sym := func() (e error) {
		tk, e := stp.token()
		if e == nil {
			if tk.isIdent {
				curr.B = &Predicate{Operator: Term, String: tk.value}
				curr = curr.B
			} else {
				e = fmt.Errorf("Expecting identifier, got '%s'", tk.value)
			}
		}
		return
	}
	branch := func(op string) {
		old := &Predicate{
			Operator: curr.Operator,
			A:        curr.A,
			B:        curr.B,
			String:   curr.String,
		}
		curr.Operator, curr.A, curr.String = op, old, ""
	}
	op := func() (string, error) {
		return moreOps(stp, []string{AndOp})
	}
	e := parseOp(op, sym, branch)
	require.NoError(t, e)
	expected := &Predicate{
		Operator: AndOp,
		A:        &Predicate{Operator: Term, String: "pepe"},
		B: &Predicate{
			Operator: AndOp,
			A:        False(),
			B:        &Predicate{Operator: Term, String: "coco"},
		},
	}
	require.Equal(t, expected, root.B)
}
