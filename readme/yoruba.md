# TesselBox - Yorùbá README
## Eré Voxel Hexagonal

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/TesselBox-game)

Eré àdúró sandbox 2D tí a gba agbára láti *Terraria*, ṣùgbọ́n tí a kọ sára **àwọn hexagonal grid**.

Ṣawari awọn agbaye, ṣe iṣẹ́ awọn orísun, kọ awọn ile-iṣẹ, ṣẹda awọn nkan, ba awọn ọta jagun ki o sì yè — gbogbo rẹ ni awọn tiles hexagonal ti o dara.

## Awọn ẹya Eré

### ✅ **Awọn ẹya Pipe**
- **Ipilẹṣẹ Agbaye Hexagonal** - Awọn agbaye ti a ṣẹda ni ilana pẹlu biomes
- **Iwakusa ati iṣẹ́-ọnà** - Iwakusa ti o da lori irinṣẹ pẹlu awọn iyara ohun elo ti o yatọ
- **Gbigbe Block** - Tẹ ọtun lati gbe awọn block pẹlu aworan iwin
- **Eto Inventory** - Inventory slot 32 pẹlu hotbar (awọn slot 9)
- **Eto Ija** - Eto ilera/ibajẹ pẹlu awọn animations ikọlu
- **Oṣu Ọjọ/Ale** - Imọlẹ ti o yipada ati ilọsiwaju akoko
- **Awọn ipa Oju ojo** - Awọn eto ojo, egbon ati iji
- **Eto Fi pamọ/Pakia** - Ipo agbaye ti o duro pẹlu fifipamọ aifọwọyi

### 🎮 **Awọn Iṣakoso**
- **WASD / Awọn itọka**: Gbigbe
- **Alafo**: Fo / Kọlu
- **Tẹ osi**: Iwakusa block
- **Tẹ ọtun**: Gbigbe block
- **E**: Ṣii akojọ aṣẹ iṣẹ́-ọnà
- **Q**: Ju nkan ti a yan silẹ
- **Kẹkẹ́ Mouse**: Yiyan hotbar
- **1-9**: Yiyan hotbar taara
- **F5**: Fifipamọ afọwọṣe
- **F9**: Pakia afọwọṣe
- **ESC**: Akojọ aṣẹ / Pa awọn akojọ aṣẹ

## Fifi sori ẹrọ ati Iṣeto

### Awọn ibeere
- **Go 1.19+** - Engine akọkọ
- **Git** - Iṣakoso ẹya

### Bibẹrẹ ni kiakia
```bash
# Ṣe ẹda repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Kọ ere
go build ./cmd/client

# Ṣiṣẹ ere
./client
```

## Awọn ibeere Eto

### O kere ju
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processor cores meji
- **RAM**: 4GB
- **GPU**: Ibamu pẹlu OpenGL 3.3+
- **Ibi ipamọ**: Aaye 500MB ọfẹ

### Iṣeduro
- **CPU**: Processor cores mẹrin
- **RAM**: 8GB+
- **GPU**: Kaadi ayaworan ti a yasọtọ
- **Ibi ipamọ**: Aaye 1GB+ ọfẹ

## Architecture

### Awọn Imọ-ẹrọ Akọkọ
- **Ede**: Go (Golang)
- **Awọn ayaworan**: Ebiten (ile-ikawe ere 2D)
- **Eto Ikọle**: Awọn modules Go

### Ilana Ise agbese
```
TesselBox/
├── cmd/client/          # Faili išẹ akọkọ ere
├── pkg/                 # Awọn idii akọkọ
│   ├── world/          # Ipilẹṣẹ ati iṣakoso agbaye
│   ├── player/         # Awọn ẹrọ orin ati fisiksi
│   ├── blocks/         # Awọn oriṣi block ati awọn abuda
│   ├── items/          # Eto awọn nkan ati iṣẹ́-ọnà
│   ├── crafting/       # Awọn ilana iṣẹ́-ọnà ati UI
│   ├── weather/        # Iṣe afihan oju ojo
│   ├── gametime/       # Oṣu ọjọ/ale
│   ├── save/           # Iṣẹ fifipamọ/pakia
│   └── render/         # Awọn eto išẹ ati UI
├── config/             # Awọn faili iṣeto
└── assets/             # Awọn ohun ini ere (ti o ba wa)
```

## Idasi

### Fun Awọn Olùṣe kọ̀dù
1. Fork repository naa
2. Ṣẹda ẹka ẹya kan (`git checkout -b feature/amazing-feature`)
3. Fi awọn ayipada rẹ pamọ (`git commit -m 'Add amazing feature'`)
4. Titari si ẹka naa (`git push origin feature/amazing-feature`)
5. Ṣii Ibeere Fa

### Awọn ilana Idagbasoke
- Tẹle awọn ilana koodu Go
- Fi awọn idanwo kun fun awọn ẹya tuntun
- Ṣe imudojuiwọn itọsọna
- Rii daju ibamu oriṣiriṣi platform

## Iwe-aṣẹ

**Iwe-aṣẹ MIT** - Wo faili [LICENSE](LICENSE) fun awọn alaye.

## Idupẹ

- **Ti gba agbára nipasẹ**: Awọn ẹrọ ere Terraria
- **Ti kọ pẹlu**: Engine ere Ebiten
- **Awọn oluranlọwọ**: Agbegbe orisun ṣiṣi

## Atilẹyin

- **Awọn ọrọ isoro**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Awọn ijiroro**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki Ise agbese](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Gbadun wiwa agbaye hexagonal ti TesselBox!*
