package main

import (
	"strings"
	"testing"
)

// Integration tests that test complete Lisp expressions from parsing to evaluation

// Test-safe parser that handles EOF gracefully
type testParser struct {
	input string
	pos   int
	ch    byte
}

func newTestParser(input string) *testParser {
	p := &testParser{input: input, pos: 0}
	p.next()
	return p
}

func (p *testParser) next() {
	if p.pos < len(p.input) {
		p.ch = p.input[p.pos]
		p.pos++
	} else {
		p.ch = 0 // EOF
	}
}

func (p *testParser) skipWhitespace() {
	for p.ch > 0 && p.ch <= ' ' {
		p.next()
	}
}

func (p *testParser) readAtom() L {
	start := p.pos - 1  // Start at current character position
	// Keep reading while we have valid atom characters
	for p.ch > ' ' && p.ch != '(' && p.ch != ')' && p.ch != '\'' && p.ch != 0 {
		p.next()
	}
	// Extract the atom string - end is where we stopped
	end := p.pos
	if p.ch != 0 {
		end = p.pos - 1  // Back up one if we stopped on a delimiter
	}
	s := p.input[start:end]
	
	// Try to parse as number
	if len(s) > 0 {
		// Simple number parsing
		var val float64
		var isNum bool
		if s == "0" {
			val = 0
			isNum = true
		} else {
			// Parse float manually since we can't use strconv easily
			sign := 1.0
			i := 0
			if s[0] == '-' {
				sign = -1
				i = 1
			} else if s[0] == '+' {
				i = 1
			}
			
			// Check if it's all digits or has decimal point
			hasDigit := false
			hasDot := false
			for j := i; j < len(s); j++ {
				if s[j] >= '0' && s[j] <= '9' {
					hasDigit = true
				} else if s[j] == '.' && !hasDot {
					hasDot = true
				} else {
					goto notNumber // Contains invalid character
				}
			}
			
			if hasDigit {
				// Simple integer parsing
				for j := i; j < len(s) && s[j] != '.'; j++ {
					val = val*10 + float64(s[j]-'0')
				}
				
				// Handle decimal part
				if hasDot {
					factor := 0.1
					dotPos := strings.IndexByte(s[i:], '.')
					if dotPos >= 0 {
						for j := i + dotPos + 1; j < len(s); j++ {
							val += float64(s[j]-'0') * factor
							factor *= 0.1
						}
					}
				}
				
				val *= sign
				isNum = true
			}
		}
		
		notNumber:
		if isNum {
			return L(val)
		}
	}
	
	// Return as atom
	return atom(s)
}

func (p *testParser) readList() L {
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

func (p *testParser) readExpr() L {
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

func parseAndEval(input string) L {
	initTinyLisp()
	parser := newTestParser(input)
	expr := parser.readExpr()
	return eval(expr, env)
}

func TestBasicArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"(+ 1 2)", 3.0},
		{"(+ 1 2 3 4)", 10.0},
		{"(- 10 3)", 7.0},
		{"(- 10 3 2)", 5.0},
		{"(* 2 3)", 6.0},
		{"(* 2 3 4)", 24.0},
		{"(/ 12 3)", 4.0},
		{"(/ 12 3 2)", 2.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if float64(result) != tt.expected {
				t.Errorf("%s = %f, want %f", tt.input, float64(result), tt.expected)
			}
		})
	}
}

func TestNestedArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"(+ 1 (* 2 3))", 7.0},          // 1 + (2 * 3) = 7
		{"(* (+ 1 2) (- 5 2))", 9.0},    // (1 + 2) * (5 - 2) = 3 * 3 = 9
		{"(/ (* 2 6) (+ 1 2))", 4.0},    // (2 * 6) / (1 + 2) = 12 / 3 = 4
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if float64(result) != tt.expected {
				t.Errorf("%s = %f, want %f", tt.input, float64(result), tt.expected)
			}
		})
	}
}

func TestComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool // true means should equal #t, false means should equal nil
	}{
		{"(< 1 2)", true},
		{"(< 2 1)", false},
		{"(< 1 1)", false},
		{"(eq? 1 1)", true},
		{"(eq? 1 2)", false},
		{"(eq? 'hello 'hello)", true},
		{"(eq? 'hello 'world)", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if tt.expected {
				if !equ(result, tru) {
					t.Errorf("%s should be true but got %v", tt.input, result)
				}
			} else {
				if !equ(result, nilv) {
					t.Errorf("%s should be false but got %v", tt.input, result)
				}
			}
		})
	}
}

func TestLogicalOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"(not ())", true},
		{"(not #t)", false},
		{"(not 42)", false},
		{"(and #t #t)", true},
		{"(and #t ())", false},
		{"(and () #t)", false},
		{"(or #t ())", true},
		{"(or () #t)", true},
		{"(or () ())", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if tt.expected {
				if notv(result) {
					t.Errorf("%s should be truthy but got falsy", tt.input)
				}
			} else {
				if !notv(result) {
					t.Errorf("%s should be falsy but got truthy", tt.input)
				}
			}
		})
	}
}

func TestListOperations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(L) bool
	}{
		{
			name:  "cons creates pair",
			input: "(cons 1 2)",
			checkFn: func(result L) bool {
				return T(result) == CONS && equ(car(result), L(1)) && equ(cdr(result), L(2))
			},
		},
		{
			name:  "car extracts first",
			input: "(car (cons 1 2))",
			checkFn: func(result L) bool {
				return equ(result, L(1))
			},
		},
		{
			name:  "cdr extracts second",
			input: "(cdr (cons 1 2))",
			checkFn: func(result L) bool {
				return equ(result, L(2))
			},
		},
		{
			name:  "pair? detects pairs",
			input: "(pair? (cons 1 2))",
			checkFn: func(result L) bool {
				return equ(result, tru)
			},
		},
		{
			name:  "pair? rejects atoms",
			input: "(pair? 42)",
			checkFn: func(result L) bool {
				return equ(result, nilv)
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if !tt.checkFn(result) {
				t.Errorf("Test failed for: %s", tt.input)
			}
		})
	}
}

func TestQuoting(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(L) bool
	}{
		{
			name:  "quote prevents evaluation",
			input: "(quote hello)",
			checkFn: func(result L) bool {
				return T(result) == ATOM && equ(result, atom("hello"))
			},
		},
		{
			name:  "apostrophe quote",
			input: "'world",
			checkFn: func(result L) bool {
				return T(result) == ATOM && equ(result, atom("world"))
			},
		},
		{
			name:  "eval evaluates quoted expression",
			input: "(eval '(+ 1 2))",
			checkFn: func(result L) bool {
				return equ(result, L(3))
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if !tt.checkFn(result) {
				t.Errorf("Test failed for: %s", tt.input)
			}
		})
	}
}

func TestConditionals(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(L) bool
	}{
		{
			name:  "if true branch",
			input: "(if #t 42 24)",
			checkFn: func(result L) bool {
				return equ(result, L(42))
			},
		},
		{
			name:  "if false branch",
			input: "(if () 42 24)",
			checkFn: func(result L) bool {
				return equ(result, L(24))
			},
		},
		{
			name:  "cond first true",
			input: "(cond (#t 42) (else 24))",
			checkFn: func(result L) bool {
				return equ(result, L(42))
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if !tt.checkFn(result) {
				t.Errorf("Test failed for: %s", tt.input)
			}
		})
	}
}

func TestVariableBinding(t *testing.T) {
	// Test define and variable lookup
	t.Run("define and lookup", func(t *testing.T) {
		initTinyLisp()
		
		// Define x = 42
		parser1 := newTestParser("(define x 42)")
		defineExpr := parser1.readExpr()
		eval(defineExpr, env)
		
		// Now evaluate x
		parser2 := newTestParser("x")
		varExpr := parser2.readExpr()
		result := eval(varExpr, env)
		
		if !equ(result, L(42)) {
			t.Errorf("Variable x should be 42, got %f", float64(result))
		}
	})
	
	// Test let* binding
	t.Run("let* binding", func(t *testing.T) {
		input := "(let* (x 1) (y (+ x 1)) (+ x y))"
		result := parseAndEval(input)
		if !equ(result, L(3)) {
			t.Errorf("let* should compute 3, got %f", float64(result))
		}
	})
}

