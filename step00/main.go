package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// результаты
type eigen struct {
	// собственные значения
	𝜦 float64

	// собственный вектор
	𝑿 []float64
}

func pm(A [][]float64) (e eigen, err error) {
	n := len(A)
	x := make([]float64, n)
	xLast := make([]float64, n)

	rand.Seed(time.Now().UnixNano())

	// инициализация произвольным вектором
	for i := range x {
		x[i] = rand.Float64() // [0.0, 1)
	}

	// точность результата
	𝛆 := 0.001

	for iter := 0; ; iter++ {
		// z(k) = A · x(k-1)
		z := make([]float64, n)
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				z[row] += A[row][col] * x[col]
			}
		}

		// x(k) = z(k) / || z(k) ||
		{
			max := z[0]
			for i := range z {
				if math.Abs(z[i]) > math.Abs(max) {
					max = z[i]
				}
			}
			for i := range x {
				x[i] = z[i] / max
			}
		}

		// ||x(k-1)-x(k-2)|| > 𝛆
		if iter > 0 {
			var max float64
			for i := range x {
				e := math.Abs(x[i] - xLast[i])
				if e > max {
					max = e
				}
			}
			if max < 𝛆 {
				// выходим из итераций
				break
			}
		}

		// устанавливаем лимит на количество итераций
		if iter > 100 {
			err = fmt.Errorf("Iteration limit")
			return
		}

		// отображаем результат каждой итерации
		fmt.Printf("iter: %2d\tx = %v\n", iter, x)
		copy(xLast, x)
	}

	// λ = (Ax , x) / (x , x)
	e.𝑿 = x

	Ax := make([]float64, n)
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			Ax[row] += A[row][col] * x[col]
		}
	}
	var Axx float64
	for i := range x {
		Axx += Ax[i] * x[i]
	}
	var xx float64
	for i := range x {
		xx += x[i] * x[i]
	}

	e.𝜦 = Axx / xx

	fmt.Println("e = ", e)

	return
}
