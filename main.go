package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fps := flag.Int("fps", 12, "Specify the name to greet")
	far := flag.Float64("far", 2.1, "Font aspect ratio")
	characterDensity := flag.String("cd", ".,:-=i|%O#@$X", "Set character density string")
	randomBlank := flag.Bool("randomBlank", false, "Set if a random character from the character density string should be pick for blank pixels")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Gif file has not been specified.")
		fmt.Println("Usage: gifterm <input.gif>")
		return
	}

	path := args[0]

	// TODO: verify file is a gif

	windowSizeChan := make(chan tea.WindowSizeMsg)
	ctx, cancel := context.WithCancel(context.Background())

	m := newModel(
		windowSizeChan,
		*fps,
		*far,
		AsciiConfig{
			CharacterDensity: *characterDensity,
			SetRandomBlank:   *randomBlank,
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

				for i := 0; i < numWorkers; i++ {
					go worker(
						ctx,
						&wg,
						c1,
						results,
						size.Width,
						size.Height,
						m.FAR,
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

					p.Send(playMsg{
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
