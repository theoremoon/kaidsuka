import std.stdio;
import std.conv;
import std.traits;
import std.exception : assertThrown;

struct MyStruct
{
public:
    int x;
    this(int x)
    {
        this.x = x;
    }

    T to(T)() if (isIntegral!(T))
    {
        return this.x.to!(T);
    }
}

void main()
{
    auto s = MyStruct(1000);
    writeln(s.to!long);
    assertThrown!ConvOverflowException(s.to!ubyte);
}
