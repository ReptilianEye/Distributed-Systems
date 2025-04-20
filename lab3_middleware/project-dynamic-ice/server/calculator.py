import Ice

import Calculations

class CalculatorI(Calculations.Calculator):
    def add(self, a, b, current:Ice.Current):
        return a + b

    def subtract(self, a, b, current:Ice.Current):
        return a - b

    def multiply(self, a, b, current:Ice.Current):
        return a * b
    
    def hello(self, name, current:Ice.Current):
        return "you found me: " + name

    def mean(self, *vals):
        vals = list(vals)[:-1]
        return sum(*vals) / len(*vals)
    
    def increment(self, val, current:Ice.Current):
        return val + 1
    
