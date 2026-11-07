package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sys/unix"
)

func main() {
	fps := flag.Int("fps", 12, "FPS")
	characterDensity := flag.String("cd", ".,:-=i|%O#@$X", "Set character density string")
	randomBlank := flag.Bool("randomBlank", false, "Set if a random character from the character density string should be pick for blank pixels")
	ofg := flag.Bool("ofg", false, "Only Foreground - Set if only the foregroud color should be set")
	characterAspectRatioOverride := flag.Float64("car", 0, "Character Aspect Ratio (cell height / width) - Overrides the value detected from the terminal")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gifterm [flags] <input.gif>")
		fmt.Fprintln(flag.CommandLine.Output(), "\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(flag.CommandLine.Output(), "Gif file has not been specified.")
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	// flag stops at the first non-flag argument, so parse again to accept flags after the file
	flag.CommandLine.Parse(args[1:])

	if flag.NArg() > 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "Unexpected arguments: %v\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	if *fps <= 0 {
		fmt.Println("-fps must be greater than 0")
		os.Exit(1)
	}

	if *characterDensity == "" {
		fmt.Println("-cd must not be empty")
		os.Exit(1)
	}

	if *characterAspectRatioOverride < 0 {
		fmt.Println("-car must be greater than 0")
		os.Exit(1)
	}

	images, err := readGif(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read GIF %q: %v\n", path, err)
		os.Exit(1)
	}

	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)

	if err != nil {
		fmt.Fprintf(os.Stderr, "TIOCGWINSZ failed: %v\n", err)
		os.Exit(1)
	}

	cw := float64(ws.Xpixel) / float64(ws.Col)
	ch := float64(ws.Ypixel) / float64(ws.Row)

	car := ch / cw

	// some terminals/multiplexers (e.g. herdr) report 0 pixel size
	if ws.Xpixel == 0 || ws.Ypixel == 0 {
		car = 2.0
	}

	if *characterAspectRatioOverride > 0 {
		car = *characterAspectRatioOverride
	}

	asciiConfig := AsciiConfig{
		CharacterDensity: *characterDensity,
		SetRandomBlank:   *randomBlank,
		OnlyForeground:   *ofg,
	}

	renderRequests := make(chan renderRequest, 1)

	p := tea.NewProgram(newModel(renderRequests, *fps, car, asciiConfig), tea.WithAltScreen())

	go renderOnResize(p, renderRequests, images, car, asciiConfig)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
