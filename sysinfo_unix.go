//go:build unix

package main

import "syscall"

func diskBytes(sys any) (int64, bool) {
	st, ok := sys.(*syscall.Stat_t)
	if !ok {
		return 0, false
	}
	return st.Blocks * 512, true
}

func fileID(sys any) (dev, ino uint64, ok bool) {
	st, ok := sys.(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}
	return uint64(st.Dev), uint64(st.Ino), true
}
