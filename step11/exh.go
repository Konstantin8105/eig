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

	// переменные для организации итераций
	var maxIteration int64 = 5000
	var iter int64 = 0

	get := func(x []float64, trans bool) (err error) {
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
					if trans {
						z[row] += A[col][row] * x[col]
						continue
					}
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
				if math.Abs((max-maxLast)/max) < 𝛆 { // eMax(x, xLast) < 𝛆
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
		return
	}

	for value := 0; value < n; value++ {

		// инициализация произвольным вектором
		u := make([]float64, n)
		initialize(u)
		err = get(u, false)
		if err != nil {
			return
		}

		l := λ(A, u)
		e = append(e, eigen{𝑿: u, 𝜦: l})

		// инициализация произвольным вектором
		v := make([]float64, n)
		initialize(v)
		err = get(v, true)
		if err != nil {
			return
		}

		// нормализация
		_, err = oneMax(u, u)
		if err != nil {
			return
		}
		_, err = oneMax(v, v)
		if err != nil {
			return
		}
		var pro float64
		for i := range u {
			pro += u[i] * v[i]
		}
		for i := range u {
			v[i] /= pro
		}

		// проверка V'*U = 1
		{
			res := 0.0
			for i := range u {
				res += v[i] * u[i]
			}
			if math.Abs(res) > 1+1e-1 || math.Abs(res) < 1-1e-1 {
				err = fmt.Errorf("check is not ok. V'*U = %.14e != 1\nu = %v\nv = %v",
					res, u, v)
				return
			}
		}

		// метод исчерпывания
		Atmp := make([][]float64, n)
		for i := 0; i < n; i++ {
			Atmp[i] = make([]float64, n)
		}

		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				Atmp[row][col] = A[row][col] - l*u[row]*v[col]
			}
		}

		A = Atmp
	}

	for i := range e {
		oneMax(e[i].𝑿, e[i].𝑿)
		if i == 0 {
			continue
		}
		if math.Abs(e[i-1].𝜦)+𝛆 < math.Abs(e[i].𝜦) {
			err = fmt.Errorf("eigen values is not less. %.14e !> %.14e",
				math.Abs(e[i-1].𝜦), math.Abs(e[i].𝜦))
		}
	}

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
var output bool = false // true

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
