# TesselBox - Polski README
## Gra Wokselska Sześciokątna

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Gra przygodowa typu sandbox 2D inspirowana *Terraria*, ale zbudowana na **siatce sześciokątnej**.

Eksploruj światy, wydobywaj zasoby, buduj struktury, twórz przedmioty, walcz z wrogami i przeżyj — wszystko w pięknych sześciokątnych kafelkach.

## Funkcje Gry

### ✅ **Pełne Funkcje**
- **Generacja Świata Sześciokątnego** - Proceduralnie generowane światy z biomami
- **Górnictwo i Tworzenie** - Górnictwo oparte na narzędziach z różnymi prędkościami materiałów
- **Umieszczanie Bloków** - Kliknięcie prawym przyciskiem myszy, aby umieścić bloki z podglądem ducha
- **System Ekwipunku** - Ekwipunek 32 slotów z paskiem szybkiego dostępu (9 slotów)
- **System Walki** - System zdrowia/obrażeń z animacjami ataków
- **Cykl Dzień/Noc** - Dynamiczne oświetlenie i postęp czasu
- **Efekty Pogody** - Systemy deszczu, śniegu i burzy
- **System Zapisywania/Wczytywania** - Trwały stan świata z automatycznym zapisem

### 🎮 **Sterowanie**
- **WASD / Strzałki**: Ruch
- **Spacja**: Skok / Atak
- **Lewy Klik**: Górnictwo bloków
- **Prawy Klik**: Umieszczanie bloków
- **E**: Otwórz menu tworzenia
- **Q**: Upuść wybrany przedmiot
- **Kółko Myszy**: Wybór paska szybkiego dostępu
- **1-9**: Bezpośredni wybór paska szybkiego dostępu
- **F5**: Ręczne zapisywanie
- **F9**: Ręczne wczytywanie
- **ESC**: Menu / Zamknij menu

## Instalacja i Konfiguracja

### Wymagania Wstępne
- **Go 1.19+** - Silnik główny
- **Git** - Kontrola wersji

### Szybki Start
```bash
# Sklonuj repozytorium
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Zbuduj grę
go build ./cmd/client

# Uruchom grę
./client
```

### Konfiguracja Rozwoju
```bash
# Zainstaluj zależności
go mod tidy

# Uruchom testy
go test ./...

# Zbuduj dla rozwoju
go build -tags debug ./cmd/client
```

## Wymagania Systemowe

### Minimalne
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Procesor dwurdzeniowy
- **RAM**: 4GB
- **GPU**: Kompatybilny z OpenGL 3.3+
- **Pamięć**: 500MB wolnego miejsca

### Zalecane
- **CPU**: Procesor czterordzeniowy
- **RAM**: 8GB+
- **GPU**: Dedykowana karta graficzna
- **Pamięć**: 1GB+ wolnego miejsca

## Architektura

### Główne Technologie
- **Język**: Go (Golang)
- **Grafika**: Ebiten (biblioteka gier 2D)
- **System Budowy**: Moduły Go

### Struktura Projektu
```
TesselBox/
├── cmd/client/          # Główny plik wykonywalny gry
├── pkg/                 # Główne pakiety
│   ├── world/          # Generacja i zarządzanie światem
│   ├── player/         # Mechaniki gracza i fizyka
│   ├── blocks/         # Typy bloków i właściwości
│   ├── items/          # System przedmiotów i tworzenia
│   ├── crafting/       # Przepisy tworzenia i interfejs
│   ├── weather/        # Symulacja pogody
│   ├── gametime/       # Cykl dzień/noc
│   ├── save/           # Funkcjonalność zapisywania/wczytywania
│   └── render/         # Systemy renderowania i interfejsu
├── config/             # Pliki konfiguracyjne
└── assets/             # Zasoby gry (jeśli istnieją)
```

## Współtworzenie

### Dla Deweloperów
1. Zrób fork repozytorium
2. Utwórz gałąź funkcji (`git checkout -b feature/fantastyczna-funkcja`)
3. Zatwierdź swoje zmiany (`git commit -m 'Dodaj fantastyczną funkcję'`)
4. Wypchnij do gałęzi (`git push origin feature/fantastyczna-funkcja`)
5. Otwórz Pull Request

### Wytyczne Rozwojowe
- Postępuj zgodnie ze standardami kodowania Go
- Dodaj testy dla nowych funkcji
- Aktualizuj dokumentację
- Zapewnij kompatybilność międzyplatformową

## Licencja

**Licencja CC BY-NC-SA 4.0** - Zobacz plik [LICENSE](LICENSE) po szczegóły.

## Podziękowania

- **Inspirowane przez**: Mechaniki gry Terraria
- **Zbudowane z**: Silnik gier Ebiten
- **Współtwórcy**: Społeczność open source

## Wsparcie

- **Problemy**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Dyskusje**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki Projektu](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Ciesz się eksploracją sześciokątnego świata TesselBox!*
