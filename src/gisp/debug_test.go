package main

import (
	"testing"
)

func TestDebugEvaluation(t *testing.T) {
	initTinyLisp()
	
	// Test looking up the + symbol
	plusSym := atom("+")
	plusVal := assoc(plusSym, env)
	
	t.Logf("+ symbol: tag=%x, ord=%d", T(plusSym), ord(plusSym))
	t.Logf("+ value from env: tag=%x, ord=%d", T(plusVal), ord(plusVal))
	t.Logf("+ should have primIndex = %d", primIndex["+"])
	
	// Check if it's a PRIM
	if T(plusVal) == PRIM {
		t.Log("+ is correctly identified as PRIM")
	} else {
		t.Errorf("+ should be PRIM, got tag %x", T(plusVal))
	}
	
	// Test the apply function manually
	args := cons(L(1), cons(L(2), nilv)) // (1 2)
	result := apply(plusVal, args, env)
	
	t.Logf("apply result: %f (tag=%x)", float64(result), T(result))
	
	// Test if it's NaN
	if result != result { // NaN check
		t.Error("Result is NaN - apply function failed")
	}
}

func TestPrimitiveStorageDebug(t *testing.T) {
	initTinyLisp()
	
	// Check how primitives are stored
	for name := range prims {
		sym := atom(name)
		val := assoc(sym, env)
		t.Logf("Primitive %s: sym tag=%x ord=%d, val tag=%x ord=%d", 
			name, T(sym), ord(sym), T(val), ord(val))
	}
}

func TestApplyDebug(t *testing.T) {
	initTinyLisp()
	
	plusVal := assoc(atom("+"), env)
	
	// Debug the apply function
	if T(plusVal) == PRIM {
		t.Logf("Found PRIM with ordinal %d, checking primitive lookup...", ord(plusVal))
		
		// Show what the primitive index mapping looks like
		for name, index := range primIndex {
			t.Logf("primIndex[%s] = %d", name, index)
		}
		
		// This is what apply() does now
		primOrd := ord(plusVal)
		for name, fn := range prims {
			if primIndex[name] == primOrd {
				t.Logf("MATCH FOUND: %s has index %d, matches prim ordinal %d", name, primIndex[name], primOrd)
				
				// Try calling the function directly
				args := cons(L(1), cons(L(2), nilv))
				result := fn(args, env)
				t.Logf("Direct call to %s function: %f", name, float64(result))
				break
			}
		}
	}
}