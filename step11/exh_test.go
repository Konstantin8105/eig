package main

import (
	"fmt"
	"testing"
)

func ExampleSimple() {

	// output true
	oldOut := output
	output = false
	defer func() {
		output = oldOut
	}()

	es := []eigen{
		{𝜦: +2.0, 𝑿: []float64{+0.5714286, +0.1428572, +1.0000000}},
		{𝜦: -5.0, 𝑿: []float64{-0.6666667, -1.0000000, -1.0000000}},
		{𝜦: -1.0, 𝑿: []float64{+0.5773503, +0.5773503, +0.5773503}},
	}

	A := generator(es)
	MatrixPrint(A)
	output = true
	fmt.Println("Before")
	PrintEigens(es)
	output = false

	// initialize
	old := initialize
	initialize = func(x []float64) {
		for i := range x {
			x[i] = 1.0 + float64(i)
		}
	}
	defer func() {
		initialize = old
	}()

	// calculate
	e, err := exh(A)
	if err != nil {
		panic(err)
	}
	output = true
	fmt.Println("After")
	PrintEigens(e)
	output = false

	// compare
	output = true
	fmt.Println(compare(e, es))
	output = false

	// Output:
	// Before
	// ---     0 ---
	// 𝜦      = +2.0000000000e+00
	// 𝑿[  0] = +5.7142860000e-01
	// 𝑿[  1] = +1.4285720000e-01
	// 𝑿[  2] = +1.0000000000e+00
	// ---     1 ---
	// 𝜦      = -5.0000000000e+00
	// 𝑿[  0] = -6.6666670000e-01
	// 𝑿[  1] = -1.0000000000e+00
	// 𝑿[  2] = -1.0000000000e+00
	// ---     2 ---
	// 𝜦      = -1.0000000000e+00
	// 𝑿[  0] = +5.7735030000e-01
	// 𝑿[  1] = +5.7735030000e-01
	// 𝑿[  2] = +5.7735030000e-01
	// After
	// ---     0 ---
	// 𝜦      = -4.9999995141e+00
	// 𝑿[  0] = +6.6666668939e-01
	// 𝑿[  1] = +9.9999990456e-01
	// 𝑿[  2] = +1.0000000000e+00
	// ---     1 ---
	// 𝜦      = +1.9999997398e+00
	// 𝑿[  0] = +5.7142183122e-01
	// 𝑿[  1] = +1.4279607663e-01
	// 𝑿[  2] = +1.0000000000e+00
	// ---     2 ---
	// 𝜦      = -1.0000000000e+00
	// 𝑿[  0] = +9.9999999640e-01
	// 𝑿[  1] = +1.0000000000e+00
	// 𝑿[  2] = +9.9999992052e-01
	// Compare [  0,  1] is same
	// Compare [  1,  0] is same
	// Compare [  2,  2] is same
	// true
}

