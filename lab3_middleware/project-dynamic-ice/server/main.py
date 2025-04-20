import Ice
import calculator
import sys


def main():
    # Create an Ice communicator. We'll use this communicator to create an object adapter.
    with Ice.initialize(sys.argv) as calc:
        # Create an object adapter that listens for incoming requests and dispatches them to servants.
        adapter = calc.createObjectAdapterWithEndpoints(
            "CalculatorAdapter", "tcp -p 4061")

        # Register the Calc servant with the adapter.
        adapter.add(calculator.CalculatorI(), Ice.stringToIdentity("calculator"))

        # Start dispatching requests.
        adapter.activate()
        print("Listening on port 4061...")

        try:
            # Wait until communicator.shutdown() is called, which never occurs in this demo.
            calc.waitForShutdown()
        except KeyboardInterrupt:
            print("Caught Ctrl+C, exiting...")


if __name__ == "__main__":
    main()
