from pathlib import Path
import zipfile

root = Path(__file__).resolve().parents[1]
dist = root / "dist"
dist.mkdir(exist_ok=True)
out = dist / "bahon-shooter.pk3"
include = ["DECORATE", "MAPINFO", "maps/MAP01.wad", "sprites/BAHOA0.png", "sprites/LEOOA0.png"]

# Store entries without DEFLATE. GZDoom can read stored ZIP entries directly and this
# avoids an unnecessary zlib layer around PNG/WAD data in the portable build.
with zipfile.ZipFile(out, "w", compression=zipfile.ZIP_STORED, allowZip64=False) as z:
    for rel in include:
        z.write(root / rel, rel)

with zipfile.ZipFile(out, "r") as z:
    bad = z.testzip()
    if bad is not None:
        raise RuntimeError(f"corrupt PK3 entry: {bad}")

print(out)
