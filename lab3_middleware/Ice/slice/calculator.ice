
#ifndef CALC_ICE
#define CALC_ICE

module Demo
{
  enum operation { MIN, MAX, AVG };
  
  exception NoInput {};
  exception EmptySequence {};

  struct A
  {
    short a;
    long b;
    float c;
    string d;
  };
  sequence<long> Seq;

  interface Calc
  {
    long add(int a, int b);
    long subtract(int a, int b);
    void op(A a1, short b1); //załóżmy, że to też jest operacja arytmetyczna ;)
    double avg(Seq seq) throws EmptySequence;
  };

};

#endif
