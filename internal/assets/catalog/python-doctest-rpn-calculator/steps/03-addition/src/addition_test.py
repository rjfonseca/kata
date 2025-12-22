"""
Tests for addition in RPN calculator.

>>> from rpn import evaluate

Simple addition:
>>> evaluate("1 2 +")
3
>>> evaluate("10 20 +")
30

Multiple additions:
>>> evaluate("1 2 + 3 +")
6
"""
