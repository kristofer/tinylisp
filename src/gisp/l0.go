package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

// NaN boxing constants
const (
	ATOM = 0x7ff8
	PRIM = 0x7ff9
	CONS = 0x7ffa
	CLOS = 0x7ffb
	NIL  = 0x7ffc
	N    = 2048
)

type L float64
type I uint64

var (
	cell [N]L
	hp   I = 0
	sp   I = N
	A      = make([]byte, N*8)
	nilv L
	tru  L
	err  L
	env  L
)

// NaN boxing helpers
func box(t, i I) L {
	return L(math.Float64frombits(uint64(t)<<48 | uint64(i)))
}

func T(x L) I {
	return I(math.Float64bits(float64(x)) >> 48)
}

func ord(x L) I {
	return I(math.Float64bits(float64(x)) & 0xFFFFFFFFFFFF)
}

func equ(x, y L) bool {
	return math.Float64bits(float64(x)) == math.Float64bits(float64(y))
}

func ifv(cond L, alt L) L {
	if notv(cond) {
		return alt
	}
	return cond
}

func bind(v, t, e L) L {
	if notv(v) {
		return e
	} else if T(v) == CONS {
		return bind(cdr(v), cdr(t), pair(car(v), car(t), e))
	}
	return pair(v, t, e)
}

func eval(x, e L) L {
	if T(x) == ATOM {
		return assoc(x, e)
	} else if T(x) == CONS {
		return apply(eval(car(x), e), cdr(x), e)
	}
	return x
}

// Atom interning
func atom(s string) L {
	i := I(0)
	for i < hp {
		// Find the null terminator to get the correct string length
		end := i
		for end < hp && A[end] != 0 {
			end++
		}
		str := string(A[i:end])
		if str == s {
			return box(ATOM, i)
		}
		i = end + 1 // Move past the null terminator
	}
	// Not found, add new atom
	copy(A[hp:], s)
	A[hp+I(len(s))] = 0
	result := box(ATOM, hp)
	hp += I(len(s) + 1)
	if hp > sp<<3 {
		panic("out of memory")
	}
	return result
}

// Cons cell creation
func cons(x, y L) L {
	cell[sp-1] = x
	cell[sp-2] = y
	sp -= 2
	if hp > sp<<3 {
		panic("out of memory")
	}
	return box(CONS, sp)
}

// car and cdr
func car(p L) L {
	if T(p)&^(CONS^CLOS) == CONS {
		return cell[ord(p)+1]
	}
	return err
}

func cdr(p L) L {
	if T(p)&^(CONS^CLOS) == CONS {
		return cell[ord(p)]
	}
	return err
}

// pair, closure, assoc
func pair(v, x, e L) L {
	return cons(cons(v, x), e)
}

func closure(v, x, e L) L {
	if equ(e, env) {
		return box(CLOS, ord(pair(v, x, nilv)))
	}
	return box(CLOS, ord(pair(v, x, e)))
}

func assoc(v, e L) L {
	for T(e) == CONS && !equ(v, car(car(e))) {
		e = cdr(e)
	}
	if T(e) == CONS {
		return cdr(car(e))
	}
	return err
}

// not and let
func notv(x L) bool {
	return T(x) == NIL
}

func letv(x L) bool {
	return !notv(x) && !notv(cdr(x))
}

// evlis
//func eval(x, e L) L // forward declaration

func evlis(t, e L) L {
	if T(t) == CONS {
		return cons(eval(car(t), e), evlis(cdr(t), e))
	} else if T(t) == ATOM {
		return assoc(t, e)
	}
	return nilv
}

// Primitives
func f_add(t, e L) L {
	t = evlis(t, e)
	n := car(t)
	for {
		t = cdr(t)
		if notv(t) {
			break
		}
		n += car(t)
	}
	return n
}

func f_sub(t, e L) L {
	t = evlis(t, e)
	n := car(t)
	for {
		t = cdr(t)
		if notv(t) {
			break
		}
		n -= car(t)
	}
	return n
}

func f_mul(t, e L) L {
	t = evlis(t, e)
	n := car(t)
	for {
		t = cdr(t)
		if notv(t) {
			break
		}
		n *= car(t)
	}
	return n
}

func f_div(t, e L) L {
	t = evlis(t, e)
	n := car(t)
	for {
		t = cdr(t)
		if notv(t) {
			break
		}
		n /= car(t)
	}
	return n
}

// Additional primitives
func f_eval(t, e L) L {
	return eval(car(evlis(t, e)), e)
}

func f_quote(t, e L) L {
	return car(t)
}

func f_cons(t, e L) L {
	t = evlis(t, e)
	return cons(car(t), car(cdr(t)))
}

func f_car(t, e L) L {
	return car(car(evlis(t, e)))
}

func f_cdr(t, e L) L {
	return cdr(car(evlis(t, e)))
}

func f_int(t, e L) L {
	n := car(evlis(t, e))
	if n < 1e16 && n > -1e16 {
		return L(int64(n))
	}
	return n
}

