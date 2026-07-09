package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func send(m tuiModel, s string) tuiModel {
	next, _ := m.Update(key(s))
	return next.(tuiModel)
}

func TestTUIDelete(t *testing.T) {
	dir := t.TempDir()
	big := filepath.Join(dir, "big.bin")
	if err := os.WriteFile(big, make([]byte, 4096), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "small.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Scan(dir, Options{})
	if err != nil {
		t.Fatal(err)
	}
	SortTree(res.Root, SortSize, false)
	totalBefore := res.Root.Size

	m := newTUIModel(res.Root, UnitAuto, Options{}, SortSize, false)

	m = send(m, "d")
	if m.mode != modeConfirm {
		t.Fatalf("pressing d should enter confirm mode, got %d", m.mode)
	}

	m = send(m, "y")

	if m.mode != modeNormal {
		t.Fatalf("after confirm should return to normal mode")
	}
	if _, err := os.Stat(big); !os.IsNotExist(err) {
		t.Fatalf("big.bin should be deleted from disk")
	}
	if len(m.current().Children) != 1 {
		t.Fatalf("expected 1 child left, got %d", len(m.current().Children))
	}
	if m.current().Children[0].Name != "small.txt" {
		t.Fatalf("wrong child remained: %s", m.current().Children[0].Name)
	}
	if m.root.Size != totalBefore-4096 {
		t.Fatalf("root size not updated: got %d want %d", m.root.Size, totalBefore-4096)
	}
}

func TestTUIDeleteCancel(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "keep.txt")
	if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, _ := Scan(dir, Options{})
	m := newTUIModel(res.Root, UnitAuto, Options{}, SortSize, false)

	m = send(m, "d")
	m = send(m, "n")

	if m.mode != modeNormal {
		t.Fatalf("cancel should return to normal mode")
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("file should still exist after cancel: %v", err)
	}
	if len(m.current().Children) != 1 {
		t.Fatalf("child should remain after cancel")
	}
}
