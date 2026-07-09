package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Unit int

const (
	UnitAuto Unit = iota
	UnitB
	UnitKB
	UnitMB
	UnitGB
	UnitTB
)

var unitNames = map[Unit]string{UnitB: "B", UnitKB: "KB", UnitMB: "MB", UnitGB: "GB", UnitTB: "TB"}
var unitDivs = map[Unit]float64{UnitB: 1, UnitKB: 1024, UnitMB: 1 << 20, UnitGB: 1 << 30, UnitTB: 1 << 40}

func parseUnit(s string) (Unit, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "auto":
		return UnitAuto, nil
	case "b":
		return UnitB, nil
	case "k", "kb":
		return UnitKB, nil
	case "m", "mb":
		return UnitMB, nil
	case "g", "gb":
		return UnitGB, nil
	case "t", "tb":
		return UnitTB, nil
	}
	return UnitAuto, fmt.Errorf("%s", T("err.unit", s))
}

func humanSize(bytes int64) string {
	return humanSizeU(bytes, UnitAuto)
}

func humanSizeU(bytes int64, u Unit) string {
	if u != UnitAuto {
		if u == UnitB {
			return fmt.Sprintf("%d B", bytes)
		}
		return fmt.Sprintf("%.1f %s", float64(bytes)/unitDivs[u], unitNames[u])
	}

	const k = 1024
	if bytes < k {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(k), 0
	for n := bytes / k; n >= k; n /= k {
		div *= k
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

func parseSize(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 0, nil
	}

	mult := float64(1)
	switch {
	case strings.HasSuffix(s, "tb"), strings.HasSuffix(s, "t"):
		mult, s = 1<<40, strings.TrimRight(s, "tb")
	case strings.HasSuffix(s, "gb"), strings.HasSuffix(s, "g"):
		mult, s = 1<<30, strings.TrimRight(s, "gb")
	case strings.HasSuffix(s, "mb"), strings.HasSuffix(s, "m"):
		mult, s = 1<<20, strings.TrimRight(s, "mb")
	case strings.HasSuffix(s, "kb"), strings.HasSuffix(s, "k"):
		mult, s = 1024, strings.TrimRight(s, "kb")
	case strings.HasSuffix(s, "b"):
		s = strings.TrimRight(s, "b")
	}

	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("%s", T("err.size", s))
	}
	return int64(val * mult), nil
}
