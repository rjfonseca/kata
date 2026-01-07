"""
Tests for addition in RPN calculator.

>>> from rpn import evaluate

Single numbers should return themselves:
>>> evaluate("5")
5
>>> evaluate("10")
10

Simple addition:
>>> evaluate("1 2 +")
3
>>> evaluate("10 20 +")
30

Multiple additions:
>>> evaluate("1 2 + 3 +")
6
"""
