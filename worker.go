package main

import (
	"context"
	"image"
	"runtime"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type renderRequest struct {
	Size       tea.WindowSizeMsg
	Generation int
}

type renderedFrame struct {
	Index int
	Ascii string
}

// renderOnResize renders every frame for each requested size, cancelling the render of the previous size.
func renderOnResize(p *tea.Program, requests <-chan renderRequest, images []image.Image, car float64, config AsciiConfig) {
	cancel := func() {}
	for request := range requests {
		cancel()
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		go renderFrames(ctx, p, request, images, car, config)
	}
	cancel()
}

func renderFrames(ctx context.Context, p *tea.Program, request renderRequest, images []image.Image, car float64, config AsciiConfig) {
	indexes := make(chan int)
	go func() {
		defer close(indexes)
		for i := range images {
			select {
			case indexes <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	results := make(chan renderedFrame)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range indexes {
				resized := resizeImage(images[i], request.Size.Width, request.Size.Height, car)
				frame := renderedFrame{Index: i, Ascii: imageToAscii(resized, config)}
				select {
				case results <- frame:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	frames := make([]string, len(images))
	renderedCount := 0
	for frame := range results {
		frames[frame.Index] = frame.Ascii
		renderedCount++
		p.Send(loadingMsg{Generation: request.Generation, Percentage: renderedCount * 100 / len(images)})
	}

	// a cancelled render sends partial frames here, the model ignores them because its generation moved on
	p.Send(processingCompletedMsg{Generation: request.Generation, Frames: frames})
}
