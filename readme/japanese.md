# TesselBox - 日本語のREADME
## 六角形ボクセルゲーム

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

*Terraria* にインスパイアされた 2D サンドボックスアドベンチャーゲームですが、**六角形グリッド** 上に構築されています。

世界を探検し、リソースを採掘し、構造物を建設し、アイテムを作成し、敵と戦い、生存しましょう — すべて美しい六角形タイルの中で。

## ゲーム機能

### ✅ **完全な機能**
- **六角形ワールド生成** - バイオーム付きの手続き型生成ワールド
- **採掘とクラフト** - 異なる素材速度を持つツールベースの採掘
- **ブロック配置** - ゴーストプレビュー付きの右クリック配置
- **インベントリシステム** - 32スロットインベントリとホットバー（9スロット）
- **戦闘システム** - 攻撃アニメーション付きのヘルス/ダメージシステム
- **昼夜サイクル** - 動的照明と時間進行
- **天候効果** - 雨、雪、嵐システム
- **保存/読み込みシステム** - 自動保存付きの永続ワールド状態

### 🎮 **コントロール**
- **WASD / 矢印キー**：移動
- **スペース**：ジャンプ / 攻撃
- **左クリック**：ブロック採掘
- **右クリック**：ブロック配置
- **E**：クラフトメニューを開く
- **Q**：選択アイテムをドロップ
- **マウスホイール**：ホットバー選択
- **1-9**：直接ホットバー選択
- **F5**：手動保存
- **F9**：手動読み込み
- **ESC**：メニュー / メニューを閉じる

## インストールとセットアップ

### 前提条件
- **Go 1.19+** - コアエンジン
- **Git** - バージョン管理

### クイックスタート
```bash
# リポジトリをクローン
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# ゲームをビルド
go build ./cmd/client

# ゲームを実行
./client
```

### 開発セットアップ
```bash
# 依存関係をインストール
go mod tidy

# テストを実行
go test ./...

# 開発用にビルド
go build -tags debug ./cmd/client
```

## システム要件

### 最低要件
- **OS**：Windows 10+、macOS 10.15+、Linux
- **CPU**：デュアルコアプロセッサ
- **RAM**：4GB
- **GPU**：OpenGL 3.3+ 対応
- **ストレージ**：500MB の空き容量

### 推奨要件
- **CPU**：クアッドコアプロセッサ
- **RAM**：8GB+
- **GPU**：専用グラフィックカード
- **ストレージ**：1GB+ の空き容量

## アーキテクチャ

### コアテクノロジー
- **言語**：Go (Golang)
- **グラフィックス**：Ebiten (2Dゲームライブラリ)
- **ビルドシステム**：Goモジュール

### プロジェクト構造
```
TesselBox/
├── cmd/client/          # メインゲーム実行ファイル
├── pkg/                 # コアパッケージ
│   ├── world/          # ワールド生成と管理
│   ├── player/         # プレイヤーメカニクスと物理
│   ├── blocks/         # ブロックタイプとプロパティ
│   ├── items/          # アイテムシステムとクラフト
│   ├── crafting/       # クラフトレシピとUI
│   ├── weather/        # 天候シミュレーション
│   ├── gametime/       # 昼夜サイクル
│   ├── save/           # 保存/読み込み機能
│   └── render/         # レンダリングとUIシステム
├── config/             # 設定ファイル
└── assets/             # ゲームアセット（存在する場合）
```

## 貢献

### 開発者向け
1. リポジトリをフォーク
2. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
3. 変更をコミット（`git commit -m 'Add amazing feature'`）
4. ブランチにプッシュ（`git push origin feature/amazing-feature`）
5. Pull Request を開く

### 開発ガイドライン
- Goコーディング標準に従う
- 新機能のテストを追加
- ドキュメントを更新
- クロスプラットフォーム互換性を確保

## ライセンス

**MITライセンス** - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## クレジット

- **インスピレーション**：Terrariaゲームメカニクス
- **構築ツール**：Ebitenゲームエンジン
- **貢献者**：オープンソースコミュニティ

## サポート

- **Issues**：[GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussions**：[GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**：[プロジェクトWiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*TesselBoxの六角形世界の探索をお楽しみください！*
