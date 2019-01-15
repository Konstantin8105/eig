package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func exh(A [][]float64) (e []eigen, err error) {
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

	// матрица А не должна состоять из одних нулей
	{
		isAllZeros := true
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				if A[row][col] != 0.0 {
					isAllZeros = false
					break
				}
			}
		}
		if isAllZeros {
			err = fmt.Errorf("all elements of matrix is zeros")
			return
		}
	}

	// для случая матрица 1х1
	if n == 1 {
		e = []eigen{
			{
				𝑿: []float64{1.0},
				𝜦: A[0][0],
			},
		}
		return
	}

	// add random seed
	rand.Seed(time.Now().UnixNano())

	// инициализация произвольным вектором
	x := make([]float64, n)
	initialize(x)

	// переменные для организации итераций
	var maxIteration int64 = 500
	var iter int64 = 0

	for max, maxLast, z := 0.0, 0.0, make([]float64, n); ; {
		// устанавливаем лимит на количество итераций
		iter++
		if iter > maxIteration {
			err = fmt.Errorf("Iteration limit")
			return
		}

		// z(k) = A · x(k-1)
		for row := 0; row < n; row++ {
			z[row] = 0.0
		}
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				z[row] += A[row][col] * x[col]
			}
		}

		// x(k) = z(k) / || z(k) ||
		max, err = oneMax(x, z)
		if err != nil {
			return
		}

		// отображаем результат каждой итерации
		if output {
			fmt.Printf("iter: %2d\tx=", iter)
			for i := range x {
				fmt.Printf("\t%10.5e", x[i])
			}
			fmt.Printf("\n")
		}

		// ||x(k-1)-x(k-2)|| > 𝛆
		if iter > 0 {
			if math.Abs((max-maxLast)/max) < 𝛆 { // eMax(x, xLast) < 𝛆 {
				if iter < 3 {
					// на случай слишком быстрой сходимости
					random(x)
					continue
				}

				// проверка результата, выходим из итераций
				break
			}
		}

		maxLast, max = max, maxLast
	}

	e = append(e, eigen{𝑿: x, 𝜦: λ(A, x)})

	return
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

// λ = (Ax , x) / (x , x)
func λ(A [][]float64, x []float64) float64 {
	n := len(A)
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
	return Axx / xx
}

// x(k) = z(k) / || z(k) ||
func oneMax(x, z []float64) (max float64, err error) {
	max = z[0]
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
	return
}

// ||x(k-1)-x(k-2)|| > 𝛆
// func eMax(x, xLast []float64) (eMax float64) {
// 	// for i := range x {
// 	// 	e := math.Abs(x[i] - xLast[i])
// 	// 	if e > eMax {
// 	// 		eMax = e
// 	// 	}
// 	// }
//
// 	// for i := range x {
// 	// 	eMax += math.Pow(x[i]-xLast[i], 2.0)
// 	// }
// 	// eMax = math.Sqrt(eMax)
//
// 	for i := range x {
// 		if math.Abs(xLast[i]) < 𝛆 && math.Abs(x[i]) < 𝛆 {
// 			continue
// 		}
// 		if xLast[i] != 0.0 {
// 			eMax += math.Abs((x[i] - xLast[i]) / xLast[i])
// 			continue
// 		}
// 		if x[i] != 0.0 {
// 			eMax += math.Abs((x[i] - xLast[i]) / x[i])
// 			continue
// 		}
// 	}
// 	return
// }
