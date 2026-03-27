# TesselBox Release Script
# Automated release with version management

set -e

VERSION=${1:-"2.0.0"}
RELEASE_BRANCH="main"

echo "🚀 Releasing TesselBox v$VERSION"

# Check if we're on the right branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "$RELEASE_BRANCH" ]; then
    echo "❌ Must be on $RELEASE_BRANCH branch to release"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Working directory is not clean"
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
go test ./...

# Build all platforms
echo "🔨 Building all platforms..."
mkdir -p bin
make release

# Create release notes
echo "📝 Creating release notes..."
cat > RELEASE_NOTES.md << EOF
# TesselBox v$VERSION

## 🎮 Features
- Perfect hexagonal grid block placement
- Improved collision detection
- Reduced player placement restriction zone
- Embedded assets for single binary distribution

## 🐛 Fixes
- Fixed block placement coordinate transformation
- Resolved hexagonal grid alignment issues
- Corrected vertical positioning of placed blocks

## 📦 Downloads
- Windows: tesselbox-windows-amd64.exe
- Linux: tesselbox-linux-amd64
- macOS Intel: tesselbox-darwin-amd64
- macOS Apple Silicon: tesselbox-darwin-arm64

## 🔧 Installation
1. Download the appropriate binary for your platform
2. Make it executable (Linux/macOS): \`chmod +x tesselbox-*\`
3. Run: \`./tesselbox-*\`

Enjoy building with TesselBox! 🎯
EOF

# Tag and push
echo "🏷️  Creating and pushing tag..."
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

echo "✅ Release v$VERSION completed!"
echo "📦 Check GitHub Actions for build status"
echo "🌐 Release will be available at: https://github.com/$(git config remote.origin.url | sed 's/.*github.com[:\/]\(.*\)\.git/\1/')/releases/tag/v$VERSION"
