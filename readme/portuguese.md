# TesselBox - README em Português
## Jogo de Voxels Hexagonais

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

Um jogo de aventura sandbox 2D inspirado no *Terraria*, mas construído em uma **grade hexagonal**.

Explore mundos, mine recursos, construa estruturas, crie itens, lute contra inimigos e sobreviva — tudo em belas telhas hexagonais.

## Recursos do Jogo

### ✅ **Recursos Completos**
- **Geração de Mundo Hexagonal** - Mundos gerados proceduralmente com biomas
- **Mineração e Criação** - Mineração baseada em ferramentas com velocidades diferentes de materiais
- **Posicionamento de Blocos** - Clique direito para posicionar blocos com visualização fantasma
- **Sistema de Inventário** - Inventário de 32 slots com barra rápida (9 slots)
- **Sistema de Combate** - Sistema de saúde/dano com animações de ataque
- **Ciclo Dia/Noite** - Iluminação dinâmica e progressão temporal
- **Efeitos Climáticos** - Sistemas de chuva, neve e tempestade
- **Sistema Salvar/Carregar** - Estado persistente do mundo com salvamento automático

### 🎮 **Controles**
- **WASD / Setas**: Movimento
- **Espaço**: Pular / Atacar
- **Clique Esquerdo**: Mineração de blocos
- **Clique Direito**: Posicionamento de blocos
- **E**: Abrir menu de criação
- **Q**: Largar item selecionado
- **Roda do Mouse**: Seleção da barra rápida
- **1-9**: Seleção direta da barra rápida
- **F5**: Salvamento manual
- **F9**: Carregamento manual
- **ESC**: Menu / Fechar menus

## Instalação e Configuração

### Pré-requisitos
- **Go 1.19+** - Motor principal
- **Git** - Controle de versão

### Início Rápido
```bash
# Clonar repositório
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Construir jogo
go build ./cmd/client

# Executar jogo
./client
```

### Configuração de Desenvolvimento
```bash
# Instalar dependências
go mod tidy

# Executar testes
go test ./...

# Construir para desenvolvimento
go build -tags debug ./cmd/client
```

## Requisitos do Sistema

### Mínimo
- **SO**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Processador dual-core
- **RAM**: 4GB
- **GPU**: Compatível com OpenGL 3.3+
- **Armazenamento**: 500MB de espaço livre

### Recomendado
- **CPU**: Processador quad-core
- **RAM**: 8GB+
- **GPU**: Placa de vídeo dedicada
- **Armazenamento**: 1GB+ de espaço livre

## Arquitetura

### Tecnologias Principais
- **Linguagem**: Go (Golang)
- **Gráficos**: Ebiten (biblioteca de jogos 2D)
- **Sistema de Construção**: Módulos Go

### Estrutura do Projeto
```
TesselBox/
├── cmd/client/          # Executável principal do jogo
├── pkg/                 # Pacotes principais
│   ├── world/          # Geração e gerenciamento do mundo
│   ├── player/         # Mecânicas do jogador e física
│   ├── blocks/         # Tipos de blocos e propriedades
│   ├── items/          # Sistema de itens e criação
│   ├── crafting/       # Receitas de criação e interface
│   ├── weather/        # Simulação do clima
│   ├── gametime/       # Ciclo dia/noite
│   ├── save/           # Funcionalidade salvar/carregar
│   └── render/         # Sistemas de renderização e interface
├── config/             # Arquivos de configuração
└── assets/             # Assets do jogo (se houver)
```

## Contribuição

### Para Desenvolvedores
1. Faça fork do repositório
2. Crie um branch de recurso (`git checkout -b feature/recurso-incrivel`)
3. Faça commit das suas mudanças (`git commit -m 'Adicionar recurso incrível'`)
4. Faça push para o branch (`git push origin feature/recurso-incrivel`)
5. Abra um Pull Request

### Diretrizes de Desenvolvimento
- Seguir padrões de codificação Go
- Adicionar testes para novos recursos
- Atualizar documentação
- Garantir compatibilidade cross-platform

## Licença

**Licença CC BY-NC-SA 4.0** - Veja o arquivo [LICENSE](LICENSE) para detalhes.

## Créditos

- **Inspirado por**: Mecânicas do jogo Terraria
- **Construído com**: Engine de jogos Ebiten
- **Contribuintes**: Comunidade open source

## Suporte

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussões**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Wiki do Projeto](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Aproveite a exploração do mundo hexagonal do TesselBox!*