func f_lt(t, e L) L {
	t = evlis(t, e)
	if car(t)-car(cdr(t)) < 0 {
		return tru
	}
	return nilv
}

func f_eq(t, e L) L {
	t = evlis(t, e)
	if equ(car(t), car(cdr(t))) {
		return tru
	}
	return nilv
}

func f_pair(t, e L) L {
	x := car(evlis(t, e))
	if T(x) == CONS {
		return tru
	}
	return nilv
}

func f_or(t, e L) L {
	x := nilv
	for !notv(t) {
		x = eval(car(t), e)
		if !notv(x) {
			break
		}
		t = cdr(t)
	}
	return x
}

func f_and(t, e L) L {
	x := tru
	for !notv(t) {
		x = eval(car(t), e)
		if notv(x) {
			break
		}
		t = cdr(t)
	}
	return x
}

func f_not(t, e L) L {
	if notv(car(evlis(t, e))) {
		return tru
	}
	return nilv
}

func f_cond(t, e L) L {
	for notv(eval(car(car(t)), e)) {
		t = cdr(t)
	}
	return eval(car(cdr(car(t))), e)
}

func f_if(t, e L) L {
	if notv(eval(car(t), e)) {
		return eval(car(cdr(cdr(t))), e) // false condition -> else branch
	}
	return eval(car(cdr(t)), e) // true condition -> then branch
}

func f_leta(t, e L) L {
	for letv(t) {
		e = pair(car(car(t)), eval(car(cdr(car(t))), e), e)
		t = cdr(t)
	}
	return eval(car(t), e)
}

func f_lambda(t, e L) L {
	return closure(car(t), car(cdr(t)), e)
}

func f_define(t, e L) L {
	env = pair(car(t), eval(car(cdr(t)), e), env)
	return car(t)
}

// Load Lisp code from a file
func loadFile(filename string, env L) L {
	content, err := os.ReadFile(filename)
	if err != nil {
		return atom("FILE-ERROR")
	}
	
	input := string(content)
	parser := newInputParser(input)
	var result L = nilv
	
	// Parse and evaluate each expression in the file
	for parser.ch != 0 {
		parser.skipWhitespace()
		if parser.ch == 0 {
			break
		}
		
		expr := parser.readExpr()
		if T(expr) == ATOM && equ(expr, atom("ERR")) {
			return atom("PARSE-ERROR")
		}
		
		result = eval(expr, env)
		if T(result) == ATOM && equ(result, atom("ERR")) {
			return atom("EVAL-ERROR")
		}
	}
	
	return result
}

// Primitive wrapper for loadFile
func f_load(t, e L) L {
	// Get the filename argument
	args := evlis(t, e)
	if notv(args) {
		return atom("MISSING-FILENAME")
	}
	
	// Extract filename string from atom
	filenameAtom := car(args)
	if T(filenameAtom) != ATOM {
		return atom("INVALID-FILENAME")
	}
	
	i := ord(filenameAtom)
	filename := ""
	for j := i; A[j] != 0; j++ {
		filename += string(A[j])
	}
	
	
	// Load and evaluate the file using global environment
	return loadFile(filename, env)
}

// Primitive table with all primitives
var prims map[string]func(L, L) L

// Store primitive index mapping
var primIndex map[string]I

// Corrected parser that handles EOF gracefully and fixes atom string extraction
type inputParser struct {
	input string
	pos   int
	ch    byte
}

func newInputParser(input string) *inputParser {
	p := &inputParser{input: input, pos: 0}
	p.next()
	return p
}

func (p *inputParser) next() {
	if p.pos < len(p.input) {
		p.ch = p.input[p.pos]
		p.pos++
	} else {
		p.ch = 0 // EOF
	}
}

func (p *inputParser) skipWhitespace() {
	for p.ch > 0 && p.ch <= ' ' {
		p.next()
	}
}

func (p *inputParser) readAtom() L {
	start := p.pos - 1 // Start at current character position
	// Keep reading while we have valid atom characters
	for p.ch > ' ' && p.ch != '(' && p.ch != ')' && p.ch != '\'' && p.ch != 0 {
		p.next()
	}
	// Extract the atom string - end is where we stopped
	end := p.pos
	if p.ch != 0 {
		end = p.pos - 1 // Back up one if we stopped on a delimiter
	}
	s := p.input[start:end]

	// Try to parse as number using strconv
	if len(s) > 0 {
		if n, err := strconv.ParseFloat(s, 64); err == nil {
			return L(n)
		}
	}

	// Return as atom
	return atom(s)
}

func (p *inputParser) readList() L {
	p.skipWhitespace()
	if p.ch == ')' {
		p.next()
		return nilv
	}
	if p.ch == 0 {
		return nilv
	}

	// Handle dot notation
	if p.ch == '.' {
		p.next()
		p.skipWhitespace()
		result := p.readExpr()
		p.skipWhitespace()
		if p.ch == ')' {
			p.next()
		}
		return result
	}

	first := p.readExpr()
	rest := p.readList()
	return cons(first, rest)
}

