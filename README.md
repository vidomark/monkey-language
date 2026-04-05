# Bytelang

A programming language implementation in Go featuring a lexer, parser, compiler, and bytecode virtual machine.

## Features

- **Data types:** integers, strings, booleans, arrays, first-class functions
- **Operators:** arithmetic (`+`, `-`, `*`, `/`), comparison (`<`, `>`, `==`, `!=`), prefix (`-`, `!`)
- **Control flow:** if/else expressions, return statements
- **Functions:** closures, recursion, higher-order functions
- **Built-ins:** `len`, `puts`, `first`, `last`, `push`
- **Macro system:** compile-time code generation with `quote`/`unquote`
- **Dual execution:** tree-walking interpreter and bytecode compiler + VM

## Quick start

```
go run main.go
```

This starts the interactive REPL where you can type expressions:

```
>> let add = fn(a, b) { a + b };
>> add(2, 3)
5
>> let map = fn(arr, f) { if (len(arr) == 0) { [] } else { push(map(last(arr), f), f(first(arr))) } };
>> map([1, 2, 3], fn(x) { x * 2 })
[2, 4, 6]
```

## Benchmark

Compare the bytecode VM against the tree-walking interpreter:

```
go run src/bytelang/benchmark/main.go --engine=vm
go run src/bytelang/benchmark/main.go --engine=eval
```

Computes `fibonacci(35)` using each engine.

## Tests

```
go test ./...
```

## Project structure

```
main.go                     Entry point (REPL)
src/bytelang/
  token/                    Token type definitions
  lexer/                    Lexical analysis
  ast/                      Abstract syntax tree nodes
  parser/                   Pratt parser
  evaluator/                Tree-walking interpreter
  object/                   Runtime object system
  code/                     Bytecode instruction set
  compiler/                 Bytecode compiler
  vm/                       Virtual machine
  repl/                     Interactive REPL (compiler mode)
  benchmark/                Performance comparison
```
