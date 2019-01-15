package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func generator(es []eigen) (A [][]float64) {
	for i := range es {
		if len(es[i].𝑿) != len(es) {
			panic(i)
		}
		for j := range es[i].𝑿 {
			if es[i].𝑿[j] > 1 || es[i].𝑿[j] < -1 {
				panic("border")
			}
		}
	}
	n := len(es)
	A = make([][]float64, n)
	for i := 0; i < n; i++ {
		A[i] = make([]float64, n)
	}

	B := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			B.Set(i, j, es[j].𝑿[i])
		}
	}
	BT := B.T()

	var lu mat.LU
	lu.Factorize(BT)
	for i := 0; i < n; i++ {
		b := mat.NewDense(n, 1, nil)
		for j := 0; j < n; j++ {
			b.Set(j, 0, es[j].𝜦*es[j].𝑿[i])
		}

		var x mat.Dense
		if err := lu.Solve(&x, false, b); err != nil {
			continue
		}

		for j := 0; j < n; j++ {
			A[i][j] = x.At(j, 0)
		}
	}
	return
}

func MatrixPrint(A [][]float64) {
	if output {
		n := len(A)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				fmt.Printf("|%+20.10f|", A[i][j])
			}
			fmt.Printf("\n")
		}
	}
}

func ExampleGenerator() {

	// output true
	oldOut := output
	output = false
	defer func() {
		output = oldOut
	}()

	es := []eigen{
		{𝜦: +2.0, 𝑿: []float64{+0.5714286, +0.1428572, +1.0000000}},
		{𝜦: -2.0, 𝑿: []float64{-0.6666667, -1.0000000, -1.0000000}},
		{𝜦: -1.0, 𝑿: []float64{+0.5773503, +0.5773503, +0.5773503}},
	}
	A := generator(es)
	output = true
	MatrixPrint(A)
	output = false

	fmt.Println("change 0 <=> 2")
	es[0], es[2] = es[2], es[0]
	for i := 0; i < 3; i++ {
		es[i].𝑿[0], es[i].𝑿[2] = es[i].𝑿[2], es[i].𝑿[0]
	}
	A = generator(es)
	output = true
	MatrixPrint(A)
	output = false

	// Output:
	// |       +1.0000003000||       -3.0000003833||       +1.0000000833|
	// |       +3.0000003000||       -3.0000003833||       -0.9999999167|
	// |       +3.0000003000||       -5.0000003833||       +1.0000000833|
	// change 0 <=> 2
	// |       +1.0000000833||       -5.0000003833||       +3.0000003000|
	// |       -0.9999999167||       -3.0000003833||       +3.0000003000|
	// |       +1.0000000833||       -3.0000003833||       +1.0000003000|
}
