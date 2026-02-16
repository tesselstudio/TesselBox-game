# TesselBox - Hausa README
## Wasan Voxel na Hexagonal

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/TesselBox-game)

Wasan sandbox na 2D wanda aka yi shi ne daga *Terraria*, amma an gina shi akan **hexagonal grid**.

Bincika duniyoyi, hakuda albarkatu, gina tsarukan, ƙirƙiri abubuwa, yaƙi da maƙiyan kuma tsira — duk a cikin kyawawan tiles na hexagonal.

## Fasali na Wasan

### ✅ **Cikakken Fasali**
- **Ƙirƙirar Duniya ta Hexagonal** - Duniyoyi da aka ƙirƙira ta hanyar tsari tare da biomes
- **Hakuda da Sana'a** - Hakuda ta hanyar kayan aiki tare da saurin kayan daban-daban
- **Sanya Block** - Danna dama don sanya blocks tare da bayanin fatalwa
- **Tsarin Inventory** - Inventory na slot 32 tare da hotbar (slots 9)
- **Tsarin Yaƙi** - Tsarin lafiya/lalacewa tare da animations na hari
- **Tsarin Rana/Dare** - Hasashen haske mai canzawa da ci gaban lokaci
- **Tasirin Yanayi** - Tsarin ruwan sama, dusar ƙanƙara da guguwa
- **Tsarin Ajiya/Load** - Tsarin duniya mai dindindin tare da ajiya ta atomatik

### 🎮 **Sarrafa**
- **WASD / Kibiya**：Girma
- **Sarari**: Tsalle / Hari
- **Danna hagu**: Hakuda block
- **Danna dama**: Sanya block
- **E**: Buɗe menu na sana'a
- **Q**: Bar abin da aka zaɓa
- **Karamar linzamin kwamfuta**: Zaɓin hotbar
- **1-9**: Zaɓin hotbar kai tsaye
- **F5**: Ajiya ta hannu
- **F9**: Load ta hannu
- **ESC**: Menu / Rufe menu

## Shigarwa da Saitawa

### Bukatu
- **Go 1.19+** - Babban injin
- **Git** - Sarrafa sigar

### Farawa mai sauri
```bash
# Kwafi repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Gina wasan
go build ./cmd/client

# Fara wasan
./client
```

## Bukatun Tsarin

### Mafi ƙanƙanta
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processor mai cores biyu
- **RAM**: 4GB
- **GPU**: Daidai da OpenGL 3.3+
- **Ajiya**: Sarari 500MB

### Shawara
- **CPU**: Processor mai cores huɗu
- **RAM**: 8GB+
- **GPU**: Katin zane na musamman
- **Ajiya**: Sarari 1GB+

## Architecture

### Babban Fasaha
- **Harshe**: Go (Golang)
- **Zane**: Ebiten (laburaree wasan 2D)
- **Tsarin Gina**: Modules na Go

### Tsarin Aikin
```
TesselBox/
├── cmd/client/          # Babban fayil na executable na wasan
├── pkg/                 # Babban fakiti
│   ├── world/          # Ƙirƙirar duniya da sarrafa
│   ├── player/         # Hanyoyin ɗan wasa da physics
│   ├── blocks/         # Nau'ikan block da kaddarorin
│   ├── items/          # Tsarin abubuwa da sana'a
│   ├── crafting/       # Girke-girke na sana'a da UI
│   ├── weather/        # Simulation na yanayi
│   ├── gametime/       # Tsarin rana/dare
│   ├── save/           # Aiki na ajiya/load
│   └── render/         # Tsarin zana da UI
├── config/             # Fayilolin saiti
└── assets/             # Abubuwan wasan (idanyesu)
```

## Bayar da gudummawa

### Ga Masu Shirin Shirye-shirye
1. Fork repository
2. Ƙirƙiri reshe na fasali (`git checkout -b feature/amazing-feature`)
3. Commit canje-canje naku (`git commit -m 'Add amazing feature'`)
4. Tura zuwa reshen (`git push origin feature/amazing-feature`)
5. Buɗe Buƙatar Ja

### Ka'idojin Haɓakawa
- Bi ka'idojin rubuta lambar Go
- Ƙara gwaje-gwaje don sabbin fasali
- Sabunta takardu
- Tabbatar da daidaitawa tsakanin dandamali

## Lasisi

**Lasisin MIT** - Duba fayil [LICENSE](LICENSE) don cikakkun bayanai.

## Godiya

- **An yi shi ne daga**: Hanyoyin wasan Terraria
- **An gina shi da**: Injin wasan Ebiten
- **Masu bayar da gudummawa**: Al'ummar buɗaɗɗen tushe

## Tallafi

- **Lamuran**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Tattaunawa**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki na Aikin](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Ji daɗin binciken duniyar hexagonal na TesselBox!*
