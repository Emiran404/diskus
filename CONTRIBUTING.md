<div align="center">

# 🤝 Contributing to diskus

**English** · [Türkçe](#-türkçe)

</div>

Thanks for your interest! Issues, ideas and pull requests are all welcome — this is a small, friendly codebase and a great first project to contribute to.

## 🛠️ Development setup

You only need [Go 1.26+](https://go.dev/dl/).

```bash
git clone https://github.com/Emiran404/diskus.git
cd diskus
make build        # builds ./diskus
make install      # installs to ~/go/bin
make vet          # go vet
make release      # cross-compiles all platforms into dist/
```

Run without installing:

```bash
go run . --top 5 ~/Downloads
```

## 🗺️ Project layout

| File | Responsibility |
|:-----|:---------------|
| [`main.go`](main.go) | Flags, validation, dispatch to renderers |
| [`scanner.go`](scanner.go) | Concurrent tree scan, symlink loop protection, sorting |
| [`render.go`](render.go) | Colored tree output |
| [`types.go`](types.go) | File-type (extension) breakdown |
| [`tui.go`](tui.go) | Interactive browser (Bubble Tea) |
| [`format.go`](format.go) | Byte ↔ human-readable conversion, size parsing |
| [`i18n.go`](i18n.go) | Language detection & selection |
| [`catalog.go`](catalog.go) | Translation catalog (all UI strings) |
| [`config.go`](config.go) | Persistent config (`--set-lang`) |
| [`sysinfo_unix.go`](sysinfo_unix.go) / [`sysinfo_other.go`](sysinfo_other.go) | Platform-specific disk/inode info (build tags) |

## 🌍 Adding a language

The most welcome contribution! Two small steps:

1. **[`i18n.go`](i18n.go)** — add a `Lang` constant, its codes in `parseLangExplicit`, a locale prefix in `langFromLocale`, and the code in `langList`.
2. **[`catalog.go`](catalog.go)** — add your language's line to every key (copy the English line and translate).

If a key is missing for your language, diskus silently falls back to English — so partial translations won't break anything.

Test it:

```bash
go run . --lang xx --help
go run . --lang xx --types .
```

## ✅ Pull request checklist

- `go build ./...` and `go vet ./...` pass
- Cross-platform check if you touched scanner/sysinfo: `GOOS=windows go build ./...`
- Keep the code comment-free (project convention) and `gofmt`-ed
- One focused change per PR — small PRs get merged fast
- UI strings go through the catalog (`T("key")`), never hardcoded

## 🐛 Reporting bugs

Open an [issue](https://github.com/Emiran404/diskus/issues) with:

- Your OS + `diskus --version` output
- The exact command you ran
- What you expected vs. what happened (paste output with `--no-color`)

<br/>

---

<br/>

<div align="center">

# 🇹🇷 Türkçe

</div>

İlgin için teşekkürler! Issue, fikir ve PR'ların hepsine açığız — küçük ve samimi bir kod tabanı, ilk açık kaynak katkın için harika bir proje.

## 🛠️ Geliştirme ortamı

Sadece [Go 1.26+](https://go.dev/dl/) gerekli.

```bash
git clone https://github.com/Emiran404/diskus.git
cd diskus
make build        # ./diskus üretir
make install      # ~/go/bin'e kurar
make vet          # go vet
make release      # tüm platformları dist/ altına derler
```

Kurmadan çalıştır:

```bash
go run . --top 5 ~/Downloads
```

## 🗺️ Proje yapısı

| Dosya | Sorumluluk |
|:------|:-----------|
| [`main.go`](main.go) | Bayraklar, doğrulama, çıktı seçimi |
| [`scanner.go`](scanner.go) | Paralel ağaç tarama, symlink döngü koruması, sıralama |
| [`render.go`](render.go) | Renkli ağaç çıktısı |
| [`types.go`](types.go) | Dosya türü (uzantı) dökümü |
| [`tui.go`](tui.go) | İnteraktif gezgin (Bubble Tea) |
| [`format.go`](format.go) | Bayt ↔ insan-okur çeviri, boyut ayrıştırma |
| [`i18n.go`](i18n.go) | Dil algılama ve seçimi |
| [`catalog.go`](catalog.go) | Çeviri kataloğu (tüm arayüz metinleri) |
| [`config.go`](config.go) | Kalıcı ayar (`--set-lang`) |
| [`sysinfo_unix.go`](sysinfo_unix.go) / [`sysinfo_other.go`](sysinfo_other.go) | Platforma özel disk/inode bilgisi (build tag) |

## 🌍 Yeni dil eklemek

En değerli katkı! İki küçük adım:

1. **[`i18n.go`](i18n.go)** — `Lang` sabiti ekle, `parseLangExplicit`'e kodlarını, `langFromLocale`'e locale önekini, `langList`'e kodu ekle.
2. **[`catalog.go`](catalog.go)** — her anahtara kendi dilinin satırını ekle (İngilizce satırı kopyalayıp çevir).

Bir anahtar eksik kalırsa diskus sessizce İngilizce'ye düşer — yarım çeviri hiçbir şeyi bozmaz.

Test et:

```bash
go run . --lang xx --help
go run . --lang xx --types .
```

## ✅ PR kontrol listesi

- `go build ./...` ve `go vet ./...` geçiyor
- scanner/sysinfo'ya dokunduysan çapraz platform kontrolü: `GOOS=windows go build ./...`
- Kod yorumsuz (proje geleneği) ve `gofmt`'li
- PR başına tek odaklı değişiklik — küçük PR hızlı merge olur
- Arayüz metinleri katalogdan geçer (`T("anahtar")`), asla elle yazılmaz

## 🐛 Hata bildirmek

[Issue](https://github.com/Emiran404/diskus/issues) açarken şunları ekle:

- İşletim sistemin + `diskus --version` çıktısı
- Çalıştırdığın komutun tamamı
- Ne bekliyordun, ne oldu (`--no-color` ile çıktıyı yapıştır)
