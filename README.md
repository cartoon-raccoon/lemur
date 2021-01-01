# monkey-jit
A JIT implementation of the Monkey programming language

Based off the Monkey programming language developed by Thorsten Ball for his book [Writing an Interpreter in Go](https://interpreterbook.com).

Currently, the language runtime (as implemented by Ball) is a tree-walking interpreter, but I plan to expand it and add a JIT native code compiler that can be swapped out with different instruction sets and architectures. I plan to implement the lexer and parser according to the book, and create a module for AST walking, but from there I will add my own code to it. Ball does add a bytecode compiler in the sequel book Writing a Compiler, but I will implement my own backend for the language, possibly using LLVM.

I have a separate programming language project involving a bytecode compiler and a register-based VM [here](https://github.com/cartoon-raccoon/verdigris).

A standard library may or may not be developed, depending on whether this language reaches maturity.

I may or may not be a bit too optimistic, but I do think that I could make this work.
