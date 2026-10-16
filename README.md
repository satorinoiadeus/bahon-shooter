# Bahon Shooter

A tiny **GZDoom / Doom II** shooter mod.

## What it does

- Opens directly into a small arena (`MAP01`).
- The boss starts as **Bahon** with 200 HP.
- When Bahon drops below 50% health, the actor is replaced by a **Leonardo DiCaprio** phase with 100 HP.
- Uses Doom II's normal weapon and projectile behavior, with custom PNG sprites.
- Includes a generated preview in `docs/preview.png`.

## Run it

You need:

1. **GZDoom**
2. A legal copy of `DOOM2.WAD`

Then run:

```bash
gzdoom -iwad DOOM2.WAD -file bahon-shooter.pk3 +map MAP01
```

On Windows you can also drag `bahon-shooter.pk3` onto `gzdoom.exe`, then select Doom II and start MAP01.

## Project layout

- `DECORATE` — Bahon and transformed boss logic.
- `MAPINFO` — game/map metadata.
- `maps/MAP01.wad` — compact UDMF arena.
- `sprites/BAHOA0.png` — Bahon sprite.
- `sprites/LEOOA0.png` — transformed sprite.
- `docs/preview.png` — concept preview.

## Build the PK3

A PK3 is just a ZIP archive with a `.pk3` extension. From the repo root:

```bash
python tools/build.py
```

The result is `dist/bahon-shooter.pk3`.

> Fan-made parody mod. No affiliation or endorsement is implied.