func TestSnippets(t *testing.T) {

	// output true
	oldOut := output
	output = true
	defer func() {
		output = oldOut
	}()

	t.Run("Fadeev: page 334", func(t *testing.T) {
		e, err := exh([][]float64{
			{1.022551, 0.116069, -0.287028, -0.429969},
			{0.228401, 0.742521, -0.176368, -0.283720},
			{0.326141, 0.097221, 0.197209, -0.216487},
			{0.433864, 0.148965, -0.193686, 0.006472},
		})
		PrintEigens(e)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("result is not checked")

		// [R,T]= spec([1.022551  0.116069  -0.287028  -0.429969 ; 0.228401  0.742521  -0.176368  -0.283720 ; 0.326141  0.097221  0.197209  -0.216487 ; 0.433864  0.148965  -0.193686  0.006472])
		//  T  =
		//
		//     0.2876392    0            0                         0
		//     0            0.3461482    0                         0
		//     0            0            0.6674828 + 8.991D-08i    0
		//     0            0            0                         0.6674828 - 8.991D-08i
		//  R  =
		//
		//   - 0.4803844    0.4256283  - 0.7074026               - 0.7074026
		//   - 0.3202569    0.2553766  - 0.3795747 + 0.1573849i  - 0.3795747 - 0.1573849i
		//   - 0.1601288    0.8512565  - 0.3652032 + 0.0187365i  - 0.3652032 - 0.0187365i
		//   - 0.8006405    0.1702519  - 0.4428447 + 0.0299782i  - 0.4428447 - 0.0299782i
	})
	t.Run("Fadeev: page 347", func(t *testing.T) {
		e, err := exh([][]float64{
			{1.00, 0.0, 1.00, 0.0},
			{1.00, 0.77777777777, 0.333333333333333, 0.3333333333333},
			{0.0, -0.02525252525, 0.555555555555555, -0.025252525252},
			{0.0, -0.88888888888, -8.64444444444444, 0.1111111111111},
		})
		PrintEigens(e)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("result is not checked")
		// [P,O] = spec([1.00, 0.0, 1.00, 0.0 ;1.00, 0.77777777777, 0.333333333333333, 0.3333333333333 ;0.0, -0.02525252525, 0.555555555555555, -0.025252525252 ;0.0, -0.88888888888, -8.64444444444444, 0.1111111111111])
		// O  =
		//
		//     1.    0            0            0
		//     0     0.3333333    0            0
		//     0     0            0.4444444    0
		//     0     0            0            0.6666667
		//  P  =
		//
		//   - 0.3656362  - 0.0509677  - 0.0618765    0.1270457
		//   - 0.6581452  - 0.5402572  - 0.6256402    0.7876832
		//   - 4.600D-12    0.0339784    0.0343758  - 0.0423486
		//     0.6581452    0.8392675    0.7768938  - 0.6013495
	})
	t.Run("D", func(t *testing.T) {
		// initialize
		old := initialize
		initialize = func(x []float64) {
			for i := range x {
				x[i] = 1.0
			}
		}
		defer func() {
			initialize = old
		}()

		eig := []eigen{
			{𝜦: 18.0, 𝑿: []float64{0.8944272, -0.4472136, 0.00000}},
			{𝜦: 18.0, 𝑿: []float64{-0.2981424, -0.5962848, 0.7453560}},
			{𝜦: +9.0, 𝑿: []float64{0.5, 1.0, 1.0}},
		}

		output = true

		// find a1,a2,a3
		{
			// y = c1*x1 + c2*x2 + c3*x3 + ...
			A := generator(eig)
			Y0 := []float64{1, 1, 1}
			AY0 := make([]float64, 3)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					AY0[i] += A[i][j] * Y0[j]
				}
			}
			CN := make([]float64, 3)
			L1 := 18.0
			for i := 0; i < 3; i++ {
				CN[i] = Y0[i] - AY0[i]/L1
			}
			oneMax(CN, CN)
			fmt.Println("CN = ", CN)
		}

		//
		{
			A := generator(eig)
			x := []float64{1, 1, 1}
			fmt.Println("A = ", A)
			for i := 0; i < 3; i++ {
				A[i][i] -= 18.0
			}
			fmt.Println("A = ", A)
			for iter := 0; iter < 40; iter++ {
				fmt.Println("iter ", iter, " -> x = ", x)
				for i := 0; i < 3; i++ {
					x[i] = 0.0
					for k := 0; k < 3; k++ {
						if i == k {
							continue
						}
						x[i] += A[i][k] * x[k]
					}
					x[i] /= -A[i][i]
				}
			}
		}

		// compare with another vectors
		MatrixPrint(generator(eig))
		MatrixPrint(generator([]eigen{
			{𝜦: 18.0, 𝑿: []float64{1.0, 0.0, -0.5}},
			{𝜦: 18.0, 𝑿: []float64{0.0, 1.0, -1.0}},
			{𝜦: +9.0, 𝑿: []float64{0.5, 1.0, 1.0}},
		}))
		MatrixPrint(generator([]eigen{
			{𝜦: 18.0, 𝑿: []float64{-1.0, 0.5, 0.0}},
			{𝜦: 18.0, 𝑿: []float64{-1.0, 0.0, 0.5}},
			{𝜦: +9.0, 𝑿: []float64{0.5, 1.0, 1.0}},
		}))
		MatrixPrint(generator([]eigen{
			{𝜦: 18.0, 𝑿: []float64{1.0, -0.5, 0.0}},
			{𝜦: 18.0, 𝑿: []float64{0.0, -1.0, 1.0}},
			{𝜦: +9.0, 𝑿: []float64{0.5, 1.0, 1.0}},
		}))

		e, err := exh([][]float64{
			{17, -2, -2},
			{-2, 14, -4},
			{-2, -4, 14},
		})
		PrintEigens(e)
		if err != nil {
			t.Fatal(err)
		}

		if !compare(e, eig) {
			t.Errorf("not same")
		}
		//
		// A=[17 -2 -2; -2 14 -4; -2 -4 14];[S,P]=spec(A)
		//  P  =
		//
		//     9.    0.     0.
		//     0.    18.    0.
		//     0.    0.     18.
		//  S  =
		//
		//     0.3333333  - 0.2981424    0.8944272
		//     0.6666667  - 0.5962848  - 0.4472136
		//     0.6666667    0.7453560    0.
		//
	})
}

// func init() {
// 	output = flag.CommandLine.Lookup("test.v").Value.String() == "true"
// }

