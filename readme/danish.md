# TesselBox - Dansk README
## Hexagonal Voxel Spil

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Et 2D sandbox eventyrspil inspireret af *Terraria*, men bygget på et **hexagonalt gitter**.

Udforsk verdener, udvind ressourcer, bygg strukturer, skab genstande, kæmp mod fjender og overlev — alt i smukke hexagonale fliser.

## Spil Funktioner

### ✅ **Komplette Funktioner**
- **Hexagonal Verdens Generering** - Procedurelt genererede verdener med biomer
- **Udnyttelse og Håndværk** - Værktøjsbaseret udnyttelse med forskellige materialehastigheder
- **Blok Placering** - Højreklik for at placere blokke med spøgelses preview
- **Inventar System** - 32-slot inventar med hotbar (9 slots)
- **Kampsystem** - Sundhed/skadesystem med angrebsanimationer
- **Dag/Nat Cyklus** - Skiftende belysning og tidsfremgang
- **Vejr Effekter** - Regn, sne og storm systemer
- **Gem/Indlæs System** - Vedvarende verdensstatus med automatisk lagring

### 🎮 **Styring**
- **WASD / Pile**: Bevægelse
- **Mellemrum**: Hop / Angrib
- **Venstre Klik**: Blok udnyttelse
- **Højre Klik**: Blok placering
- **E**: Åbn håndværksmenu
- **Q**: Slip valgte genstand
- **Musehjul**: Hotbar valg
- **1-9**: Direkte hotbar valg
- **F5**: Manuel lagring
- **F9**: Manuel indlæsning
- **ESC**: Menu / Luk menuer

## Installation og Opsætning

### Krav
- **Go 1.19+** - Kærnemotor
- **Git** - Versionskontrol

### Hurtig Start
```bash
# Klon repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Byg spil
go build ./cmd/client

# Start spil
./client
```

## System Krav

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: To-kærne processor
- **RAM**: 4GB
- **GPU**: Kompatibel med OpenGL 3.3+
- **Lager**: 500MB ledig plads

### Anbefalet
- **CPU**: Fire-kærne processor
- **RAM**: 8GB+
- **GPU**: Dedikeret grafik kort
- **Lager**: 1GB+ ledig plads

## Arkitektur

### Kærne Teknologier
- **Sprog**: Go (Golang)
- **Grafik**: Ebiten (2D spil bibliotek)
- **Bygge System**: Go moduler

### Projekt Struktur
```
TesselBox/
├── cmd/client/          # Hoved spil eksekverbar fil
├── pkg/                 # Kærne pakker
│   ├── world/          # Verdens generering og styring
│   ├── player/         # Spiller mekanikker og fysik
│   ├── blocks/         # Blok typer og egenskaber
│   ├── items/          # Genstands system og håndværk
│   ├── crafting/       # Håndværks opskrifter og UI
│   ├── weather/        # Vejr simulering
│   ├── gametime/       # Dag/nat cyklus
│   ├── save/           # Gem/indlæs funktionalitet
│   └── render/         # Rendering og UI systemer
├── config/             # Konfigurationsfiler
└── assets/             # Spil aktiver (hvis nogen)
```

## Bidrag

### For Udviklere
1. Fork repository
2. Opret en funktion gren (`git checkout -b feature/amazing-feature`)
3. Commit dine ændringer (`git commit -m 'Add amazing feature'`)
4. Push til grenen (`git push origin feature/amazing-feature`)
5. Åbn en Pull Request

### Udviklings Retningslinjer
- Følg Go kode standarder
- Tilføj tests til nye funktioner
- Opdater dokumentation
- Sørg for tværsplatform kompatibilitet

## Licens

**MIT Licens** - Se [LICENSE](LICENSE) filen for detaljer.

## Tak

- **Inspireret af**: Terraria spil mekanikker
- **Bygget med**: Ebiten spil motor
- **Bidragydere**: Open source fællesskab

## Support

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Diskussioner**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projekt Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Nyd udforskningen af TesselBox's hexagonale verden!*
