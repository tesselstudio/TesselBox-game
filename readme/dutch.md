# TesselBox - Nederlands README
## Hexagonale Voxel Game

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Een 2D sandbox avonturenspel geïnspireerd door *Terraria*, maar gebouwd op een **hexagonale grid**.

Verken werelden, mijn hulpbronnen, bouw structuren, creëer items, vecht tegen vijanden en overleef — alles in mooie hexagonale tegels.

## Game Kenmerken

### ✅ **Volledige Kenmerken**
- **Hexagonale Wereld Generatie** - Procedureel gegenereerde werelden met biomen
- **Mijnbouw en Crafting** - Tool-gebaseerde mijnbouw met verschillende materiaal snelheden
- **Blok Plaatsing** - Rechtsklik om blokken te plaatsen met ghost preview
- **Inventaris Systeem** - 32-slot inventaris met hotbar (9 slots)
- **Gevechts Systeem** - Gezondheid/schade systeem met aanval animaties
- **Dag/Nacht Cyclus** - Dynamische verlichting en tijd voortgang
- **Weer Effecten** - Regen, sneeuw en storm systemen
- **Opslaan/Laden Systeem** - Persistente wereld status met auto-save

### 🎮 **Besturing**
- **WASD / Pijltjes**: Beweging
- **Spatie**: Springen / Aanvallen
- **Linker Klik**: Blok mijnbouw
- **Rechter Klik**: Blok plaatsing
- **E**: Open crafting menu
- **Q**: Laat geselecteerd item vallen
- **Muiswiel**: Hotbar selectie
- **1-9**: Directe hotbar selectie
- **F5**: Handmatig opslaan
- **F9**: Handmatig laden
- **ESC**: Menu / Sluit menu's

## Installatie en Setup

### Vereisten
- **Go 1.19+** - Core engine
- **Git** - Versie beheer

### Snelle Start
```bash
# Kloon repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Bouw game
go build ./cmd/client

# Start game
./client
```

### Ontwikkel Setup
```bash
# Installeer afhankelijkheden
go mod tidy

# Voer tests uit
go test ./...

# Bouw voor ontwikkeling
go build -tags debug ./cmd/client
```

## Systeem Vereisten

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dual-core processor
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ compatible
- **Opslag**: 500MB vrije ruimte

### Aanbevolen
- **CPU**: Quad-core processor
- **RAM**: 8GB+
- **GPU**: Dedicated graphics card
- **Opslag**: 1GB+ vrije ruimte

## Architectuur

### Kern Technologieën
- **Taal**: Go (Golang)
- **Graphics**: Ebiten (2D game library)
- **Build Systeem**: Go modules

### Project Structuur
```
TesselBox/
├── cmd/client/          # Hoofd game executable
├── pkg/                 # Kern pakketten
│   ├── world/          # Wereld generatie en beheer
│   ├── player/         # Speler mechanics en physics
│   ├── blocks/         # Blok types en eigenschappen
│   ├── items/          # Item systeem en crafting
│   ├── crafting/       # Crafting recepten en UI
│   ├── weather/        # Weer simulatie
│   ├── gametime/       # Dag/nacht cyclus
│   ├── save/           # Opslaan/laden functionaliteit
│   └── render/         # Rendering en UI systemen
├── config/             # Configuratie bestanden
└── assets/             # Game assets (indien aanwezig)
```

## Bijdragen

### Voor Ontwikkelaars
1. Fork de repository
2. Maak een feature branch (`git checkout -b feature/geweldige-feature`)
3. Commit je wijzigingen (`git commit -m 'Add geweldige feature'`)
4. Push naar de branch (`git push origin feature/geweldige-feature`)
5. Open een Pull Request

### Ontwikkel Richtlijnen
- Volg Go coding standaarden
- Voeg tests toe voor nieuwe features
- Update documentatie
- Zorg voor cross-platform compatibiliteit

## Licentie

**MIT Licentie** - Zie [LICENSE](LICENSE) bestand voor details.

## Credits

- **Geïnspireerd door**: Terraria game mechanics
- **Gebouwd met**: Ebiten game engine
- **Bijdragers**: Open source community

## Support

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussies**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Project Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Geniet van het verkennen van de hexagonale wereld van TesselBox!*
