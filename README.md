# eig

Построение алгоритма для расчета собственных значений и собственного вектора 
итерационном методов.

Причина: при расчете частот свободных колебаний или устойчивости 
методом конечных элементов требуется иметь лишь несколько собсвенных 
значений что и обеспечивается итерационными методами, а не всех 
как при расчете прямыми методами.

Общий вид уравнения:
```
A · x = λ · B · x, 
где:
	A,B - симметричные матрицы
	x   - собственный вектор
	λ   - собственное значение
```

Алгоритм буду писать на языке программирования `Golang` в связи с тем,
что язык уже знаком и обладает необходимым инструментарием для 
исследования.

## step00 : Степенной метод(power method)

Начнём с упрощенного уравнения
```
A · x = λ · x, 
где:
	A   - матрица
	x   - собственный вектор
	λ   - собственное значение
```

Начнем с рассмотрения степенного метода, как с наиболее простого с точки
зрения псевдокода, который представлен ниже:

```
Выбираем произвольный вектор x(0)
for k = 1,2,...		(while )
	z(k) = A · x(k-1)
	x(k) = z(k) / || z(k) ||
	if ||q(k-1)-q(k-2)|| < 𝛆 then break
end

λ = (Ax , x) / (x , x)

```

Используемая литература:
* http://ergodic.ugr.es/cphys/LECCIONES/FORTRAN/power_method.pdf
* http://www.cs.huji.ac.il/~csip/tirgul2.pdf

Сразу стоит определить ограничения данного метода, который выражается в 
` | λ1 | > | λi | `, то есть первое собственное значение должно доминировать, 
что не всегда выполнимо, к примеру:

```
A = 
[1  0]
[0 -1]

λ1 = 1
λ2 = -1
```

```
A = 
[ 2 0 0 ]
[ 0 2 0 ]
[ 0 0 1 ]

λ1 = 2
λ2 = 2
λ3 = 1
```

Но зная это надо обязательно посмотреть как ведёт себя алгоритм.

Начнем программирование нашего алгоритма с заранее заданных ограничений
задаваемых на текущий этап, а именно:
* алгоритм должен быть работающим
* в качестве входных данных будут использоваться матрицы 3х3 и 2х2
* тест не должен занимать более 1 секунды
* вопросы об оптимизации производительности, распараллеливании игнорируются
* каждый этап со своим характерным кодом храниться отдельно от других - разбито на папки
* расчет производиться только одного собственного значений

Рассмотрим исходный код `step00`:
```golang
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
```

Тестировать будет следующим образом:

```golang
package main

import "testing"

func Test(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		e, err := pm([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("2", func(t *testing.T) {
		e, err := pm([][]float64{
			{4, 5},
			{6, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
}

```

Посмотрим результат:

```
=== RUN   Test/1
iter:  0	x = [1 0.4039100393253347]
iter:  1	x = [1 0.3581238769008534]
iter:  2	x = [1 0.3441236223389014]
iter:  3	x = [1 0.3384004255025737]
iter:  4	x = [1 0.33579212584816726]
iter:  5	x = [1 0.3345448562892963]
e =  {-2.00703469759771 [1 0.3339347232253568]}
```
Полученное значение собственного значение `-2.00703469759771` близко к 
точному значению `-2`.

```
=== RUN   Test/2
iter:  0	x = [0.8391201728306759 1]
iter:  1	x = [0.8327566516820931 1]
e =  {10.000034031736032 [0.8333910214590675 1]}
--- PASS: Test (0.00s)
    --- PASS: Test/1 (0.00s)
    --- PASS: Test/2 (0.00s)
PASS
ok  	github.com/Konstantin8105/eig/step00	0.003s
```
Полученное значение собственного значение `10.000034031736032` близко к 
точному значению `10`.


Давайте проверим что будет происходить если в качестве матрицы взять, матрицу 
не имеющией доминантного собственного значения, которые были описаны ранее.

```go
	t.Run("E3", func(t *testing.T) {
		e, err := pm([][]float64{
			{1, 0},
			{0, -1},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("E4", func(t *testing.T) {
		e, err := pm([][]float64{
			{2, 0, 0},
			{0, 2, 0},
			{0, 0, 1},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
```

