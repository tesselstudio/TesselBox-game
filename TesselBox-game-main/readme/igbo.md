# TesselBox - Igbo README
## Egwuregwu Voxel Hexagonal

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Egwuregwu sandbox 2D nke sitere na *Terraria*, mana ewuru ya na **hexagonal grid**.

Chọpụta ụwa, kụpụ ihe onwunwe, wuo ihe owuwu, mepụta ihe, buso ndị iro agha ma dị ndụ — niile na ezigbo tiles hexagonal.

## Atụmatụ Egwuregwu

### ✅ **Atụmatụ zuru ezu**
- **Mmepụta Ụwa Hexagonal** - Ụwa ndị a na-emepụta n'usoro nwere biomes
- **Ịkụpụ na Ịkpụzi** - Ịkụpụ dabere na ngwaọrụ nwere ọsọ ihe dị iche iche
- **Idobe Block** - Pịa aka nri iji dobere blocks nwere ihe ngosi mmụọ
- **Usoro Inventory** - Inventory slot 32 nwere hotbar (slots 9)
- **Usoro Ọgụ** - Usoro ahụike/mmebi nwere animations ọgụ
- **Usoro Ụbọchị/Abalị** - Ìhè na-agbanwe agbanwe na ọganihu oge
- **Mmetụta Ihu igwe** - Usoro mmiri ozuzo, snow na oké ifufe
- **Usoro Chekwaa/Load** - Ọnọdụ ụwa na-adịgide adịgide nwere nchekwa akpaaka

### 🎮 **Njikwa**
- **WASD / Akara ngosi**: Mfegharị
- **Oghere**: Wụlikwa / Wakpo
- **Pịa aka ekpe**: Ịkụpụ block
- **Pịa aka nri**: Idobe block
- **E**: Mepee menu ịkpụzi
- **Q**: Hapu ihe ahọpụtara
- **Ụgbọ mmiri ozi**: Nhọrọ hotbar
- **1-9**: Nhọrọ hotbar ozugbo
- **F5**: Nchekwa aka
- **F9**: Load aka
- **ESC**: Menu / Mechie menu

## Nrụnye na Nhazi

### Ihe achọrọ
- **Go 1.19+** - Engine isi
- **Git** - Njikwa ụdị

### Mbido ngwa ngwa
```bash
# Detuo ebe nchekwa
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Wuo egwuregwu
go build ./cmd/client

# Malite egwuregwu
./client
```

## Ihe achọrọ Sistemụ

### Kacha nta
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processor cores abụọ
- **RAM**: 4GB
- **GPU**: Dakọtara na OpenGL 3.3+
- **Nchekwa**: Oghere 500MB efu

### Atụ aro
- **CPU**: Processor cores anọ
- **RAM**: 8GB+
- **GPU**: Kaadị eserese akọwapụtara
- **Nchekwa**: Oghere 1GB+ efu

## Architecture

### Teknụzụ Isi
- **Asụsụ**: Go (Golang)
- **Eserese**: Ebiten (ọbá akwụkwọ egwuregwu 2D)
- **Sistemụ Wu**：Modules Go

### Nhazi Ọrụ
```
TesselBox/
├── cmd/client/          # Isi faịlụ egwuregwu executable
├── pkg/                 # Isi ngwugwu
│   ├── world/          # Mmepụta na njikwa ụwa
│   ├── player/         # Ụzọ egwuregwu na physics
│   ├── blocks/         # Ụdị blocks na akụrụngwa
│   ├── items/          # Usoro ihe na ịkpụzi
│   ├── crafting/       # Ntuziaka ịkpụzi na UI
│   ├── weather/        # Simulation ihu igwe
│   ├── gametime/       # Usoro ụbọchị/abali
│   ├── save/           # Ọrụ chekwaa/load
│   └── render/         # Usoro eserese na UI
├── config/             # Faịlụ nhazi
└── assets/             # Ihe onwunwe egwuregwu (ọ bụrụ adị)
```

## Onyinye

### Maka Ndị Mepụta Kọdụ
1. Fork ebe nchekwa
2. Mepụta alaka atụmatụ (`git checkout -b feature/amazing-feature`)
3. Kọmiti mgbanwe gị (`git commit -m 'Add amazing feature'`)
4. Bugharịa na alaka (`git push origin feature/amazing-feature`)
5. Mepee Arịrịọ Mbugharị

### Ụkpụrụ Mmepe
- Gbaso ụkpụrụ kọdụ Go
- Tinye ule maka atụmatụ ọhụrụ
- Melite akwụkwọ ntuziaka
- Gbaa mbọ hụ na ndakọrịta n'etiti ikpo okwu

## Ikike

**Ikike CC BY-NC-SA 4.0** - Lee faịlụ [LICENSE](LICENSE) maka nkọwa.

## Ekele

- **Sitere na**: Ụzọ egwuregwu Terraria
- **Wuru na**: Engine egwuregwu Ebiten
- **Ndị nyere aka**: Obodo e wepụtara

## Nkwado

- **Nsogbu**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Mkparịta ụka**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki Ọrụ](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Nwee obi ụtọ na nyocha ụwa hexagonal nke TesselBox!*
