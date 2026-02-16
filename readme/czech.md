# TesselBox - Česká README
## Hexagonální Voxelová Hra

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

2D sandbox adventurní hra inspirovaná *Terraria*, ale postavená na **hexagonální mřížce**.

Prozkoumejte světy, těžte zdroje, stavějte struktury, vytvářejte předměty, bojujte s nepřáteli a přežijte — vše v krásných hexagonálních dlaždicích.

## Funkce Hry

### ✅ **Kompletní Funkce**
- **Hexagonální Generování Světa** - Procedurálně generované světy s biomi
- **Těžba a Výroba** - Těžba založená na nástrojích s různými rychlostmi materiálů
- **Umístění Bloků** - Pravé kliknutí pro umístění bloků s náhledem ducha
- **Systém Inventáře** - 32-slotový inventář s rychlým panelem (9 slotů)
- **Bojový Systém** - Systém zdraví/poškození s animacemi útoků
- **Cyklus Den/Noc** - Dynamické osvětlení a časový postup
- **Efekty Počasí** - Systémy deště, sněhu a bouře
- **Systém Uložení/Nahrání** - Trvalý stav světa s automatickým ukládáním

### 🎮 **Ovládání**
- **WASD / Šipky**: Pohyb
- **Mezerník**: Skok / Útok
- **Levé Kliknutí**: Těžba bloků
- **Pravé Kliknutí**: Umístění bloků
- **E**: Otevřít menu výroby
- **Q**: Položit vybraný předmět
- **Kolečko Myši**: Výběr rychlého panelu
- **1-9**: Přímý výběr rychlého panelu
- **F5**: Manuální uložení
- **F9**: Manuální nahrání
- **ESC**: Menu / Zavřít menu

## Instalace a Nastavení

### Předpoklady
- **Go 1.19+** - Hlavní engine
- **Git** - Správa verzí

### Rychlý Start
```bash
# Klonovat repozitář
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Sestavit hru
go build ./cmd/client

# Spustit hru
./client
```

### Nastavení Vývoje
```bash
# Nainstalovat závislosti
go mod tidy

# Spustit testy
go test ./...

# Sestavit pro vývoj
go build -tags debug ./cmd/client
```

## Systémové Požadavky

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dvoujádrový procesor
- **RAM**: 4GB
- **GPU**: Kompatibilní s OpenGL 3.3+
- **Úložiště**: 500MB volného místa

### Doporučeno
- **CPU**: Čtyřjádrový procesor
- **RAM**: 8GB+
- **GPU**: Vyhrazená grafická karta
- **Úložiště**: 1GB+ volného místa

## Architektura

### Klíčové Technologie
- **Jazyk**: Go (Golang)
- **Grafika**: Ebiten (2D herní knihovna)
- **Sestavovací Systém**: Go moduly

### Struktura Projektu
```
TesselBox/
├── cmd/client/          # Hlavní spustitelný soubor hry
├── pkg/                 # Klíčové balíčky
│   ├── world/          # Generování a správa světa
│   ├── player/         # Mechaniky hráče a fyzika
│   ├── blocks/         # Typy bloků a vlastnosti
│   ├── items/          # Systém předmětů a výroba
│   ├── crafting/       # Recepty výroby a UI
│   ├── weather/        # Simulace počasí
│   ├── gametime/       # Cyklus den/noc
│   ├── save/           # Funkce uložení/nahrání
│   └── render/         # Systémy vykreslování a UI
├── config/             # Konfigurační soubory
└── assets/             # Herní assety (pokud existují)
```

## Přispívání

### Pro Vývojáře
1. Forknout repozitář
2. Vytvořit větev funkce (`git checkout -b feature/úžasná-funkce`)
3. Commitnout změny (`git commit -m 'Přidat úžasnou funkci'`)
4. Pushnout do větve (`git push origin feature/úžasná-funkce`)
5. Otevřít Pull Request

### Pokyny pro Vývoj
- Dodržovat standardy kódování Go
- Přidat testy pro nové funkce
- Aktualizovat dokumentaci
- Zajistit kompatibilitu napříč platformami

## Licence

**MIT Licence** - Podrobnosti viz soubor [LICENSE](LICENSE).

## Poděkování

- **Inspirováno**: Mechanikami hry Terraria
- **Postaveno s**: Herním enginem Ebiten
- **Přispěvatelé**: Komunita open source

## Podpora

- **Problémy**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Diskuze**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki Projektu](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Užijte si průzkum hexagonálního světa TesselBox!*
