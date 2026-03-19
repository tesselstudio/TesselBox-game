# TesselBox - Deutsche README
## Hexagonales Voxel-Spiel

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Ein 2D-Sandbox-Abenteuerspiel inspiriert von *Terraria*, aber auf einem **hexagonalen Gitter** aufgebaut.

Erkunde Welten, baue Ressourcen ab, errichte Strukturen, stelle Gegenstände her, kämpfe gegen Feinde und überlebe — alles in wunderschönen hexagonalen Kacheln.

## Spiel-Features

### ✅ **Vollständige Features**
- **Hexagonale Weltgenerierung** - Prozedural generierte Welten mit Biomen
- **Abbau und Herstellung** - Werkzeugbasierter Abbau mit unterschiedlichen Materialgeschwindigkeiten
- **Blockplatzierung** - Rechtsklick zum Platzieren von Blöcken mit Geistervorschau
- **Inventarsystem** - 32-Slot-Inventar mit Schnellleiste (9 Slots)
- **Kampfsystem** - Gesundheit/Schaden-System mit Angriffsanimationen
- **Tag/Nacht-Zyklus** - Dynamische Beleuchtung und Zeitfortschritt
- **Wettereffekte** - Regen-, Schnee- und Sturmsysteme
- **Speichern/Laden-System** - Persistenter Weltzustand mit automatischer Speicherung

### 🎮 **Steuerung**
- **WASD / Pfeiltasten**: Bewegung
- **Leertaste**: Springen / Angreifen
- **Linksklick**: Blöcke abbauen
- **Rechtsklick**: Blöcke platzieren
- **E**: Herstellungsmenü öffnen
- **Q**: Ausgewählten Gegenstand fallen lassen
- **Mausrad**: Schnellleistenauswahl
- **1-9**: Direkte Schnellleistenauswahl
- **F5**: Manuelle Speicherung
- **F9**: Manuelles Laden
- **ESC**: Menü / Menüs schließen

## Installation und Einrichtung

### Voraussetzungen
- **Go 1.19+** - Hauptantrieb
- **Git** - Versionskontrolle

### Schnellstart
```bash
# Repository klonen
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Spiel erstellen
go build ./cmd/client

# Spiel ausführen
./client
```

### Entwicklungseinrichtung
```bash
# Abhängigkeiten installieren
go mod tidy

# Tests ausführen
go test ./...

# Für Entwicklung erstellen
go build -tags debug ./cmd/client
```

## Systemanforderungen

### Minimum
- **Betriebssystem**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dual-Core-Prozessor
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ kompatibel
- **Speicher**: 500MB freier Speicherplatz

### Empfohlen
- **CPU**: Quad-Core-Prozessor
- **RAM**: 8GB+
- **GPU**: Dedizierte Grafikkarte
- **Speicher**: 1GB+ freier Speicherplatz

## Architektur

### Kerntechnologien
- **Sprache**: Go (Golang)
- **Grafik**: Ebiten (2D-Spielebibliothek)
- **Build-System**: Go-Module

### Projektstruktur
```
TesselBox/
├── cmd/client/          # Hauptspiel-Executable
├── pkg/                 # Kernpakete
│   ├── world/          # Weltgenerierung & -verwaltung
│   ├── player/         # Spieler-Mechaniken & Physik
│   ├── blocks/         # Blocktypen & Eigenschaften
│   ├── items/          # Gegenstandssystem & Herstellung
│   ├── crafting/       # Herstellungsrezepte & Benutzeroberfläche
│   ├── weather/        # Wetter-Simulation
│   ├── gametime/       # Tag/Nacht-Zyklus
│   ├── save/           # Speichern/Laden-Funktionalität
│   └── render/         # Rendering & Benutzeroberflächen-Systeme
├── config/             # Konfigurationsdateien
└── assets/             # Spiel-Assets (falls vorhanden)
```

## Mitwirken

### Für Entwickler
1. Repository forken
2. Feature-Branch erstellen (`git checkout -b feature/tolle-funktion`)
3. Änderungen committen (`git commit -m 'Tolle Funktion hinzufügen'`)
4. Auf Branch pushen (`git push origin feature/tolle-funktion`)
5. Pull Request öffnen

### Entwicklungsrichtlinien
- Go-Coding-Standards befolgen
- Tests für neue Features hinzufügen
- Dokumentation aktualisieren
- Plattformübergreifende Kompatibilität sicherstellen

## Lizenz

**CC BY-NC-SA 4.0 Lizenz** - Siehe [LICENSE](LICENSE)-Datei für Details.

## Danksagungen

- **Inspiriert von**: Terraria-Spielmechaniken
- **Erstellt mit**: Ebiten-Spiele-Engine
- **Mitwirkende**: Open-Source-Community

## Support

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Diskussionen**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Projekt-Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Genieße die Erkundung der hexagonalen Welt von TesselBox!*
