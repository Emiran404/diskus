package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type typeStat struct {
	Ext   string
	Size  int64
	Count int64
}

func aggregateTypes(n *Node, acc map[string]*typeStat) {
	if !n.IsDir {
		ext := strings.ToLower(filepath.Ext(n.Name))
		if ext == "" {
			ext = T("types.noext")
		}
		s := acc[ext]
		if s == nil {
			s = &typeStat{Ext: ext}
			acc[ext] = s
		}
		s.Size += n.Size
		s.Count++
		return
	}
	for _, c := range n.Children {
		aggregateTypes(c, acc)
	}
}

func RenderTypes(root string, res *ScanResult, top int, unit Unit) {
	acc := make(map[string]*typeStat)
	aggregateTypes(res.Root, acc)

	stats := make([]*typeStat, 0, len(acc))
	for _, s := range acc {
		stats = append(stats, s)
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Size > stats[j].Size })

	fmt.Println()
	fmt.Println(titleStyle.Render("📊 " + root + T("types.title")))
	fmt.Println()

	if len(stats) == 0 {
		fmt.Println(pctStyle.Render(T("types.nofiles")))
		fmt.Println()
		return
	}

	total := res.Root.Size
	limit := len(stats)
	if top > 0 && top < limit {
		limit = top
	}

	nameW := 0
	for i := 0; i < limit; i++ {
		if l := len([]rune(stats[i].Ext)); l > nameW {
			nameW = l
		}
	}

	for i := 0; i < limit; i++ {
		s := stats[i]
		pct := 0.0
		if total > 0 {
			pct = float64(s.Size) / float64(total) * 100
		}
		pad := strings.Repeat(" ", nameW-len([]rune(s.Ext))+1)
		fmt.Printf("  %s%s%s  %s  %s  %s\n",
			dirStyle.Render(s.Ext),
			pad,
			renderBar(pct),
			sizeStyle.Render(fmt.Sprintf("%10s", humanSizeU(s.Size, unit))),
			pctStyle.Render(fmt.Sprintf("%5.1f%%", pct)),
			pctStyle.Render(T("types.files", s.Count)),
		)
	}

	if limit < len(stats) {
		fmt.Println(pctStyle.Render(T("types.more_types", len(stats)-limit)))
	}

	fmt.Println()
	fmt.Printf("  %s %s  %s\n",
		totalStyle.Render(T("render.total")),
		sizeStyle.Render(humanSizeU(total, unit)),
		pctStyle.Render(T("types.summary", res.Root.Count, len(stats))),
	)
	fmt.Println()
}
