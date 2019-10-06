import std.stdio;
import asdf;

struct ChipInfo
{
    bool rigid;
    double friction;
}

struct MapChipInfo
{
    uint chipsize;
    ChipInfo[] chips;
}

const json = `
{
	"chipsize": 16,
	"chips": [
		{
			"rigid": false,
			"friction": 0
		},
		{
			"rigid": true,
			"friction": 1
		},
		{
			"rigid": true,
			"friction": 0.85
		}
	]
}
`;
void main()
{
    auto info = MapChipInfo(16, [
            ChipInfo(false, 0.0), ChipInfo(true, 1.0), ChipInfo(true, 0.85)
            ]);
    writeln(info.serializeToJsonPretty());
    writeln(json.deserialize!MapChipInfo);
}