```
=== RUN   Test/E3
iter:  0	x = [-0.2758362805320346 1]
iter:  1	x = [0.2758362805320346 1]
iter:  2	x = [-0.2758362805320346 1]
iter:  3	x = [0.2758362805320346 1]
iter:  4	x = [-0.2758362805320346 1]
iter:  5	x = [0.2758362805320346 1]
...
iter: 93	x = [0.2758362805320346 1]
iter: 94	x = [-0.2758362805320346 1]
iter: 95	x = [0.2758362805320346 1]
iter: 96	x = [-0.2758362805320346 1]
iter: 97	x = [0.2758362805320346 1]
iter: 98	x = [-0.2758362805320346 1]
iter: 99	x = [0.2758362805320346 1]
iter: 100	x = [-0.2758362805320346 1]
        main_test.go:33: Iteration limit
=== RUN   Test/E4
iter:  0	x = [0.38516596921127894 1 0.4163313768328755]
iter:  1	x = [0.38516596921127894 1 0.20816568841643776]
iter:  2	x = [0.38516596921127894 1 0.10408284420821888]
iter:  3	x = [0.38516596921127894 1 0.05204142210410944]
iter:  4	x = [0.38516596921127894 1 0.02602071105205472]
iter:  5	x = [0.38516596921127894 1 0.01301035552602736]
iter:  6	x = [0.38516596921127894 1 0.00650517776301368]
iter:  7	x = [0.38516596921127894 1 0.00325258888150684]
iter:  8	x = [0.38516596921127894 1 0.00162629444075342]
e =  {1.999999424211786 [0.38516596921127894 1 0.00081314722037671]}
```

Как мы видем из результатов, тест Е3 не сходиться, а тест Е4 сходиться, но
собственные вектор не корректен.

## step01 : автоматическая проверка

```golang
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
```

Увеличим точность:

```golang

// точность результата
var 𝛆 float64 = 1e-6

```

посмотрим результат тестов:

```
=== RUN   Test/1
iter:  0	x = [1 0.410232090428574]
iter:  1	x = [1 0.35964342964733637]
iter:  2	x = [1 0.34469484646217846]
...
iter: 13	x = [1 0.33333821562899224]
iter: 14	x = [1 0.33333577440965456]
iter: 15	x = [1 0.33333455385361765]
e =  {-2.0000071399897674 [1 0.33333394358900653]}
Delta = 5.7894349817e-07
=== RUN   Test/2
iter:  0	x = [0.869999217333778 1]
iter:  1	x = [0.8297456718022086 1]
iter:  2	x = [0.8336928734313528 1]
iter:  3	x = [0.8332973870780034 1]
iter:  4	x = [0.833336928036396 1]
iter:  5	x = [0.8333329738638025 1]
e =  {10.000000021214595 [0.8333333692802942 1]}
Delta = 3.0376742543e-07
=== RUN   Test/E3
iter:  0	x = [-0.8239081044647364 1]
iter:  1	x = [0.8239081044647364 1]
iter:  2	x = [-0.8239081044647364 1]
...
iter: 98	x = [-0.8239081044647364 1]
iter: 99	x = [0.8239081044647364 1]
iter: 100	x = [-0.8239081044647364 1]
=== RUN   Test/E4
iter:  0	x = [0.15070757797200232 0.32302703977622665 1]
iter:  1	x = [0.30141515594400464 0.6460540795524533 1]
iter:  2	x = [0.4665478718945736 1 0.7739290189861031]
...
iter: 19	x = [0.4665478718945736 1 5.9046098250282525e-06]
iter: 20	x = [0.4665478718945736 1 2.9523049125141263e-06]
iter: 21	x = [0.4665478718945736 1 1.4761524562570631e-06]
e =  {1.9999999999995528 [0.4665478718945736 1 7.380762281285316e-07]}
Delta = 7.3807622813e-07
```

Результат не изменился, тест Е3 не сходиться, а тест Е4 сходиться, но
собственные вектор не корректен.

Добавим обработку ошибок для теста Е3:
```golang
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
```

Добавим обработку ошибок для теста Е4:
```golang
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
```

