# Changelog

All notable changes to **diskus** are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.0.5] - 2026-07-09

### Added
- **Interactive delete** in the TUI — press `d` on any item to remove it, with a `y/n` confirmation prompt. Freed space is reflected instantly in the tree without a rescan.
- **Reveal in file manager** — press `o` to open the selected item in Finder (macOS), Explorer (Windows) or the default file manager (Linux).
- **Rescan** — press `r` to re-read the current tree from disk.
- Status line in the TUI showing the result of the last action (deleted / rescanned / errors).
- Unit tests covering the delete and cancel flows.

### Changed
- The TUI footer now lists all available keys.
- Entering a directory no longer requires it to be non-empty.

## [1.0.0] - 2026-07-09

### Added
- Concurrent folder size scanning with a colorful tree view (size bars + percentages).
- Interactive browser (`--tui`) built with Bubble Tea.
- File-type breakdown (`--types`).
- Filters and views: `--top`, `--depth`, `--min`, `--exclude`, `--all`, `--sort`, `--reverse`, `--unit`.
- Real disk usage (`--disk`) and loop-safe symlink following (`--follow`).
- JSON output (`--json`), `--no-color` / `NO_COLOR`, `--verbose`.
- 10 UI languages with system-language auto-detection and a persistent `--set-lang` setting.
- Cross-platform binaries for macOS, Linux and Windows.

[1.0.5]: https://github.com/Emiran404/diskus/releases/tag/v1.0.5
[1.0.0]: https://github.com/Emiran404/diskus/releases/tag/v1.0.0
