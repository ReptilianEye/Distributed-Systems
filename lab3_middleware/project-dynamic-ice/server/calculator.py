import Ice

import Calculations

class CalculatorI(Calculations.Calculator):
    def add(self, a, b, current:Ice.Current):
        return a + b

    def subtract(self, a, b, current:Ice.Current):
        return a - b

    def multiply(self, a, b, current:Ice.Current):
        return a * b
    
    def dynamichello(self,name, current:Ice.Current):
        return "you found me: " + name
