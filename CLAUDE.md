# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is tinylisp, a collection of minimal Lisp interpreters written in C. The project demonstrates how to implement a complete Lisp interpreter in just 99 lines of C using NaN boxing and functional programming techniques. The codebase includes multiple variants optimized for different use cases and platforms.

## Build Commands

### Basic Compilation
```bash
# Compile the standard optimized version (recommended)
cc -o tinylisp src/tinylisp-opt.c

# Compile the basic version
cc -o tinylisp src/tinylisp.c

# Compile the extras version with readline support
cc -o tinylisp-extras src/tinylisp-extras.c -lreadline

# Compile floating point versions
cc -o tinylisp-float src/tinylisp-float.c
cc -o tinylisp-float-opt src/tinylisp-float-opt.c

# Compile Sharp PC-G850 versions
cc -o lisp850 src/lisp850.c
cc -o lisp850-opt src/lisp850-opt.c
```

### Running with Standard Libraries
```bash
# Load common Lisp functions
cat src/common.lisp src/list.lisp src/math.lisp | ./tinylisp

# Run unit tests
./tinylisp < tests/dotcall.lisp
./tinylisp-extras < tests/dotcall-extras.lisp
```

## Code Architecture

### Core Implementation Strategy
- **NaN Boxing**: All Lisp values are stored as IEEE 754 double-precision floats, with non-numeric values encoded in the NaN space using tag bits
- **Unified Memory Model**: Uses a single `cell[]` array that serves as both a stack for cons cells (growing down from top) and a heap for atom strings (growing up from bottom)
- **Tags**: ATOM (0x7ff8), PRIM (0x7ff9), CONS (0x7ffa), CLOS (0x7ffb), NIL (0x7ffc), MACR (0x7ffc in extras)
- **Garbage Collection**: Simple stack-based allocation with automatic reclamation when stack unwinds

### Key Files and Their Purpose
- `tinylisp.c`: Original 99-line implementation with full Lisp functionality
- `tinylisp-opt.c`: Optimized version with tail-call optimization and improved performance
- `tinylisp-extras.c`: Extended version with 16 additional primitives including file I/O, exceptions, macros, and tracing
- `tinylisp-float.c/tinylisp-float-opt.c`: Single-precision variants
- `lisp850.c/lisp850-opt.c`: BCD boxing variants for Sharp PC-G850 pocket computer
- `gisp/l0.go`: Go port of the interpreter
- `common.lisp`: Standard Lisp library functions implemented in tinylisp
- `list.lisp`: List manipulation functions
- `math.lisp`: Mathematical functions

### Memory Management
- Stack grows downward from `cell[N]` with `sp` (stack pointer)
- Atom heap grows upward from `cell[0]` with `hp` (heap pointer)
- Safety invariant: `hp <= sp << 3` (atoms take 8x less space than cells)
- Memory size controlled by `N` constant (default 1024 for 8KB, or 8192 for 64KB in extras)

### Evaluation Model
- `eval(expr, env)`: Main evaluation function
- `evlis(list, env)`: Evaluates argument lists
- `evarg(arg, env, flag)`: Optimized argument evaluation with dot-operator support
- Environments are association lists of variable bindings
- Supports lexical scoping with closures

### Built-in Primitives (21 in basic version)
Core: `eval`, `quote`, `cons`, `car`, `cdr`, `pair?`
Arithmetic: `+`, `-`, `*`, `/`, `int`, `<`
Logic: `eq?`, `or`, `and`, `not`
Control: `cond`, `if`
Binding: `let*`, `lambda`, `define`

### Additional Primitives in Extras Version (+16)
File I/O: `load`, `read`, `print`, `println`
Exceptions: `catch`, `throw`
Mutation: `setq`, `set-car!`, `set-cdr!`
Introspection: `assoc`, `env`, `trace`
Extended binding: `let`, `letrec*`, `letrec`
Macros: `macro` with backquote support

### Testing
The test suite in `tests/` verifies:
- Dot operator functionality in lambda parameters and function calls
- Arithmetic and logic operations
- List operations and memory management
- Static scoping and closures
- All language features including macros (in extras version)

### Go Implementation
The `src/gisp/` directory contains a Go port that implements the same NaN boxing strategy and core Lisp functionality, demonstrating the portability of the design.