package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

var version = "dev"

func main() {
	currentLang = detectLangFromArgs(os.Args[1:])

	showVersion := flag.Bool("version", false, T("flag.version"))

	top := flag.Int("top", 0, T("flag.top"))
	depth := flag.Int("depth", 1, T("flag.depth"))
	all := flag.Bool("all", false, T("flag.all"))
	excludeArg := flag.String("exclude", "node_modules,.git,vendor", T("flag.exclude"))
	asJSON := flag.Bool("json", false, T("flag.json"))
	sortArg := flag.String("sort", "size", T("flag.sort"))
	reverse := flag.Bool("reverse", false, T("flag.reverse"))
	unitArg := flag.String("unit", "auto", T("flag.unit"))
	minArg := flag.String("min", "", T("flag.min"))
	tui := flag.Bool("tui", false, T("flag.tui"))
	follow := flag.Bool("follow", false, T("flag.follow"))
	disk := flag.Bool("disk", false, T("flag.disk"))
	noColor := flag.Bool("no-color", false, T("flag.nocolor"))
	verbose := flag.Bool("verbose", false, T("flag.verbose"))
	types := flag.Bool("types", false, T("flag.types"))
	langArg := flag.String("lang", "auto", T("flag.lang")+" ("+langList+")")
	setLangArg := flag.String("set-lang", "", T("flag.setlang")+" ("+langList+")")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, T("usage.tagline"))
		fmt.Fprintln(os.Stderr, T("usage.usage"))
		fmt.Fprintln(os.Stderr, T("usage.usageline"))
		fmt.Fprintln(os.Stderr, T("usage.options"))
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, T("usage.examples"))
		fmt.Fprintln(os.Stderr, "  diskus")
		fmt.Fprintln(os.Stderr, "  diskus --tui ~/Downloads")
		fmt.Fprintln(os.Stderr, "  diskus --top 10 --depth 2 ~/Downloads")
		fmt.Fprintln(os.Stderr, "  diskus --sort count --min 10mb ~/Documents")
		fmt.Fprintln(os.Stderr, "  diskus --types --top 10 ~/Downloads")
		fmt.Fprintln(os.Stderr, T("usage.ex.depth0"))
		fmt.Fprintln(os.Stderr, T("usage.ex.setlang"))
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("diskus", version)
		return
	}

	if *setLangArg != "" {
		handleSetLang(*setLangArg)
		return
	}

	if lang, ok := parseLang(*langArg); ok {
		currentLang = lang
	} else {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+T("err.invalidlang", *langArg)+" ("+langList+")"))
		os.Exit(1)
	}

	if *noColor || os.Getenv("NO_COLOR") != "" {
		lipgloss.SetColorProfile(termenv.Ascii)
	}

	root := "."
	if args := flag.Args(); len(args) > 0 {
		root = args[0]
	}

	sortMode, err := parseSort(*sortArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+err.Error()))
		os.Exit(1)
	}
	unit, err := parseUnit(*unitArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+err.Error()))
		os.Exit(1)
	}
	minBytes, err := parseSize(*minArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+err.Error()))
		os.Exit(1)
	}

	opts := Options{
		ShowHidden: *all,
		Exclude:    parseExclude(*excludeArg),
		Follow:     *follow,
		Disk:       *disk,
	}

	var counter int64
	stopProgress := func() {}
	if !*asJSON && isatty.IsTerminal(os.Stderr.Fd()) {
		opts.Progress = &counter
		stopProgress = startProgress(&counter)
	}

	res, err := Scan(root, opts)
	stopProgress()
	if err != nil {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+err.Error()))
		os.Exit(1)
	}

	SortTree(res.Root, sortMode, *reverse)

	maxDepth := *depth
	if maxDepth <= 0 {
		maxDepth = math.MaxInt
	}

	switch {
	case *tui:
		if err := RunTUI(res.Root, unit); err != nil {
			fmt.Fprintln(os.Stderr, errStyle.Render(T("err.tui")+err.Error()))
			os.Exit(1)
		}
	case *types:
		RenderTypes(root, res, *top, unit)
	case *asJSON:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res.Root); err != nil {
			fmt.Fprintln(os.Stderr, T("err.json"), err)
			os.Exit(1)
		}
	default:
		Render(root, res, RenderOptions{
			Top:      *top,
			MaxDepth: maxDepth,
			Min:      minBytes,
			Unit:     unit,
		})
		if *verbose && len(res.Errors) > 0 {
			printErrors(res.Errors)
		}
	}
}

func handleSetLang(val string) {
	v := strings.ToLower(strings.TrimSpace(val))

	var toSave config
	if v != "auto" {
		if _, ok := parseLangExplicit(v); !ok {
			currentLang = resolveAutoLang()
			fmt.Fprintln(os.Stderr, errStyle.Render(T("err.prefix")+T("err.invalidlang", val)+" ("+langList+")"))
			os.Exit(1)
		}
		toSave.Lang = v
	}

	path, err := saveConfig(toSave)
	if err != nil {
		fmt.Fprintln(os.Stderr, errStyle.Render(T("setlang.savefail")+err.Error()))
		os.Exit(1)
	}

	currentLang = resolveAutoLang()
	if v == "auto" {
		fmt.Println(totalStyle.Render("✓ ") + T("setlang.auto"))
	} else {
		fmt.Println(totalStyle.Render("✓ ") + T("setlang.done", v))
	}
	fmt.Println(pctStyle.Render(T("setlang.configfile") + path))
}

func printErrors(errs []error) {
	const maxShown = 50
	fmt.Fprintln(os.Stderr, errStyle.Render(T("err.inaccessible", len(errs))))
	for i, e := range errs {
		if i >= maxShown {
			fmt.Fprintln(os.Stderr, errStyle.Render(T("err.more", len(errs)-maxShown)))
			break
		}
		fmt.Fprintln(os.Stderr, errStyle.Render("  • "+e.Error()))
	}
}

func parseExclude(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func startProgress(counter *int64) func() {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				n := atomic.LoadInt64(counter)
				fmt.Fprintf(os.Stderr, T("scan.progress"), frames[i%len(frames)], n)
				i++
			}
		}
	}()
	return func() { close(done); time.Sleep(10 * time.Millisecond) }
}
