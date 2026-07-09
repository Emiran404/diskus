package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

func revealPath(path string, isDir bool) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-R", path).Start()
	case "windows":
		return exec.Command("explorer", "/select,"+filepath.Clean(path)).Start()
	default:
		target := path
		if !isDir {
			target = filepath.Dir(path)
		}
		return exec.Command("xdg-open", target).Start()
	}
}