func Test(t *testing.T) {
	// initialize
	old := initialize
	initialize = func(x []float64) {
		for i := range x {
			x[i] = 1.0 //+ float64(i)
		}
	}
	defer func() {
		initialize = old
	}()

	tcs := []struct {
		es   []eigen
		name string
	}{
		{
			name: "simple",
			es: []eigen{
				{𝜦: +2.0, 𝑿: []float64{+0.5714286, +0.1428572, +1.0000000}},
				{𝜦: -5.0, 𝑿: []float64{-0.6666667, -1.0000000, -1.0000000}},
				{𝜦: -1.0, 𝑿: []float64{+0.5773503, +0.5773503, +0.5773503}},
			},
		},
		{
			name: "Матрица 2х2 с одним собственным значением",
			es: []eigen{
				{𝜦: -6.0, 𝑿: []float64{1.0, 0.0}},
				{𝜦: -6.0, 𝑿: []float64{0.0, 1.0}},
			},
		},
		{
			name: "Доминанирование l1 и l2 == l3. Собственные вектора разные",
			es: []eigen{
				{𝜦: -1.0, 𝑿: []float64{0.6, 1.0, 1.0}},
				{𝜦: -5.0, 𝑿: []float64{0.5, 0.2, 1.0}},
				{𝜦: -1.0, 𝑿: []float64{1.0, 1.0, 1.0}},
			},
		},
		{
			name: "Большие числа",
			es: []eigen{
				{𝜦: 1e01, 𝑿: []float64{0.0, 0.0, 0.0, 1.0}},
				{𝜦: 1e04, 𝑿: []float64{0.0, 0.0, 1.0, 0.0}},
				{𝜦: 1e08, 𝑿: []float64{0.0, 1.0, 0.0, 0.0}},
				{𝜦: 1e12, 𝑿: []float64{1.0, 0.0, 0.0, 0.0}},
			},
		},
		{
			name: "Нет доминантной l1 = l2 > 0. Собственные вектора разные",
			es: []eigen{
				{𝜦: +5.0, 𝑿: []float64{0.5, 0.2, 1.0}},
				{𝜦: +5.0, 𝑿: []float64{0.6, 1.0, 1.0}},
				{𝜦: -1.0, 𝑿: []float64{1.0, 1.0, 1.0}},
			},
			//
			// -->[U,I] = spec([-10 1.875 7.125; -15 6.875 7.125; -15 1.875 12.125])
			//  I  =
			//
			//   - 1.    0     0
			//     0     5.    0
			//     0     0     5.
			//  U  =
			//
			//   - 0.5773503    0.3905667  - 0.0576896
			//   - 0.5773503    0.6509446  - 0.9886500
			//   - 0.5773503    0.6509446    0.1387193
			//
		},
		{
			name: "Нет доминантной l1 = l2 < 0. Собственные вектора разные",
			es: []eigen{
				{𝜦: -5.0, 𝑿: []float64{0.5, 0.2, 1.0}},
				{𝜦: -5.0, 𝑿: []float64{0.6, 1.0, 1.0}},
				{𝜦: -1.0, 𝑿: []float64{1.0, 1.0, 1.0}},
			},
			//
			// -->[P,O]=spec([5 -1.25 -4.75; 10 -6.25 -4.75; 10 -1.25 -9.75])
			//  O  =
			//
			//   - 1.    0     0
			//     0   - 5.    0
			//     0     0   - 5.
			//  P  =
			//
			//     0.5773503  - 0.3905667    0.2640036
			//     0.5773503  - 0.6509446  - 0.6377152
			//     0.5773503  - 0.6509446    0.7236169
			//
		},
		{
			name: "Нет доминантной l1 = - l2. Собственные вектора разные",
			es: []eigen{
				{𝜦: +5.0, 𝑿: []float64{1.0, 0.4, 0.0}},
				{𝜦: -5.0, 𝑿: []float64{0.0, 0.2, 1.0}},
				{𝜦: -1.0, 𝑿: []float64{1.0, 1.0, 1.0}},
			},
			//
			// -->[G,J]=spec([11 -15 3; 4 -5 0; -4 10 -7])
			// J  =
			//
			//          5.    0     0
			//          0   - 1.    0
			//          0     0   - 5.
			// G  =
			//
			//   0.9284767  - 0.5773503    5.459D-16
			//   0.3713907  - 0.5773503    0.1961161
			// - 4.538D-17  - 0.5773503    0.9805807
			//
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if output {
				fmt.Println(tc.name)
			}

			// generate
			A := generator(tc.es)
			MatrixPrint(A)
			PrintEigens(tc.es)

			// calculate
			e, err := exh(A)
			if err != nil {
				t.Fatal(err)
			}
			PrintEigens(e)

			// compare
			if !compare(e, tc.es) {
				t.Errorf("not same")
			}
		})
	}

	// Output:
}

func TestGauss(t *testing.T) {

	// output true
	oldOut := output
	output = true
	defer func() {
		output = oldOut
	}()

	// A := [][]float64{
	// {-1, -3, -1, 2, 2},
	// {-1, -7, -2, 1, 2},
	// {-1, -11, -3, 0, 2},
	// {-2, -2, -1, 5, 4},
	// }

	// A := [][]float64{
	// {-1, -2, -2},
	// {-2, -4, -4},
	// {-2, -4, -4},
	// }

	A := [][]float64{
		{+8, -2, -2},
		{-2, +5, -4},
		{-2, -4, +5},
	}

	MatrixPrint(A)
	fmt.Println("--------")

	for k := 0; k < len(A)-1; k++ {
		for i := k + 1; i < len(A); i++ {
			factor := (A[i][k] / A[k][k])
			if A[k][k] == 0.0 {
				fmt.Println("break")
				break
			}
			fmt.Println("factor = ", factor)
			for col := k; col < len(A[i]); col++ {
				A[i][col] -= A[k][col] * factor
			}
		}
		MatrixPrint(A)
		fmt.Println("--------")
	}
}