func (p *inputParser) readExpr() L {
	p.skipWhitespace()

	switch p.ch {
	case 0:
		return nilv
	case '(':
		p.next()
		return p.readList()
	case '\'':
		p.next()
		return cons(atom("quote"), cons(p.readExpr(), nilv))
	default:
		return p.readAtom()
	}
}

// Error handling for primitives
func apply(f, t, e L) L {
	if T(f) == PRIM {
		// Use the ordinal to look up the primitive directly
		primOrd := ord(f)
		for name, fn := range prims {
			if primIndex[name] == primOrd {
				return fn(t, e)
			}
		}
		return err
	} else if T(f) == CLOS {
		return eval(cdr(car(f)), bind(car(car(f)), evlis(t, e), ifv(cdr(f), env)))
	}
	return err
}

// Print function with type detection
func printExpr(x L) {
	switch T(x) {
	case NIL:
		fmt.Print("()")
	case ATOM:
		i := ord(x)
		s := ""
		for j := i; A[j] != 0; j++ {
			s += string(A[j])
		}
		fmt.Print(s)
	case PRIM:
		fmt.Print("<primitive>")
	case CONS:
		printlist(x)
	case CLOS:
		fmt.Printf("{closure %d}", ord(x))
	default:
		fmt.Printf("%.10g", float64(x))
	}
}

func printlist(t L) {
	fmt.Print("(")
	first := true
	for {
		if !first {
			fmt.Print(" ")
		}
		printExpr(car(t))
		t = cdr(t)
		if notv(t) {
			break
		}
		if T(t) != CONS {
			fmt.Print(" . ")
			printExpr(t)
			break
		}
		first = false
	}
	fmt.Print(")")
}

func gc() {
	sp = ord(env)
}

var (
	buf      = make([]byte, 40)
	see byte = ' '
	rdr      = bufio.NewReader(os.Stdin)
)

func look() {
	c, err := rdr.ReadByte()
	if err != nil {
		os.Exit(0)
	}
	see = c
}

func seeing(c byte) bool {
	if c == ' ' {
		return see > 0 && see <= c
	}
	return see == c
}

func get() byte {
	c := see
	look()
	return c
}

func scan() byte {
	i := 0
	for seeing(' ') {
		look()
	}
	if seeing('(') || seeing(')') || seeing('\'') {
		buf[i] = get()
		i++
	} else {
		for i < 39 && !seeing('(') && !seeing(')') && !seeing(' ') {
			buf[i] = get()
			i++
		}
	}
	buf[i] = 0
	return buf[0]
}

func Read() L {
	scan()
	return parse()
}

func list() L {
	var x L
	if scan() == ')' {
		return nilv
	}
	if string(buf[:1]) == "." {
		x = Read()
		scan()
		return x
	}
	x = parse()
	return cons(x, list())
}

func quoteExpr() L {
	return cons(atom("quote"), cons(Read(), nilv))
}

func atomic() L {
	s := string(buf[:])
	n, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return L(n)
	}
	return atom(s)
}

func parse() L {
	switch buf[0] {
	case '(':
		return list()
	case '\'':
		return quoteExpr()
	default:
		return atomic()
	}
}

// Example usage
func main() {
	prims = map[string]func(L, L) L{
		"eval":   f_eval,
		"quote":  f_quote,
		"cons":   f_cons,
		"car":    f_car,
		"cdr":    f_cdr,
		"+":      f_add,
		"-":      f_sub,
		"*":      f_mul,
		"/":      f_div,
		"int":    f_int,
		"<":      f_lt,
		"eq?":    f_eq,
		"pair?":  f_pair,
		"or":     f_or,
		"and":    f_and,
		"not":    f_not,
		"cond":   f_cond,
		"if":     f_if,
		"let*":   f_leta,
		"lambda": f_lambda,
		"define": f_define,
		"load":   f_load,
	}
	fmt.Println("tinylisp")
	nilv = box(NIL, 0)
	err = atom("ERR")
	tru = atom("#t")
	env = pair(tru, tru, nilv)

	// Initialize primitive index mapping with deterministic order
	primIndex = make(map[string]I)
	primOrd := I(0)

	// Use a slice to ensure deterministic order
	primNames := []string{
		"eval", "quote", "cons", "car", "cdr", "+", "-", "*", "/", "int",
		"<", "eq?", "pair?", "or", "and", "not", "cond", "if", "let*", "lambda", "define", "load",
	}

	for _, name := range primNames {
		if _, exists := prims[name]; exists {
			primIndex[name] = primOrd
			env = pair(atom(name), box(PRIM, primOrd), env)
			primOrd++
		}
	}
	// REPL using safer input handling
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("\n%d> ", int(sp)-int(hp)/8)
		if !scanner.Scan() {
			break // EOF or error
		}

		input := scanner.Text()
		if input == "" {
			continue // Skip empty lines
		}

		// Parse using our corrected parser logic (same as tests)
		parser := &inputParser{input: input, pos: 0}
		parser.next()
		expr := parser.readExpr()

		// Evaluate and print
		result := eval(expr, env)
		printExpr(result)
		gc()
	}
}
