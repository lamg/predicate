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
	"fmt"
	"github.com/lamg/predicate"
	"log"
	"os"
	"strings"
)

func main() {
	var e error
	if len(os.Args) != 2 {
		e = fmt.Errorf("Need one argument, not %d", len(os.Args))
	}
	if e == nil {
		var p *predicate.Predicate
		p, e = predicate.Parse(strings.NewReader(os.Args[1]))
		if e == nil {
			np := predicate.Reduce(p,
				func(name string) (b, def bool) {
					b, def = name == "true",
						name == "true" || name == "false"
					return
				})
			fmt.Println(predicate.String(np))
		}
	}
	if e != nil {
		log.Fatal(e)
	}
}
