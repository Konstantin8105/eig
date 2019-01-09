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

// точность результата
var 𝛆 float64 = 1e-6

func pm(A [][]float64) (e eigen, err error) {

	defer func() {
		fmt.Println("e = ", e)
	}()
	n := len(A)

	// масштабирование
	scale := 250.0 // пока это произвольное число
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			A[row][col] *= scale
		}
	}
	defer func() {
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				A[row][col] /= scale
			}
		}
		if err == nil {
			e.𝜦 /= scale
		}
	}()

	x := make([]float64, n)
	xLast := make([]float64, n)

	rand.Seed(time.Now().UnixNano())

	// инициализация произвольным вектором
	for i := range x {
		x[i] = rand.Float64() // [0.0, 1)
	}

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
		if iter > 500 {
			err = fmt.Errorf("Iteration limit")
			return
		}

		// отображаем результат каждой итерации
		fmt.Printf("iter: %2d\tx = %v\n", iter, x)

		// обработка повторения значений х, но с изменением знака кроме 1.0
		isSame := true
		for i := range x {
			if x[i] == 1.0 {
				continue
			}
			if x[i] != -xLast[i] {
				isSame = false
				break
			}
		}
		if isSame {
			err = fmt.Errorf("Loop values x")
			return
		}

		// значение х не изменяется кроме 1.0
		isSame = false
		for i := range x {
			if x[i] == 1.0 && xLast[i] == 1.0 {
				continue
			}
			if x[i] == xLast[i] {
				isSame = true
				break
			}
		}
		if isSame {
			err = fmt.Errorf("one or more values of eigenvector is not change")
			return
		}

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

	return
}
