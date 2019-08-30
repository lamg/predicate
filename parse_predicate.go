package predicate

import (
	"fmt"
	"io"
)

/*
Grammar in EBNF syntax

predicate = term {('≡'|'≢') term}.
term = implication | consequence.
implication = junction {'⇒' junction}.
consequence = junction {'⇐' junction}.
junction = disjunction | conjunction.
disjunction = factor {'∨' factor}.
conjunction = factor {'∧' factor}.
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.

*/

func Parse(rd io.Reader) (p *Predicate, e error) {
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

	st := &predState{
		stp:  &scanStatePreserver{tkf: tokens(rd, ss)},
		root: new(Predicate),
	}
	st.curr = st.root
	println("****")

	e = st.predicate()
	p = st.root.B
	if len(st.stp.stored) != 0 {
		println("value:", st.stp.stored[len(st.stp.stored)-1].value)
	}
	return
}

type predState struct {
	stp  *scanStatePreserver
	root *Predicate
	curr []*Predicate
}

func (s *predState) predicate() (e error) {
	factor := factor0(s.stp, s.predicate, s.insert)
	conjunction := s.ops(factor, AndOp)
	disjunction := s.ops(factor, OrOp)
	junction := s.alt(disjunction, conjunction)
	implication := s.ops(junction, ImpliesOp)
	consequence := s.ops(junction, FollowsOp)
	term := s.alt(implication, consequence)
	e = s.ops(term, EquivalesOp, NotEquivalesOp)()
	return
}

func factor0(stp *scanStatePreserver, predicate func() error,
	insert func(*Predicate)) func() error {
	return func() (e error) {
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
func (s *predState) insert(p *Predicate) {
	s.curr = append(s.curr, p)
}

func (s *predState) ops(f func() error,
	os ...string) func() error {
	ms := func() (string, bool, error) {
		return moreOps(s.stp, os)
	}
	return func() error { return parseOp(ms, f, s.branch) }
}

func (s *predState) branch(op string) {

	old := &Predicate{
		Operator: s.curr.Operator,
		A:        s.curr.A,
		B:        s.curr.B,
		String:   s.curr.String,
	}
	s.curr.Operator, s.curr.A, s.curr.String = op, old, ""

}

func (s *predState) save() {
	s.stp.saveState()
}

func (s *predState) back() {
	s.stp.backToSaved()
	if len(s.curr) != 0 {
		s.curr = s.curr[:len(s.curr)-1]
	} else {
		panic("No available subtrees")
	}
}

func (s *predState) alt(fs ...func() error) func() error {
	return func() error {
		return alternative(s.save, s.back, fs)
	}
}
