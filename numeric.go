package predicate

import (
	"unicode"
)

const (
	EqualOp          = "="
	NotEqualOp       = "≠" // C-k !=
	GreaterOp        = ">"
	LesserOp         = "<"
	AtMostOp         = "≤" // C-k =<
	AtLeastOp        = "≥" // C-k >=
	SumOp            = "+"
	SubstractionOp   = "-"
	MultiplicationOp = "×" // C-k *X
	DivisionOp       = "÷" // C-k -:
)

/*
numeric_predicate = numeric_expression
	("≠" numeric_expression | operator_expr_chain ).
operator_expr_chain = at_most_chain | at_least_chain | equal_chain.
at_most_chain = at_most_elemet {at_most_element}.
at_most_element = (">"|"≥") numeric_expression.
at_least_chain = at_least_element {at_least_element}.
at_least_element = ("<"|"≤") numeric_expression.
numeric_expression = product_chain {("+"|"-") product_chain}.
product_chain = number {("×"|"÷") number}.
*/

func parseNumeric(ns []*Numeric) parser {
	return func(ts []token, i int,
		ext parser) (p *Predicate, n int, e error) {
		return
	}
}

type Numeric struct {
	Operator string
	A        *Numeric
	B        *Numeric
	Value    string
}

func parseNumericExpr(ts []token, i int) (p *Numeric, n int,
	e error) {
	p, n, e = parseProduct(ts, i)
	var op string
	if n != len(ts) {
		op = moreOps(ts, []string{SumOp, SubstractionOp})
		if op == "" {
			e = notFound(SumOp + " or " + SubstractionOp)
		}
	}
	if e == nil {
		if op == "" {

		}
	}
	if e == nil {

	}
	return
}

func parseProduct(ts []token, i int) (p *Numeric, n int, e error) {
	return
}

func numberScan() func(rune) (token, bool, bool) {
	var val string
	return func(rn rune) (t token, cont, prod bool) {
		cont = unicode.IsDigit(rn)
		if cont {
			val = val + string(rn)
		}
		prod = val != "" && !cont
		if prod {
			t = token{value: val, isNumber: true}
		}
		return
	}
}
