package main

import (
	"testing"
)

func TestParserAtomDebug(t *testing.T) {
	initTinyLisp()
	
	// Test parsing just "x"
	parser := newTestParser("x")
	result := parser.readExpr()
	
	t.Logf("Parsed 'x': tag=%x ord=%d", T(result), ord(result))
	
	// Extract the atom string to see what it actually is
	if T(result) == ATOM {
		i := ord(result)
		end := i
		for end < hp && A[end] != 0 {
			end++
		}
		str := string(A[i:end])
		t.Logf("Atom string: '%s' (length=%d)", str, len(str))
		
		// Compare to direct atom call
		direct := atom("x")
		t.Logf("Direct atom('x'): ord=%d", ord(direct))
		
		// Show hex bytes to see if there are any hidden characters
		for j := i; j < end; j++ {
			t.Logf("Byte at %d: 0x%02x ('%c')", j-i, A[j], A[j])
		}
	}
}