"""
Tests for subtraction in RPN calculator.

>>> from rpn import evaluate

Simple subtraction:
>>> evaluate("5 3 -")
2
>>> evaluate("10 20 -")
-10

Mixed operations:
>>> evaluate("10 2 - 3 +")
11
>>> evaluate("10 2 + 3 -")
9
"""
