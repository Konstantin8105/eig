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

	if output {
		fmt.Printf("Delta = %.10e\n", delta)
	}

	if delta > 𝛆*10.0 {
		err = fmt.Errorf("Precition is not ok. %.5e > %.5e", delta, 𝛆)
		return
	}

	return
}

func Test(t *testing.T) {
	t.Run("example: 1", func(t *testing.T) {
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("example: 2", func(t *testing.T) {
		e, err := check([][]float64{
			{4, 5},
			{6, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("No dominant: 1", func(t *testing.T) {
		e, err := check([][]float64{
			{1, 0},
			{0, -1},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})
	t.Run("No dominant: 2", func(t *testing.T) {
		e, err := check([][]float64{
			{2, 0, 0},
			{0, 2, 0},
			{0, 0, 1},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})
	t.Run("No dominant: 3", func(t *testing.T) {
		e, err := check([][]float64{
			{-3, 0},
			{1, 3},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})

	t.Run("Low ratio : |𝜦2|/|𝜦1| = 0.1", func(t *testing.T) {
		e, err := check([][]float64{
			{4, 5},
			{6, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("Big ratio : |𝜦2|/|𝜦1| = 0.9", func(t *testing.T) {
		e, err := check([][]float64{
			{-4, 10},
			{7, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("matrix size: zero", func(t *testing.T) {
		e, err := check([][]float64{})
		if err == nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("matrix size: nil", func(t *testing.T) {
		e, err := check(nil)
		if err == nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("matrix size: rectangle", func(t *testing.T) {
		e, err := check([][]float64{
			{4, 12, 23, 34},
			{2, 34},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})

	t.Run("matrix size: one", func(t *testing.T) {
		e, err := check([][]float64{
			{4},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("initialize by zeros", func(t *testing.T) {
		old := initialize
		initialize = func(x []float64) {
			for i := range x {
				x[i] = 0.0
			}
		}
		defer func() {
			initialize = old
		}()
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})
	t.Run("initialize by eigenvector1", func(t *testing.T) {
		old := initialize
		initialize = func(x []float64) {
			x[0] = 1.0
			x[1] = 0.3333333333333333
		}
		defer func() {
			initialize = old
		}()
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("initialize by eigenvector2", func(t *testing.T) {
		old := initialize
		initialize = func(x []float64) {
			x[0] = 1.00
			x[1] = 0.25
		}
		defer func() {
			initialize = old
		}()
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if math.Abs(e.𝜦+2) > 1e-4 {
			t.Fatalf("result is not correct: %.14e ---> prec = %.14e", e.𝜦, e.𝜦+2)
		}
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("initialize specific : 1", func(t *testing.T) {
		old := initialize
		initialize = func(x []float64) {
			x[0] = 5.0
			x[1] = 2.0
		}
		defer func() {
			initialize = old
		}()
		e, err := check([][]float64{
			{4, -5},
			{2, -3},
		})
		if math.Abs(e.𝜦-2) > 1e-4 {
			t.Fatalf("result is not correct: %.14e ---> prec = %.14e", e.𝜦, e.𝜦+2)
		}
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("matrix with zeros", func(t *testing.T) {
		e, err := check([][]float64{
			{0.0, 0.0},
			{0.0, 0.0},
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
		_ = e
	})

	t.Run("lower triangle matrix", func(t *testing.T) {
		e, err := check([][]float64{
			{2, 1},
			{0, -4},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("upper triangle matrix", func(t *testing.T) {
		e, err := check([][]float64{
			{2, 3, 1},
			{0, -1, 2},
			{0, 0, 3},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
}

func ExampleInitByEigenvector1and2() {
	old := initialize
	defer func() {
		initialize = old
	}()
	oldO := output
	defer func() {
		output = oldO
	}()
	// eigenvector 1 : [1 0.333333]
	// eigenvector 2 : [1 0.25]
	n := 50000
	output = false
	for i := int(n * 9999.0 / 10000.0); i < n; i++ {
		value := float64(i) / float64(n-1)
		initialize = func(x []float64) {
			x[0] = 1.00
			x[1] = 0.25*value + 0.33333333333333333*(1.0-value)
		}
		e, err := check([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("ratio: %8.7f x: [%3.2f %8.7f]. Result: 𝜦=%6.4f 𝑿=[%6.4f %6.4f]\n",
			value, 1.0, 0.25*value+0.33333333333333333*(1.0-value),
			e.𝜦, e.𝑿[0], e.𝑿[1])
	}

	// Output:
	// ratio: 0.9999200 x: [1.00 0.2500067]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999400 x: [1.00 0.2500050]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999600 x: [1.00 0.2500033]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999800 x: [1.00 0.2500017]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 1.0000000 x: [1.00 0.2500000]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
}
