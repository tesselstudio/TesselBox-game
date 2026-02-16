# TesselBox - Suomi README
## Hexagoninen Voxel -peli

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

2D hiekkalaatikko-seikkailupeli, joka on saanut inspiraation *Terraria*sta, mutta rakennettu **hexagoniselle ruudukolle**.

Tutki maailmoja, louhi resursseja, rakenna rakenteita, luo esineitä, taistele vihollisia vastaan ja selviydy — kaikki kauniissa hexagonisissa laatoissa.

## Pelin Ominaisuudet

### ✅ **Täydelliset Ominaisuudet**
- **Hexagoninen Maailman Generointi** - Menettelyllisesti generoidut maailmat biomeineen
- **Louhinta ja Ammattitaito** - Työkaluun perustuva louhinta erilaisilla materiaali nopeuksilla
- **Blokin Sijoittaminen** - Oikea klikkaus blokkien sijoittamiseen haamu-esikatselulla
- **Inventaariojärjestelmä** - 32-paikkainen inventaario hotbarilla (9 paikkaa)
- **Taistelujärjestelmä** - Terveys/vahinko järjestelmä hyökkäysanimaatioilla
- **Päivä/Yö Syklit** - Vaihtuva valaistus ja ajan eteneminen
- **Säävaikutukset** - Sade, lumi ja myrsky järjestelmät
- **Tallenna/Lataa Järjestelmä** - Pysyvä maailman tila automaattisella tallennuksella

### 🎮 **Ohjaus**
- **WASD / Nuolet**: Liikkuminen
- **Välilyönti**: Hyppy / Hyökkäys
- **Vasen Klikkaus**: Blokin louhinta
- **Oikea Klikkaus**: Blokin sijoittaminen
- **E**: Avaa ammattitaito valikko
- **Q**: Pudota valittu esine
- **Hiiren Rulla**: Hotbar valinta
- **1-9**: Suora hotbar valinta
- **F5**: Manuaalinen tallennus
- **F9**: Manuaalinen lataus
- **ESC**: Valikko / Sulje valikot

## Asennus ja Asetukset

### Vaatimukset
- **Go 1.19+** - Päämoottori
- **Git** - Versiohallinta

### Pika-aloitus
```bash
# Kloonaa repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Rakenna peli
go build ./cmd/client

# Käynnistä peli
./client
```

## Järjestelmävaatimukset

### Minimivaatimukset
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Kaksiytiminen prosessori
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ yhteensopiva
- **Tallennustila**: 500MB vapaata tilaa

### Suositellut
- **CPU**: Neliytiminen prosessori
- **RAM**: 8GB+
- **GPU**: Oma näytönohjain
- **Tallennustila**: 1GB+ vapaata tilaa

## Arkkitehtuuri

### Ydinteknologiat
- **Kieli**: Go (Golang)
- **Grafiikka**: Ebiten (2D pelikirjasto)
- **Rakennejärjestelmä**: Go-modulit

### Projektin Rakenne
```
TesselBox/
├── cmd/client/          # Pään pelin suoritettava tiedosto
├── pkg/                 # Ydinkirjastot
│   ├── world/          # Maailman generointi ja hallinta
│   ├── player/         # Pelaajan mekaniikat ja fysiikka
│   ├── blocks/         # Blokki tyypit ja ominaisuudet
│   ├── items/          # Esine järjestelmä ja ammattitaito
│   ├── crafting/       # Ammattitaito reseptit ja UI
│   ├── weather/        # Sää simulointi
│   ├── gametime/       # Päivä/yö sykli
│   ├── save/           # Tallenna/lataa toiminnallisuus
│   └── render/         # Renderöinti ja UI järjestelmät
├── config/             # Konfiguraatiotiedostot
└── assets/             # Pelin resurssit (jos olemassa)
```

## Osallistuminen

### Kehittäjille
1. Forkkaa repository
2. Luo ominaisuus haara (`git checkout -b feature/amazing-feature`)
3. Commitoi muutoksesi (`git commit -m 'Add amazing feature'`)
4. Puskaa haaraan (`git push origin feature/amazing-feature`)
5. Avaa Pull Request

### Kehitys Ohjeet
- Noudata Go-koodausstandardeja
- Lisää testejä uusiin ominaisuuksiin
- Päivitä dokumentaatio
- Varmista monialustaisuus yhteensopivuus

## Lisenssi

**MIT-lisenssi** - Katso [LICENSE](LICENSE) tiedosto yksityiskohdista.

## Kiitokset

- **Inspiroitu**: Terraria peli mekaniikoista
- **Rakennettu**: Ebiten pelimoottorilla
- **Avustajat**: Avoimen lähdekoodin yhteisö

## Tuki

- **Ongelmat**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Keskustelut**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projektin Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Nauti TesselBoxin hexagonisen maailman tutkimisesta!*
