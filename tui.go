package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tuiHeader   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	tuiSelected = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#7D56F4")).Foreground(lipgloss.Color("#FFFFFF"))
	tuiFooter   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	tuiCrumb    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	tuiWarn     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#febc2e"))
	tuiOk       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50FA7B"))
)

const (
	modeNormal = iota
	modeConfirm
)

type tuiModel struct {
	root        *Node
	stack       []*Node
	cursorStack []int
	cursor      int
	offset      int
	height      int
	width       int
	unit        Unit
	opts        Options
	sortMode    SortMode
	reverse     bool
	mode        int
	status      string
	statusErr   bool
}

func newTUIModel(root *Node, unit Unit, opts Options, sortMode SortMode, reverse bool) tuiModel {
	return tuiModel{
		root:     root,
		stack:    []*Node{root},
		cursor:   0,
		height:   20,
		width:    80,
		unit:     unit,
		opts:     opts,
		sortMode: sortMode,
		reverse:  reverse,
	}
}

func (m tuiModel) current() *Node { return m.stack[len(m.stack)-1] }

func (m tuiModel) Init() tea.Cmd { return nil }

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 5
		if m.height < 3 {
			m.height = 3
		}
	case tea.KeyMsg:
		key := msg.String()

		if m.mode == modeConfirm {
			switch key {
			case "y", "Y", "enter":
				m = m.doDelete()
			case "n", "N", "esc", "q":
				m.mode = modeNormal
			}
			return m, nil
		}

		m.status = ""
		switch key {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.current().Children)-1 {
				m.cursor++
			}
		case "g", "home":
			m.cursor = 0
		case "G", "end":
			m.cursor = len(m.current().Children) - 1
		case "enter", "right", "l":
			children := m.current().Children
			if len(children) > 0 {
				sel := children[m.cursor]
				if sel.IsDir {
					m.stack = append(m.stack, sel)
					m.cursorStack = append(m.cursorStack, m.cursor)
					m.cursor = 0
					m.offset = 0
				}
			}
		case "left", "h", "backspace":
			if len(m.stack) > 1 {
				m.stack = m.stack[:len(m.stack)-1]
				m.cursor = m.cursorStack[len(m.cursorStack)-1]
				m.cursorStack = m.cursorStack[:len(m.cursorStack)-1]
				m.offset = 0
			}
		case "d":
			if len(m.current().Children) > 0 {
				m.mode = modeConfirm
			}
		case "o":
			children := m.current().Children
			if len(children) > 0 {
				sel := children[m.cursor]
				if err := revealPath(sel.Path, sel.IsDir); err != nil {
					m.status, m.statusErr = T("tui.open_failed", err.Error()), true
				}
			}
		case "r":
			m = m.rescan()
		}
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	return m, nil
}

func (m tuiModel) doDelete() tuiModel {
	m.mode = modeNormal
	children := m.current().Children
	if m.cursor >= len(children) {
		return m
	}
	sel := children[m.cursor]

	if err := os.RemoveAll(sel.Path); err != nil {
		m.status, m.statusErr = T("tui.delete_failed", err.Error()), true
		return m
	}

	for _, n := range m.stack {
		n.Size -= sel.Size
		n.Count -= sel.Count
	}
	cur := m.current()
	cur.Children = append(children[:m.cursor], children[m.cursor+1:]...)
	if m.cursor >= len(cur.Children) && m.cursor > 0 {
		m.cursor--
	}
	m.status, m.statusErr = T("tui.deleted", sel.Name), false
	return m
}

func (m tuiModel) rescan() tuiModel {
	res, err := Scan(m.root.Path, m.opts)
	if err != nil {
		m.status, m.statusErr = err.Error(), true
		return m
	}
	SortTree(res.Root, m.sortMode, m.reverse)
	m.root = res.Root
	m.stack = []*Node{res.Root}
	m.cursorStack = nil
	m.cursor = 0
	m.offset = 0
	m.status, m.statusErr = T("tui.rescanned"), false
	return m
}

func (m tuiModel) View() string {
	var b strings.Builder

	crumbs := make([]string, len(m.stack))
	for i, n := range m.stack {
		name := n.Name
		if name == "" || name == "." {
			name = n.Path
		}
		crumbs[i] = name
	}
	cur := m.current()
	b.WriteString(tuiHeader.Render("📁 "))
	b.WriteString(tuiCrumb.Render(strings.Join(crumbs, " / ")))
	b.WriteString("  ")
	b.WriteString(sizeStyle.Render(humanSizeU(cur.Size, m.unit)))
	b.WriteString(tuiFooter.Render(T("tui.files", cur.Count)))
	b.WriteString("\n\n")

	children := cur.Children
	if len(children) == 0 {
		b.WriteString(tuiFooter.Render(T("tui.empty")))
	}

	end := m.offset + m.height
	if end > len(children) {
		end = len(children)
	}
	nameW := 0
	for i := m.offset; i < end; i++ {
		if l := displayLen(children[i]); l > nameW {
			nameW = l
		}
	}
	if nameW > 40 {
		nameW = 40
	}

	for i := m.offset; i < end; i++ {
		e := children[i]
		pct := 0.0
		if cur.Size > 0 {
			pct = float64(e.Size) / float64(cur.Size) * 100
		}

		name := e.Name
		if e.IsDir {
			name += "/"
		}
		runes := []rune(name)
		if len(runes) > nameW {
			name = string(runes[:nameW-1]) + "…"
		}
		pad := strings.Repeat(" ", nameW-len([]rune(name))+1)

		row := fmt.Sprintf(" %s%s%s  %10s  %5.1f%%",
			name, pad, renderBar(pct), humanSizeU(e.Size, m.unit), pct)

		if i == m.cursor {
			b.WriteString(tuiSelected.Render("›" + row))
		} else {
			line := " "
			if e.IsDir {
				line += dirStyle.Render(row)
			} else {
				line += fileStyle.Render(row)
			}
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	if len(children) > m.height {
		b.WriteString(tuiFooter.Render(fmt.Sprintf("\n  %d-%d / %d", m.offset+1, end, len(children))))
	}

	b.WriteString("\n")
	switch {
	case m.mode == modeConfirm:
		selName := ""
		if m.cursor < len(children) {
			selName = children[m.cursor].Name
		}
		b.WriteString(tuiWarn.Render(T("tui.confirm", selName)))
	case m.status != "":
		if m.statusErr {
			b.WriteString(tuiWarn.Render("⚠ " + m.status))
		} else {
			b.WriteString(tuiOk.Render("✓ " + m.status))
		}
	default:
		b.WriteString(tuiFooter.Render(T("tui.footer")))
	}
	return b.String()
}

func RunTUI(root *Node, unit Unit, opts Options, sortMode SortMode, reverse bool) error {
	p := tea.NewProgram(newTUIModel(root, unit, opts, sortMode, reverse), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
