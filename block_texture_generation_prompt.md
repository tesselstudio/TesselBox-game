# Block Texture Generation Prompt for AI Image Generator

## Overview
Generate high-quality PNG textures for all blocks in TesselBox, a hexagonal voxel-based sandbox game. Each texture should be a 16x16 pixel image that replaces the current procedural color-based appearance with detailed, realistic textures suitable for hexagonal block faces.

## General Guidelines
- **Resolution**: 37 pixels per texture
- **Format**: PNG with transparency support
- **Block Shape**: Hexagonal geometry - textures should work well when mapped to hexagonal faces
- **Style**: Pixel art with clear, recognizable features optimized for hexagonal tiles
- **Lighting**: Assume neutral lighting, no shadows unless part of the material
- **Perspective**: Top-down view for most blocks, designed to fit hexagonal boundaries
- **Color Reference**: Use the provided RGBA values as base colors, but enhance with realistic details and variations



## Output Requirements
- Generate one 16x16 PNG texture per block type
- Name files using block ID (e.g., `dirt.png`, `stone.png`, `grass.png`)
- Include multiple variations where specified (e.g., `dirt_dry.png`, `dirt_normal.png`, `dirt_wet.png`)
- Ensure textures tile seamlessly for adjacent hexagonal blocks
- Design textures to work well with hexagonal geometry and texture mapping
- Maintain appropriate transparency for transparent blocks
- Use the reference colors as base but add realistic details and variations

## Quality Standards
- Sharp pixel art with clear details
- Consistent style across all textures
- Appropriate level of detail for 37 cell resolution
- Colors should match the game's aesthetic while being more realistic than flat colors