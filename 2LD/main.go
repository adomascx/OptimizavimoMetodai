package main

import (
	"fmt"
	"math"
	"os"
)

// mano studento skaitmenys
var A float64 = 3
var B float64 = 9

// 2D koordinaciu taskas
type point struct {
	x float64
	y float64
}

// Xn tasko reprezentacija
type startPoint struct {
	name string
	p    point
}

// algoritmo isvesties info
type runSummary struct {
	end  point
	dev  float64
	iter int
	fn   int
	grad int
}

// iskvietimus laikome globaliai
// kai f-ja kvieciama, ji pati atnaujina counter
type iskvietimuSk struct {
	fnCalls   int
	gradCalls int
}

var iskvietimai iskvietimuSk

type funkcija func(p point) (z float64)

type gradientas func(p point) (g point)

// musu tikslo f-ja (turis kvadratu)
func tiksloFunkcija(p point) float64 {
	iskvietimai.fnCalls++
	return -(p.x * p.y * (1 - p.x - p.y)) / 8
}

// is tikslo f-jos gauta gradiento f-ja
func gradientoFunkcija(p point) (g point) {
	iskvietimai.gradCalls++
	g.x = p.y * (p.y + 2*p.x - 1) / 8
	g.y = p.x * (2*p.y + p.x - 1) / 8
	return g
}

// helper funkcija, jei naujas taskas nepriklauso funkcijos riboms
func project(p point) point {
	if p.x < 0 {
		p.x = 0
	}
	if p.y < 0 {
		p.y = 0
	}
	if s := p.x + p.y; s > 1 {
		p.x /= s
		p.y /= s
	}
	return p
}

// helper funkcija tvarkanti floating point precision klaidas
func snap(v float64) float64 {
	if math.Abs(v) < 5e-13 {
		return 0
	}
	return v
}

