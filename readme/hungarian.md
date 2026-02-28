# TesselBox - Magyar README
## Hexagonális Voxel Játék

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Egy 2D sandbox kalandjáték, amely a *Terraria* ihlette, de **hexagonális rácson** épült.

Fedezd fel a világokat, bányássz erőforrásokat, építs struktúrákat, készíts tárgyakat, harcolj az ellenségekkel és élj túl — minden a gyönyörű hexagonális csempékben.

## Játék Jellemzők

### ✅ **Teljes Jellemzők**
- **Hexagonális Világ Generálás** - Eljárásosan generált világok biomokkal
- **Bányászat és Készítés** - Eszközalapú bányászat különböző anyaghastásokkal
- **Blokk Elhelyezés** - Jobb klikk a blokkok elhelyezéséhez szellem előnézettel
- **Leltár Rendszer** - 32-slot leltár hotbar-ral (9 slot)
- **Harc Rendszer** - Egészség/sérülés rendszer támadási animációkkal
- **Nap/Éjszaka Ciklus** - Változó megvilágítás és idő előrehaladás
- **Időjárás Hatások** - Eső, hó és vihar rendszerek
- **Mentés/Betöltés Rendszer** - Állandó világállapot automatikus mentéssel

### 🎮 **Irányítás**
- **WASD / Nyilak**: Mozgás
- **Szóköz**: Ugrás / Támadás
- **Bal Klikk**: Blokk bányászat
- **Jobb Klikk**: Blokk elhelyezés
- **E**: Készítés menü megnyitása
- **Q**: Kiválasztott tárgy elengedése
- **Egér Görgő**: Hotbar kiválasztás
- **1-9**: Közvetlen hotbar kiválasztás
- **F5**: Manuális mentés
- **F9**: Manuális betöltés
- **ESC**: Menü / Menük bezárása

## Telepítés és Beállítás

### Követelmények
- **Go 1.19+** - Fő motor
- **Git** - Verziókezelés

### Gyors Kezdés
```bash
# Repository klónozása
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Játék építése
go build ./cmd/client

# Játék indítása
./client
```

## Rendszer Követelmények

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Kétmagos processzor
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ kompatibilis
- **Tárhely**: 500MB szabad hely

### Ajánlott
- **CPU**: Négy magos processzor
- **RAM**: 8GB+
- **GPU**: Dedikált grafikus kártya
- **Tárhely**: 1GB+ szabad hely

## Architektúra

### Fő Technológiák
- **Nyelv**: Go (Golang)
- **Grafika**: Ebiten (2D játék könyvtár)
- **Építési Rendszer**: Go modulok

### Projekt Struktúra
```
TesselBox/
├── cmd/client/          # Fő játék futtatható fájl
├── pkg/                 # Fő csomagok
│   ├── world/          # Világ generálás és kezelés
│   ├── player/         # Játékos mechanikák és fizika
│   ├── blocks/         # Blokk típusok és tulajdonságok
│   ├── items/          # Tárgy rendszer és készítés
│   ├── crafting/       # Készítési receptek és UI
│   ├── weather/        # Időjárás szimuláció
│   ├── gametime/       # Nap/éjszaka ciklus
│   ├── save/           # Mentés/betöltés funkcionalitás
│   └── render/         # Renderelés és UI rendszerek
├── config/             # Konfigurációs fájlok
└── assets/             # Játék eszközök (ha vannak)
```

## Közreműködés

### Fejlesztőknek
1. Forkold a repositoryt
2. Hozz létre egy feature ágat (`git checkout -b feature/amazing-feature`)
3. Commitold a változásaidat (`git commit -m 'Add amazing feature'`)
4. Pushold az ágba (`git push origin feature/amazing-feature`)
5. Nyiss egy Pull Requestet

### Fejlesztési Irányelvek
- Kövesd a Go kódolási szabványokat
- Adj hozzá teszteket az új funkciókhoz
- Frissítsd a dokumentációt
- Biztosítsd a keresztplatform kompatibilitást

## Licenc

**CC BY-NC-SA 4.0 Licenc** - Lásd a [LICENSE](LICENSE) fájlt a részletekért.

## Köszönet

- **Inspirálta**: A Terraria játék mechanikái
- **Épült**: Ebiten játék motorral
- **Közreműködők**: Nyílt forráskódú közösség

## Támogatás

- **Problémák**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Fórum**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projekt Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Élvezd a TesselBox hexagonális világának felfedezését!*
