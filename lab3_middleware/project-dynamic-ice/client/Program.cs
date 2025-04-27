using System;
using Calculations;

namespace Client
{
    public class Program
    {
        public static int Main(string[] args)
        {
            try
            {
                using (Ice.Communicator communicator = Ice.Util.initialize(ref args))
                {
                    var obj = communicator.stringToProxy("calculator:default -h localhost -p 4061");
                    var calculator = CalculatorPrxHelper.checkedCast(obj);
                    if (calculator == null)
                    {
                        throw new ApplicationException("Invalid proxy");
                    }

                    calculator.add(1, 2);
                    Console.WriteLine("Addition result: " + calculator.add(1, 2));
                    Ice.OutputStream outStream = new(communicator);
                    byte[] inParams = [];
                    byte[] outParams = [];
                    outStream.startEncapsulation();
                    outStream.writeString("Piotr");
                    outStream.endEncapsulation();
                    inParams = outStream.finished();
                    if (calculator.ice_invoke("hello", Ice.OperationMode.Normal, inParams, out outParams))
                    {
                        Ice.InputStream inStream = new(communicator, outParams);
                        inStream.startEncapsulation();
                        string result = inStream.readString();
                        inStream.endEncapsulation();
                        Console.WriteLine(result);
                    }
                    else
                    {
                        Console.WriteLine("Dynamic 'hello' invocation failed.");
                    }
                    outStream.reset();
                    outStream.startEncapsulation();
                    outStream.writeIntSeq([1,2,3,4,5]);
                    outStream.endEncapsulation();
                    inParams = outStream.finished();
                    if (calculator.ice_invoke("mean", Ice.OperationMode.Normal, inParams, out outParams))
                    {

                        // Handle success
                        Ice.InputStream inStream = new(communicator, outParams);
                        inStream.startEncapsulation();
                        double result = inStream.readDouble();
                        inStream.endEncapsulation();
                        Console.WriteLine(result);
                    }
                    else
                    {
                        Console.WriteLine("Dynamic 'mean' invocation failed.");
                    }
                    outStream.reset();
                    outStream.startEncapsulation();
                    outStream.writeInt(5);
                    outStream.endEncapsulation();
                    inParams = outStream.finished();
                    if (calculator.ice_invoke("increment", Ice.OperationMode.Normal, inParams, out outParams))
                    {

                        Ice.InputStream inStream = new(communicator, outParams);
                        inStream.startEncapsulation();
                        double result = inStream.readInt();
                        inStream.endEncapsulation();
                        Console.WriteLine(result);
                    }
                    else
                    {
                        Console.WriteLine("Dynamic 'increment' invocation failed.");
                    }

                }
            }
            catch (Exception e)
            {
                Console.Error.WriteLine(e);
                return 1;
            }
            return 0;
        }
    }
}