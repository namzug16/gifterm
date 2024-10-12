package main

// import (
// 	"fmt"
// 	"os"
// 	"sync"
//
// 	tea "github.com/charmbracelet/bubbletea"
// 	"golang.org/x/term"
// )
//
// func startProcessingWithoutEngine(path string) {
// 	fmt.Println("Start program")
// 	windowSizeChan := make(chan tea.WindowSizeMsg)
//
// 	m := newModel(
// 		windowSizeChan,
//     1,
//     1,
// 	)
//
// 	images, err := readGif(path)
// 	if err != nil {
// 		fmt.Println("Error getting terminal size:", err)
// 		return
// 	}
//
// 	width, height, err := term.GetSize(int(os.Stdin.Fd()))
// 	if err != nil {
// 		fmt.Println("Error getting terminal size:", err)
// 		return
// 	}
//
// 	c1 := chanFromImages(images)
// 	results := make(chan job)
//
// 	numWorkers := 10
// 	var wg sync.WaitGroup
//
// 	wg.Add(numWorkers)
//
// 	for i := 0; i < numWorkers; i++ {
// 		go worker(&wg, c1, results, width, height)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()
//
// 	frames := make(map[int]string)
//
// 	for j := range results {
// 		frames[j.Index] = j.Ascii
// 		pe := int((float32(len(frames)) / float32(len(images))) * 100)
// 		fmt.Println("Progress: ", pe)
// 	}
//
// 	m.Frames = frames
//
// 	for i := 0; i < len(m.Frames); i++ {
// 		fmt.Println(m.Frames[i])
// 	}
// }
