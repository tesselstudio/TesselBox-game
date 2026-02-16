# TesselBox - Afrikaans README
## Hexagonale Voxel Speletjie

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

'n 2D sandbox avontuurspel geïnspireer deur *Terraria*, maar gebou op 'n **hexagonale rooster**.

Verken wêrelde, ontgin hulpbronne, bou strukture, skep items, veg teen vyande en oorleef — alles in pragtige hexagonale teëls.

## Speletjie Kenmerke

### ✅ **Volledige Kenmerke**
- **Hexagonale Wêreld Generering** - Prosedureel gegenereerde wêrelde met biom
- **Ontginning en Vervaardiging** - Gereedskap-gebaseerde ontginning met verskillende materiaal spoed
- **Blok Plasing** - Regsklik om blokke te plaas met spook voorskou
- **Inventaris Stelsel** - 32-slot inventaris met hotbar (9 slotte)
- **Geveg Stelsel** - Gesondheid/skade stelsel met aanval animasies
- **Dag/Nag Siklus** - Dinamiese beligting en tyd vordering
- **Weer Effekte** - Reën, sneeu en storm stelsels
- **Stoor/Laai Stelsel** - Bestendige wêreld status met outomatiese stoor

### 🎮 **Beheer**
- **WASD / Pyle**: Beweging
- **Spasie**: Spring / Aanval
- **Linker Klik**: Blok ontginning
- **Regter Klik**: Blok plasing
- **E**: Maak vervaardiging kieslys oop
- **Q**: Laat geselekteerde item val
- **Muiswiel**: Hotbar seleksie
- **1-9**: Direkte hotbar seleksie
- **F5**: Handmatige stoor
- **F9**: Handmatige laai
- **ESC**: Kieslys / Sluit kieslyste

## Installasie en Opstelling

### Vereistes
- **Go 1.19+** - Kern enjin
- **Git** - Weergawe beheer

### Vinnige Begin
```bash
# Klooneer repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Bou speletjie
go build ./cmd/client

# Begin speletjie
./client
```

### Ontwikkeling Opstelling
```bash
# Installeer afhanklikhede
go mod tidy

# Voer toetse uit
go test ./...

# Bou vir ontwikkeling
go build -tags debug ./cmd/client
```

## Stelsel Vereistes

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dubbelkern verwerker
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ versoenbaar
- **Berging**: 500MB vrye ruimte

### Aanbeveel
- **CPU**: Vierkern verwerker
- **RAM**: 8GB+
- **GPU**: Toegewyde grafiese kaart
- **Berging**: 1GB+ vrye ruimte

## Argitektuur

### Kern Tegnologieë
- **Taal**: Go (Golang)
- **Grafika**: Ebiten (2D speletjie biblioteek)
- **Boustelsel**: Go modules

### Projek Struktuur
```
TesselBox/
├── cmd/client/          # Hoof speletjie uitvoerbare lêer
├── pkg/                 # Kern pakkette
│   ├── world/          # Wêreld generering en bestuur
│   ├── player/         # Speler meganika en fisika
│   ├── blocks/         # Blok tipes en eienskappe
│   ├── items/          # Item stelsel en vervaardiging
│   ├── crafting/       # Vervaardiging resepte en UI
│   ├── weather/        # Weer simulasie
│   ├── gametime/       # Dag/nag siklus
│   ├── save/           # Stoor/laai funksionaliteit
│   └── render/         # Rendering en UI stelsels
├── config/             # Konfigurasie lêers
└── assets/             # Speletjie bates (indien enige)
```

## Bydra

### Vir Ontwikkelaars
1. Fork die repository
2. Skep 'n kenmerk tak (`git checkout -b feature/amazing-feature`)
3. Commit jou veranderinge (`git commit -m 'Add amazing feature'`)
4. Stuur na die tak (`git push origin feature/amazing-feature`)
5. Maak 'n Pull Request oop

### Ontwikkeling Riglyne
- Volg Go kodering standaarde
- Voeg toetse by vir nuwe kenmerke
- Dateer dokumentasie op
- Verseker kruisplatform verenigbaarheid

## Lisensie

**MIT Lisensie** - Sien [LICENSE](LICENSE) lêer vir besonderhede.

## Erkennings

- **Geïnspireer deur**: Terraria speletjie meganika
- **Gebou met**: Ebiten speletjie enjin
- **Bydraers**: Oopbron gemeenskap

## Ondersteuning

- **Kwessies**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Besprekings**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projek Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Geniet die verkenning van TesselBox se hexagonale wêreld!*
