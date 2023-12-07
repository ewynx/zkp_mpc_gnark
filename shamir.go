package main

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// p = 21888242871839275222246405745257275088548364400416034343698204186575808495617

type Polynomial struct {
	Coefficients []fr.Element
	Degree       int
}

type Point struct {
	X fr.Element
	Y fr.Element
}

// return result after evaluation polynomial at x = xval
func (p Polynomial) Eval(x fr.Element) fr.Element {
	var res, xval, temp fr.Element
	res.SetZero()
	xval.SetOne()

	for _, coeff := range p.Coefficients {
		temp.Mul(&coeff, &xval)
		res.Add(&res, &temp)
		xval.Mul(&xval, &x)
	}
	return res
}

func GenerateSecret() fr.Element {
	var s fr.Element
	s.SetRandom()
	return s
}

func GetSharesSecret(s fr.Element, inputs []uint64, t int) []Point {
	// 1. Generate polynomial of degree t
	p := CreatePol(s, t)

	// 2. Evaluate the polynomial at the given inputs
	res := make([]Point, len(inputs))
	for i, input := range inputs {
		var input_el fr.Element
		input_el.SetUint64(input)
		y := p.Eval(input_el)
		res[i] = Point{X: input_el, Y: y}
	}
	return res
}

// returns a polynomial of  degree t with random coefficients c_i, except for c_0=s
func CreatePol(s fr.Element, t int) Polynomial {
	p := Polynomial{Coefficients: make([]fr.Element, t), Degree: t}
	// Set the secret at the constant coefficient
	p.Coefficients[0].Set(&s)
	// Then set the other coefficients with random values
	for i := 1; i < t; i++ {
		p.Coefficients[i].SetRandom()
	}
	return p
}

func Interpolate(points []Point) fr.Element {
	var result fr.Element
	result.SetZero()

	for i := 0; i < len(points); i++ {
		var term fr.Element
		term.Set(&points[i].Y)

		for j := 0; j < len(points); j++ {
			if i != j {
				// Ensure that the x values are distinct
				if points[i].X.Equal(&points[j].X) {
					panic("Interpolation error: duplicate x values")
				}

				// Calculate the denominator (xs[i] - xs[j])
				var denominator fr.Element
				denominator.Sub(&points[i].X, &points[j].X)

				// Calculate the inverse of the denominator
				var inv fr.Element
				inv.Inverse(&denominator)

				// Calculate the numerator (-xs[j])
				var numerator fr.Element
				numerator.Neg(&points[j].X)

				// Multiply term by (numerator * inv)
				var temp fr.Element
				temp.Mul(&numerator, &inv)
				term.Mul(&term, &temp)

			}
		}
		result.Add(&result, &term)
	}

	return result
}
