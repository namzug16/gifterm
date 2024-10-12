package main

import (
	"context"
	"image"
	"image/color"
	"sync"

	"golang.org/x/image/draw"
)

type job struct {
	Image image.Image
	Ascii string
	Index int
}

func chanFromImages(imgs []*image.Paletted) <-chan job {
	out := make(chan job)
	go func() {
		correctedImgs := make([]image.Image, len(imgs))

		bounds := imgs[0].Bounds()
		savedPreviousImg := image.NewRGBA(bounds)
		draw.Draw(savedPreviousImg, bounds, image.NewUniform(color.Transparent), image.Point{}, draw.Src)

		for i := 0; i < len(imgs); i++ {
			cumulativeImage := image.NewRGBA(bounds)
			draw.Draw(cumulativeImage, bounds, savedPreviousImg, image.Point{}, draw.Src)
			frame := imgs[i]
			draw.Draw(cumulativeImage, bounds, frame, image.Point{}, draw.Over)
			correctedImgs[i] = cumulativeImage
			savedPreviousImg = cumulativeImage
		}

		for i, img := range correctedImgs {
			job := job{
				Index: i,
				Image: img,
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan job,
	result chan<- job,
	w,
	h int,
	far float64,
  aConfig AsciiConfig,
) {
	defer wg.Done()
	c2 := resizeImages(ctx, jobs, w, h, far)
	c3 := imagesToAscii(ctx, c2, aConfig)
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-c3:
			if !ok {
				return
			}
			result <- j
		}
	}
}

func resizeImages(
	ctx context.Context,
	input <-chan job,
	w,
	h int,
	far float64,
) <-chan job {
	out := make(chan job)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-input:
				if !ok {
					return
				}

				job.Image = resizeImage(job.Image, w, h, far)
				out <- job
			}
		}
	}()
	return out
}

func imagesToAscii(
	ctx context.Context,
	input <-chan job,
  aConfig AsciiConfig,
) <-chan job {
	out := make(chan job)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-input:
				if !ok {
					return
				}

				job.Ascii = imageToAscii(job.Image, aConfig)
				out <- job
			}
		}
	}()
	return out
}
