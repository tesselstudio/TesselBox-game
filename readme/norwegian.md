# TesselBox - Norsk README
## Hexagonal Voxel-spill

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Et 2D sandbox eventyrspill inspirert av *Terraria*, men bygget på et **hexagonalt rutenett**.

Utforsk verdener, utvin ressurser, bygg strukturer, lag gjenstander, kjemp mot fiender og overlev — alt i vakre hexagonale fliser.

## Spillfunksjoner

### ✅ **Komplette Funksjoner**
- **Hexagonal Verdensgenerering** - Prosedyremessig genererte verdener med biomer
- **Utvinning og Håndverk** - Verktøybasert utvinning med forskjellige materialhastigheter
- **Blokkplassering** - Høyreklikk for å plassere blokker med spøkelses forhåndsvisning
- **Inventarsystem** - 32-slot inventar med hotbar (9 spor)
- **Kampsystem** - Helse/skade system med angrepsanimasjoner
- **Dag/Natt Syklus** - Skiftende belysning og tidsfremgang
- **Værefekter** - Regn, snø og storm systemer
- **Lagre/Laste System** - Vedvarende verdensstatus med automatisk lagring

### 🎮 **Kontroll**
- **WASD / Piler**: Bevegelse
- **Mellomrom**: Hopp / Angrip
- **Venstre Klikk**: Blokk utvinning
- **Høyre Klikk**: Blokk plassering
- **E**: Åpne håndverksmeny
- **Q**: Slipp valgte gjenstand
- **Musehjul**: Hotbar valg
- **1-9**: Direkte hotbar valg
- **F5**: Manuell lagring
- **F9**: Manuell lasting
- **ESC**: Meny / Lukk menyer

## Installasjon og Oppsett

### Krav
- **Go 1.19+** - Kjerne motor
- **Git** - Versjonskontroll

### Hurtigstart
```bash
# Klon repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Bygg spill
go build ./cmd/client

# Start spill
./client
```

## Systemkrav

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Tokerne prosessor
- **RAM**: 4GB
- **GPU**: Kompatibel med OpenGL 3.3+
- **Lagring**: 500MB ledig plass

### Anbefalt
- **CPU**: Firekerne prosessor
- **RAM**: 8GB+
- **GPU**: Dedikert grafikkort
- **Lagring**: 1GB+ ledig plass

## Arkitektur

### Kjerneteknologier
- **Språk**: Go (Golang)
- **Grafikk**: Ebiten (2D spillbibliotek)
- **Byggesystem**: Go moduler

### Prosjektstruktur
```
TesselBox/
├── cmd/client/          # Hoved spill kjørbar fil
├── pkg/                 # Kjernespakker
│   ├── world/          # Verdensgenerering og styring
│   ├── player/         # Spiller mekanikker og fysikk
│   ├── blocks/         # Blokktyper og egenskaper
│   ├── items/          # Gjenstandssystem og håndverk
│   ├── crafting/       # Håndverksoppskrifter og UI
│   ├── weather/        # Vær simulering
│   ├── gametime/       # Dag/natt syklus
│   ├── save/           # Lagre/laste funksjonalitet
│   └── render/         # Rendering og UI systemer
├── config/             # Konfigurasjonsfiler
└── assets/             # Spillressurser (hvis noen)
```

## Bidra

### For Utviklere
1. Fork repository
2. Opprett en funksjonsgren (`git checkout -b feature/amazing-feature`)
3. Commit endringene dine (`git commit -m 'Add amazing feature'`)
4. Push til grenen (`git push origin feature/amazing-feature`)
5. Åpne en Pull Request

### Utviklings Retningslinjer
- Følg Go-kode standarder
- Legg til tester for nye funksjoner
- Oppdater dokumentasjon
- Sørg for tverrplattform kompatibilitet

## Lisens

**MIT Lisens** - Se [LICENSE](LICENSE) filen for detaljer.

## Takk

- **Inspirert av**: Terraria spillmekanikker
- **Bygget med**: Ebiten spillmotor
- **Bidragsytere**: Åpen kildekode fellesskap

## Støtte

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Diskusjoner**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Prosjekt Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

* Nyt utforskningen av TesselBox' hexagonale verden! *
