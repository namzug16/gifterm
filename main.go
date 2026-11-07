package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	if *characterAspectRatioOverride < 0 {
		fmt.Println("-car must be greater than 0")
		os.Exit(1)
	}

	if !gifExists(path) {
		fmt.Printf("GIF %q does not exist or it is not a GIF\n", path)
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

	windowSizeChan := make(chan tea.WindowSizeMsg)

	ctx, cancel := context.WithCancel(context.Background())

	m := newModel(
		windowSizeChan,
		*fps,
		car,
		AsciiConfig{
			CharacterDensity: *characterDensity,
			SetRandomBlank:   *randomBlank,
			OnlyForeground:   *ofg,
		},
	)

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func(m *model) {
		images, err := readGif(path)
		// FIX: handle error
		if err != nil {
			return
		}

		for size := range windowSizeChan {
			m.reset()
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			go func(size tea.WindowSizeMsg) {
				c1 := chanFromImages(images)
				results := make(chan job)

				numWorkers := 10
				var wg sync.WaitGroup

				wg.Add(numWorkers)

				for range numWorkers {
					go worker(
						ctx,
						&wg,
						c1,
						results,
						size.Width,
						size.Height,
						m.CAR,
						m.AsciiConfiguration,
					)
				}

				go func() {
					wg.Wait()
					close(results)
				}()

				go func() {
					frames := make(map[int]string)

					for j := range results {
						select {
						case <-ctx.Done():
							return
						default:
							frames[j.Index] = j.Ascii
							pe := int((float32(len(frames)) / float32(len(images))) * 100)
							p.Send(loadingMsg{
								p: pe,
							})
						}
					}

					m.Frames = frames

					p.Send(processingCompletedMsg{
						m: *m,
					})
				}()
			}(size)
		}
	}(&m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func gifExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return false
	}
	if info.IsDir() {
		return false
	}

	if strings.ToLower(filepath.Ext(path)) != ".gif" {
		return false
	}

	//NOTE: header check and signature
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 6)
	if _, err := io.ReadFull(f, buf); err != nil {
		return false
	}
	sig := string(buf)
	return sig == "GIF87a" || sig == "GIF89a"
}
