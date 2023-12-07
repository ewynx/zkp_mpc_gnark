package main

import (
	"github.com/consensys/gnark/frontend"
)

// Have to prove the final shared value is indeed the sum of the received shares

type Shares struct {
	X1 frontend.Variable
	X2 frontend.Variable
	X3 frontend.Variable
}

type MpcAddCircuit struct {
	// The shares that the party received
	C Shares // this is private
	// The final value the party will distribute
	Y frontend.Variable `gnark:",public"` // (default is private)
}

// At compile time, frontend.Compile(...) recursively parses the struct fields that contains frontend.Variable to build the frontend.constraintSystem.
func (circuit *MpcAddCircuit) Define(api frontend.API) error {
	sum := api.Add(circuit.C.X1, circuit.C.X2, circuit.C.X3)
	api.AssertIsEqual(circuit.Y, sum)
	return nil
}
