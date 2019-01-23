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
			if output && iter > 0 {
				fmt.Printf("iter: %2d\tx=", iter)
				for i := range x {
					fmt.Printf("\t%10.5e", x[i])
				}
				fmt.Printf("\t𝛆 = %10.5e\n", math.Abs((max-maxLast)/max))
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

	for value := 0; value < n; {
		if output {
			fmt.Println("Input A. value = ", value)
			MatrixPrint(A)
		}

		// инициализация произвольным вектором
		u := make([]float64, n)
		initialize(u)
		err = get(u, false)
		if err != nil {
			return
		}

		l := λ(A, u)

		value += Gauss(A, l)

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
var 𝛆 float64 = 1e-15

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

func Gauss(A [][]float64, l float64) int {
	n := len(A)
	U := make([][]float64, n)
	for i := 0; i < n; i++ {
		U[i] = make([]float64, n)
	}

	// copying
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			U[i][j] = A[i][j]
		}
		U[i][i] -= l
	}

	// gauss decomposition to U triangle matrix
	for k := 0; k < len(A)-1; k++ {
		for i := k + 1; i < len(A); i++ {
			factor := (U[i][k] / U[k][k])
			if U[k][k] == 0.0 {
				break
			}
			for col := k; col < len(U[i]); col++ {
				U[i][col] -= U[k][col] * factor
			}
		}
	}

	if output {
		fmt.Println("++++++++++++++")
		fmt.Println("l = ", l)
		MatrixPrint(A)
		MatrixPrint(U)
	}

	var amount int
	for i := 0; i < n; i++ {
		var sum float64
		var sumC float64
		for j := 0; j < n; j++ {
			sum += math.Abs(U[i][j])
			sumC += math.Abs(U[j][i])
		}

		if output {
			fmt.Printf("row = %d\tsum = %.14e\n", i, sum)
			fmt.Printf("col = %d\tsum = %.14e\n", i, sum)
		}

		if math.Abs(sum) < 𝛆 {
			amount++
		}
		if math.Abs(sumC) < 𝛆 {
			amount++
		}
	}

	if output {
		fmt.Println("Count = ", amount)
		fmt.Println("++++++++++++++")
	}

	fmt.Println("// TODO : ADD IMPLEMENTATION FOR КРАТНЫЕ СОБСТВЕННЫЕ ЗНАЧЕНИЯ")

	return amount
}
