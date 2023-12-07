package main

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/test"
	"testing"
)

func TestMpcAddCircuit(t *testing.T) {
	var x1, x2, x3, temp, y fr.Element
	x1.SetRandom()
	x2.SetRandom()
	x3.SetRandom()
	temp.Add(&x1, &x2)
	y.Add(&x3, &temp)

	assert := test.NewAssert(t)
	shares := Shares{
		X1: x1,
		X2: x2,
		X3: x3,
	}
	var mpcAddCircuit MpcAddCircuit
	assert.ProverSucceeded(&mpcAddCircuit, &MpcAddCircuit{
		C: shares,
		Y: y,
	}, test.WithCurves(ecc.BN254))
}
