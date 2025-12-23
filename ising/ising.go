package ising

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Params задаёт параметры моделирования.
type Params struct {
	L      int     // размер решётки LxL
	T1     float64 // начальная температура
	T2     float64 // конечная температура
	Tcount int     // количество точек по температуре

	ASteps int // шаги термализации
	MSteps int // шаги измерения
	Copies int // количество независимых копий решётки
}

// ResultRow содержит усреднённые по шагам и копиям величины при одной температуре.
type ResultRow struct {
	T float64
	U float64
	M float64
	C float64
	X float64
}

type array2d [][]int

// pbc — периодические граничные условия по размеру L.
func pbc(x, L int) int {
	if x < 0 {
		return x + L
	}
	return x % L
}

// calcParameters считает полную энергию и намагниченность решётки.

func calcParameters(lattice array2d, L int, J, h float64, energy, moment *float64) {
	*energy = 0
	*moment = 0
	for x := 0; x < L; x++ {
		for y := 0; y < L; y++ {
			S := lattice[x][y]
			Sr := lattice[pbc(x+1, L)][y]
			Sb := lattice[x][pbc(y+1, L)]
			*energy += -J * float64(S) * float64(Sr)
			*energy += -J * float64(S) * float64(Sb)
			*energy += -h * float64(S)
			*moment += float64(S)
		}
	}
}

// mcStep — один шаг алгоритма Метрополиса в выбранной точке решётки.
func mcStep(lattice array2d, L int, J, h, T float64, x, y int) {
	S0 := lattice[x][y]
	S1 := -S0
	Sr := lattice[pbc(x+1, L)][y] // правый
	Sb := lattice[x][pbc(y+1, L)] // нижний
	Sl := lattice[pbc(x-1, L)][y] // левый
	St := lattice[x][pbc(y-1, L)] // верхний
	dE := float64(S1-S0) * (-h - J*float64(Sl+Sr+St+Sb))
	if rand.Float64() < math.Exp(-dE/T) {
		lattice[x][y] = S1
	}
}

// nextStep выполняет один "маятник" — N случайных попыток 
func nextStep(lattice array2d, L int, J, h, T float64) {
	N := L * L
	for i := 0; i < N; i++ {
		x := rand.Intn(L)
		y := rand.Intn(L)
		mcStep(lattice, L, J, h, T, x, y)
	}
}

// RunIsing запускает моделирование Изинга для заданных J, h и параметров.

func RunIsing(J, h float64, params Params) ([]ResultRow, error) {
	if params.L <= 0 {
		return nil, fmt.Errorf("L must be > 0")
	}
	if params.Tcount <= 1 {
		return nil, fmt.Errorf("Tcount must be > 1")
	}
	if params.MSteps <= 0 {
		return nil, fmt.Errorf("MSteps must be > 0")
	}
	if params.Copies <= 0 {
		return nil, fmt.Errorf("Copies must be > 0")
	}

	L := params.L
	T1 := params.T1
	T2 := params.T2
	Tcount := params.Tcount
	aSteps := params.ASteps
	mSteps := params.MSteps
	copies := params.Copies

	// Инициализация генератора случайных чисел.
	rand.Seed(time.Now().UnixNano())

	N := L * L

	// Создаём копии решётки, изначально все спины = +1.
	lattices := make([]array2d, copies)
	for k := 0; k < copies; k++ {
		lattice := make(array2d, 0, L)
		for i := 0; i < L; i++ {
			row := make([]int, L)
			for j := 0; j < L; j++ {
				row[j] = 1
			}
			lattice = append(lattice, row)
		}
		lattices[k] = lattice
	}

	results := make([]ResultRow, 0, Tcount)

	for tIdx := 0; tIdx < Tcount; tIdx++ {
		T := T1 + float64(tIdx)*(T2-T1)/float64(Tcount-1)

		E := 0.0
		E2 := 0.0
		M := 0.0
		M2 := 0.0

		// Для каждой копии проводим термализацию и измерение.
		for copyIdx := 0; copyIdx < copies; copyIdx++ {
			lattice := lattices[copyIdx]

			// Термализация.
			for s := 0; s < aSteps; s++ {
				nextStep(lattice, L, J, h, T)
			}

			// Измерения.
			for s := 0; s < mSteps; s++ {
				nextStep(lattice, L, J, h, T)

				energy := 0.0
				moment := 0.0
				calcParameters(lattice, L, J, h, &energy, &moment)

				// Усреднение по шагам и копиям.
				E += energy / float64(mSteps) / float64(copies)
				E2 += energy * energy / float64(mSteps) / float64(copies)
				M += math.Abs(moment) / float64(mSteps) / float64(copies)
				M2 += moment * moment / float64(mSteps) / float64(copies)
			}
		}

		U := E / float64(N)
		C := (E2 - E*E) / (T * T * float64(N))
		m := M / float64(N)
		X := (M2 - M*M) / (T * float64(N))

		results = append(results, ResultRow{
			T: T,
			U: U,
			M: m,
			C: C,
			X: X,
		})
	}

	return results, nil
}
