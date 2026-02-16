# TesselBox - README ya Kiswahili
## Mchezo wa Voxel ya Hexagonal

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Mchezo wa adventure ya sandbox ya 2D uliochochewa na *Terraria*, lakini uliojengwa kwenye **gridi ya hexagonal**.

Chunguza dunia, chimbua rasilimali, jenga miundo, tengeneza vitu, pigana na maadui na kuishi — yote katika tiles nzuri za hexagonal.

## Vipengele vya Mchezo

### ✅ **Vipengele Kamili**
- **Uzazi wa Ulimwengu wa Hexagonal** - Ulimwengu uliotengenezwa kwa utaratibu na biom
- **Uchimbaji na Ufundi** - Uchimbaji unaotegemea zana na kasi tofauti za nyenzo
- **Kuweka Block** - Bonyeza kulia kuweka block na onyesho la mzimu
- **Mfumo wa Inventory** - Inventory ya slot 32 na hotbar (slot 9)
- **Mfumo wa Mapambano** - Mfumo wa afya/ubaya na michoro ya mashambulizi
- **Mzunguko wa Mchana/Usiku** - Taa inayobadilika na maendeleo ya wakati
- **Madoido ya Hali ya Hewa** - Mifumo ya mvua, theluji na dhoruba
- **Mfumo wa Hifadhi/Pakia** - Hali ya ulimwengu inayodumu na hifadhi otomatiki

### 🎮 **Udhibiti**
- **WASD / Mishale**: Mwendo
- **Nafasi**: Ruka / Shambulia
- **Bonyeza Kushoto**: Uchimbaji wa block
- **Bonyeza Kulia**: Kuweka block
- **E**: Fungua menyu ya ufundi
- **Q**: Achia kitu kilichochaguliwa
- **Kgurumo ya Panya**: Uchaguzi wa hotbar
- **1-9**: Uchaguzi wa moja kwa moja wa hotbar
- **F5**: Hifadhi ya mwongozo
- **F9**: Upakiaji wa mwongozo
- **ESC**: Menyu / Funga menyu

## Usakinishaji na Kuweka

### Vigezo
- **Go 1.19+** - Injini kuu
- **Git** - Udhibiti wa toleo

### Kuanza Haraka
```bash
# Nakala hazina
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Jenga mchezo
go build ./cmd/client

# Endesha mchezo
./client
```

## Mahitaji ya Mfumo

### Kiwango cha Chini
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processor ya cores mbili
- **RAM**: 4GB
- **GPU**: Inayolingana na OpenGL 3.3+
- **Hifadhi**: Nafasi 500MB ya bure

### Inapendekezwa
- **CPU**: Processor ya cores nne
- **RAM**: 8GB+
- **GPU**: Kadi ya picha iliyotengwa
- **Hifadhi**: Nafasi 1GB+ ya bure

## Usanifu

### Teknolojia kuu
- **Lugha**: Go (Golang)
- **Michoro**: Ebiten (maktaba ya michezo ya 2D)
- **Mfumo wa Ujenzi**: Moduli za Go

### Muundo wa Mradi
```
TesselBox/
├── cmd/client/          # Faili kuu ya kutekeleza mchezo
├── pkg/                 # Pakiti kuu
│   ├── world/          # Uzazi na usimamizi wa ulimwengu
│   ├── player/         # Mechanics za mchezaji na fizikia
│   ├── blocks/         # Aina za block na sifa
│   ├── items/          # Mfumo wa vitu na ufundi
│   ├── crafting/       # Mapishi ya ufundi na UI
│   ├── weather/        # Uigaji wa hali ya hewa
│   ├── gametime/       # Mzunguko wa mchana/usiku
│   ├── save/           # Utendaji wa hifadhi/pakia
│   └── render/         # Mifumo ya utekelezaji na UI
├── config/             # Faili za usanidi
└── assets/             # Mali za mchezo (kama zipo)
```

## Kuchangia

### Kwa Wasanidi Programu
1. Fork hazina
2. Tengeneza tawi la kipengele (`git checkout -b feature/amazing-feature`)
3. Commit mabadiliko yako (`git commit -m 'Add amazing feature'`)
4. Push kwenye tawi (`git push origin feature/amazing-feature`)
5. Fungua Ombi la Kuvuta

### Miongozo ya Maendeleo
- Fuata viwango vya kuandika misimbo ya Go
- Ongeza majaribio kwa vipengele vipya
- Sasisha nyaraka
- Hakikisha uwezo wa kupita kwenye majukwaa

## Leseni

**Leseni ya MIT** - Angalia faili ya [LICENSE](LICENSE) kwa maelezo.

## Shukrani

- **Imechochewa na**: Mechanics za mchezo wa Terraria
- **Imejengwa kwa**: Injini ya michezo ya Ebiten
- **Wachangiaji**: Jumuiya ya chanzo wazi

## Usaidizi

- **Masuala**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Majadiliano**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki ya Mradi](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Furahia kuchunguza ulimwengu wa hexagonal wa TesselBox!*