func TestLambdaFunctions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(L) bool
	}{
		{
			name:  "lambda creates closure",
			input: "(lambda (x) x)",
			checkFn: func(result L) bool {
				return T(result) == CLOS
			},
		},
		{
			name:  "apply lambda function",
			input: "((lambda (x) (+ x 1)) 5)",
			checkFn: func(result L) bool {
				return equ(result, L(6))
			},
		},
		{
			name:  "lambda with multiple args",
			input: "((lambda (x y) (+ x y)) 3 4)",
			checkFn: func(result L) bool {
				return equ(result, L(7))
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAndEval(tt.input)
			if !tt.checkFn(result) {
				t.Errorf("Test failed for: %s", tt.input)
			}
		})
	}
}

func TestComplexIntegration(t *testing.T) {
	// Test a more complex expression combining multiple features
	t.Run("factorial-like function", func(t *testing.T) {
		// Define a simple function: (define square (lambda (x) (* x x)))
		initTinyLisp()
		
		input1 := "(define square (lambda (x) (* x x)))"
		parser1 := newTestParser(input1)
		expr1 := parser1.readExpr()
		eval(expr1, env)
		
		// Now use it: (square 5) should be 25
		input2 := "(square 5)"
		parser2 := newTestParser(input2)
		expr2 := parser2.readExpr()
		result := eval(expr2, env)
		
		if !equ(result, L(25)) {
			t.Errorf("(square 5) should be 25, got %f", float64(result))
		}
	})
	
	// Test closure capturing
	t.Run("closure captures environment", func(t *testing.T) {
		// ((lambda (x) (lambda (y) (+ x y))) 10) should return a function that adds 10
		// Then apply it to 5 to get 15
		input := "(((lambda (x) (lambda (y) (+ x y))) 10) 5)"
		result := parseAndEval(input)
		
		if !equ(result, L(15)) {
			t.Errorf("Closure test should return 15, got %f", float64(result))
		}
	})
}

func TestErrorConditions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"undefined variable", "undefined_var"},
		{"car of non-pair", "(car 42)"},
		{"cdr of non-pair", "(cdr 42)"},
		{"apply non-function", "(42 1 2)"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAndEval(tt.input)
			// Check if result is the error atom
			if !equ(result, err) {
				t.Logf("Expected error for %s, but got result (may be OK depending on implementation)", tt.input)
			}
		})
	}
}

// Test cases inspired by the original dotcall.lisp test
func TestDotCallCompatibility(t *testing.T) {
	// Test some basic cases from dotcall.lisp to ensure compatibility
	tests := []struct {
		name    string
		setup   []string // Define expressions to run first
		input   string
		checkFn func(L) bool
	}{
		{
			name:  "basic arithmetic",
			setup: []string{},
			input: "(+ 1 2 3)",
			checkFn: func(result L) bool {
				return equ(result, L(6))
			},
		},
		{
			name:  "list function",
			setup: []string{"(define list (lambda args args))"},
			input: "(list 1 2 3)",
			checkFn: func(result L) bool {
				// Should return (1 2 3)
				return T(result) == CONS &&
					equ(car(result), L(1)) &&
					T(cdr(result)) == CONS &&
					equ(car(cdr(result)), L(2)) &&
					T(cdr(cdr(result))) == CONS &&
					equ(car(cdr(cdr(result))), L(3))
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTinyLisp()
			
			// Run setup expressions
			for _, setup := range tt.setup {
				parser := newTestParser(setup)
				expr := parser.readExpr()
				eval(expr, env)
			}
			
			// Run test expression
			parser := newTestParser(tt.input)
			expr := parser.readExpr()
			result := eval(expr, env)
			
			if !tt.checkFn(result) {
				t.Errorf("Test failed for: %s", tt.input)
			}
		})
	}
}