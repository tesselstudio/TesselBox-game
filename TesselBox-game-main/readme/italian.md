# TesselBox - Italiano README
## Gioco di Voxel Esagonali

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Un gioco di avventura sandbox 2D ispirato a *Terraria*, ma costruito su una **griglia esagonale**.

Esplora mondi, estrai risorse, costruisci strutture, crea oggetti, combatti nemici e sopravvivi — tutto in bellissime tessere esagonali.

## Caratteristiche del Gioco

### ✅ **Caratteristiche Complete**
- **Generazione Mondo Esagonale** - Mondi generati proceduralmente con biomi
- **Estrazione e Creazione** - Estrazione basata su strumenti con diverse velocità materiali
- **Posizionamento Blocchi** - Clic destro per posizionare blocchi con anteprima fantasma
- **Sistema Inventario** - Inventario 32 slot con barra veloce (9 slot)
- **Sistema Combattimento** - Sistema salute/danno con animazioni attacco
- **Ciclo Giorno/Notte** - Illuminazione dinamica e progressione temporale
- **Effetti Meteo** - Sistemi pioggia, neve e tempesta
- **Sistema Salva/Carica** - Stato mondo persistente con salvataggio automatico

### 🎮 **Controlli**
- **WASD / Frecce**: Movimento
- **Spazio**: Salto / Attacco
- **Clic Sinistro**: Estrazione blocchi
- **Clic Destro**: Posizionamento blocchi
- **E**: Apri menu creazione
- **Q**: Rilascia oggetto selezionato
- **Rotella Mouse**: Selezione barra veloce
- **1-9**: Selezione diretta barra veloce
- **F5**: Salvataggio manuale
- **F9**: Caricamento manuale
- **ESC**: Menu / Chiudi menu

## Installazione e Configurazione

### Prerequisiti
- **Go 1.19+** - Motore principale
- **Git** - Controllo versione

### Avvio Rapido
```bash
# Clona repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Costruisci gioco
go build ./cmd/client

# Avvia gioco
./client
```

### Configurazione Sviluppo
```bash
# Installa dipendenze
go mod tidy

# Esegui test
go test ./...

# Costruisci per sviluppo
go build -tags debug ./cmd/client
```

## Requisiti di Sistema

### Minimo
- **SO**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processore dual-core
- **RAM**: 4GB
- **GPU**: Compatibile OpenGL 3.3+
- **Archiviazione**: 500MB spazio libero

### Consigliato
- **CPU**: Processore quad-core
- **RAM**: 8GB+
- **GPU**: Scheda video dedicata
- **Archiviazione**: 1GB+ spazio libero

## Architettura

### Tecnologie Principali
- **Linguaggio**: Go (Golang)
- **Grafica**: Ebiten (libreria giochi 2D)
- **Sistema Build**: Moduli Go

### Struttura Progetto
```
TesselBox/
├── cmd/client/          # Eseguibile principale gioco
├── pkg/                 # Pacchetti principali
│   ├── world/          # Generazione e gestione mondo
│   ├── player/         # Meccaniche giocatore e fisica
│   ├── blocks/         # Tipi blocchi e proprietà
│   ├── items/          # Sistema oggetti e creazione
│   ├── crafting/       # Ricette creazione e interfaccia
│   ├── weather/        # Simulazione meteo
│   ├── gametime/       # Ciclo giorno/notte
│   ├── save/           # Funzionalità salva/carica
│   └── render/         # Sistemi rendering e interfaccia
├── config/             # File configurazione
└── assets/             # Asset gioco (se presenti)
```

## Contributi

### Per Sviluppatori
1. Fork il repository
2. Crea un branch feature (`git checkout -b feature/fantastica-feature`)
3. Committa i tuoi cambiamenti (`git commit -m 'Add fantastica feature'`)
4. Push al branch (`git push origin feature/fantastica-feature`)
5. Apri una Pull Request

### Linee Guida Sviluppo
- Seguire standard codifica Go
- Aggiungere test per nuove feature
- Aggiornare documentazione
- Garantire compatibilità cross-platform

## Licenza

**Licenza CC BY-NC-SA 4.0** - Vedi file [LICENSE](LICENSE) per dettagli.

## Ringraziamenti

- **Ispirato da**: Meccaniche gioco Terraria
- **Costruito con**: Motore giochi Ebiten
- **Collaboratori**: Comunità open source

## Supporto

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussioni**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki Progetto](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Goditi l'esplorazione del mondo esagonale di TesselBox!*
