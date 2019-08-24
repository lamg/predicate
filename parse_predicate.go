package predicate

import (
	"fmt"
	"io"
)

/*
Grammar in EBNF syntax

predicate = term {('≡'|'≢') term}.
term = implication | consequence | junction.
implication = junction {'⇒' junction}.
consequence = junction {'⇐' junction}.
junction = disjunction | conjunction | factor.
disjunction = factor {'∨' factor}.
conjunction = factor {'∧' factor}.
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.

*/

func Parse(rd io.Reader) (p *Predicate, e error) {
	root := new(Predicate)
	curr := root
	ss := []scanner{
		identScan,
		spaceScan,
		strScan(NotOp),
		strScan(AndOp),
		strScan(OrOp),
		strScan(EquivalesOp),
		strScan(NotEquivalesOp),
		strScan(ImpliesOp),
		strScan(OPar),
		strScan(CPar),
	}

	insert := func(p *Predicate) { curr.B = p; curr = curr.B }

	stp := &scanStatePreserver{tkf: tokens(rd, ss)}
	branch := func(op string) {
		old := &Predicate{
			Operator: curr.Operator,
			A:        curr.A,
			B:        curr.B,
			String:   curr.String,
		}
		curr.Operator, curr.A, curr.String = op, old, ""
	}
	ops := operators0(stp, branch)
	alt := alt0(stp)
	var predicate func() error

	factor := factor0(stp, predicate, insert)
	conjunction := ops(factor, AndOp)
	disjunction := ops(factor, OrOp)
	junction := alt(disjunction, conjunction, factor)

	implication := ops(junction, ImpliesOp)
	consequence := ops(junction, FollowsOp)
	term := alt(implication, consequence, junction)
	predicate = ops(term, EquivalesOp, NotEquivalesOp)

	e = predicate()
	p = root.B
	return
}

func factor0(stp *scanStatePreserver, predicate func() error,
	insert func(*Predicate)) func() error {
	return func() (e error) {
		stp.saveState()
		t, e := stp.token()
		if e == nil {
			if t.value == NotOp {
				insert(&Predicate{Operator: t.value})
				t, e = stp.token()
			}
			if e == nil {
				if t.isIdent {
					insert(&Predicate{Operator: Term, String: t.value})
				} else if t.value == OPar {
					e = predicate()
					if e == nil {
						t, e = stp.token()
						if e == nil && t.value != CPar {
							e = fmt.Errorf("Expecting closing parenthesis")
						}
					}
				} else {
					e = fmt.Errorf("Expecting identifier or opening " +
						"parenthesis")
				}
			}
		}
		return
	}
}

func operators0(s *scanStatePreserver, branch func(string)) func(func() error, ...string) func() error {
	return func(sym func() error, ops ...string) func() error {
		ms := func() (string, error) { return moreOps(s, ops) }
		return func() error { return parseOp(ms, sym, branch) }
	}
}

func alt0(s *scanStatePreserver) func(...func() error) func() error {
	return func(fs ...func() error) func() error {
		return func() error {
			return alternative(
				s.saveState,
				s.backToSaved,
				fs,
			)
		}
	}
}
