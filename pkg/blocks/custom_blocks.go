package blocks

import (
    "image/color"
)

// Custom block definitions from block designer

// Air block definition
func init() {
    BlockDefinitions["air"] = &BlockProperties{
        ID: 0,
        Name: "Air",
        Color: color.RGBA{128, 128, 128, 255},
        Hardness: 0.0,
        Transparent: true,
        Solid: false,
        Collectible: false,
        Flammable: false,
        LightLevel: 0,
        Gravity: false,
        Viscosity: 0.0,
        Pattern: "solid",
    }
}

// Dirt block definition
func init() {
    BlockDefinitions["dirt"] = &BlockProperties{
        ID: 1,
        Name: "Dirt",
        Color: color.RGBA{139, 69, 19, 255},
        Hardness: 1.0,
        Transparent: false,
        Solid: true,
        Collectible: true,
        Flammable: false,
        LightLevel: 0,
        Gravity: false,
        Viscosity: 0.0,
        Pattern: "solid",
    }
}
