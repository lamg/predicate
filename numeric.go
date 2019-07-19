package predicate

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

func parseNumeric(ts []token, i int) (p *Predicate, n int,
	e error) {
	return
}
