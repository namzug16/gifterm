package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sys/unix"
)

func main() {
	fps := flag.Int("fps", 12, "FPS")
	characterDensity := flag.String("cd", ".,:-=i|%O#@$X", "Set character density string")
	randomBlank := flag.Bool("randomBlank", false, "Set if a random character from the character density string should be pick for blank pixels")
	ofg := flag.Bool("ofg", false, "Only Foreground - Set if only the foregroud color should be set")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Gif file has not been specified.")
		fmt.Println("Usage: gifterm <input.gif>")
		os.Exit(1)
	}

	//TODO: verify file exists

	path := args[0]

	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)

	if err != nil {
		fmt.Fprintf(os.Stderr, "TIOCGWINSZ failed: %v\n", err)
		os.Exit(1)
	}

	cw := float64(ws.Xpixel) / float64(ws.Col)
	ch := float64(ws.Ypixel) / float64(ws.Row)

	car := ch / cw;

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
