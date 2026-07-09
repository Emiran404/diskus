package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type Options struct {
	ShowHidden bool
	Exclude    []string
	Follow     bool
	Disk       bool
	Progress   *int64
}

type Node struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Size     int64   `json:"size"`
	Count    int64   `json:"count"`
	IsDir    bool    `json:"is_dir"`
	Children []*Node `json:"children,omitempty"`
}

type ScanResult struct {
	Root   *Node
	Errors []error
}

type errCollector struct {
	mu   sync.Mutex
	errs []error
}

func (c *errCollector) add(err error) {
	if err == nil {
		return
	}
	c.mu.Lock()
	c.errs = append(c.errs, err)
	c.mu.Unlock()
}

type fileIDKey struct{ dev, ino uint64 }

type scanCtx struct {
	opts    Options
	ec      *errCollector
	visited sync.Map
}

func Scan(root string, opts Options) (*ScanResult, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	c := &scanCtx{opts: opts, ec: &errCollector{}}

	if !info.IsDir() {
		countOne(opts.Progress)
		return &ScanResult{Root: &Node{
			Name: info.Name(), Path: root, Size: sizeOf(info, opts.Disk), Count: 1,
		}}, nil
	}

	rootNode := &Node{Name: filepath.Base(root), Path: root, IsDir: true}
	if key, ok := fileKey(info); ok {
		c.visited.Store(key, struct{}{})
	}

	children, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, child := range children {
		if skip(child.Name(), opts) {
			continue
		}
		childPath := filepath.Join(root, child.Name())

		wg.Add(1)
		go func(name, path string) {
			defer wg.Done()
			node := c.walk(name, path)
			mu.Lock()
			rootNode.Children = append(rootNode.Children, node)
			rootNode.Size += node.Size
			rootNode.Count += node.Count
			mu.Unlock()
		}(child.Name(), childPath)
	}
	wg.Wait()

	return &ScanResult{Root: rootNode, Errors: c.ec.errs}, nil
}

func (c *scanCtx) walk(name, path string) *Node {
	countOne(c.opts.Progress)
	node := &Node{Name: name, Path: path}

	fi, err := os.Lstat(path)
	if err != nil {
		c.ec.add(err)
		return node
	}
	mode := fi.Mode()

	if mode&os.ModeSymlink != 0 {
		if !c.opts.Follow {
			return node
		}
		ti, err := os.Stat(path)
		if err != nil {
			c.ec.add(err)
			return node
		}
		if ti.IsDir() {
			return c.walkDir(path, ti, node)
		}
		if ti.Mode().IsRegular() {
			node.Size, node.Count = sizeOf(ti, c.opts.Disk), 1
		}
		return node
	}

	if mode.IsDir() {
		return c.walkDir(path, fi, node)
	}

	if mode.IsRegular() {
		node.Size, node.Count = sizeOf(fi, c.opts.Disk), 1
	}
	return node
}

func (c *scanCtx) walkDir(path string, fi os.FileInfo, node *Node) *Node {
	node.IsDir = true

	if key, ok := fileKey(fi); ok {
		if _, seen := c.visited.LoadOrStore(key, struct{}{}); seen {
			return node
		}
	}

	children, err := os.ReadDir(path)
	if err != nil {
		c.ec.add(err)
		return node
	}

	for _, child := range children {
		if skip(child.Name(), c.opts) {
			continue
		}
		cn := c.walk(child.Name(), filepath.Join(path, child.Name()))
		node.Children = append(node.Children, cn)
		node.Size += cn.Size
		node.Count += cn.Count
	}
	return node
}

func sizeOf(fi os.FileInfo, disk bool) int64 {
	if disk {
		if b, ok := diskBytes(fi.Sys()); ok {
			return b
		}
	}
	return fi.Size()
}

func fileKey(fi os.FileInfo) (fileIDKey, bool) {
	dev, ino, ok := fileID(fi.Sys())
	return fileIDKey{dev, ino}, ok
}

func skip(name string, opts Options) bool {
	if !opts.ShowHidden && strings.HasPrefix(name, ".") {
		return true
	}
	for _, ex := range opts.Exclude {
		if name == ex {
			return true
		}
	}
	return false
}

type SortMode int

const (
	SortSize SortMode = iota
	SortName
	SortCount
)

func parseSort(s string) (SortMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "size":
		return SortSize, nil
	case "name":
		return SortName, nil
	case "count":
		return SortCount, nil
	}
	return SortSize, fmt.Errorf("%s", T("err.sort", s))
}

func SortTree(n *Node, mode SortMode, reverse bool) {
	less := func(a, b *Node) bool {
		switch mode {
		case SortName:
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case SortCount:
			return a.Count > b.Count
		default:
			return a.Size > b.Size
		}
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		if reverse {
			return less(n.Children[j], n.Children[i])
		}
		return less(n.Children[i], n.Children[j])
	})
	for _, c := range n.Children {
		if c.IsDir {
			SortTree(c, mode, reverse)
		}
	}
}

func countOne(counter *int64) {
	if counter != nil {
		atomic.AddInt64(counter, 1)
	}
}
