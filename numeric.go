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
product_chain = numeric {("×"|"÷") numeric}.
numeric = identifier | number | '(' predicate ')'.
*/

type Numeric struct {
	Operator string
	A        *Numeric
	B        *Numeric
	Value    string
}

func numberScan() func(rune) (*token, bool, bool) {
	var val string
	return func(rn rune) (t *token, cont, prod bool) {
		cont = unicode.IsDigit(rn)
		if cont {
			val = val + string(rn)
		}
		prod = val != "" && !cont
		if prod {
			t = &token{value: val, isNumber: true}
		}
		return
	}
}
