package main

import (
	"testing"
)

func TestEnvironmentDebug(t *testing.T) {
	initTinyLisp()
	
	// Test atom interning first
	t.Log("Testing atom interning:")
	plus1 := atom("+")
	plus2 := atom("+")
	t.Logf("atom('+') first call: tag=%x ord=%d", T(plus1), ord(plus1))
	t.Logf("atom('+') second call: tag=%x ord=%d", T(plus2), ord(plus2))
	t.Logf("Are they equal? %v", equ(plus1, plus2))
	t.Logf("Current hp: %d", hp)
	
	// Test specific lookups  
	symbols := []string{"+", "-", "*", "/", "define", "car", "cdr"}
	for _, sym := range symbols {
		atomSym := atom(sym)
		value := assoc(atomSym, env)
		if T(value) != NIL {
			t.Logf("assoc(%s): atom ord=%d -> tag=%x ord=%d", sym, ord(atomSym), T(value), ord(value))
		} else {
			t.Logf("assoc(%s): NOT FOUND", sym)
		}
	}
}