Посмотрим на результаты:
```
=== RUN   Test
=== RUN   Test/1
iter:  0	x = [1 0.4043840695707527]
iter:  1	x = [1 0.3582406165785594]
iter:  2	x = [1 0.3441678279922198]
...
iter: 13	x = [1 0.33333801520286443]
iter: 14	x = [1 0.333335674202341]
iter: 15	x = [1 0.3333345037513984]
e =  {-2.0000068468961194 [1 0.3333339185382563]}
Delta = 5.5517793764e-07
=== RUN   Test/2
iter:  0	x = [0.9675894666177765 1]
iter:  1	x = [0.820908579637794 1]
iter:  2	x = [0.8345851407421243 1]
iter:  3	x = [0.8332082465431966 1]
iter:  4	x = [0.8333458429512196 1]
iter:  5	x = [0.833332082380934 1]
iter:  6	x = [0.8333334584286672 1]
e =  {9.999999992617324 [0.8333333208238008 1]}
Delta = 1.0571097932e-07
=== RUN   Test/E3
iter:  0	x = [-0.006308903007061117 1]
iter:  1	x = [0.006308903007061117 1]
=== RUN   Test/E4
iter:  0	x = [1 0.5858060710718106 0.16525381480819515]
iter:  1	x = [1 0.5858060710718106 0.08262690740409757]
--- PASS: Test (0.00s)
    --- PASS: Test/1 (0.00s)
    --- PASS: Test/2 (0.00s)
    --- PASS: Test/E3 (0.00s)
        main_test.go:72: Loop values x
    --- PASS: Test/E4 (0.00s)
        main_test.go:84: one or more values of eigenvector is not change
PASS
ok  	github.com/Konstantin8105/eig/step01	0.003s
```

## step02 : Тест с различными значениями |𝜦2|/|𝜦1|

Переименуем тесты на более подходящие и добавим тесты с разным соотношением
свободных значений.
```golang
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
```

В результате видно что чем больше соотношение |𝜦2|/|𝜦1|, тем больше необходимо 
итераций или по другому, чем дальше друг от друга |𝜦1| и |𝜦2| тем быстрее будет
сходиться.
```
=== RUN   Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1
iter:  0	x = [0.7845213075471055 1]
iter:  1	x = [0.8383618058706501 1]
iter:  2	x = [0.8328319986482193 1]
iter:  3	x = [0.8333834818865702 1]
iter:  4	x = [0.8333283186288978 1]
iter:  5	x = [0.8333338348052857 1]
e =  {9.999999970404932 [0.8333332831861533 1]}
Delta = 4.2376543762e-07
=== RUN   Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9
iter:  0	x = [0.26898382935618304 1]
iter:  1	x = [1 0.7712726263551744]
iter:  2	x = [0.34198618989492297 1]
...
iter: 130	x = [0.7142851144879747 1]
iter: 131	x = [0.7142862541039066 1]
iter: 132	x = [0.7142852284495247 1]
e =  {9.99999913140197 [0.7142861515384336 1]}
Delta = 6.7603359794e-06
    --- PASS: Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1 (0.00s)
    --- FAIL: Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9 (0.00s)
        main_test.go:104: Precition is not ok
```


## step03 : Тест с различными значениями |𝜦2|/|𝜦1| с масштабированием

Давайте рассмотрим масштабирование:
```
A · x = λ · x
B = A · n
B · x = λB · x
A · n · x = λ · n · x

где:
	n   - постоянная
```

Добавим код:
```golang
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
```

Сравнение результатов:

Было
```
=== RUN   Test
=== RUN   Test/example:_1
iter:  0	x = [1 0.41505506722929053]
...
iter: 15	x = [1 0.3333345925025327]
e =  {-2.0000073660824516 [1 0.3333339629131764]}
Delta = 5.9727628148e-07
=== RUN   Test/example:_2
iter:  0	x = [0.8677381595362804 1]
...
iter:  5	x = [0.833332995622703 1]
e =  {10.000000019930463 [0.8333333671044032 1]}
Delta = 2.8538019963e-07
=== RUN   Test/No_dominant:_1
iter:  0	x = [1 -0.2887417204267915]
iter:  1	x = [1 0.2887417204267915]
=== RUN   Test/No_dominant:_2
iter:  0	x = [0.0036879944609551373 1 0.14299247703760867]
iter:  1	x = [0.0036879944609551373 1 0.07149623851880434]
=== RUN   Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1
iter:  0	x = [0.8262774870764508 1]
...
iter:  4	x = [0.833332625022946 1]
e =  {10.00000004180192 [0.833333404164402 1]}
Delta = 5.9855327917e-07
=== RUN   Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9
iter:  0	x = [1 0.8347170334810676]
...
iter: 129	x = [0.7142852006468684 1]
e =  {9.999999081695126 [0.7142861765608418 1]}
Delta = 7.1472057939e-06
--- FAIL: Test (0.00s)
    --- PASS: Test/example:_1 (0.00s)
    --- PASS: Test/example:_2 (0.00s)
    --- PASS: Test/No_dominant:_1 (0.00s)
        main_test.go:72: Loop values x
    --- PASS: Test/No_dominant:_2 (0.00s)
        main_test.go:84: one or more values of eigenvector is not change
    --- PASS: Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1 (0.00s)
    --- FAIL: Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9 (0.00s)
        main_test.go:104: Precition is not ok
```

