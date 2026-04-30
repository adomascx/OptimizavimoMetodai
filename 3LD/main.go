package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// mano studento skaitmenys
var A float64 = 9
var B float64 = 3
var C float64 = 9

// kintamuju vektorius
type kintamieji struct {
	x float64
	y float64
	z float64
}

// iskvietimus laikome globaliai
// kai f-ja kvieciama, ji pati atnaujina counter
var iskvietimai int

type funkcija func(X kintamieji) float64
type gradientas func(X kintamieji) kintamieji

// musu tikslo f-ja
func tiksloFunkcija(X kintamieji) float64 {
	iskvietimai++
	return -X.x * X.y * X.z
}

// lygybinis apribojimas
func g(X kintamieji) float64 {
	return 2*X.x*X.y + 2*X.x*X.z + 2*X.y*X.z - 1
}

// nelygybiniai apribojimai
func h1(X kintamieji) float64 {
	return -X.x
}

func h2(X kintamieji) float64 {
	return -X.y
}

func h3(X kintamieji) float64 {
	return -X.z
}

func baudosDalis(X kintamieji) float64 {
	gX := g(X)

	h1X := math.Max(0, h1(X))
	h2X := math.Max(0, h2(X))
	h3X := math.Max(0, h3(X))

	return gX*gX + h1X*h1X + h2X*h2X + h3X*h3X
}

func baudosFunkcija(X kintamieji, r float64) float64 {
	return tiksloFunkcija(X) + (1/r)*baudosDalis(X)

}

func gradientoBaudosFunkcija(X kintamieji, r float64) kintamieji {
	grad := kintamieji{
		x: -(X.y * X.z),
		y: -(X.x * X.z),
		z: -(X.x * X.y),
	}

	gX := g(X)
	dg := kintamieji{
		x: 2*X.y + 2*X.z,
		y: 2*X.x + 2*X.z,
		z: 2*X.x + 2*X.y,
	}

	h1X := math.Max(0, h1(X))
	h2X := math.Max(0, h2(X))
	h3X := math.Max(0, h3(X))

	grad.x += (1 / r) * (2*gX*dg.x - 2*h1X)
	grad.y += (1 / r) * (2*gX*dg.y - 2*h2X)
	grad.z += (1 / r) * (2*gX*dg.z - 2*h3X)

	return grad
}

func greiciausiojoNusileidimoAlgo(f funkcija, gf gradientas, p0 kintamieji, eps, epsGamma, initialStep, maxGamma float64, maxIter int) (p kintamieji, value float64, k int) {
	p = p0

	value = f(p)

	for k = 0; k < maxIter; k++ {
		g := gf(p)
		gNorm := math.Sqrt(g.x*g.x + g.y*g.y + g.z*g.z)
		if gNorm <= eps {
			break
		}

		s := kintamieji{x: -g.x, y: -g.y, z: -g.z}
		phi := func(gamma float64) float64 {
			cand := kintamieji{x: p.x + gamma*s.x, y: p.y + gamma*s.y, z: p.z + gamma*s.z}
			return f(cand)
		}

		a := 0.0
		b := initialStep
		fa := phi(a)
		fb := phi(b)

		for b < maxGamma && fb < fa {

			a = b
			fa = fb

			b *= 2
			if b > maxGamma {
				b = maxGamma
			}

			fb = phi(b)
		}

		if 0 == b {
			b = maxGamma
		}

		tau := (math.Sqrt(5) - 1) / 2
		c := b - tau*(b-a)
		d := a + tau*(b-a)
		fc := phi(c)
		fd := phi(d)

		for i := 0; i < 200 && (b-a) > epsGamma; i++ {
			if fc < fd {
				b = d
				d = c
				fd = fc
				c = b - tau*(b-a)
				fc = phi(c)
			} else {
				a = c
				c = d
				fc = fd
				d = a + tau*(b-a)
				fd = phi(d)
			}
		}

		gamma := (a + b) / 2
		p = kintamieji{x: p.x + gamma*s.x, y: p.y + gamma*s.y, z: p.z + gamma*s.z}

		value = f(p)
	}

	return p, value, k
}

