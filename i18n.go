package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Lang int

const (
	LangTR Lang = iota
	LangEN
	LangDE
	LangFR
	LangES
	LangIT
	LangPT
	LangRU
	LangZH
	LangJA
)

const langList = "auto, tr, en, de, fr, es, it, pt, ru, zh, ja"

var currentLang = LangTR

func T(key string, args ...any) string {
	m, ok := catalog[key]
	if !ok {
		return key
	}
	s, ok := m[currentLang]
	if !ok || s == "" {
		s = m[LangEN]
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

func parseLang(s string) (Lang, bool) {
	if v := strings.ToLower(strings.TrimSpace(s)); v == "" || v == "auto" {
		return resolveAutoLang(), true
	}
	return parseLangExplicit(s)
}

func parseLangExplicit(s string) (Lang, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "tr", "turkish", "türkçe", "turkce":
		return LangTR, true
	case "en", "english", "ingilizce":
		return LangEN, true
	case "de", "german", "deutsch", "almanca":
		return LangDE, true
	case "fr", "french", "français", "francais", "fransızca":
		return LangFR, true
	case "es", "spanish", "español", "espanol", "ispanyolca":
		return LangES, true
	case "it", "italian", "italiano", "italyanca":
		return LangIT, true
	case "pt", "portuguese", "português", "portugues", "portekizce":
		return LangPT, true
	case "ru", "russian", "русский", "rusça", "rusca":
		return LangRU, true
	case "zh", "chinese", "中文", "çince", "cince", "cn":
		return LangZH, true
	case "ja", "japanese", "日本語", "japonca", "jp":
		return LangJA, true
	}
	return LangTR, false
}

func resolveAutoLang() Lang {
	if v := os.Getenv("DISKUS_LANG"); v != "" {
		if l, ok := parseLangExplicit(v); ok {
			return l
		}
	}
	if v := loadConfig().Lang; v != "" {
		if l, ok := parseLangExplicit(v); ok {
			return l
		}
	}
	return detectSystemLang()
}

func detectSystemLang() Lang {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if l, ok := langFromLocale(os.Getenv(key)); ok {
			return l
		}
	}
	if runtime.GOOS == "darwin" {
		if l, ok := langFromLocale(macOSLanguage()); ok {
			return l
		}
	}
	return LangEN
}

func langFromLocale(v string) (Lang, bool) {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" || v == "c" || v == "posix" {
		return LangEN, false
	}
	prefixes := []struct {
		p string
		l Lang
	}{
		{"tr", LangTR}, {"en", LangEN}, {"de", LangDE}, {"fr", LangFR},
		{"es", LangES}, {"it", LangIT}, {"pt", LangPT}, {"ru", LangRU},
		{"zh", LangZH}, {"ja", LangJA},
	}
	for _, x := range prefixes {
		if strings.HasPrefix(v, x.p) {
			return x.l, true
		}
	}
	return LangEN, false
}

func macOSLanguage() string {
	out, err := exec.Command("defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.Trim(strings.TrimSpace(line), "\",")
		if line == "" || line == "(" || line == ")" {
			continue
		}
		return strings.ToLower(line)
	}
	return ""
}

func detectLangFromArgs(args []string) Lang {
	for i, a := range args {
		switch {
		case a == "-lang" || a == "--lang":
			if i+1 < len(args) {
				if l, ok := parseLang(args[i+1]); ok {
					return l
				}
			}
		case strings.HasPrefix(a, "-lang="), strings.HasPrefix(a, "--lang="):
			if l, ok := parseLang(a[strings.IndexByte(a, '=')+1:]); ok {
				return l
			}
		}
	}
	return resolveAutoLang()
}
