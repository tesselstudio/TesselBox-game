#!/bin/bash

# TesselBox Icon Generator
# Creates placeholder icons for all platforms

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ASSETS_DIR="$(dirname "$SCRIPT_DIR")/assets"
ICONS_DIR="$ASSETS_DIR/icons"

echo "🎨 Creating TesselBox icons..."

# Create icons directory
mkdir -p "$ICONS_DIR"

# Create a simple SVG icon (hexagon theme)
cat > "$ICONS_DIR/icon.svg" << 'EOF'
<svg width="256" height="256" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg">
  <!-- Background -->
  <rect width="256" height="256" fill="#2C3E50"/>
  
  <!-- Hexagon Grid Pattern -->
  <defs>
    <pattern id="hexPattern" x="0" y="0" width="60" height="52" patternUnits="userSpaceOnUse">
      <polygon points="30,5 55,20 55,40 30,52 5,40 5,20" fill="#34495E" stroke="#ECF0F1" stroke-width="1"/>
    </pattern>
  </defs>
  
  <!-- Main hexagon -->
  <rect width="256" height="256" fill="url(#hexPattern)"/>
  
  <!-- Center hexagon (logo) -->
  <polygon points="128,40 198,80 198,160 128,200 58,160 58,80" 
           fill="#E74C3C" stroke="#C0392B" stroke-width="3"/>
  
  <!-- Letter T -->
  <rect x="108" y="80" width="40" height="15" fill="white"/>
  <rect x="118" y="95" width="20" height="60" fill="white"/>
  
  <!-- Game title -->
  <text x="128" y="230" font-family="Arial, sans-serif" font-size="18" 
        font-weight="bold" text-anchor="middle" fill="white">TESSELBOX</text>
</svg>
EOF

echo "✅ Created SVG icon"

# Convert to different formats if ImageMagick is available
if command -v convert &> /dev/null; then
    echo "🔄 Converting to platform-specific formats..."
    
    # Windows ICO (multiple sizes)
    convert "$ICONS_DIR/icon.svg" -resize 16x16 "$ICONS_DIR/icon_16.png"
    convert "$ICONS_DIR/icon.svg" -resize 32x32 "$ICONS_DIR/icon_32.png"
    convert "$ICONS_DIR/icon.svg" -resize 48x48 "$ICONS_DIR/icon_48.png"
    convert "$ICONS_DIR/icon.svg" -resize 256x256 "$ICONS_DIR/icon_256.png"
    
    # Create ICO file (requires icotool or convert with ICO support)
    if command -v icotool &> /dev/null; then
        icotool -c -o "$ICONS_DIR/tesselbox.ico" \
            "$ICONS_DIR/icon_16.png" \
            "$ICONS_DIR/icon_32.png" \
            "$ICONS_DIR/icon_48.png" \
            "$ICONS_DIR/icon_256.png"
        echo "✅ Created Windows ICO file"
    else
        echo "⚠️  icotool not found, creating placeholder ICO"
        cp "$ICONS_DIR/icon_256.png" "$ICONS_DIR/tesselbox.ico"
    fi
    
    # macOS ICNS
    if command -v iconutil &> /dev/null; then
        mkdir -p "$ICONS_DIR/tesselbox.iconset"
        convert "$ICONS_DIR/icon.svg" -resize 16x16 "$ICONS_DIR/tesselbox.iconset/icon_16x16.png"
        convert "$ICONS_DIR/icon.svg" -resize 32x32 "$ICONS_DIR/tesselbox.iconset/icon_16x16@2x.png"
        convert "$ICONS_DIR/icon.svg" -resize 32x32 "$ICONS_DIR/tesselbox.iconset/icon_32x32.png"
        convert "$ICONS_DIR/icon.svg" -resize 64x64 "$ICONS_DIR/tesselbox.iconset/icon_32x32@2x.png"
        convert "$ICONS_DIR/icon.svg" -resize 128x128 "$ICONS_DIR/tesselbox.iconset/icon_128x128.png"
        convert "$ICONS_DIR/icon.svg" -resize 256x256 "$ICONS_DIR/tesselbox.iconset/icon_256x256.png"
        convert "$ICONS_DIR/icon.svg" -resize 256x256 "$ICONS_DIR/tesselbox.iconset/icon_128x128@2x.png"
        convert "$ICONS_DIR/icon.svg" -resize 512x512 "$ICONS_DIR/tesselbox.iconset/icon_256x256@2x.png"
        
        iconutil -c icns "$ICONS_DIR/tesselbox.iconset" -o "$ICONS_DIR/TesselBox.icns"
        echo "✅ Created macOS ICNS file"
    else
        echo "⚠️  iconutil not found, creating placeholder ICNS"
        cp "$ICONS_DIR/icon_256.png" "$ICONS_DIR/TesselBox.icns"
    fi
    
    # Linux PNG
    cp "$ICONS_DIR/icon_256.png" "$ICONS_DIR/tesselbox.png"
    echo "✅ Created Linux PNG file"
    
    # Clean up temporary files
    rm -f "$ICONS_DIR/icon_*.png"
    rm -rf "$ICONS_DIR/tesselbox.iconset"
    
else
    echo "⚠️  ImageMagick not found, keeping SVG only"
    cp "$ICONS_DIR/icon.svg" "$ICONS_DIR/tesselbox.png"
fi

echo "🎉 Icon generation completed!"
echo "📁 Icons located in: $ICONS_DIR"
