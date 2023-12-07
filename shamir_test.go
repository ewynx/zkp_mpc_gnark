package main

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"testing"
)

func TestSecretRecovery1(t *testing.T) {
	// 3 shares with threshold 1
	secret := GenerateSecret()

	n := 3         // Number of-  shares
	threshold := 1 // Threshold (degree of polynomial)

	// Generate inputs (1..n)
	inputs := make([]uint64, n)
	for i := range inputs {
		inputs[i] = uint64(i + 1)
	}

	// Obtain shares
	shares := GetSharesSecret(secret, inputs, threshold)
	for _, share := range shares {
		fmt.Printf("share: %v\n", share)
	}

	// Recover the polynomial using interpolation
	recoveredVal := Interpolate(shares)
	fmt.Printf("Interpolation result: %s\n", recoveredVal.String())

	// Check if the constant term of the recovered polynomial is the secret
	if !recoveredVal.Equal(&secret) {
		t.Errorf("The recovered secret does not match the original secret")
	}
}

func TestSecretRecovery2(t *testing.T) {
	// 15 shares with threshold 4
	secret := GenerateSecret()

	n := 15        // Number of shares
	threshold := 4 // Threshold (degree of polynomial)

	inputs := make([]uint64, n)
	for i := range inputs {
		inputs[i] = uint64(i + 1)
	}

	// Obtain shares
	shares := GetSharesSecret(secret, inputs, threshold)

	// Recover the polynomial
	recoveredVal := Interpolate(shares)

	// Check if the constant term of the recovered polynomial is the secret
	if !recoveredVal.Equal(&secret) {
		t.Errorf("The recovered secret does not match the original secret")
	}
}

func TestInterpolate(t *testing.T) {
	points := []Point{
		{X: newFrElement(1), Y: newFrElement(12)},
		{X: newFrElement(2), Y: newFrElement(7)},
		{X: newFrElement(3), Y: newFrElement(9)},
	}

	result := Interpolate(points)
	expectedResult := newFrElement(24)

	if !result.Equal(&expectedResult) {
		t.Errorf("Interpolation result is incorrect")
	}
}

// Helper function to create a new Element with a given uint64 value
func newFrElement(val uint64) fr.Element {
	var el fr.Element
	el.SetUint64(val)
	return el
}
