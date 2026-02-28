# TesselBox - README en Français
## Jeu de Voxels Hexagonaux

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Un jeu d'aventure sandbox 2D inspiré par *Terraria*, mais construit sur une **grille hexagonale**.

Explorez des mondes, minez des ressources, construisez des structures, fabriquez des objets, combattez des ennemis et survivez — tout dans de magnifiques tuiles hexagonales.

## Fonctionnalités du Jeu

### ✅ **Fonctionnalités Complètes**
- **Génération de Monde Hexagonal** - Mondes générés procéduralement avec biomes
- **Minage et Artisanat** - Minage basé sur les outils avec différentes vitesses de matériau
- **Placement de Blocs** - Clic droit pour placer des blocs avec aperçu fantôme
- **Système d'Inventaire** - Inventaire 32 emplacements avec barre rapide (9 emplacements)
- **Système de Combat** - Système santé/dégâts avec animations d'attaque
- **Cycle Jour/Nuit** - Éclairage dynamique et progression temporelle
- **Effets Météorologiques** - Systèmes pluie, neige et tempête
- **Système Sauvegarde/Chargement** - État persistant du monde avec sauvegarde automatique

### 🎮 **Contrôles**
- **WASD / Flèches** : Déplacement
- **Espace** : Sauter / Attaquer
- **Clic Gauche** : Miner des blocs
- **Clic Droit** : Placer des blocs
- **E** : Ouvrir le menu d'artisanat
- **Q** : Lâcher l'objet sélectionné
- **Molette Souris** : Sélection barre rapide
- **1-9** : Sélection directe barre rapide
- **F5** : Sauvegarde manuelle
- **F9** : Chargement manuel
- **Échap** : Menu / Fermer les menus

## Installation et Configuration

### Prérequis
- **Go 1.19+** - Moteur principal
- **Git** - Contrôle de version

### Démarrage Rapide
```bash
# Cloner le dépôt
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Construire le jeu
go build ./cmd/client

# Lancer le jeu
./client
```

### Configuration Développement
```bash
# Installer les dépendances
go mod tidy

# Lancer les tests
go test ./...

# Construire pour le développement
go build -tags debug ./cmd/client
```

## Configuration Système Requise

### Minimum
- **OS** : Windows 10+, macOS 10.15+, Linux
- **CPU** : Processeur double-cœur
- **RAM** : 4GB
- **GPU** : Compatible OpenGL 3.3+
- **Stockage** : 500MB d'espace libre

### Recommandé
- **CPU** : Processeur quadricœur
- **RAM** : 8GB+
- **GPU** : Carte graphique dédiée
- **Stockage** : 1GB+ d'espace libre

## Architecture

### Technologies Principales
- **Langage** : Go (Golang)
- **Graphismes** : Ebiten (bibliothèque de jeu 2D)
- **Système de Construction** : Modules Go

### Structure du Projet
```
TesselBox/
├── cmd/client/          # Exécutable principal du jeu
├── pkg/                 # Paquets principaux
│   ├── world/          # Génération et gestion du monde
│   ├── player/         # Mécaniques joueur et physique
│   ├── blocks/         # Types de blocs et propriétés
│   ├── items/          # Système d'objets et artisanat
│   ├── crafting/       # Recettes artisanat et interface
│   ├── weather/        # Simulation météo
│   ├── gametime/       # Cycle jour/nuit
│   ├── save/           # Fonctionnalité sauvegarde/chargement
│   └── render/         # Systèmes rendu et interface
├── config/             # Fichiers de configuration
└── assets/             # Ressources du jeu (si présentes)
```

## Contribuer

### Pour les Développeurs
1. Fork le dépôt
2. Créez une branche de fonctionnalité (`git checkout -b feature/fonctionnalite-geniale`)
3. Commitez vos changements (`git commit -m 'Ajouter fonctionnalité géniale'`)
4. Poussez vers la branche (`git push origin feature/fonctionnalite-geniale`)
5. Ouvrez une Pull Request

### Directives de Développement
- Suivre les standards de codage Go
- Ajouter des tests pour les nouvelles fonctionnalités
- Mettre à jour la documentation
- Assurer la compatibilité multiplateforme

## Licence

**Licence CC BY-NC-SA 4.0** - Voir le fichier [LICENSE](LICENSE) pour les détails.

## Crédits

- **Inspiré par** : Mécaniques du jeu Terraria
- **Construit avec** : Moteur de jeu Ebiten
- **Contributeurs** : Communauté open source

## Support

- **Issues** : [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussions** : [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki** : [Wiki du Projet](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Profitez de l'exploration du monde hexagonal de TesselBox !*
