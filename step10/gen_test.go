package main

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// результаты
type eigen struct {
	// собственные значения
	𝜦 float64

	// собственный вектор
	𝑿 []float64
}

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

func TestGenerator(t *testing.T) {
	A := generator([]eigen{
		{
			𝜦: 2.0,
			𝑿: []float64{
				0.5714286,
				0.1428572,
				1.0000000,
			},
		},
		{
			𝜦: -2.0,
			𝑿: []float64{
				-0.6666667,
				-1.0000000,
				-1.0000000,
			},
		},
		{
			𝜦: -1.0,
			𝑿: []float64{
				0.5773503,
				0.5773503,
				0.5773503,
			},
		},
	})
	n := 3
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%+8.6e ", A[i][j])
		}
		fmt.Print("\n")
	}
}
