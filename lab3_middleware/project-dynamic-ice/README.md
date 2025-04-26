# Project: Dynamic ICE Middleware

This project demonstrates the use of the ICE middleware for dynamic invocation of remote methods. It includes both client and server implementations for a calculator service.

## Project Structure

- **client/**: Contains the client-side implementation for invoking remote methods.
  - `Program.cs`: Entry point for the client application.
  - `Operations.cs`: Auto-generated code for proxy and interface definitions.
- **server/**: Contains the server-side implementation of the calculator service.
  - `main.py`: Entry point for the server application.
  - `calculator.py`: Implementation of the calculator service.
  - `Operations_ice.py`: Auto-generated code for server-side operations.
- **Calculations/**: Contains additional generated files for the service.

## Features

The calculator service supports the following operations:
- `add(a, b)`: Adds two integers.
- `subtract(a, b)`: Subtracts the second integer from the first.
- `multiply(a, b)`: Multiplies two integers.
- `hello(name)`: Returns a greeting message.
- `mean(seq)`: Calculates the mean of a sequence of integers.
- `increment(a)`: Increments an integer by 1.

## Prerequisites

- **Python**: Required for running the server.
- **.NET**: Required for running the client.
- **ZeroC ICE**: Ensure ICE is installed and configured.

## Running the Server

Note: Python environment is using [uv](https://docs.astral.sh/uv/) for running the server. You need to install it first.

1. Navigate to the `server/` directory.
2. Run the server using Python:
```bash
uv run main.py
```
3. The server will start and listen for incoming requests.
4. Ensure the server is running before starting the client.

## Running the Client
1. Navigate to the `client/` directory.
2. Build the client application:
```bash
dotnet build
```
3. Run the client application:
```bash
dotnet run
```