Стало
```
=== RUN   Test
=== RUN   Test/example:_1
iter:  0	x = [1 0.3970094718669067]
...
iter: 15	x = [1 0.3333344348877632]
e =  {-2.000006444049512 [1 0.33333388410690795]}
Delta = 5.2251305819e-07
=== RUN   Test/example:_2
iter:  0	x = [0.9665420787664261 1]
...
iter:  6	x = [0.8333334575188529 1]
e =  {9.999999992671016 [0.8333333209147823 1]}
Delta = 1.0494214707e-07
=== RUN   Test/No_dominant:_1
iter:  0	x = [1 -0.6300067252143378]
iter:  1	x = [1 0.6300067252143378]
e =  {0 []}
=== RUN   Test/No_dominant:_2
iter:  0	x = [0.44959179856179154 1 0.8284975407886491]
iter:  1	x = [0.44959179856179154 1 0.41424877039432456]
e =  {0 []}
=== RUN   Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1
iter:  0	x = [0.8526777191589066 1]
...
iter:  5	x = [0.8333331419093029 1]
e =  {10.000000011297157 [0.8333333524757386 1]}
Delta = 1.6176163487e-07
=== RUN   Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9
iter:  0	x = [1 0.726848993174966]
...
iter: 132	x = [0.7142861993845608 1]
e =  {10.000000867276368 [0.7142852776969006 1]}
Delta = 6.7500719937e-06
--- FAIL: Test (0.00s)
    --- PASS: Test/example:_1 (0.00s)
    --- PASS: Test/example:_2 (0.00s)
    --- PASS: Test/No_dominant:_1 (0.00s)
        main_test.go:72: Loop values x
    --- PASS: Test/No_dominant:_2 (0.00s)
        main_test.go:84: one or more values of eigenvector is not change
    --- PASS: Test/Low_ratio_:_|𝜦2|/|𝜦1|_=_0.1 (0.00s)
    --- FAIL: Test/Big_ratio_:_|𝜦2|/|𝜦1|_=_0.9 (0.00s)
        main_test.go:104: Precition is not ok
```

Результаты не изменились, поэтому данная модификация бесполезна.

## step04: Размер входной матрицы

Добавим тестов для проверки ситуации входной матрицы размером 0х0 или
матрица `nil`.
```golang
t.Run("matrix size: zero", func(t *testing.T) {
	e, err := check([][]float64{})
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
t.Run("matrix size: nil", func(t *testing.T) {
	e, err := check(nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
```

Добавим следующий код для отработки матриц нулевого размера:
```golang
n := len(A)
if n == 0 {
	err = fmt.Errorf("matrix size is zero")
	return
}
```

Проверим случай с прямоугольной матрицей - этот случай является ошибкой.
```golang
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
```
Добавим следующий код для отработки данного случая:
```golang
// проверка на квадратность входной матрицы
for row := 0; row < len(A); row++ {
	if len(A[row]) != n {
		err = fmt.Errorf("input matrix is not square in row %d: [%d,%d]", row, n, len(A[row]))
		return
	}
}
```

Проверим случай с размером матрица 1х1.
```golang
t.Run("matrix size: one", func(t *testing.T) {
	e, err := check([][]float64{
		{4},
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
```
Это простой частный случай и для его отработки добавим следующее:
```golang
// для случая матрица 1х1
if n == 1 {
	e.𝑿 = []float64{1.0}
	e.𝜦 = A[0][0]
	return
}
```

## step05: исследование начального вектора `х`

У нас изначально реализован рандомизированная задача начального вектора `x`.
```golang
// инициализация произвольным вектором
rand.Seed(time.Now().UnixNano())
for i := range x {
	x[i] = rand.Float64() // [0.0, 1)
}
```

Модифицируем код для задания требуемых значений начального вектора.
```golang
func random(x []float64) {
	// инициализация произвольным вектором
	rand.Seed(time.Now().UnixNano())
	for i := range x {
		x[i] = rand.Float64() // [0.0, 1)
	}
}

var initialize func([]float64) = random
```
и
```golang
	...
	// инициализация произвольным вектором
	initialize(x)

	...
```

