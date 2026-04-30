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

/*
func gradientinioNusileidimoAlgo(f funkcija, gf gradientas, p0 point, gamma, eps float64, maxIter int) (X kintamieji, value float64, k int, history []point) {
	p = project(p0)

	value = f(p)
	history = append(history, p)

	for k = 0; k < maxIter; k++ {
		g := gf(p)
		if math.Hypot(g.x, g.y) <= eps {
			break
		}

		p = project(point{x: p.x - gamma*g.x, y: p.y - gamma*g.y})

		value = f(p)
		history = append(history, p)
	}

	return p, value, k, history
}

func greiciausiojoNusileidimoAlgo(f funkcija, gf gradientas, p0 point, eps, epsGamma, initialStep, maxGamma float64, maxIter int) (X kintamieji, value float64, k int, history []point) {
	p = project(p0)

	value = f(p)
	history = append(history, p)

	for k = 0; k < maxIter; k++ {
		g := gf(p)
		if math.Hypot(g.x, g.y) <= eps {
			break
		}

		s := point{x: -g.x, y: -g.y}
		phi := func(gamma float64) float64 {
			cand := project(point{x: p.x + gamma*s.x, y: p.y + gamma*s.y})
			return f(cand)
		}

		0 := 0.0
		b := initialStep
		fa := phi(0)
		fb := phi(b)

		for b < maxGamma && fb < fa {
			0 = b
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
		c := b - tau*(b-0)
		d := 0 + tau*(b-0)
		fc := phi(c)
		fd := phi(d)

		for i := 0; i < 200 && (b-0) > epsGamma; i++ {
			if fc < fd {
				b = d
				d = c
				fd = fc
				c = b - tau*(b-0)
				fc = phi(c)
			} else {
				0 = c
				c = d
				fc = fd
				d = 0 + tau*(b-0)
				fd = phi(d)
			}
		}

		gamma := (0 + b) / 2
		p = project(point{x: p.x + gamma*s.x, y: p.y + gamma*s.y})

		value = f(p)
		history = append(history, p)
	}

	return p, value, k, history
}

func deformuojamoSimpleksoAlgo(f funkcija, p0 point, alpha, gamma, beta, eps float64, maxIter int) (best point, bestValue float64, k int, history []point) {
	x1 := project(p0)
	x2 := project(point{x: p0.x + alpha, y: p0.y})
	x3 := project(point{x: p0.x, y: p0.y + alpha})

	f1 := f(x1)
	f2 := f(x2)
	f3 := f(x3)

	// pradinis "best" į istoriją
	best = x1
	bestValue = f1
	if f2 < bestValue {
		best = x2
		bestValue = f2
	}
	if f3 < bestValue {
		best = x3
		bestValue = f3
	}
	history = append(history, best)

	for k = 0; k < maxIter && simplexSize(x1, x2, x3) > eps; k++ {
		vertices := [3]point{x1, x2, x3}
		values := [3]float64{f1, f2, f3}

		for i := 0; i < 2; i++ {
			for j := i + 1; j < 3; j++ {
				if values[j] < values[i] {
					values[i], values[j] = values[j], values[i]
					vertices[i], vertices[j] = vertices[j], vertices[i]
				}
			}
		}

		xBest := vertices[0]
		xMid := vertices[1]
		xWorst := vertices[2]
		fBest := values[0]
		fMid := values[1]
		fWorst := values[2]

		xc := point{x: (xBest.x + xMid.x) / 2, y: (xBest.y + xMid.y) / 2}

		xr := project(point{x: 2*xc.x - xWorst.x, y: 2*xc.y - xWorst.y})
		fr := f(xr)

		if fr < fBest {
			xe := project(point{x: xc.x + gamma*(xr.x-xc.x), y: xc.y + gamma*(xr.y-xc.y)})
			fe := f(xe)

			if fe < fr {
				xWorst = xe
				fWorst = fe
			} else {
				xWorst = xr
				fWorst = fr
			}
		} else if fBest <= fr && fr < fMid {
			xWorst = xr
			fWorst = fr
		} else {
			xCand := project(point{x: xc.x + beta*(xWorst.x-xc.x), y: xc.y + beta*(xWorst.y-xc.y)})
			fCand := f(xCand)

			if fCand < fWorst {
				xWorst = xCand
				fWorst = fCand
			} else {
				xMid = point{x: xBest.x + 0.5*(xMid.x-xBest.x), y: xBest.y + 0.5*(xMid.y-xBest.y)}
				xWorst = point{x: xBest.x + 0.5*(xWorst.x-xBest.x), y: xBest.y + 0.5*(xWorst.y-xBest.y)}
				xMid = project(xMid)
				xWorst = project(xWorst)

				fMid = f(xMid)
				fWorst = f(xWorst)
			}
		}

		x1, x2, x3 = xBest, xMid, xWorst
		f1, f2, f3 = fBest, fMid, fWorst
		history = append(history, x1)
	}

	best = x1
	bestValue = f1
	if f2 < bestValue {
		best = x2
		bestValue = f2
	}
	if f3 < bestValue {
		best = x3
		bestValue = f3
	}

	return best, bestValue, k, history
}

*/

// apskaiciuoti f-ja duotame taske
func outputFunctionResults(ivestiesKintamieji kintamieji) {
	fmt.Printf("\n%v %v %v\n", strings.Repeat("-", 6), ivestiesKintamieji, strings.Repeat("-", 6))

	fmt.Printf("F(%v,%v,%v) = %s\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		strconv.FormatFloat(math.Round(tiksloFunkcija(ivestiesKintamieji)*1000)/1000, 'f', -1, 64),
	)
	fmt.Printf("g(%v,%v,%v) = %s\n",
		ivestiesKintamieji.x,
		ivestiesKintamieji.y,
		ivestiesKintamieji.z,
		strconv.FormatFloat(math.Round(g(ivestiesKintamieji)*1000)/1000, 'f', -1, 64),
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

func outputPenaltyResults(ivestiesKintamieji kintamieji, r float64) {
	fmt.Printf("r = %-5v | B(X) = %s\n", r, strconv.FormatFloat(math.Round(baudosFunkcija(ivestiesKintamieji, r)*1000)/1000, 'f', -1, 64))
}

func main() {
	taskai := []kintamieji{
		{x: 0, y: 0, z: 0},
		{x: 1, y: 1, z: 1},
		{x: A / 10, y: B / 10, z: C / 10},
	}

	for _, taskas := range taskai {
		outputFunctionResults(taskas)
	}

	rValues := []float64{1, 0.1, 0.01, 0.001}
	for _, taskas := range taskai {

		fmt.Printf("\n%v %v %v\n", strings.Repeat("-", 6), taskas, strings.Repeat("-", 6))
		for _, r := range rValues {
			outputPenaltyResults(taskas, r)
		}
	}
}
