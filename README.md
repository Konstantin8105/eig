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

## Степенной метод(power method)

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

