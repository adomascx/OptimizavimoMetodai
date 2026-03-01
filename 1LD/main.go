package main

import (
	"fmt"
	"log"
	"math"
)

// mano studento skaitmenys
var A float64 = 3
var B float64 = 9

type tiksloFunkcija func(x float64) float64

// Tikslo funkcija + jos 1osios ir 2osios eilės isvestines
func funkcija(x float64) (y float64) {
	return math.Pow(x*x-A, 2) / (B - 1)
}

func isvestine_1(x float64) (y float64) {
	return (x * 4 * (x*x - A)) / (B - 1)
}

func isvestine_2(x float64) (y float64) {
	return (x*x*12 - A*4) / (B - 1)
}

// Algoritmu implementacijos
func intervaloDalijimoPusiauAlgo(f tiksloFunkcija, l, r, epsilon float64, maxIter int) (xMin, fMin float64, iteracijuSk, funkcijosSkaiciavimuSk int, err error) {

	funkcijosSkaiciavimuSk = 0

	xm := (l + r) / 2

	for i := 0; i < maxIter; i++ {

		if (r - l) <= epsilon {
			x := (l + r) / 2
			return x, f(x), i, funkcijosSkaiciavimuSk + 1, nil
		}

		L := r - l

		x1 := l + L/4
		x2 := r - L/4

		if f(x1) < f(xm) {
			r = xm
			xm = x1

		} else if f(x2) < f(xm) {
			l = xm
			xm = x2

		} else {
			l = x1
			r = x2

		}
		funkcijosSkaiciavimuSk += 2

	}

	x := (l + r) / 2

	return x, f(x), maxIter, funkcijosSkaiciavimuSk + 1, fmt.Errorf("pasiektas maksimalus iteracijų skaičius.")

}

func auksinioPjuvioAlgo(f tiksloFunkcija, l, r, epsilon float64, maxIter int) (xMin, fMin float64, iteracijuNr, funkcijosSkaiciavimuSk int, err error) {

	funkcijosSkaiciavimuSk = 0

	tau := (math.Sqrt(5.0) - 1.0) / 2.0 // ~0.618

	x1 := r - tau*(r-l)
	x2 := l + tau*(r-l)

	for i := 0; i < maxIter; i++ {

		if (r - l) <= epsilon {
			x := (l + r) / 2
			return x, f(x), i, funkcijosSkaiciavimuSk + 1, nil
		}

		if f(x2) < f(x1) {
			l = x1
			x1 = x2
			x2 = l + tau*(r-l)

		} else {
			r = x2
			x2 = x1
			x1 = r - tau*(r-l)

		}
		funkcijosSkaiciavimuSk += 1
	}

	x := (l + r) / 2
	return x, f(x), maxIter, funkcijosSkaiciavimuSk + 1, fmt.Errorf("pasiektas maksimalus iteracijų skaičius.")
}

func niutonoAlgo(f, df, d2f tiksloFunkcija, x0, epsilon float64, maxIter int) (xMin, fMin float64, iteracijuNr, funkcijosSkaiciavimuSk int, err error) {

	funkcijosSkaiciavimuSk = 0

	x := x0

	for i := 0; i < maxIter; i++ {

		if math.Abs(df(x)) <= epsilon {
			return x, f(x), i, funkcijosSkaiciavimuSk + 1, nil
		}

		xNext := x - df(x)/d2f(x)
		funkcijosSkaiciavimuSk += 2

		if math.Abs(xNext-x) <= epsilon {
			x = xNext
			return x, f(x), i + 1, funkcijosSkaiciavimuSk, nil
		}
		x = xNext
	}

	return x, f(x), maxIter, funkcijosSkaiciavimuSk + 1, fmt.Errorf("pasiektas maksimalus iteracijų skaičius.")
}

func main() {

	xMin, fMin, iter, funk, err := intervaloDalijimoPusiauAlgo(funkcija, 1, 10, 1e-4, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n**Intervalo Dalijimo Algoritmas**\nxMin: %v\nfMin: %.1e\niteraciju: %v\nfunkcijos skaiciavimu: %v\n", xMin, fMin, iter, funk)

	xMin, fMin, iter, funk, err = auksinioPjuvioAlgo(funkcija, 1, 10, 1e-4, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n**Auksinio Pjuvio Algoritmas**\nxMin: %v\nfMin: %.1e\niteraciju: %v\nfunkcijos skaiciavimu: %v\n", xMin, fMin, iter, funk)

	xMin, fMin, iter, funk, err = niutonoAlgo(funkcija, isvestine_1, isvestine_2, 5, 1e-4, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n**Niutono Algoritmas**\nxMin: %v\nfMin: %.1e\niteraciju: %v\nfunkcijos skaiciavimu: %v\n", xMin, fMin, iter, funk)
}
