# Reverse Polish Notation (RPN) Calculator Kata

This kata is designed to help you practice Test-Driven Development (TDD) by building a simple command-line calculator that evaluates expressions written in Reverse Polish Notation (RPN).

## What is RPN?

RPN, also known as postfix notation, is a mathematical notation in which every operator follows all of its operands. It is stack-based and does not require any parentheses to specify the order of operations.

For example, the infix expression `(1 + 2) * 3` would be written in RPN as `1 2 + 3 *`.

## The Goal

Your goal is to implement a function, `Evaluate(string) (int, error)`, that takes an RPN expression as a string and returns the calculated integer result.

## The Steps

The kata is divided into several steps, each adding a new layer of functionality. You should follow the TDD cycle for each step:
1.  **Red:** Write a failing test that defines the new behavior.
2.  **Green:** Write the simplest possible code to make the test pass.
3.  **Refactor:** Clean up your code while keeping the tests green.

Good luck!
