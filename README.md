# diskus

A fast, colorful folder size analyzer for your terminal — written in Go.

Bir klasörün ne kadar yer kapladığını renkli, hızlı ve interaktif olarak gösteren terminal aracı.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)](#installation)

> **English** · [Türkçe](#türkçe)

```
📁 ~/Downloads

  ├─ videos/       ███████████████░░░░░     4.2 GB   68.1%
  ├─ archives/     █████░░░░░░░░░░░░░░░░     1.4 GB   22.7%
  ├─ images/       █░░░░░░░░░░░░░░░░░░░░   380.2 MB    6.0%
  └─ documents/    ░░░░░░░░░░░░░░░░░░░░░   190.5 MB    3.1%

  Total: 6.1 GB  (12480 files)
```

## Features

- 🎨 **Colorful tree view** with size bars and percentages
- 🕹️ **Interactive browser** (`--tui`) — walk into folders with arrow keys
- ⚡ **Fast** — top-level directories are scanned concurrently
- 📊 **File-type breakdown** (`--types`) — see which extensions eat your disk
- 🔎 **Filters** — `--top`, `--depth`, `--min`, `--exclude`, `--all`
- 💽 **Real disk usage** (`--disk`) — allocated blocks, like `du`
- 🔗 **Loop-safe symlinks** (`--follow`)
- 📄 **JSON output** (`--json`) for scripts
- 🌍 **10 languages** — auto-detects your system language
- 📦 **Single binary** — no runtime, works on macOS, Linux and Windows

## Installation

### With Go

```bash
go install github.com/Emiran404/diskus@latest
```

### Download a binary

Grab a prebuilt binary for your platform from the [Releases](https://github.com/Emiran404/diskus/releases) page, then put it somewhere on your `PATH`.

### Build from source

```bash
git clone https://github.com/Emiran404/diskus.git
cd diskus
make install      # or: go install .
```

> If `diskus: command not found`, make sure `~/go/bin` is on your `PATH`:
> `export PATH="$HOME/go/bin:$PATH"`

## Usage

```bash
diskus                              # analyze current folder
diskus ~/Downloads                  # analyze a specific path
diskus --top 10                     # only the 10 largest items
diskus --depth 2                    # show 2 nested levels (tree)
diskus --depth 0                    # unlimited depth
diskus --tui ~/Downloads            # interactive browser
diskus --types --top 10             # size breakdown by file type
diskus --sort count --min 10mb      # sort by file count, hide < 10 MB
diskus --disk                       # real allocated disk space (like du)
diskus --json . > report.json       # machine-readable output
```

### Options

| Flag | Description |
|------|-------------|
| `--top N` | Show only the largest N items per level (0 = all) |
| `--depth N` | Nested levels to show (0 = unlimited) |
| `--all` | Include hidden files (starting with `.`) |
| `--exclude a,b` | Names to skip (default `node_modules,.git,vendor`) |
| `--sort` | `size` (default), `name`, `count` |
| `--reverse` | Reverse the sort order |
| `--unit` | `auto`, `b`, `kb`, `mb`, `gb`, `tb` |
| `--min` | Hide items below this size (e.g. `10mb`) |
| `--tui` | Launch interactive browser |
| `--types` | Size breakdown by file type |
| `--follow` | Follow symlinks (loop-safe) |
| `--disk` | Use allocated disk space instead of logical size |
| `--json` | Print result as JSON |
| `--no-color` | Disable colored output (also honors `NO_COLOR`) |
| `--verbose` | List inaccessible paths |
| `--lang` | Interface language for this command |
| `--set-lang` | Set language permanently |
| `--version` | Show version |

### Interactive mode

```
↑/↓  move    →  enter folder    ←  back    q  quit
```

## Languages

Turkish, English, German, French, Spanish, Italian, Portuguese, Russian, Chinese, Japanese.

The language is chosen in this order: `--lang` flag → `DISKUS_LANG` env var → saved setting → system language.

```bash
diskus --lang de .        # German, this command only
diskus --set-lang en      # set English permanently
diskus --set-lang auto    # go back to system language
```

Codes: `tr en de fr es it pt ru zh ja` (and `auto`).

## License

[MIT](LICENSE) © Emirhan Gök

---

## Türkçe

Terminalde klasör boyutlarını renkli, hızlı ve interaktif gösteren bir araç. Tek binary, kurulum kolay, macOS · Linux · Windows üzerinde çalışır.

### Özellikler

- 🎨 **Renkli ağaç görünümü** — çubuklar ve yüzdelerle
- 🕹️ **İnteraktif gezgin** (`--tui`) — ok tuşlarıyla klasörlere gir/çık
- ⚡ **Hızlı** — üst seviye klasörler eşzamanlı taranır
- 📊 **Dosya türü dökümü** (`--types`) — hangi uzantı ne kadar yer kaplıyor
- 🔎 **Filtreler** — `--top`, `--depth`, `--min`, `--exclude`, `--all`
- 💽 **Gerçek disk kullanımı** (`--disk`) — `du` gibi ayrılan blok
- 🔗 **Döngü korumalı symlink** (`--follow`)
- 📄 **JSON çıktı** (`--json`)
- 🌍 **10 dil** — sistem dilini otomatik algılar

### Kurulum

```bash
# Go ile
go install github.com/Emiran404/diskus@latest

# Kaynaktan
git clone https://github.com/Emiran404/diskus.git
cd diskus
make install
```

> `diskus: command not found` alırsan `~/go/bin`'i PATH'e ekle:
> `export PATH="$HOME/go/bin:$PATH"`

Ayrıca [Releases](https://github.com/Emiran404/diskus/releases) sayfasından platformuna uygun hazır binary'yi indirebilirsin.

### Kullanım

```bash
diskus                              # bulunduğun klasör
diskus ~/Downloads                  # belirli bir yol
diskus --top 10                     # en büyük 10 öğe
diskus --depth 2                    # 2 seviye iç içe (ağaç)
diskus --tui ~/Downloads            # interaktif gezgin
diskus --types --top 10             # dosya türü dökümü
diskus --sort count --min 10mb      # sayıya göre sırala, 10 MB altını gizle
diskus --disk                       # gerçek disk kullanımı
diskus --json . > rapor.json        # makine-okur çıktı
```

### Seçenekler

| Bayrak | Açıklama |
|--------|----------|
| `--top N` | Her seviyede en büyük N öğe (0 = hepsi) |
| `--depth N` | İç içe seviye sayısı (0 = sınırsız) |
| `--all` | Gizli dosyaları dahil et |
| `--exclude a,b` | Atlanacak adlar (varsayılan `node_modules,.git,vendor`) |
| `--sort` | `size` (varsayılan), `name`, `count` |
| `--reverse` | Sıralamayı ters çevir |
| `--unit` | `auto`, `b`, `kb`, `mb`, `gb`, `tb` |
| `--min` | Bu boyutun altını gizle (ör. `10mb`) |
| `--tui` | İnteraktif gezgini başlat |
| `--types` | Dosya türü dökümü |
| `--follow` | Symlink'leri takip et (döngü korumalı) |
| `--disk` | Mantıksal boyut yerine ayrılan disk alanı |
| `--json` | JSON olarak yazdır |
| `--no-color` | Renkleri kapat (`NO_COLOR` da geçerli) |
| `--verbose` | Erişilemeyen yolları listele |
| `--lang` | Bu komut için arayüz dili |
| `--set-lang` | Dili kalıcı ayarla |
| `--version` | Sürümü göster |

### Diller

Türkçe, İngilizce, Almanca, Fransızca, İspanyolca, İtalyanca, Portekizce, Rusça, Çince, Japonca.

Dil önceliği: `--lang` bayrağı → `DISKUS_LANG` ortam değişkeni → kayıtlı ayar → sistem dili.

```bash
diskus --lang de .        # sadece bu komut Almanca
diskus --set-lang en      # kalıcı İngilizce
diskus --set-lang auto    # sistem diline dön
```

### Lisans

[MIT](LICENSE) © Emirhan Gök
