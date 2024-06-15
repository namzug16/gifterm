package main

import (
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestIndexCharacter(t *testing.T) {
  r := uint8(200)
  g := uint8(120)
  b := uint8(200)

  c := characterFromRgb(r, g, b)

  t.Log("CHARACTER : ", c)
}

func ProcessingSingleImage(t *testing.T) {
	start := time.Now()

  t.Log("STARTED")

	img := readImage("input/frame_0001.png")

  stop := time.Since(start)
  start = time.Now()

  t.Log("Image read; ", stop)

	w := 144
	h := 144

	resizedImg := resizeImage(img, w, h)

  stop = time.Since(start)
  start = time.Now()

  t.Log("Image resized; ", stop)

	imageToAscii1(resizedImg, w, h)

  stop = time.Since(start)
  start = time.Now()

  t.Log("Image ascii 1; ", stop)

	imageToAscii2(resizedImg)

  stop = time.Since(start)
  start = time.Now()

  t.Log("Image ascii 2; ", stop)

	imageToAscii3(resizedImg)

  stop = time.Since(start)

  t.Log("Image ascii 3; ", stop)
}

func ProcessingSingleImage2(t *testing.T) {
	start := time.Now()

  t.Log("STARTED")

	img := readImage("input/frame_0001.png")

  stop := time.Since(start)
  start = time.Now()

  t.Log("Image read; ", stop)

	w := 144
	h := 144

	resizedImg := resizeImage(img, w, h)

  stop = time.Since(start)
  start = time.Now()

  t.Log("Image resized; ", stop)

	imageToAscii3(resizedImg)

  stop = time.Since(start)

  t.Log("Image ascii 3; ", stop)
}

func ProcessingImageTime(t *testing.T) {
	m := newModel()

	start := time.Now()

	files, _ := m.readFiles()

	var fileNames []string

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			fileNames = append(fileNames, m.Dir+"/"+file.Name())
		}
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i] < fileNames[j]
	})

	stop := time.Since(start)
	start = time.Now()

	t.Log("Got files after: ", stop)
	t.Log("Files count: ", len(fileNames))

	m.Files = fileNames

	m.Width = 144
	m.Height = 144

	c1 := m.loadImages(m.Files)
	results := make(chan job)

	numWorkers := 10
	var wg sync.WaitGroup

	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go m.worker(i, &wg, c1, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for j := range results {
		m.Frames[j.InputPath] = j.Ascii
	}

	t.Log("Proccessed them: ", time.Since(start))
}
