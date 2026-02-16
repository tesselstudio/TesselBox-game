# TesselBox - Svenska README
## Hexagonal Voxel-spel

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Ett 2D sandbox-äventyrsspel inspirerat av *Terraria*, men byggt på ett **hexagonalt rutnät**.

Utforska världar, utvinna resurser, bygg strukturer, skapa föremål, kämpa mot fiender och överlev — allt i vackra hexagonala plattor.

## Spel Funktioner

### ✅ **Kompletta Funktioner**
- **Hexagonal Världs Generering** - Procedurellt genererade världar med biom
- **Brytning och Tillverkning** - Verktygsbaserad brytning med olika materialhastigheter
- **Block Placering** - Högerklicka för att placera block med spökförhandsvisning
- **Inventarie System** - 32-slotars inventarie med snabbverktygsfält (9 slots)
- **Strids System** - Hälsa/skadesystem med attackanimeringar
- **Dag/Natt Cykel** - Dynamisk belysning och tidsframsteg
- **Väder Effekter** - Regn, snö och storms system
- **Spara/Ladda System** - Beständig världstillstånd med autosparande

### 🎮 **Kontroller**
- **WASD / Piltangenter**: Rörelse
- **Mellanslag**: Hoppa / Attackera
- **Vänster Klick**: Block brytning
- **Höger Klick**: Block placering
- **E**: Öppna tillverkning meny
- **Q**: Släpp valt föremål
- **Mushjul**: Snabbverktygsfält val
- **1-9**: Direkt snabbverktygsfält val
- **F5**: Manuell sparande
- **F9**: Manuell laddning
- **ESC**: Meny / Stäng menyer

## Installation och Setup

### Förutsättningar
- **Go 1.19+** - Kärnmotor
- **Git** - Versionshantering

### Snabb Start
```bash
# Klona repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Bygg spel
go build ./cmd/client

# Starta spel
./client
```

### Utvecklings Setup
```bash
# Installera beroenden
go mod tidy

# Kör tester
go test ./...

# Bygg för utveckling
go build -tags debug ./cmd/client
```

## Systemkrav

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dubbelkärnig processor
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ kompatibel
- **Lagring**: 500MB ledigt utrymme

### Rekommenderat
- **CPU**: Fyrkärnig processor
- **RAM**: 8GB+
- **GPU**: Dedikerat grafikkort
- **Lagring**: 1GB+ ledigt utrymme

## Arkitektur

### Kärn Teknologier
- **Språk**: Go (Golang)
- **Grafik**: Ebiten (2D-spelbibliotek)
- **Byggsystem**: Go-moduler

### Projekt Struktur
```
TesselBox/
├── cmd/client/          # Huvudspel körbar fil
├── pkg/                 # Kärnpaket
│   ├── world/          # Världs generering och hantering
│   ├── player/         # Spelar mekaniker och fysik
│   ├── blocks/         # Block typer och egenskaper
│   ├── items/          # Föremålssystem och tillverkning
│   ├── crafting/       # Tillverknings recept och UI
│   ├── weather/        # Väder simulering
│   ├── gametime/       # Dag/natt cykel
│   ├── save/           # Spara/ladda funktionalitet
│   └── render/         # Rendering och UI system
├── config/             # Konfigurationsfiler
└── assets/             # Spel tillgångar (om några finns)
```

## Bidra

### För Utvecklare
1. Forka repository
2. Skapa en funktionsgren (`git checkout -b feature/amazing-feature`)
3. Committa dina ändringar (`git commit -m 'Add amazing feature'`)
4. Pusha till grenen (`git push origin feature/amazing-feature`)
5. Öppna en Pull Request

### Utvecklings Riktlinjer
- Följ Go-kodningsstandarder
- Lägg till tester för nya funktioner
- Uppdatera dokumentation
- Säkerställ plattformsoberoende kompatibilitet

## Licens

**MIT Licens** - Se [LICENSE](LICENSE)-filen för detaljer.

## Credits

- **Inspirerad av**: Terraria spelmekaniker
- **Byggd med**: Ebiten spel motor
- **Bidragsgivare**: Open source community

## Support

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Diskussioner**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projekt Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Njut av att utforska TesselBox hexagonala värld!*
