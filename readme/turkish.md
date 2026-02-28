# TesselBox - Türkçe README
## Altıgen Voxel Oyunu

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

*Terraria*'dan ilham alan 2D kum havuzu macera oyunu, ancak **altıgen ızgara** üzerinde inşa edilmiştir.

Dünyaları keşfedin, kaynakları çıkarın, yapılar inşa edin, eşyalar yaratın, düşmanlarla savaşın ve hayatta kalın — hepsi güzel altıgen karolarda.

## Oyun Özellikleri

### ✅ **Tam Özellikler**
- **Altıgen Dünya Oluşturma** - Biyomlarla prosedürel oluşturulan dünyalar
- **Madencilik ve Üretim** - Farklı malzeme hızlarına sahip araç tabanlı madencilik
- **Blok Yerleştirme** - Hayalet önizleme ile blok yerleştirmek için sağ tıklama
- **Envanter Sistemi** - Hızlı çubuklu 32 slot envanter (9 slot)
- **Savaş Sistemi** - Saldırı animasyonlarına sahip sağlık/hasar sistemi
- **Gündüz/Gece Döngüsü** - Dinamik aydınlatma ve zaman ilerlemesi
- **Hava Etkileri** - Yağmur, kar ve fırtına sistemleri
- **Kaydet/Yükle Sistemi** - Otomatik kaydetme ile kalıcı dünya durumu

### 🎮 **Kontroller**
- **WASD / Ok Tuşları**: Hareket
- **Boşluk**: Zıplama / Saldırma
- **Sol Tık**: Blok madenciliği
- **Sağ Tık**: Blok yerleştirme
- **E**: Üretim menüsünü aç
- **Q**: Seçili eşyayı bırak
- **Fare Tekerleği**: Hızlı çubuk seçimi
- **1-9**: Doğrudan hızlı çubuk seçimi
- **F5**: Manuel kaydetme
- **F9**: Manuel yükleme
- **ESC**: Menü / Menüleri kapat

## Kurulum ve Ayarlar

### Ön Gereksinimler
- **Go 1.19+** - Ana motor
- **Git** - Sürüm kontrolü

### Hızlı Başlangıç
```bash
# Depoyu klonla
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Oyunu derle
go build ./cmd/client

# Oyunu çalıştır
./client
```

### Geliştirme Kurulumu
```bash
# Bağımlılıkları yükle
go mod tidy

# Testleri çalıştır
go test ./...

# Geliştirme için derle
go build -tags debug ./cmd/client
```

## Sistem Gereksinimleri

### Minimum
- **İşletim Sistemi**: Windows 10+, macOS 10.15+, Linux
- **İşlemci**: Çift çekirdekli işlemci
- **RAM**: 4GB
- **Ekran Kartı**: OpenGL 3.3+ uyumlu
- **Depolama**: 500MB boş alan

### Önerilen
- **İşlemci**: Dört çekirdekli işlemci
- **RAM**: 8GB+
- **Ekran Kartı**: Özel grafik kartı
- **Depolama**: 1GB+ boş alan

## Mimarisi

### Ana Teknolojiler
- **Dil**: Go (Golang)
- **Grafikler**: Ebiten (2D oyun kütüphanesi)
- **Derleme Sistemi**: Go modülleri

### Proje Yapısı
```
TesselBox/
├── cmd/client/          # Ana oyun çalıştırılabilir dosyası
├── pkg/                 # Ana paketler
│   ├── world/          # Dünya oluşturma ve yönetimi
│   ├── player/         # Oyuncu mekanikleri ve fizik
│   ├── blocks/         # Blok türleri ve özellikleri
│   ├── items/          # Eşya sistemi ve üretim
│   ├── crafting/       # Üretim tarifleri ve arayüz
│   ├── weather/        # Hava durumu simülasyonu
│   ├── gametime/       # Gündüz/gece döngüsü
│   ├── save/           # Kaydet/yükle işlevselliği
│   └── render/         # Oluşturma ve kullanıcı arayüzü sistemleri
├── config/             # Yapılandırma dosyaları
└── assets/             # Oyun varlıkları (varsa)
```

## Katkıda Bulunma

### Geliştiriciler İçin
1. Depoyu çatalla
2. Özellik dalı oluştur (`git checkout -b feature/harika-ozellik`)
3. Değişikliklerini işle (`git commit -m 'Harika özellik ekle'`)
4. Dale gönder (`git push origin feature/harika-ozellik`)
5. Çekme İsteği Aç

### Geliştirme Yönergeleri
- Go kodlama standartlarını takip et
- Yeni özellikler için testler ekle
- Dokümantasyonu güncelle
- Çoklu platform uyumluluğu sağla

## Lisans

**CC BY-NC-SA 4.0 Lisansı** - Ayrıntılar için [LICENSE](LICENSE) dosyasına bakın.

## Teşekkürler

- **İlhami**: Terraria oyun mekanikleri
- **İle İnşa Edildi**: Ebiten oyun motoru
- **Katkıda Bulunanlar**: Açık kaynak topluluğu

## Destek

- **Sorunlar**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Tartışmalar**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Proje Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*TesselBox'un altıgen dünyasının keşfini keyfini çıkarın!*
