package main

import (
	"math"
	"testing"
)

func TestNaNBoxingDebug(t *testing.T) {
	numbers := []float64{0, 1, -1, 3.14159, 1e10, -1e10}
	
	t.Log("Analyzing IEEE 754 bit patterns:")
	for _, n := range numbers {
		bits := math.Float64bits(n)
		tag := bits >> 48
		
		t.Logf("Number %f:", n)
		t.Logf("  Bits: %016x", bits)
		t.Logf("  Tag (top 16): %04x", tag)
		t.Logf("  ATOM constant: %04x", ATOM)
		t.Logf("  Tag >= ATOM: %t", tag >= uint64(ATOM))
		t.Logf("  Is finite: %t", math.IsInf(n, 0) == false && math.IsNaN(n) == false)
		t.Log("")
	}
	
	// Test what makes a valid NaN for boxing
	t.Log("NaN patterns for boxing:")
	nanAtom := box(ATOM, 123)
	nanPrim := box(PRIM, 456)
	
	t.Logf("box(ATOM, 123): bits=%016x, tag=%04x", math.Float64bits(float64(nanAtom)), T(nanAtom))
	t.Logf("box(PRIM, 456): bits=%016x, tag=%04x", math.Float64bits(float64(nanPrim)), T(nanPrim))
}