// apskaiciuoti tikslo f-ja duotame taske
func outputFunctionResults(ivestiesKintamieji kintamieji) {
	fmt.Printf("\n%v %v %v\n", strings.Repeat("-", 6), ivestiesKintamieji, strings.Repeat("-", 6))

	fmt.Printf("F(%v,%v,%v) = %s\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		strconv.FormatFloat(math.Round(tiksloFunkcija(ivestiesKintamieji)*1e3)/1e3, 'f', -1, 64),
	)
	fmt.Printf("g(%v,%v,%v) = %s\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		strconv.FormatFloat(math.Round(g(ivestiesKintamieji)*1e3)/1e3, 'f', -1, 64),
	)
	fmt.Printf("h1(%v,%v,%v) = %v\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		h1(ivestiesKintamieji),
	)
	fmt.Printf("h2(%v,%v,%v) = %v\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		h2(ivestiesKintamieji),
	)
	fmt.Printf("h3(%v,%v,%v) = %v\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		h3(ivestiesKintamieji),
	)
}

// apskaiciuoti baudos f-ja duotame taske
func outputPenaltyResults(ivestiesKintamieji kintamieji, r float64) {
	fmt.Printf("r = %-5v | B(X) = %s\n", r, strconv.FormatFloat(math.Round(baudosFunkcija(ivestiesKintamieji, r)*1e3)/1e3, 'f', -1, 64))
}

func outputOptimizationResults(start kintamieji, eps, epsGamma, initialStep, maxGamma float64, maxIter int, verbose bool) {

	fmt.Printf("\n%v %v %v\n", strings.Repeat("-", 24), start, strings.Repeat("-", 24))
	iskvietimai = 0
	trueVal := 1 / math.Sqrt(6)
	p := start
	totalIter := 0
	for r := 1.0; r >= 1e-3; r *= 0.1 {
		penaltyFunc := func(X kintamieji) float64 {
			return baudosFunkcija(X, r)
		}
		penaltyGrad := func(X kintamieji) kintamieji {
			return gradientoBaudosFunkcija(X, r)
		}

		var k int
		p, _, k = greiciausiojoNusileidimoAlgo(penaltyFunc, penaltyGrad, p, eps, epsGamma, initialStep, maxGamma, maxIter)
		totalIter += k
		errorDist := math.Sqrt(
			(p.x-trueVal)*(p.x-trueVal) +
				(p.y-trueVal)*(p.y-trueVal) +
				(p.z-trueVal)*(p.z-trueVal),
		)
		if verbose {
			fmt.Printf("r = %.1e | skirtumas nuo tikro = %-8s | k = %d\n",
				r,
				strconv.FormatFloat(math.Round(errorDist*1e4)/1e4, 'f', -1, 64),
				k,
			)
		}
	}

	if verbose {
		fmt.Printf("\n")
	}

	errorDist := math.Sqrt(
		(p.x-trueVal)*(p.x-trueVal) +
			(p.y-trueVal)*(p.y-trueVal) +
			(p.z-trueVal)*(p.z-trueVal),
	)
	functionEvals := iskvietimai
	fX := tiksloFunkcija(p)
	iskvietimai = functionEvals
	fmt.Printf("X = %+-20v\nskirtumas nuo tikro = %-8s | f(X) = %-9s | k = %d | evals = %d\n",
		p,
		strconv.FormatFloat(math.Round(errorDist*1e4)/1e4, 'f', -1, 64),
		strconv.FormatFloat(math.Round(fX*1e3)/1e3, 'f', -1, 64),
		totalIter,
		functionEvals,
	)
}

func main() {
	const (
		eps         = 1e-6
		epsGamma    = 1e-6
		initialStep = 0.1
		maxGamma    = 10.0
		maxIter     = 1e6
	)

	taskai := []kintamieji{
		{x: 0, y: 0, z: 0},
		{x: 1, y: 1, z: 1},
		{x: A / 10, y: B / 10, z: C / 10},
	}

	// tikslo f-jos reiksmes taskuose
	fmt.Printf("\n%v\n", strings.Repeat("=", 38))
	fmt.Println("Tikslo funkcijos reikšmės")
	for _, taskas := range taskai {
		outputFunctionResults(taskas)
	}

	// baudos f-jos reiksmes taskuose
	rValues := []float64{1, 0.1, 0.01, 0.001}
	fmt.Printf("\n%v\n", strings.Repeat("=", 38))
	fmt.Println("Baudos funkcijos reikšmės")
	for _, taskas := range taskai {

		fmt.Printf("\n%v %v %v\n", strings.Repeat("-", 6), taskas, strings.Repeat("-", 6))
		for _, r := range rValues {
			outputPenaltyResults(taskas, r)
		}
	}

	fmt.Printf("\n%v\n", strings.Repeat("=", 38))
	fmt.Println("Baudos funkcijos minimizavimas")
	for _, taskas := range taskai {
		outputOptimizationResults(taskas, eps, epsGamma, initialStep, maxGamma, maxIter, true)
	}
}