func gradientinioNusileidimoAlgo(f funkcija, gf gradientas, p0 point, gamma, eps float64, maxIter int) (p point, value float64, k int, history []point) {
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

func greiciausiojoNusileidimoAlgo(f funkcija, gf gradientas, p0 point, eps, epsGamma, initialStep, maxGamma float64, maxIter int) (p point, value float64, k int, history []point) {
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

		if a == b {
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

func simplexSize(a, b, c point) float64 {
	dAB := math.Hypot(a.x-b.x, a.y-b.y)
	dAC := math.Hypot(a.x-c.x, a.y-c.y)
	dBC := math.Hypot(b.x-c.x, b.y-c.y)
	return math.Max(dAB, math.Max(dAC, dBC))
}

// helper funkcija duomenu eksportavimui i Python grafiku generavimui
func exportPlotData(file string, points []startPoint, trueMin float64, gdHistory, sdHistory, simplexHistory [][]point, gdSummary, sdSummary, simplexSummary []runSummary, gridN int) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Fprintln(out, "# META key value")
	fmt.Fprintf(out, "META A %.10f\n", A)
	fmt.Fprintf(out, "META B %.10f\n", B)
	fmt.Fprintf(out, "META TRUE_MIN %.12f\n", trueMin)
	fmt.Fprintf(out, "GRID_N %d\n", gridN)

	fmt.Fprintln(out, "# SUMMARY algo start endX endY dev iter fn grad")
	for i, sp := range points {
		s := gdSummary[i]
		fmt.Fprintf(out, "SUMMARY GD %s %.10f %.10f %.10e %d %d %d\n", sp.name, s.end.x, s.end.y, s.dev, s.iter, s.fn, s.grad)
	}
	for i, sp := range points {
		s := sdSummary[i]
		fmt.Fprintf(out, "SUMMARY SD %s %.10f %.10f %.10e %d %d %d\n", sp.name, s.end.x, s.end.y, s.dev, s.iter, s.fn, s.grad)
	}
	for i, sp := range points {
		s := simplexSummary[i]
		fmt.Fprintf(out, "SUMMARY SIMPLEX %s %.10f %.10f %.10e %d %d %d\n", sp.name, s.end.x, s.end.y, s.dev, s.iter, s.fn, s.grad)
	}

	fmt.Fprintln(out, "# PATH algo start idx x y")
	for i, sp := range points {
		for j, p := range gdHistory[i] {
			fmt.Fprintf(out, "PATH GD %s %d %.10f %.10f\n", sp.name, j, p.x, p.y)
		}
	}
	for i, sp := range points {
		for j, p := range sdHistory[i] {
			fmt.Fprintf(out, "PATH SD %s %d %.10f %.10f\n", sp.name, j, p.x, p.y)
		}
	}
	for i, sp := range points {
		for j, p := range simplexHistory[i] {
			fmt.Fprintf(out, "PATH SIMPLEX %s %d %.10f %.10f\n", sp.name, j, p.x, p.y)
		}
	}

	fmt.Fprintln(out, "# GRID i j x y gradMag")
	for j := 0; j < gridN; j++ {
		y := float64(j) / float64(gridN-1)
		for i := 0; i < gridN; i++ {
			x := float64(i) / float64(gridN-1)
			gx := y * (y + 2*x - 1) / 8
			gy := x * (2*y + x - 1) / 8
			gm := math.Hypot(gx, gy)
			fmt.Fprintf(out, "GRID %d %d %.10f %.10f %.12f\n", i, j, x, y, gm)
		}
	}

	return nil
}

func runAndReport(title string, points []startPoint, trueMin float64, run func(point) (point, float64, int, []point), history [][]point, summary []runSummary) {
	fmt.Printf("\n***%s***\n", title)
	for i, start := range points {
		iskvietimai = iskvietimuSk{}
		res, value, iter, hist := run(start.p)
		history[i] = hist
		dev := math.Abs(value - trueMin)
		summary[i] = runSummary{end: res, dev: dev, iter: iter, fn: iskvietimai.fnCalls, grad: iskvietimai.gradCalls}
		res.x = snap(res.x)
		res.y = snap(res.y)
		fmt.Printf("%s X=(%.6f,%.6f) dev=%.3e iter=%d fn=%d grad=%d\n",
			start.name, res.x, res.y, dev, iter, iskvietimai.fnCalls, iskvietimai.gradCalls)
	}
}

func main() {
	gammaGD := 0.2
	eps := 1e-6
	epsGamma := 1e-6
	maxIter := 5000
	initialStep := 1e-3
	maxGamma := 16.0
	alpha := 0.1
	gammaSimplex := 2.0
	beta := 0.5
	trueMin := -1.0 / 216.0

	points := []startPoint{
		{name: "X0", p: point{x: 0, y: 0}},
		{name: "X1", p: point{x: 1, y: 1}},
		{name: "Xm", p: point{x: A / 10, y: B / 10}},
	}

	gdHistory := make([][]point, len(points))
	sdHistory := make([][]point, len(points))
	simplexHistory := make([][]point, len(points))
	gdSummary := make([]runSummary, len(points))
	sdSummary := make([]runSummary, len(points))
	simplexSummary := make([]runSummary, len(points))

	runAndReport("Gradientinio nusileidimo algoritmas", points, trueMin, func(p point) (point, float64, int, []point) {
		return gradientinioNusileidimoAlgo(tiksloFunkcija, gradientoFunkcija, p, gammaGD, eps, maxIter)
	}, gdHistory, gdSummary)

	runAndReport("Greiciausiojo nusileidimo algoritmas", points, trueMin, func(p point) (point, float64, int, []point) {
		return greiciausiojoNusileidimoAlgo(tiksloFunkcija, gradientoFunkcija, p, eps, epsGamma, initialStep, maxGamma, maxIter)
	}, sdHistory, sdSummary)

	runAndReport("Deformuojamo simplekso algoritmas", points, trueMin, func(p point) (point, float64, int, []point) {
		return deformuojamoSimpleksoAlgo(tiksloFunkcija, p, alpha, gammaSimplex, beta, eps, maxIter)
	}, simplexHistory, simplexSummary)

	if err := exportPlotData("plot_data.txt", points, trueMin, gdHistory, sdHistory, simplexHistory, gdSummary, sdSummary, simplexSummary, 301); err != nil {
		fmt.Printf("Nepavyko eksportuoti plot_data.txt: %v\n", err)
	} else {
		fmt.Println("\nDuomenys eksportuoti i plot_data.txt")
	}
}
