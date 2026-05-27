package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"text/tabwriter"
)

var kintamieji = []string{"x1", "x2", "x3", "x4", "s1", "s2", "s3"}

// Sprendzia tiesinio programavimo uzdavini simplex metodu: max z = c^T x

// A yra m x n matrica (jau su laisvumo kintamaisiais), b – m ilgio, c – n ilgio
// basis0 nurodo pradine baze (kintamuju vardais). hook kvieciamas pradinei lentelei
// (iter=0) ir po kiekvienos iteracijos
func simplex(A [][]float64, b []float64, c []float64, basis0 []string, hook func(iter int, tableau [][]float64, basis []string)) (zMax float64, x []float64, basis []string, tableau [][]float64) {
	const eps = 1e-9
	m := len(A)
	n := len(A[0])
	basis = append([]string(nil), basis0...)

	// Sukuriame pradine simplex lentele (tableau)
	// - Eilutes: m apribojimu + 1 tikslo funkcijos eilute
	// - Stulpeliai: n kintamuju + 1 desines puses (RHS) stulpelis
	tableau = make([][]float64, m+1)
	for i := range m {
		row := make([]float64, n+1)
		copy(row, A[i])
		row[n] = b[i]
		tableau[i] = row
	}

	// Tikslo funkcija: i paskutine eilute dedame -c
	// Tada optimalumo paieskai ieskome neigiamu koeficientu tikslo eiluteje
	obj := make([]float64, n+1)
	for j := range n {
		obj[j] = -c[j]
	}
	tableau[m] = obj

	// isvedam pradine lentele
	hook(0, tableau, basis)

	for iter := 1; ; iter++ {
		// 1) Parenkame i baze ieinanti kintamaji
		// Pasirenkame stulpeli su labiausiai neigiamu koeficientu tikslo eiluteje
		enter := -1
		best := -eps
		for j := range n {
			if tableau[m][j] < best {
				best = tableau[m][j]
				enter = j
			}
		}

		// Jei nebera neigiamu koeficientu, pasiekta optimumo salyga
		if enter == -1 {
			break
		}

		// 2) Parenkame is bazes iseinanti kintamaji
		// Ziurime tik i teigiamus a_ij, kad isliktume leistinoje srityje
		leave := -1
		minRatio := 0.0
		for i := range m {
			a := tableau[i][enter]
			if a > eps {
				r := tableau[i][n] / a
				if leave == -1 || r < minRatio-eps {
					minRatio = r
					leave = i
				}
			}
		}

		// Pastaba: jei leave liktu -1, uzdavinys butu neapibreztas (unbounded)
		pivot := tableau[leave][enter]

		// 3) Pivot operacijos
		// Normalizuojame pivot eilute ir nunuliname enter stulpeli kitose eilutese
		for j := range tableau[leave] {
			tableau[leave][j] /= pivot
		}

		for i := range tableau {
			if i == leave {
				continue
			}
			factor := tableau[i][enter]
			if math.Abs(factor) < eps {
				continue
			}
			for j := range tableau[i] {
				tableau[i][j] -= factor * tableau[leave][j]
			}
		}

		// Atnaujiname baze: leave eilutes bazini kintamaji pakeiciame enter
		basis[leave] = kintamieji[enter]
		if hook != nil {
			hook(iter, tableau, basis)
		}
	}

	// Is tableles istraukiame sprendini: baziniai kintamieji lygus RHS, kiti – 0
	// Cia turime ir slack kintamuosius, bet graziname tik x1..x4
	xAll := make([]float64, n)
	for i := range m {

		col := -1
		for j, name := range kintamieji {
			if name == basis[i] {
				col = j
				break
			}
		}

		if col >= 0 {
			xAll[col] = tableau[i][n]
		}
	}

	x = append([]float64(nil), xAll[:4]...)
	return tableau[m][n], x, basis, tableau
}

