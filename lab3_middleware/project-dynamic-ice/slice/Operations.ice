#pragma once

module Calculations
{
    sequence<int> Seq;

    /// Represents a simple calculator.
    interface Calculator
    {
        /// Adds two integers.
        /// @param a The first integer.
        /// @param b The second integer.
        /// @return The sum of a and b.
        int add(int a, int b);

        /// Subtracts two integers.
        /// @param a The first integer.
        /// @param b The second integer.
        /// @return The difference of a and b.
        int subtract(int a, int b);

        /// Multiplies two integers.
        /// @param a The first integer.
        /// @param b The second integer.
        /// @return The product of a and b.
        int multiply(int a, int b);

        string hello(string name);

        double mean(Seq seq);

        int increment(int a);

        }

}