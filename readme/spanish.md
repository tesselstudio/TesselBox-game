# TesselBox - README en Español
## Juego de Voxels Hexagonales

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Un juego de aventura sandbox 2D inspirado en *Terraria*, pero construido en una **cuadrícula hexagonal**.

Explora mundos, mina recursos, construye estructuras, fabrica objetos, lucha contra enemigos y sobrevive — todo en hermosas baldosas hexagonales.

## Características del Juego

### ✅ **Características Completas**
- **Generación de Mundos Hexagonales** - Mundos generados proceduralmente con biomas
- **Minería y Fabricación** - Minería basada en herramientas con diferentes velocidades de material
- **Colocación de Bloques** - Clic derecho para colocar bloques con vista previa fantasma
- **Sistema de Inventario** - Inventario de 32 ranuras con barra rápida (9 ranuras)
- **Sistema de Combate** - Sistema de salud/daño con animaciones de ataque
- **Ciclo Día/Noche** - Iluminación dinámica y progresión del tiempo
- **Efectos Climáticos** - Sistemas de lluvia, nieve y tormenta
- **Sistema Guardar/Cargar** - Estado persistente del mundo con autoguardado

### 🎮 **Controles**
- **WASD / Flechas**: Movimiento
- **Espacio**: Saltar / Atacar
- **Clic Izquierdo**: Minar bloques
- **Clic Derecho**: Colocar bloques
- **E**: Abrir menú de fabricación
- **Q**: Soltar objeto seleccionado
- **Rueda del Ratón**: Selección de barra rápida
- **1-9**: Selección directa de barra rápida
- **F5**: Guardado manual
- **F9**: Carga manual
- **ESC**: Menú / Cerrar menús

## Instalación y Configuración

### Prerrequisitos
- **Go 1.19+** - Motor principal
- **Git** - Control de versiones

### Inicio Rápido
```bash
# Clonar el repositorio
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Construir el juego
go build ./cmd/client

# Ejecutar el juego
./client
```

### Configuración de Desarrollo
```bash
# Instalar dependencias
go mod tidy

# Ejecutar pruebas
go test ./...

# Construir para desarrollo
go build -tags debug ./cmd/client
```

## Requisitos del Sistema

### Mínimos
- **SO**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Procesador de doble núcleo
- **RAM**: 4GB
- **GPU**: Compatible con OpenGL 3.3+
- **Almacenamiento**: 500MB de espacio libre

### Recomendados
- **CPU**: Procesador de cuatro núcleos
- **RAM**: 8GB+
- **GPU**: Tarjeta gráfica dedicada
- **Almacenamiento**: 1GB+ de espacio libre

## Arquitectura

### Tecnologías Principales
- **Lenguaje**: Go (Golang)
- **Gráficos**: Ebiten (biblioteca de juegos 2D)
- **Sistema de Construcción**: Módulos Go

### Estructura del Proyecto
```
TesselBox/
├── cmd/client/          # Ejecutable principal del juego
├── pkg/                 # Paquetes principales
│   ├── world/          # Generación y gestión del mundo
│   ├── player/         # Mecánicas del jugador y física
│   ├── blocks/         # Tipos de bloques y propiedades
│   ├── items/          # Sistema de objetos y fabricación
│   ├── crafting/       # Recetas de fabricación e interfaz
│   ├── weather/        # Simulación del clima
│   ├── gametime/       # Ciclo día/noche
│   ├── save/           # Funcionalidad guardar/cargar
│   └── render/         # Sistemas de renderizado e interfaz
├── config/             # Archivos de configuración
└── assets/             # Recursos del juego (si los hay)
```

## Contribuir

### Para Desarrolladores
1. Haz fork del repositorio
2. Crea una rama de característica (`git checkout -b feature/caracteristica-increible`)
3. Confirma tus cambios (`git commit -m 'Agregar característica increíble'`)
4. Sube a la rama (`git push origin feature/caracteristica-increible`)
5. Abre un Pull Request

### Directrices de Desarrollo
- Seguir estándares de codificación Go
- Agregar pruebas para nuevas características
- Actualizar documentación
- Asegurar compatibilidad multiplataforma

## Licencia

**Licencia CC BY-NC-SA 4.0** - Ver archivo [LICENSE](LICENSE) para detalles.

## Créditos

- **Inspirado por**: Mecánicas del juego Terraria
- **Construido con**: Motor de juegos Ebiten
- **Contribuyentes**: Comunidad de código abierto

## Soporte

- **Problemas**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discusiones**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki del Proyecto](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*¡Disfruta explorando el mundo hexagonal de TesselBox!*
