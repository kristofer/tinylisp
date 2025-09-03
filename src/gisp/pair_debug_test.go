package main

import (
	"testing"
)

func TestPairDebug(t *testing.T) {
	initTinyLisp()
	
	// Create a pair
	pair := cons(L(1), L(2))
	t.Logf("Created pair: tag=%x ord=%d", T(pair), ord(pair))
	
	// Test the car and cdr
	first := car(pair)
	second := cdr(pair)
	t.Logf("car(pair): %f", float64(first))
	t.Logf("cdr(pair): %f", float64(second))
	
	// Create the arguments for f_pair
	args := cons(pair, nilv)
	t.Logf("Args: tag=%x ord=%d", T(args), ord(args))
	
	// Test evlis on the args
	evaled := evlis(args, env)
	t.Logf("evlis(args): tag=%x ord=%d", T(evaled), ord(evaled))
	
	// Test car of evaled
	first_evaled := car(evaled)
	t.Logf("car(evlis(args)): tag=%x ord=%d", T(first_evaled), ord(first_evaled))
	
	// Test the f_pair function
	result := f_pair(args, env)
	t.Logf("f_pair result: tag=%x ord=%d", T(result), ord(result))
	t.Logf("tru: tag=%x ord=%d", T(tru), ord(tru))
	t.Logf("nilv: tag=%x ord=%d", T(nilv), ord(nilv))
	t.Logf("result == tru: %t", equ(result, tru))
	t.Logf("result == nilv: %t", equ(result, nilv))
}