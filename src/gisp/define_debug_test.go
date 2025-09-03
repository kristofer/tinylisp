package main

import (
	"testing"
)

func TestDefineDebug(t *testing.T) {
	initTinyLisp()
	
	// Test atom interning consistency
	t.Log("Testing atom interning consistency:")
	x1 := atom("x")
	x2 := atom("x")
	t.Logf("atom('x') first call: ord=%d", ord(x1))
	t.Logf("atom('x') second call: ord=%d", ord(x2))
	t.Logf("Are they equal? %v", equ(x1, x2))
	
	// Parse and evaluate (define x 42)
	parser1 := newTestParser("(define x 42)")
	defineExpr := parser1.readExpr()
	t.Logf("Parsed define expression: tag=%x", T(defineExpr))
	
	// Check what x is in the parsed expression
	defineList := defineExpr
	if T(defineList) == CONS {
		second := car(cdr(defineList)) // Should be x
		t.Logf("x in define expression: tag=%x ord=%d", T(second), ord(second))
	}
	
	result1 := eval(defineExpr, env)
	t.Logf("Define result: tag=%x ord=%d", T(result1), ord(result1))
	
	// Check if x is now in the environment using the SAME atom
	value := assoc(x1, env)
	t.Logf("Looking up x with atom ord=%d: result tag=%x value=%f", ord(x1), T(value), float64(value))
	
	// Parse and evaluate just x
	parser2 := newTestParser("x")
	varExpr := parser2.readExpr()
	t.Logf("Parsed variable expression: tag=%x ord=%d", T(varExpr), ord(varExpr))
	
	result2 := eval(varExpr, env)
	t.Logf("Variable evaluation result: tag=%x, value=%f", T(result2), float64(result2))
}