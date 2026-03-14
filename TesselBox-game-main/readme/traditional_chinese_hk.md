# TesselBox - 繁體中文（香港）README
## 六角網格像素遊戲

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

一款受 *Terraria* 啟發嘅 2D 沙盒冒險遊戲，但建立喺**六角網格**上。

探索世界、採礦資源、建造建築、製作物品、戰鬥敵人同生存 — 一切都喺美麗嘅六角磚塊中。

## 遊戲特色

### ✅ **完整特色**
- **六角世界生成** - 帶有生態嘅程序生成世界
- **採礦同製作** - 基於工具嘅採礦，具有唔同材料速度
- **方塊放置** - 右鍵放置方塊同有預覽效果
- **物品欄系統** - 32 格物品欄同快捷欄（9 格）
- **戰鬥系統** - 生命值/傷害系統同攻擊動畫
- **晝夜循環** - 動態照明同時間進程
- **天氣效果** - 雨、雪同風暴系統
- **保存/載入系統** - 持久世界狀態同自動保存
- **創造模式** - 無限資源同方塊庫
- **飛行模式** - 按 F 鍵切換飛行
- **指令系統** - 聊天指令（/help, /give, /creative, /survival, /tp）
- **選單系統** - 主選單同方塊庫界面

### 🎮 **控制方法**
- **WASD / 方向鍵**：移動
- **空白鍵**：跳躍 / 攻擊
- **左鍵**：採礦方塊
- **右鍵**：放置方塊
- **E**：開啟製作選單
- **Q**：丟棄揀選物品
- **滑鼠滾輪**：快捷欄選擇
- **1-9**：直接快捷欄選擇
- **B**：開啟方塊庫（創造模式）
- **F**：切換飛行模式
- **/**：開啟指令模式
- **F5**：手動保存
- **F9**：手動載入
- **ESC**：選單 / 關閉選單

## 安裝同設定

### 必要條件
- **Go 1.19+** - 核心引擎
- **Git** - 版本控制

### 快速開始
```bash
# 複製儲存庫
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# 建置遊戲
go build ./cmd/client

# 執行遊戲
./client
```

### 開發設定
```bash
# 安裝依賴
go mod tidy

# 執行測試
go test ./...

# 開發建置
go build -tags debug ./cmd/client
```

## 系統需求

### 最低需求
- **作業系統**：Windows 10+, macOS 10.15+, Linux
- **CPU**：雙核心處理器
- **RAM**：4GB
- **GPU**：支援 OpenGL 3.3+
- **儲存空間**：500MB 可用空間

### 建議需求
- **CPU**：四核心處理器
- **RAM**：8GB+
- **GPU**：獨立顯示卡
- **儲存空間**：1GB+ 可用空間

## 架構

### 核心技術
- **語言**：Go (Golang)
- **圖形**：Ebiten (2D 遊戲程式庫)
- **建置系統**：Go 模組

### 專案結構
```
TesselBox/
├── cmd/client/          # 主遊戲執行檔
├── pkg/                 # 核心套件
│   ├── world/          # 世界生成同管理
│   ├── player/         # 玩家機制同物理
│   ├── blocks/         # 方塊類型同屬性
│   ├── items/          # 物品系統同製作
│   ├── crafting/       # 製作配方同介面
│   ├── weather/        # 天氣模擬
│   ├── gametime/       # 晝夜循環
│   ├── save/           # 保存/載入功能
│   └── render/         # 渲染同介面系統
├── config/             # 設定檔案
└── assets/             # 遊戲資源（如有）
```

## 指令系統

### 可用指令
- **/help** - 顯示所有可用指令
- **/give [物品名稱] [數量]** - 畀物品到物品欄
- **/creative** - 切換到創造模式
- **/survival** - 切換到生存模式
- **/tp [x] [y]** - 傳送到指定座標

### 創造模式特色
- 無限資源
- 方塊庫（按 B 開啟）
- 即刻破壞方塊
- 飛行移動

## 貢獻

### 開發者
1. Fork 此儲存庫
2. 建立功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交變更 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

### 開發指南
- 遵循 Go 編碼標準
- 為新功能新增測試
- 更新文件
- 確保跨平台相容性

## 授權

**CC BY-NC-SA 4.0 授權** - 詳見 [LICENSE](LICENSE) 檔案。

## 鳴謝

- **靈感來源**：Terraria 遊戲機制
- **建置工具**：Ebiten 遊戲引擎
- **貢獻者**：開源社群

## 支援

- **問題回報**：[GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **討論區**：[GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**：[專案 Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*享受探索 TesselBox 嘅六角世界！*
