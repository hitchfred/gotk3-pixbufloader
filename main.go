package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gotk3/gotk3/gdk"
)

var iterations = flag.Int("iterations", 1, "Number of iterations")

func loadPixbuf(wg *sync.WaitGroup, file string, out chan<- *gdk.Pixbuf) {
	defer wg.Done()

	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		log.Printf("PixBufLoader: %q\n", err)
		return
	}

	fp, err := os.Open(file)
	if err != nil {
		log.Printf("Open: %q\n", err)
		return
	}
	defer fp.Close()

	_, err = io.Copy(loader, fp)
	if err != nil {
		log.Printf("Copy: %q\n", err)
		return
	}

	err = loader.Close()
	if err != nil {
		log.Printf("Close: %q\n", err)
		return
	}

	pb, err := loader.GetPixbuf()
	if err != nil {
		log.Printf("GetPixBuf: %q\n", err)
		return
	}
	out <- pb
}

func loadPixbufs(files []string) <-chan *gdk.Pixbuf {
	rv := make(chan *gdk.Pixbuf, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, f := range files {
		go loadPixbuf(&wg, f, rv)
	}
	go func() {
		wg.Wait()
		close(rv)
	}()
	return rv
}

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) < 1 {
		log.Printf("usage: pixbufloader pic1 [pic2 ... picN]\n")
		return
	}

	totalGood := 0
	total := 0
	for i := 0; i < *iterations; i++ {
		good := 0
		for pb := range loadPixbufs(files) {
			if pb != nil {
				good++
			}
		}
		log.Printf("got %d/%d pixbufs\n", good, len(files))
		totalGood += good
		total += len(files)
	}
	log.Printf("after %d iterations: got %d/%d pixbufs (failed %d)\n",
		*iterations, totalGood, total, total-totalGood)
}
