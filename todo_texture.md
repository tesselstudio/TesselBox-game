# TesselBox Blocks Reference

## All Blocks

| # | Block Name | ID | Hardness | Properties |
|---|------------|-----|----------|------------|
| 1 | Air | `air` | 0 | Transparent, Non-solid |
| 2 | Dirt | `dirt` | 0.5 | Solid, Collectible, Gravity |
| 3 | Grass | `grass` | 0.6 | Solid, Collectible, Gravity |
| 4 | Stone | `stone` | 1.5 | Solid, Collectible, Gravity |
| 5 | Sand | `sand` | 0.5 | Solid, Collectible, Gravity |
| 6 | Water | `water` | - | Liquid |
| 7 | Log | `log` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 8 | Leaves | `leaves` | 0.2 | Transparent, Solid, Collectible, Flammable |
| 9 | Tropical Log | `tropical_log` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 10 | Temperate Log | `temperate_log` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 11 | Pine Log | `pine_log` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 12 | Tropical Leaves | `tropical_leaves` | 0.2 | Transparent, Solid, Collectible, Flammable |
| 13 | Temperate Leaves | `temperate_leaves` | 0.2 | Transparent, Solid, Collectible, Flammable |
| 14 | Pine Leaves | `pine_leaves` | 0.2 | Transparent, Solid, Collectible, Flammable |
| 15 | Coal Ore | `coal_ore` | 1.5 | Solid, Collectible, Gravity |
| 16 | Iron Ore | `iron_ore` | 2.0 | Solid, Collectible, Gravity |
| 17 | Gold Ore | `gold_ore` | 2.0 | Solid, Collectible, Gravity |
| 18 | Diamond Ore | `diamond_ore` | 3.0 | Solid, Collectible, Gravity |
| 19 | Bedrock | `bedrock` | -1 | Unbreakable, Solid |
| 20 | Glass | `glass` | - | Transparent |
| 21 | Brick | `brick` | - | Solid |
| 22 | Plank | `plank` | - | Solid, Flammable |
| 23 | Cactus | `cactus` | - | Solid, Collectible |
| 24 | Workbench | `workbench` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 25 | Furnace | `furnace` | 1.5 | Solid, Collectible, Gravity |
| 26 | Anvil | `anvil` | 2.0 | Solid, Collectible, Gravity |
| 27 | Gravel | `gravel` | 0.6 | Solid, Collectible, Gravity |
| 28 | Sandstone | `sandstone` | 0.8 | Solid, Collectible, Gravity |
| 29 | Obsidian | `obsidian` | 5.0 | Solid, Collectible, Gravity |
| 30 | Ice | `ice` | 0.5 | Transparent, Solid, Collectible |
| 31 | Snow | `snow` | 0.2 | Solid, Collectible |
| 32 | Torch | `torch` | 0.1 | Transparent, Non-solid, Light Level 14 |
| 33 | Crafting Table | `crafting_table` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 34 | Chest | `chest` | 1.0 | Solid, Collectible, Flammable, Gravity |
| 35 | Ladder | `ladder` | - | Transparent |
| 36 | Fence | `fence` | - | Solid |
| 37 | Gate | `gate` | - | Solid |
| 38 | Door | `door` | - | Solid |
| 39 | Window | `window` | - | Transparent |
| 40 | Flower | `flower` | 0.1 | Transparent, Non-solid, Collectible |
| 41 | Tall Grass | `tall_grass` | - | Transparent, Non-solid |
| 42 | Red Mushroom | `mushroom_red` | - | Collectible |
| 43 | Brown Mushroom | `mushroom_brown` | - | Collectible |
| 44 | Wool | `wool` | 0.3 | Solid, Collectible, Flammable, Gravity |
| 45 | Bookshelf | `bookshelf` | - | Solid, Flammable |
| 46 | Jukebox | `jukebox` | - | Solid |
| 47 | Note Block | `note_block` | - | Solid |
| 48 | Pumpkin | `pumpkin` | 0.5 | Solid, Collectible, Gravity |
| 49 | Melon | `melon` | - | Solid, Collectible |
| 50 | Hay Bale | `hay_bale` | - | Solid, Flammable |
| 51 | Cobblestone | `cobblestone` | 1.5 | Solid, Collectible, Gravity |
| 52 | Mossy Cobblestone | `mossy_cobblestone` | - | Solid |
| 53 | Stone Bricks | `stone_bricks` | - | Solid |
| 54 | Chiseled Stone | `chiseled_stone` | - | Solid |
| 55 | Randomland Portal | `randomland_portal` | -1 | Unbreakable, Transparent, Light Level 15 |

## Liquids

| Name | Type | Properties |
|------|------|------------|
| Water | `water` | Density 1.0, Viscosity 0.3, Transparent |
| Lava | `lava` | Density 3.0, Viscosity 0.8, Light Level 12 |

## Block Properties Explained

- **Hardness**: Time in seconds to break (0 = instant, -1 = unbreakable)
- **Transparent**: Can see through/light passes through
- **Solid**: Player can stand on it
- **Collectible**: Can be picked up and added to inventory
- **Flammable**: Can be burned by fire
- **Gravity**: Falls like sand/gravel
- **Light Level**: 0-15 (0 = no light, 15 = full brightness)