package main

import (
	"fmt"
	"math"
)

// результаты
type eigen struct {
	// собственные значения
	𝜦 float64

	// собственный вектор
	𝑿 []float64
}

func (e eigen) String() (out string) {
	if output {
		out += fmt.Sprintf("𝜦      = %+14.10e\n", e.𝜦)
		for i := range e.𝑿 {
			out += fmt.Sprintf("𝑿[%3d] = %+14.10e\n", i, e.𝑿[i])
		}
	}
	return
}

func PrintEigens(e []eigen) {
	if output {
		for i := range e {
			fmt.Printf("--- %5d ---\n%s", i, e[i])
		}
	}
}

func compare(e, es []eigen) bool {
	isSame := func(e1, e2 eigen) bool {
		if math.Abs(e1.𝜦-e2.𝜦) > 𝛆*100 {
			return false
		}
		oneMax(e1.𝑿, e1.𝑿)
		oneMax(e2.𝑿, e2.𝑿)
		for i := range e1.𝑿 {
			if math.Abs(e1.𝑿[i]-e2.𝑿[i]) > 𝛆*100 {
				return false
			}
		}
		return true
	}
	var found int
	for i := range e {
		for j := range es {
			if isSame(e[i], es[j]) {
				found++
				if output {
					fmt.Printf("Compare [%3d,%3d] is same\n", i, j)
				}
				break
			}
			if output {
				fmt.Printf("Compare [%3d,%3d] is not same\n", i, j)
			}
		}
	}
	isSameEigen := found == len(e) && found == len(es)
	return isSameEigen
}
