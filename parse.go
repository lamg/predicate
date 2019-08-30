package predicate

import (
	"bufio"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"unicode"
	"unicode/utf8"
)

func alternative(save, back, syms []func() error) (e error) {
	save()
	bf := func(i int) (b bool) {
		e = syms[i]()
		b = e == nil
		if !b {
			back()
		}
		return
	}
	ok, _ := alg.BLnSrch(bf, len(syms))
	if !ok {
		e = errorAlt()
	} else {
		e = nil
	}
	return
}

func parseOp(op func() (string, bool, error), sym func() error,
	branch func(string)) (e error) {
	e = sym()
	o, end := "", false
	for e == nil && !end {
		if o != "" {
			println("branch:", o)
			branch(o)
			e, o = sym(), ""
		} else {
			o, end, e = op()
		}
	}
	return
}

func moreOps(s *scanStatePreserver, ops []string) (
	op string, end bool, e error,
) {
	t, e := s.token()
	if e == nil {
		ib := func(i int) bool { return ops[i] == t.value }
		ok, _ := alg.BLnSrch(ib, len(ops))
		if ok {
			op = t.value
		} else {
			e = fmt.Errorf("Not found operator")
		}
	}
	end = e != nil || t.value == CPar
	if e == io.EOF {
		e = nil
	}
	return
}

type scanStatePreserver struct {
	tkf     func() (*token, error)
	stored  []*token
	restore bool
	save    bool
}

func (s *scanStatePreserver) token() (t *token, e error) {
	if s.restore {
		t, s.stored = s.stored[0], s.stored[1:]
		s.restore = len(s.stored) != 0
	} else {
		t, e = s.tkf()
		if e == nil && t.value == "" {
			t, e = s.tkf()
		}
	}
	if s.save && e == nil {
		s.stored = append(s.stored, t)
	}
	return
}

func (s *scanStatePreserver) saveState() {
	s.save, s.stored, s.restore = true, make([]*token, 0), false
}

func (s *scanStatePreserver) backToSaved() {
	s.save, s.restore = false, len(s.stored) != 0
}

const (
	OPar = "("
	CPar = ")"
)

func notRec(t string) (e error) {
	e = fmt.Errorf("Not recognized symbol \"%s\"", t)
	return
}

func errorAlt() (e error) {
	return fmt.Errorf("Error parsing alternative")
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

func tokens(source io.Reader, ss []scanner) (
	tf func() (*token, error)) {
	rd := bufio.NewReader(source)
	var rn rune
	var sc func(rune) (*token, bool, bool)
	n, end, read, search, scan := 0, false, true, false, false
	var err error
	tf = func() (t *token, e error) {
		if end {
			e = err
		}
		for !end {
			if read {
				rn, _, e = rd.ReadRune()
				if e == io.EOF {
					rn = 0x3 // 0x3 is the end of file character
				} else if e != nil {
					end = true
				}
				read, search = false, !scan
			} else if search {
				if n == len(ss) {
					e, end = fmt.Errorf("Not recognized '%s'", string(rn)),
						true
				} else {
					sc, n, search = ss[n](), n+1, false
				}
			} else if !search {
				t, read, end = sc(rn)
				search = !read || end
				scan = !search
			}
		}
		n, end, err = 0, e != nil, e
		if t != nil {
			e = nil
		}
		return
	}
	return
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
