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
		e.𝑿 = []float64{1.0}
		e.𝜦 = A[0][0]
		return
	}

	var (
		x     = make([]float64, n)
		xLast = make([]float64, n)
	)

	// инициализация произвольным вектором
	initialize(x)

	// add random seed
	rand.Seed(time.Now().UnixNano())

	// переменные для организации итераций
	var maxIteration int64 = 500
	var iter int64 = 0

	for {

		// устанавливаем лимит на количество итераций
		iter++
		if iter > maxIteration {
			err = fmt.Errorf("Iteration limit")
			return
		}

		// z(k) = A · x(k-1)
		z := make([]float64, n)
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				z[row] += A[row][col] * x[col]
			}
		}

		// x(k) = z(k) / || z(k) ||
		err = oneMax(x, z)
		if err != nil {
			return
		}

		// проверка на парность
		if iter%3 == 0 {
			lambda := λ(A, x)
			for i := range x {
				x[i] = x[i] + lambda*xLast[i]
			}
			err = oneMax(x, x)
			if err != nil {
				return
			}
			continue
		}

		// отображаем результат каждой итерации
		if output {
			fmt.Printf("iter: %2d\tx = %v\n", iter, x)
		}

		// ||x(k-1)-x(k-2)|| > 𝛆
		if iter > 0 {
			if eMax(x, xLast) < 𝛆 {
				// на случай слишком быстрой сходимости,
				// добавим возмущения
				if iter < 3 {
					// добавляем возмужение
					perturbation := 0.02 * (1 + rand.Float64())
					offset := 0.005
					for i := range x {
						// x[i] = [-1.0,...,1.0]
						factor := math.Abs(x[i])
						if factor > 0.5 {
							factor = 1.0 - factor
						}
						// factor graph
						// x[i]    : -1.0  -0.75  -0.5  -0.25  0.0  0.25  0.5  0.75  1.0
						// factor  :  0.0   0.25   0.5   0.25  0.0  0.25  0.5  0.25  0.0
						x[i] += perturbation*factor*factor + offset*float64(i)/float64(n)
					}
					continue
				}

				// проверка результата
				// выходим из итераций
				break
			}
		}

		copy(xLast, x)
	}

	e.𝑿 = x
	e.𝜦 = λ(A, x)

	if output {
		fmt.Println("e = ", e)
	}

	return
}

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
func oneMax(x, z []float64) (err error) {
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
	return
}

// ||x(k-1)-x(k-2)|| > 𝛆
func eMax(x, xLast []float64) (eMax float64) {
	for i := range x {
		e := math.Abs(x[i] - xLast[i])
		if e > eMax {
			eMax = e
		}
	}
	return
}
