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

func random(x []float64) {
	// инициализация произвольным вектором
	rand.Seed(time.Now().UnixNano())
	for i := range x {
		x[i] = rand.Float64() // [0.0, 1)
	}
	// проверка чтобы все элементы не нулевые
	for i := range x {
		if x[i] != 0.0 {
			return
		}
	}
	random(x)
}

var initialize func([]float64) = random

// выводить на экран
var output bool = true

func pm(A [][]float64) (e eigen, err error) {
	n := len(A)
	if n == 0 {
		err = fmt.Errorf("matrix size is zero")
		return
	}

	// проверка на квадратность входной матрицы
	for row := 0; row < len(A); row++ {
		if len(A[row]) != n {
			err = fmt.Errorf("input matrix is not square in row %d: [%d,%d]", row, n, len(A[row]))
			return
		}
	}

	// для случая матрица 1х1
	if n == 1 {
		e.𝑿 = []float64{1.0}
		e.𝜦 = A[0][0]
		return
	}

	x := make([]float64, n)
	xLast := make([]float64, n)

	// инициализация произвольным вектором
	initialize(x)

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
			if max == 0.0 {
				err = fmt.Errorf("all values of eigenvector is zeros")
				return
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

			// на случай слишком быстрой сходимости,
			// добавим возмущения
			if iter < 3 {
				// добавляем возмужение
				perturbation := 0.002
				for i := range x {
					// x[i] = [-1.0,...,1.0]
					factor := math.Abs(x[i])
					if factor > 0.5 {
						factor = 1.0 - factor
					}
					// factor graph
					// x[i]    tmp
					// -1.0    0.0
					// -0.75   0.25
					// -0.5    0.5
					// -0.25   0.25
					//  0.0    0.0
					//  0.25   0.25
					//  0.5    0.5
					//  0.75   0.25
					//  1.0    0.0
					x[i] += perturbation * factor * factor
				}
				continue
			}

			// проверка резульатата
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
		if output {
			fmt.Printf("iter: %2d\tx = %v\n", iter, x)
		}

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

	if output {
		fmt.Println("e = ", e)
	}

	return
}
