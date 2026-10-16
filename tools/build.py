from pathlib import Path
import zipfile

root = Path(__file__).resolve().parents[1]
dist = root / "dist"
dist.mkdir(exist_ok=True)
out = dist / "bahon-shooter.pk3"
include = ["DECORATE", "MAPINFO", "maps/MAP01.wad", "sprites/BAHOA0.png", "sprites/LEOOA0.png"]
with zipfile.ZipFile(out, "w", zipfile.ZIP_DEFLATED) as z:
    for rel in include:
        z.write(root / rel, rel)
print(out)
