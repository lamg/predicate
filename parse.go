package predicate

import (
	"bufio"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"unicode"
	"unicode/utf8"
)

/*
Grammar in EBNF syntax

predicate = term {('≡'|'≢') term}.
term = junction ({'⇒' junction} | {'⇐' junction}).
junction = factor ({'∨' factor} | {'∧' factor}).
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
		&scanStatePreserver{tkf: tokens(rd, ss)},
	}

	p, e = st.predicate()
	return
}

type predState struct {
	*scanStatePreserver
}

func (s *predState) predicate() (p *Predicate, e error) {
	factor := s.factor()
	conjunction := s.parseOp(factor, AndOp)
	disjunction := s.parseOp(factor, OrOp)
	junction := s.alternative("junction", disjunction, conjunction)
	implication := s.parseOp(junction, ImpliesOp)
	consequence := s.parseOp(junction, FollowsOp)
	term := s.alternative("term", implication, consequence)
	p, e = s.parseOp(term, EquivalesOp, NotEquivalesOp)()
	return
}

func (s *predState) alternative(name string,
	xs ...func() (*Predicate, error)) func() (*Predicate, error) {
	return func() (p *Predicate, e error) {
		s.save()
		bf := func(i int) (b bool) {
			p, e = xs[i]()
			if e != nil {
				s.back()
				println("back:", s.curr)
			}
			b = e == nil
			return
		}
		ok, _ := alg.BLnSrch(bf, len(xs))
		if !ok {
			e = errorAlt(name, e)
		} else {
			e = nil
		}
		s.drop()
		return
	}
}

func (s *predState) parseOp(sym func() (*Predicate, error),
	ops ...string) func() (*Predicate, error) {
	return func() (p *Predicate, e error) {
		p, e = sym()
		curr := p
		o, end := "", false
		for !end {
			if o != "" {
				var b *Predicate
				b, e = sym()
				if e == nil {
					old := &Predicate{
						Operator: curr.Operator,
						A:        curr.A,
						B:        curr.B,
						String:   curr.String,
					}
					curr.Operator, curr.A, curr.B = o, old, b
					curr = curr.B
				}
				o = ""
			} else {
				// first time it could fail
				o, e = moreOps(s.scanStatePreserver, ops)
				end = o == "" || e != nil
			}
		}
		return
	}
}

func moreOps(s *scanStatePreserver, ops []string) (
	op string, e error,
) {
	t, e := s.token()
	if e == nil {
		ib := func(i int) bool { return ops[i] == t.value }
		ok, _ := alg.BLnSrch(ib, len(ops))
		if ok {
			op = t.value
		} else if !(t.value == CPar || t.value == string(eof)) {
			e = fmt.Errorf("Not found operator %s in %v", t.value, ops)
		}
	}
	return
}

func (s *predState) factor() func() (*Predicate, error) {
	return func() (p *Predicate, e error) {
		t, e := s.token()
		var nt *Predicate
		if e == nil {
			if t.value == NotOp {
				nt = &Predicate{Operator: NotOp}
				t, e = s.token()
			}
			if e == nil {
				if t.isIdent {
					p = &Predicate{Operator: Term, String: t.value}
				} else if t.value == OPar {
					p, e = s.predicate()
					if e == nil {
						t, e = s.token()
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
		if e == nil && nt != nil {
			nt.B = p
			p = nt
		}
		return
	}
}

type scanStatePreserver struct {
	tkf    func() (*token, error)
	stored []*token
	curr   int
	saved  []int
}

func (s *scanStatePreserver) token() (t *token, e error) {
	if s.curr == len(s.stored) {
		t, e = s.tkf()
		if e == nil && t.value == "" {
			t, e = s.tkf()
		}
		if e == nil {
			s.stored = append(s.stored, t)
		}
	} else {
		t = s.stored[s.curr]
	}
	if e == nil {
		s.curr = s.curr + 1
	}
	return
}

func (s *scanStatePreserver) save() {
	s.saved = append(s.saved, s.curr)
}

func (s *scanStatePreserver) back() {
	l := len(s.saved) - 1
	s.curr = s.saved[l]
}

func (s *scanStatePreserver) drop() {
	s.saved = s.saved[:len(s.saved)-1]
}

const (
	OPar = "("
	CPar = ")"
)

func errorAlt(name string, d error) (e error) {
	return fmt.Errorf("Error parsing %s (%v)", name, d)
}

func unexpEnd() (e error) {
	e = fmt.Errorf("Unexpected end of input")
	return
}

type token struct {
	value    string
	isIdent  bool
	isNumber bool
}

const (
	// 0x3 is the end of file character
	eof = 0x3
)

func tokens(source io.Reader, ss []scanner) (
	tf func() (*token, error)) {
	rd := bufio.NewReader(source)
	ss = append(ss, eofScan)
	var rn rune
	var sc func(rune) (*token, bool, bool)
	n, end, read, search, scan := 0, false, true, false, false
	tf = func() (t *token, e error) {
		if end {
			e = io.EOF
		}
		for !end {
			if read {
				rn, _, e = rd.ReadRune()
				if e != nil {
					rn = eof
					if e == io.EOF {
						e = nil
					}
				}
				read, search = false, !scan
			} else if search {
				if n == len(ss) {
					e, end = notRec(string(rn)), true
				} else {
					sc, n, search = ss[n](), n+1, false
				}
			} else if !search {
				t, read, end = sc(rn)
				search = !read || end
				scan = !search
			}
		}
		n, end = 0, e != nil || t.value == string(eof)
		return
	}
	return
}

func equalErr(a, b error) (ok bool) {
	ok = (a == nil) == (b == nil) &&
		(a == nil || a.Error() == b.Error())
	return
}

func notRec(s string) (e error) {
	return fmt.Errorf("Not recognized '%s'", s)
}

type scanner func() func(rune) (*token, bool, bool)

func identScan() func(rune) (*token, bool, bool) {
	var ident string
	return func(rn rune) (t *token, cont, prod bool) {
		cont = unicode.IsLetter(rn) ||
			(ident != "" && unicode.IsDigit(rn))
		if cont {
			ident = ident + string(rn)
		} else if ident != "" {
			t, prod = &token{value: ident, isIdent: true}, true
		}
		return
	}
}

func strScan(strScan string) (s scanner) {
	s = func() func(rune) (*token, bool, bool) {
		str := strScan
		return func(rn rune) (t *token, cont, prod bool) {
			sr, size := utf8.DecodeRuneInString(str)
			cont = sr != utf8.RuneError && sr == rn
			if cont {
				str = str[size:]
			}
			prod = len(str) == 0
			if prod {
				t, cont = &token{value: strScan}, true
			}
			return
		}
	}
	return
}

func spaceScan() func(rune) (*token, bool, bool) {
	start := false
	return func(rn rune) (t *token, cont, prod bool) {
		cont = unicode.IsSpace(rn)
		if cont {
			start = true
		}
		prod = start && !cont
		if prod {
			t, start = new(token), false
		}
		return
	}
}

func eofScan() func(rune) (*token, bool, bool) {
	return func(r rune) (t *token, cont, prod bool) {
		if r == eof {
			t, prod = &token{value: string(r)}, true
		}
		return
	}
}
