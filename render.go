package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	dirStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Bold(true)
	fileStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	sizeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
	pctStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	barStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	totalStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50FA7B"))
	treeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
)

const barWidth = 20

type RenderOptions struct {
	Top      int
	MaxDepth int
	Min      int64
	Unit     Unit
}

func Render(root string, res *ScanResult, ro RenderOptions) {
	fmt.Println()
	fmt.Println(titleStyle.Render("📁 " + root))
	fmt.Println()

	if len(res.Root.Children) == 0 {
		fmt.Println(pctStyle.Render(T("render.empty")))
		fmt.Println()
		return
	}

	renderLevel(res.Root, ro, 1, "")

	fmt.Println()
	fmt.Printf("  %s %s  %s\n",
		totalStyle.Render(T("render.total")),
		sizeStyle.Render(humanSizeU(res.Root.Size, ro.Unit)),
		pctStyle.Render(T("render.files", res.Root.Count)),
	)

	if len(res.Errors) > 0 {
		fmt.Println(errStyle.Render(T("render.inacc_warn", len(res.Errors))))
	}
	fmt.Println()
}

func renderLevel(parent *Node, ro RenderOptions, depth int, prefix string) {
	children := parent.Children
	if ro.Min > 0 {
		filtered := make([]*Node, 0, len(children))
		for _, c := range children {
			if c.Size >= ro.Min {
				filtered = append(filtered, c)
			}
		}
		children = filtered
	}

	total := len(children)
	limit := total
	if ro.Top > 0 && ro.Top < limit {
		limit = ro.Top
	}

	nameW := 0
	for i := 0; i < limit; i++ {
		if l := displayLen(children[i]); l > nameW {
			nameW = l
		}
	}
	if nameW > 40 {
		nameW = 40
	}

	for i := 0; i < limit; i++ {
		e := children[i]
		last := i == limit-1

		pct := 0.0
		if parent.Size > 0 {
			pct = float64(e.Size) / float64(parent.Size) * 100
		}

		connector := "├─ "
		if last {
			connector = "└─ "
		}

		name := e.Name
		if e.IsDir {
			name += "/"
		}
		runes := []rune(name)
		if len(runes) > nameW {
			name = string(runes[:nameW-1]) + "…"
		}
		styled := fileStyle.Render(name)
		if e.IsDir {
			styled = dirStyle.Render(name)
		}
		pad := strings.Repeat(" ", nameW-len([]rune(name))+1)

		fmt.Printf("  %s%s%s%s  %s  %s\n",
			treeStyle.Render(prefix+connector),
			styled,
			pad,
			renderBar(pct),
			sizeStyle.Render(fmt.Sprintf("%10s", humanSizeU(e.Size, ro.Unit))),
			pctStyle.Render(fmt.Sprintf("%5.1f%%", pct)),
		)

		if e.IsDir && depth < ro.MaxDepth && len(e.Children) > 0 {
			childPrefix := prefix + "│  "
			if last {
				childPrefix = prefix + "   "
			}
			renderLevel(e, ro, depth+1, childPrefix)
		}
	}

	if limit < total {
		more := treeStyle.Render(prefix + "   ")
		fmt.Printf("  %s%s\n", more, pctStyle.Render(T("render.more_items", total-limit)))
	}
}

func displayLen(n *Node) int {
	l := len([]rune(n.Name))
	if n.IsDir {
		l++
	}
	return l
}

func renderBar(pct float64) string {
	filled := int(pct / 100 * barWidth)
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return barStyle.Render(bar)
}
