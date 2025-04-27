import Ice
import calculator
import sys


def main():
    with Ice.initialize(sys.argv) as calc:
        adapter = calc.createObjectAdapterWithEndpoints(
            "CalculatorAdapter", "tcp -p 4061")

        adapter.add(calculator.CalculatorI(), Ice.stringToIdentity("calculator"))

        adapter.activate()
        print("Listening on port 4061...")

        try:
            calc.waitForShutdown()
        except KeyboardInterrupt:
            print("Caught Ctrl+C, exiting...")


if __name__ == "__main__":
    main()
