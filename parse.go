package predicate

import (
	"bufio"
	"fmt"
	alg "github.com/lamg/algorithms"
	"io"
	"unicode"
	"unicode/utf8"
)

func optional(back func(), sym func() error) {
	e := sym()
	if e != nil {
		back()
	}
}

func alternative(back func(), syms []func() error) (e error) {
	bf := func(i int) (b bool) {
		e = syms[i]()
		b = e == nil
		if !b {
			back()
		}
		return
	}
	alg.BLnSrch(bf, len(syms))
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

func tokens(source io.Reader, ss []scanner) (
	tf func() (*token, error)) {
	rd := bufio.NewReader(source)
	tf = func() (t *token, e error) {
		var rn rune
		var sc func(rune) (*token, bool, bool)
		n, end, read, search, scan := 0, false, true, false, false
		for !end {
			if read {
				rn, _, e = rd.ReadRune()
				if e == io.EOF {
					rn, e = 0x3, nil
				}
				read, search = false, !scan
			} else if search {
				if n == len(ss) {
					e = fmt.Errorf("Not recognized '%s'", string(rn))
				} else {
					sc, n, search = ss[n](), n+1, false
				}
			} else if !search {
				t, scan, end = sc(rn)
				println("rune: ", string(rn), "end:", end, "t = nil:", t == nil)
				search, read, end = !scan, scan, end && t.value != ""
			}
			end = end || e != nil
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
			t = new(token)
		}
		return
	}
}