Добавим тест с инициализатором нулями:
```golang
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
```
Его результат предсказуем некорректен:
```
=== RUN   Test/initialize_by_zeros
iter:  0	x = [NaN NaN]
e =  {NaN [NaN NaN]}
Delta = NaN
```
Это легко отработать в ошибку:
```golang
if max == 0.0 {
	err = fmt.Errorf("all values of eigenvector is zeros")
	return
}
```
также стоит внести это в рандимизиронную функцию:
```golang
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
```

Что если инициализировать точным значением свободного вектора 1 формы
```golang
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
```
Результат предсказуемо положительный:
```
=== RUN   Test/initialize_by_eigenvector1
iter:  0	x = [1 0.33333333333333326]
e =  {-2 [1 0.3333333333333333]}
Delta = 1.1102230246e-16
```

Что если инициализировать точным значением свободного вектора 2 формы
```golang
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
	if math.Abs(e.𝜦-(-2)) > 1e-5 {
		t.Fatal("result is not correct")
	}
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
```
Результат предсказуемый не неприемлемый, так как результат сходиться ко второму
собственному значению.
```
=== RUN   Test/initialize_by_eigenvector2
iter:  0	x = [1 0.25]
e =  {-1 [1 0.25]}
Delta = 0.0000000000e+00
    --- FAIL: Test/initialize_by_eigenvector2 (0.00s)
        main_test.go:198: result is not correct
```

Требуется отработка данной ситуации. Для изучения создадим пример, который 
будет показывать насколько близко надо приблизиться к собственному 
вектору 2 форму, чтобы получить некорректный ответ.

```golang
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
		fmt.Printf("ratio: %6.4f x: [%6.4f %6.4f]. Result: 𝜦=%6.4f 𝑿=[%6.4f %6.4f]\n",
			value, 1.0, 0.25*value+0.33333333333333333*(1.0-value),
			e.𝜦, e.𝑿[0], e.𝑿[1])
	}

	// Output:
	// ratio: 0.9999200 x: [1.00 0.2500067]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999400 x: [1.00 0.2500050]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999600 x: [1.00 0.2500033]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 0.9999800 x: [1.00 0.2500017]. Result: 𝜦=-2.0000 𝑿=[1.0000 0.3333]
	// ratio: 1.0000000 x: [1.00 0.2500000]. Result: 𝜦=-1.0000 𝑿=[1.0000 0.2500]
}
```

Данный пример показал что даже приблизив начальный собственный вектор 
на 99.998% к собственному вектору 2 формы, результат всё равно сходиться
к корректному.
Схожее поведение стоит ожидать и в случаи, если когда матрица имеет `n`
собственных значений и в качестве начального используются точное
значение некорректного собственного значения.
Для решения этой проблемы будем следовать следующей стратегии - если
алгоритм слишком быстро сходиться то добавляется возмущение. Но при этом
необходимо не создавать вечного цикла.

```golang
if max < 𝛆 {
	// на случай слишком быстрой сходимости,
	// добавим возмущения
	if iter < 3 {
		// добавляем возмужение
		perturbation := 0.02 * (1 + rand.Float64())
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
		iter = -1 // сброс количества итераций
		continue
	}
...
```

## step06: матрица из нулей но не нулевого размера

Добавим тест:

```golang
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
```

```
=== RUN   Test/matrix_with_zeros
--- PASS: Test (0.00s)
    --- PASS: Test/matrix_with_zeros (0.00s)
        main_test.go:216: all values of eigenvector is zeros
```

Необходимо добавить обработку этой ошибки, следующим образом:

```golang
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
```

Теперь результат описания ошибки корректнее:
```
=== RUN   Test/matrix_with_zeros
--- PASS: Test (0.00s)
    --- PASS: Test/matrix_with_zeros (0.00s)
        main_test.go:216: all elements of matrix is zeros
```

## step07: треугольные матрицы

Добавим тесты с треугольными входными матрицами.
Особенность треугольных матриц в том что диагональные элементы матрицы
будут являться собственными значениями.

```golang
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
```

Результат сходится до требуемой точности:

