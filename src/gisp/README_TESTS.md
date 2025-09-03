# Go TinyLisp Test Suite

This directory contains comprehensive tests for the Go implementation of tinylisp (`l0.go`). The test suite is designed to help debug and validate the interpreter implementation.

## Test Files

### 1. `l0_test.go` - Core Unit Tests
Tests the fundamental building blocks of the interpreter:

- **NaN Boxing**: Tests the box/unbox operations for different data types
- **Memory Management**: Tests atom interning, cons cell allocation, and memory safety
- **Data Structures**: Tests cons, car, cdr operations and list construction
- **Primitive Functions**: Tests arithmetic (+, -, *, /), comparison (<, eq?), and logic (and, or, not)
- **Environment Operations**: Tests variable binding, lookup, and environment management
- **Basic Evaluation**: Tests evaluation of atoms, numbers, and simple expressions

### 2. `parser_test.go` - Parser and Scanner Tests
Tests the parsing and tokenization components:

- **Number Parsing**: Tests parsing of integers, floats, scientific notation
- **Atom Parsing**: Tests parsing of symbols, operators, and special atoms
- **List Parsing**: Tests parsing of empty lists, simple lists, and dotted pairs
- **Quote Parsing**: Tests parsing of quoted expressions ('x, '(1 2 3))
- **Complex Expressions**: Tests parsing of nested lists and function calls
- **Scanner Tests**: Tests the tokenization of different input types
- **Error Handling**: Tests parser behavior with malformed input

### 3. `integration_test.go` - End-to-End Integration Tests
Tests complete Lisp expressions from parsing through evaluation:

- **Basic Arithmetic**: Simple arithmetic operations
- **Nested Expressions**: Complex nested arithmetic and function calls
- **Comparisons**: Equality and ordering tests
- **Logical Operations**: Boolean logic with proper truthiness
- **List Operations**: cons, car, cdr, pair? operations
- **Quoting and Evaluation**: Quote and eval functionality
- **Conditionals**: if and cond expressions
- **Variable Binding**: define and let* binding
- **Lambda Functions**: Function definition and application
- **Complex Integration**: Multi-step expressions combining multiple features
- **Error Conditions**: Testing error cases and edge conditions
- **Compatibility Tests**: Tests based on the original dotcall.lisp test suite

## Running the Tests

### Run All Tests
```bash
cd src/gisp
go test -v
```

### Run Specific Test Files
```bash
# Run only unit tests
go test -v -run "Test.*" l0_test.go l0.go

# Run only parser tests  
go test -v -run "Test.*" parser_test.go l0.go

# Run only integration tests
go test -v -run "Test.*" integration_test.go l0.go
```

### Run Specific Test Functions
```bash
# Test only NaN boxing
go test -v -run "TestNaNBoxing"

# Test only arithmetic operations
go test -v -run "TestArithmetic"

# Test only parsing
go test -v -run "TestParsing.*"
```

## Test Structure

Each test follows Go testing conventions:

- Tests are functions starting with `Test`
- Use `t.Errorf()` for test failures with descriptive messages
- Use `t.Run()` for subtests to organize related test cases
- Include both positive and negative test cases
- Test edge cases and error conditions

## Debugging with Tests

### Finding Issues
1. **Start with Unit Tests**: If basic operations fail, check `l0_test.go`
2. **Check Parsing**: If expressions aren't parsed correctly, check `parser_test.go`  
3. **Integration Issues**: If parsing works but evaluation fails, check `integration_test.go`

### Common Issues to Look For
1. **NaN Boxing Errors**: Tag bits not set correctly, ordinals corrupted
2. **Memory Management**: Stack/heap collision, incorrect pointer arithmetic
3. **Environment Issues**: Variable lookup failures, binding problems
4. **Parser Issues**: Incorrect tokenization, malformed parse trees
5. **Evaluation Issues**: Incorrect primitive implementations, closure problems

### Adding New Tests
When adding features or fixing bugs:

1. **Add Unit Tests**: Test the individual function in `l0_test.go`
2. **Add Parser Tests**: If syntax changes, add tests in `parser_test.go`
3. **Add Integration Tests**: Test the complete feature in `integration_test.go`

## Test Coverage

The test suite covers:
- ✅ Core data structures (atoms, numbers, lists, closures)
- ✅ All primitive functions (arithmetic, logic, list operations)
- ✅ Parser for all syntax forms (atoms, lists, quotes, dots)
- ✅ Variable binding and environments
- ✅ Lambda functions and closures
- ✅ Control flow (if, cond)
- ✅ Error conditions and edge cases
- ✅ Memory management and safety invariants
- ✅ Compatibility with original test cases

The tests are designed to catch the most common implementation issues in Lisp interpreters and provide clear diagnostic information when failures occur.