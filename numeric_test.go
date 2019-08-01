package predicate

import (
	"testing"
)

func TestNumberScan(t *testing.T) {
	text := " 3939 3939    000"
	tks := []string{"3939", "3939", "000"}
	ss := []scanner{
		spaceScan,
		numberScan,
	}
	testScan(t, ss, text, tks)
}
