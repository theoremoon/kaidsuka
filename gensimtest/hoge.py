import json

with open("lyrics.json") as f:
    f = json.load(f)
xs = []
for k, v in f.items():
    xs.append({"title": k, "lyrics": v})

with open("lyrics.json", "w") as f:
    json.dump(xs, f, indent=2, ensure_ascii=False)