```
=== RUN   Test/lower_triangle_matrix
iter:  0	x = [-0.5376658386946178 1]
iter:  1	x = [0.018832919347308907 1]
...
iter: 19	x = [-0.16666595904190312 1]
iter: 20	x = [-0.16666702047904844 1]
=== RUN   Test/upper_triangle_matrix
iter:  0	x = [1 0.3176319230555149 0.7494395060708893]
iter:  1	x = [1 0.3190545969693437 0.6072703715544727]
...
iter: 27	x = [1 0.20000174862493353 0.40000349724985906]
iter: 28	x = [1 0.20000116574655347 0.4000023314931096]
e =  {3.0000032381743904 [1 0.20000077716286063 0.4000015543257204]}
Delta = 1.5863821838e-06
--- PASS: Test (0.00s)
    --- PASS: Test/lower_triangle_matrix (0.00s)
    --- PASS: Test/upper_triangle_matrix (0.00s)
```

## step08: специфический случай

Тест:

```golang
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
t.Run("initialize specific : 2", func(t *testing.T) {
	old := initialize
	initialize = func(x []float64) {
		x[0] = -3.0
		x[1] = 2.0
	}
	defer func() {
		initialize = old
	}()
	e, err := check([][]float64{
		{2, 3},
		{1, 4},
	})
	if math.Abs(e.𝜦-5) > 1e-4 {
		t.Fatalf("result is not correct: %.14e ---> prec = %.14e", e.𝜦, e.𝜦+2)
	}
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
```

Пофиксим тест "initialize specific : 2":

```golang
// значение х не изменяется кроме 1.0
if iter > 0 {
	isSame = false
	...
}
```

Но вот ещё тест:

```golang
t.Run("initialize specific : 3", func(t *testing.T) {
	old := initialize
	initialize = func(x []float64) {
		x[0] = 1.0
		x[1] = 1.0
	}
	defer func() {
		initialize = old
	}()
	e, err := check([][]float64{
		{2, 3},
		{1, 4},
	})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(e.𝜦-5) > 1e-4 {
		t.Fatalf("result is not correct: %.14e ---> prec = %.14e", e.𝜦, e.𝜦+2)
	}
	_ = e
})
```

В результате ошибка на корректные данные:
```
=== RUN   Test/initialize_specific_:_3
iter:  0	x = [1 1]
iter:  0	x = [1 1]
--- FAIL: Test (0.00s)
    --- FAIL: Test/initialize_specific_:_3 (0.00s)
        main_test.go:274: Loop values x
```

Опишим проблему:
- Если инициализация вектора х совпадает со собственным вектором 2 формы и
  при этом собственное значение 2 формы равен единице и его 
  собственный вектор состоит из единиц, то процесс зацикливается.

Для решения этой проблемы изменим алгоритм возмущения:

```golang
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
iter = 0 // сброс количества итераций
continue
```

Какова причина причина того, что используется возмущение, а не генерация
новых рандомных значений:
>  предполагается, что если все же иницированные значения имеют точное значение
>  необходимой нам собственного значения

## step09: специфический случай без доминированного собственного значения

Тест:

```golang
t.Run("initialize specific : 4", func(t *testing.T) {
	old := initialize
	initialize = func(x []float64) {
		x[0] = 3.0
		x[1] = 0.0
		x[2] = 1.0
	}
	defer func() {
		initialize = old
	}()
	e, err := check([][]float64{
		{3, 2, -3},
		{-3, -4, 9},
		{-1, -2, 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = e
})
```

Точное значение собственных векторов: `2, 2, 2e-16`.

В настоящий момент вылетает с ошибкой по количеству итераций.

```
...
iter:  1	x = [0.5922532264568454 1 0.8640844088189488]
iter:  2	x = [0.592253226456844  1 0.8640844088189479]
iter:  1	x = [0.5918284041467258 1 0.8639428013822417]
iter:  2	x = [0.5918284041467267 1 0.8639428013822422]
iter:  1	x = [0.5900439234481327 1 0.8633479744827108]
iter:  2	x = [0.590043923448133  1 0.863347974482711 ]
...
--- FAIL: Test (0.00s)
    --- FAIL: Test/initialize_specific_:_4 (0.00s)
        main_test.go:297: global iteration limit. please send to developer
```

Далее будем использовать следующую литературу:
* Д.К.Фадеев, В.Н.Фадеева "Вычислительные методы линейной алгебры"

Список сделанного:
* Добавил примеры из книги Фадеев
* Выделение расчета собственного значения в отдельную функцию
* Добавлена отработка парных свободных значений
* Удален код для отслеживания повторяющихся собственных вектором. Это было 
  сделано, так как это работало не для всех ситуаций, а также только 
  когда матрица имела парные собственные числа.

