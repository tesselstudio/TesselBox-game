# TesselBox - Zulu README
## Umdlalo we-Voxel we-Hexagonal

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Umdlalo we-sandbox we-2D otholakala ku-*Terraria*, kodwa wakhiwe ku-**hexagonal grid**.

Hlola imihlaba, zimbela izinsiza, akha izakhiwo, dala izinto, ulwe nezitha futhi uphile — konke kuma-tiles amahle e-hexagonal.

## Izici Zomdlalo

### ✅ **Izici Ezigcwele**
- **Ukwakhiwa Komhlaba we-Hexagonal** - Imihlaba eyakhiwe ngokuququzelelayo enama-biomes
- **Ukumba nokwakha** - Ukumba okusekelwe kumathuluzi ngesivinini sezinto ezahlukene
- **Ukubeka i-Block** - Chofoza kwesokudla ukubeka ama-blocks ngemboniswano ye-spirit
- **Uhlelo lwe-Inventory** - I-inventory ye-slot engu-32 ne-hotbar (ama-slots angu-9)
- **Uhlelo Lokulwa** - Uhlelo lwezempilo/ukulimala ngama-animations okuhlasela
- **Uhlelo Losuku/Ebusuku** - Ukukhanya okushintshayo nentuthuko yesikhathi
- **Imithelela Yesimo Sezulu** - Izinhlelo zemvula, iqhwa nesiphepho
- **Uhlelo Lokulondoloza/Ukulayisha** - Isimo somhlaba esihlala njalo ngokulondoloza okuzenzakalelayo

### 🎮 **Ukuphathwa**
- **WASD / Imicibisholo**: Ukunyakaza
- **Isikhala**: Gxuma / Hlasela
- **Chofoza kwesobunxele**: Ukumba i-block
- **Chofoza kwesokudla**: Ukubeka i-block
- **E**: Vula imenyu yokwakha
- **Q**: Dedela into ekhethiwe
- **I-wheel yemouse**: Ukukhetha i-hotbar
- **1-9**: Ukukhetha i-hotbar okuqondile
- **F5**: Ukulondoloza ngesandla
- **F9**: Ukulayisha ngesandla
- **ESC**: Imenyu / Vala amamenyu

## Ukufakwa nokusetha

### Izidingo
- **Go 1.19+** - Injini eyinhloko
- **Git** - Ukulawula inguqulo

### Ukuqala Okusheshayo
```bash
# Klona i-repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Yakha umdlalo
go build ./cmd/client

# Qalisa umdlalo
./client
```

## Izidingo Zohlelo

### Ezincane kakhulu
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Iprosesa enamakhodi amabili
- **RAM**: 4GB
- **GPU**: Ehambisana ne-OpenGL 3.3+
- **Isitoreji**: Isikhala esingama-500MB esikhululekile

### Okunconywayo
- **CPU**: Iprosesa enamakhodi amane
- **RAM**: 8GB+
- **GPU**: Ikhadi lemifanekiso elikhethekile
- **Isitoreji**: Isikhala esingama-1GB+ esikhululekile

## Architecture

### Ubuchwepheshe Obuyinhloko
- **Ulimi**: Go (Golang)
- **Imidwebo**: Ebiten (umtapo wezincwadi womdlalo we-2D)
- **Uhlelo Lokwakha**: Ama-modules e-Go

### Isakhiwo Somsebenzi
```
TesselBox/
├── cmd/client/          # Ifayela eliyinhloko lokusebenza komdlalo
├── pkg/                 # Ama-package ayinhloko
│   ├── world/          # Ukwakhiwa nokulawulwa komhlaba
│   ├── player/         # Izindlela zomdlali ne-physics
│   ├── blocks/         # Izinhlobo zama-blocks nezici
│   ├── items/          # Uhlelo lwezinto nokwakha
│   ├── crafting/       # Imiyalo yokwakha ne-UI
│   ├── weather/        # Ukulingisa kwesimo sezulu
│   ├── gametime/       # Uhlelo losuku/ebusuku
│   ├── save/           # Umsebenzi wokulondoloza/ukulayisha
│   └── render/         # Izinhlelo zokudweba ne-UI
├── config/             # Amafayela wokusetha
└── assets/             # Impahla yomdlalo (uma ikhona)
```

## Ukunikela

### Kubathuthukisi Bekhodi
1. Fork i-repository
2. Dala igatsha lesici (`git checkout -b feature/amazing-feature`)
3. Commit izinguquko zakho (`git commit -m 'Add amazing feature'`)
4. Phusha egatsheni (`git push origin feature/amazing-feature`)
5. Vula Isicelo Sokudonsa

### Imigomo Yokuthuthukisa
- Landela imigomo yekhodi ye-Go
- Engeza ukuhlola kwezici ezintsha
- Buyekeza imibhalo
- Qinisekisa ukuhambisana kwamapulatifomu ahlukene

## Ilayisense

**Ilayisense ye-MIT** - Buka ifayela [LICENSE](LICENSE) lemininingwane.

## Ukubonga

- **Kutholakala ku**: Izindlela zomdlalo we-Terraria
- **Yakhiwe nge**: Injini yomdlalo ye-Ebiten
- **Abanikeli**: Umphakathi womthombo ovulekile

## Ukwesekwa

- **Izinkinga**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Izingxoxo**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **I-Wiki**: [I-Wiki Yomsebenzi](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Jabulela ukuhlola umhlaba we-hexagonal we-TesselBox!*
