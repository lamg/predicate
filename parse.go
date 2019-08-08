package predicate

import (
	"bufio"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"unicode"
	"unicode/utf8"
)

func predicateSym(curr func() *Predicate) (r func() *Predicate,
	b func(string)) {
	p := curr()
	c := p
	b = func(s string) {
		n := curr()
		c.Operator, c.A, c.B = s, c, n
		c = n
	}
	r = func() *Predicate { return p }
	return
}

func parseOp0(op func() (string, error), sym func() error,
	branch func(string), notFoundOps func() error) (e error) {
	var o string
	e = sym()
	if e == nil {
		o, e = op()
	}
	for e == nil && o != "" {
		branch(o)
		e = sym()
		o, e = op()
	}
	return
}

func moreOps(ts func() (*token, bool), next func(),
	ops []string) (op string, e error) {
	t, ok := ts()
	if ok {
		ib := func(i int) bool { return ops[i] == t.value }
		ok, _ = alg.BLnSrch(ib, len(ops))
		if ok {
			op = t.value
			next()
		}
	} else {
		e = unexpEnd()
	}
	return
}

const (
	OPar = "("
	CPar = ")"
)

func notRec(t string) (e error) {
	e = fmt.Errorf("Not recognized symbol \"%s\"", t)
	return
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

func tokens(source io.Reader, ss []scanner) (ts []token,
	e error) {
	var rn rune
	var sc func(rune) (token, bool, bool)
	var t token
	rd := bufio.NewReader(source)
	n, prod, read, search, scan := 0, false, true, false, false
	for prod || scan || e == nil {
		if read {
			rn, _, e = rd.ReadRune()
			if e != nil {
				rn = ' '
			}
			read, search = false, !scan
		} else if search {
			if n == len(ss) {
				e = fmt.Errorf("Not recognized '%s'", string(rn))
			} else {
				sc, n, scan, search = ss[n](), n+1, true, false
			}
		} else if scan && !prod {
			t, scan, prod = sc(rn)
			search, read = !scan && !prod, scan
		} else if prod {
			if t.value != "" {
				ts = append(ts, t)
			}
			n, prod, search = 0, false, true
		}
	}
	if e == io.EOF {
		e = nil
	}
	return
}

type scanner func() func(rune) (token, bool, bool)

func identScan() func(rune) (token, bool, bool) {
	var ident string
	return func(rn rune) (t token, cont, prod bool) {
		cont = unicode.IsLetter(rn) ||
			(ident != "" && unicode.IsDigit(rn))
		if cont {
			ident = ident + string(rn)
		} else if ident != "" {
			t, prod = token{value: ident, isIdent: true}, true
		}
		return
	}
}

func strScan(strScan string) (s scanner) {
	s = func() func(rune) (token, bool, bool) {
		str := strScan
		return func(rn rune) (t token, cont, prod bool) {
			sr, size := utf8.DecodeRuneInString(str)
			cont = sr != utf8.RuneError && sr == rn
			if cont {
				str = str[size:]
			}
			prod = len(str) == 0
			if prod {
				t, cont = token{value: strScan}, true
			}
			return
		}
	}
	return
}

func spaceScan() func(rune) (token, bool, bool) {
	start := false
	return func(rn rune) (t token, cont, prod bool) {
		cont = unicode.IsSpace(rn)
		if cont {
			start = true
		}
		prod = start && !cont
		return
	}
}
