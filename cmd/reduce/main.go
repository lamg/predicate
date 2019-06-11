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

package main

import (
	"bufio"
	"fmt"
	pred "github.com/lamg/predicate"
	"log"
	"os"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		t := sc.Text()
		p, e := pred.Parse(strings.NewReader(t))
		if e == nil {
			stdInterp := func(name string) (val, def bool) {
				val, def = name == pred.TrueStr,
					name == pred.TrueStr || name == pred.FalseStr
				return
			}
			np := pred.Reduce(p, stdInterp)
			fmt.Println(pred.String(np))
		} else {
			log.Println(e.Error())
		}
	}
	e := sc.Err()
	if e != nil {
		log.Fatal(e)
	}
}
