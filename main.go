package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/gotk3/gotk3/gdk"
)

func loadPixbuf(wg *sync.WaitGroup, pb **gdk.Pixbuf, file string) {
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

	*pb, err = loader.GetPixbuf()
	if err != nil {
		log.Printf("GetPixBuf: %q\n", err)
		return
	}
}

func loadPixbufs(files []string) (rv []*gdk.Pixbuf) {
	rv = make([]*gdk.Pixbuf, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))
	for i, f := range files {
		var pb *gdk.Pixbuf
		go loadPixbuf(&wg, &pb, f)
		rv[i] = pb
	}
	wg.Wait()
	return
}

func main() {
	if len(os.Args) < 2 {
		log.Printf("usage: pixbufloader pic1 [pic2 ... picN]\n")
	}
	pbs := loadPixbufs(os.Args[1:])
	good := 0
	for _, pb := range pbs {
		if pb != nil {
			good++
		}
	}
	log.Printf("got %d/%d pixbufs\n", good, len(pbs))
}
