//go:build !unix

package main

func diskBytes(sys any) (int64, bool) { return 0, false }

func fileID(sys any) (dev, ino uint64, ok bool) { return 0, 0, false }
