package main

import (
	"fmt"
	"math"
	"testing"
)

func check(A [][]float64) (e eigen, err error) {
	e, err = pm(A)
	if err != nil {
		return
	}
	// Ax=lx
	// Ax-lx=0
	n := len(A)
	res := make([]float64, n)

	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			res[row] += A[row][col] * e.𝑿[col]
		}
		res[row] -= e.𝜦 * e.𝑿[row]
	}

	var delta float64
	for i := range res {
		delta += math.Pow(res[i], 2.0)
	}
	delta = math.Sqrt(delta)

	fmt.Printf("Delta = %.10e\n", delta)

	if delta > 𝛆 {
		err = fmt.Errorf("Precition is not ok")
		return
	}

	return
}

func Test(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("2", func(t *testing.T) {
		e, err := check([][]float64{
			{4, 5},
			{6, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("E3", func(t *testing.T) {
		e, err := check([][]float64{
			{1, 0},
			{0, -1},
		})
		if err == nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("E4", func(t *testing.T) {
		e, err := check([][]float64{
			{2, 0, 0},
			{0, 2, 0},
			{0, 0, 1},
		})
		if err == nil {
			t.Fatal(err)
		}
		_ = e
	})
}
