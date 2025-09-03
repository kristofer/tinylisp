# TinyLisp Go Implementation - Current State

## Project Overview
This is a Go port of tinylisp, a minimal Lisp interpreter originally written in C. The Go implementation is located in `/Users/kristofer/LocalProjects/tinylisp/src/gisp/l0.go` and includes comprehensive tests.

## Current Status
The Go implementation is **92% functional** with 38/41 tests passing. The core interpreter works correctly for arithmetic, lists, variables, functions, closures, and complex expressions.

## Major Issues Identified and Fixed

### 1. Atom Interning Bug (FIXED)
**Problem**: The original `atom()` function in l0.go had a string extraction bug:
```go
str := string(A[i:])  // BUG: reads to end of slice instead of to null terminator
```

**Fix Applied**: Modified to find null terminator correctly:
```go
// Find the null terminator to get the correct string length
end := i
for end < hp && A[end] != 0 {
    end++
}
str := string(A[i:end])
```

### 2. Primitive Lookup Bug (FIXED) 
**Problem**: All primitives were stored with ordinal 0, making them indistinguishable in the `apply()` function.

**Fix Applied**: 
- Added `primIndex` map for deterministic primitive indexing
- Modified initialization to assign unique ordinals to each primitive
- Updated `apply()` function to use ordinal-based lookup

### 3. Parser String Extraction Bug (FIXED)
**Problem**: Test parsers had incorrect string slicing in `readAtom()` causing empty strings.

**Fix Applied**: Fixed the string slicing logic:
```go
end := p.pos
if p.ch != 0 {
    end = p.pos - 1  // Back up one if we stopped on a delimiter
}
s := p.input[start:end]
```

### 4. EOF Handling Bug (FIXED)
**Problem**: Original `look()` function calls `os.Exit(0)` on EOF, causing tests to panic.

**Fix Applied**: Created test-safe parsers that handle EOF gracefully by returning 0 instead of exiting.

### 5. NaN Boxing Misunderstanding (FIXED)
**Problem**: Test incorrectly assumed negative numbers should have tags < ATOM.

**Fix Applied**: Updated test to check that numbers remain finite rather than checking tag bits, since negative numbers naturally have high bits set in IEEE 754.

## Current REPL Issue
**Active Problem**: The REPL in l0.go returns `ERR` for expressions like `(+ 1 2)` because it still uses the original buggy parser functions.

**Attempted Fix**: Added a `correctedParser` type to l0.go and modified main() to use it, but there's a compilation error on line 659 referencing undefined `testParser` instead of `correctedParser`.

## Test Files Structure

### Core Test Files
- `l0_test.go` - Unit tests for core functions (NaN boxing, arithmetic, etc.)
- `parser_test.go` - Parser and scanner tests  
- `integration_test.go` - End-to-end integration tests
- `debug_test.go`, `env_debug_test.go`, etc. - Debug helper tests

### Test Coverage
- ✅ Basic arithmetic operations
- ✅ Nested arithmetic expressions  
- ✅ Comparison operations
- ✅ Logical operations
- ✅ List operations (cons, car, cdr, pair?)
- ✅ Quoting and evaluation
- ✅ Variable binding (define, let*)
- ✅ Lambda functions and closures
- ✅ Complex integration (recursive functions, currying)
- ✅ Error conditions
- ✅ Memory management
- ✅ NaN boxing (corrected)

### Failing Tests (3/41)
1. **TestConditionals** - `if` statements failing (2 sub-tests)
2. **TestParsingNumbers** - Negative number parsing (`-5`)
3. **TestParsingAtoms** - EOF issue with original parser

## Key Architecture Details

### NaN Boxing Implementation
- Uses IEEE 754 double precision with tag bits in NaN space
- Tags: ATOM (0x7ff8), PRIM (0x7ff9), CONS (0x7ffa), CLOS (0x7ffb), NIL (0x7ffc)
- Regular numbers (positive/negative) stored as-is
- Only NaN values use tag space for type information

### Memory Layout
- `cell[N]` array serves as both stack (grows down) and atom heap (grows up)
- Stack pointer `sp` starts at N, heap pointer `hp` starts at 0
- Safety invariant: `hp <= sp << 3`

### Environment Structure
- Environments are association lists of variable bindings
- Global environment `env` contains primitive function mappings
- Lexical scoping implemented via closures

## File Locations
- Main implementation: `/Users/kristofer/LocalProjects/tinylisp/src/gisp/l0.go`
- Tests: `/Users/kristofer/LocalProjects/tinylisp/src/gisp/*_test.go`
- Documentation: `/Users/kristofer/LocalProjects/tinylisp/CLAUDE.md`

## Compilation Error to Fix
In l0.go line 659, change `testParser` to `correctedParser`:
```go
// Current (broken):
parser := &testParser{input: input, pos: 0}

// Should be:
parser := &correctedParser{input: input, pos: 0}
```

## Test Commands
- Run all tests: `go test -v`
- Run specific test: `go test -v -run TestName`
- Run REPL: `go run l0.go` (currently broken due to parser issue)

## Context
The user reported that the REPL returns `ERR` for inputs like `(+ 1 2)`. This led to investigation revealing multiple parser and atom interning bugs that have been systematically fixed in the test suite, but the main REPL code still uses the original buggy functions.