func printTableau(iter int, tableau [][]float64, basis []string) {
	// Patogus simplex lenteles spausdinimas: tabwriter sulygiuoja stulpelius
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\niteracija %d\n", iter)

	fmt.Fprint(w, "\t")
	for _, name := range kintamieji {
		fmt.Fprintf(w, "%s\t", name)
	}
	fmt.Fprintln(w, "desine puse\t")

	m := len(tableau) - 1
	for i := range m {
		fmt.Fprintf(w, "%s\t", basis[i])
		for j := range tableau[i] {
			fmt.Fprintf(w, "%.6g\t", tableau[i][j])
		}
		fmt.Fprintln(w)
	}

	fmt.Fprint(w, "z\t")
	for j := range tableau[m] {
		fmt.Fprintf(w, "%.6g\t", tableau[m][j])
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w)
}

func main() {
	// Studento skaitmenys (naudojami 2 uzdavinyje kaip b vektorius)
	const A, B, C float64 = 9, 3, 9

	// Apribojimu matrica su laisvumo (slack) kintamaisiais
	// Stulpeliai atitinka: x1 x2 x3 x4 s1 s2 s3
	A_slack := [][]float64{
		{-1, 1, -1, -1, 1, 0, 0},
		{2, 4, 0, 0, 0, 1, 0},
		{0, 0, 1, 1, 0, 0, 1},
	}
	// Tikslo funkcijos koeficientai (slack kintamiesiems koeficientai = 0)
	c_slack := []float64{-2, 3, 0, 5, 0, 0, 0}
	// 1 uzdavinio desine puse
	b_orig := []float64{8, 10, 3}
	// 2 uzdavinio desine puse (pagal studento skaicius)
	b_ind := []float64{A, B, C}
	// Pradine baze: pradedame nuo slack kintamuju
	basis0 := []string{"s1", "s2", "s3"}

	// Kosmetika isvedimui (antrastes atskyrimui)

	l1padding, l2padding := strings.Repeat("=", 16), strings.Repeat("-", 16)

	fmt.Println(l1padding, "1 uzdavinys", l1padding)
	z1, x1, basis1, _ := simplex(
		A_slack,
		b_orig,
		c_slack,
		basis0,
		func(iter int, tableau [][]float64, basis []string) {
			printTableau(iter, tableau, basis)
		},
	)
	f1 := -z1
	fmt.Println(l2padding, "resultatas", l2padding)
	fmt.Printf("z_max = %.6g\n", z1)
	fmt.Printf("f_min = %.6g\n", f1)
	fmt.Printf("x* = [x1=%.6g, x2=%.6g, x3=%.6g, x4=%.6g]\n", x1[0], x1[1], x1[2], x1[3])
	fmt.Printf("galutine baze = %v\n\n", basis1)

	fmt.Println(l1padding, "2 uzdavinys", l1padding)
	z2, x2, basis2, _ := simplex(
		A_slack,
		b_ind,
		c_slack,
		basis0,
		func(iter int, tableau [][]float64, basis []string) {
			printTableau(iter, tableau, basis)
		},
	)
	f2 := -z2
	fmt.Println(l2padding, "resultatas", l2padding)
	fmt.Printf("z_max = %.6g\n", z2)
	fmt.Printf("f_min = %.6g\n", f2)
	fmt.Printf("x* = [x1=%.6g, x2=%.6g, x3=%.6g, x4=%.6g]\n", x2[0], x2[1], x2[2], x2[3])
	fmt.Printf("galutine baze = %v\n\n", basis2)

	fmt.Println(l1padding, "palyginimas", l1padding)
	fmt.Printf("Pirmas uzdavimys:  f_min=%.6g, x=[%.6g %.6g %.6g %.6g], baze=%v\n", f1, x1[0], x1[1], x1[2], x1[3], basis1)
	fmt.Printf("Antras uzdavinys:  f_min=%.6g, x=[%.6g %.6g %.6g %.6g], baze=%v\n", f2, x2[0], x2[1], x2[2], x2[3], basis2)
}
