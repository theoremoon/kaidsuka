import std.stdio;
import std.file;
import asdf;

struct MapLayer
{
public:
    int[] data;

    int width;
    int height;
}

struct MapInfo
{
public:
    uint width;
    uint height;
    uint tilewidth;
    uint tileheight;
    MapLayer[] layers;
}

void main()
{
    auto mapinfo = readText("map1.json").deserialize!MapInfo();
	writeln(mapinfo);

}
