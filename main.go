package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	// Define a flag to specify the output filename
	outputFileName := flag.String("filename", "output.zip", "Name of the output zip file")
	flag.Parse()

	// Create a new zip archive
	zipWriter := zip.NewWriter(os.Stdout)
	defer zipWriter.Close()

	// Create a new file inside the zip
	fileInZip, err := zipWriter.Create(*outputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ch := make(chan []byte)

	wg := sync.WaitGroup{}

	// Read the content from the channel
	go func() {
		wg.Add(1)
		for content := range ch {
			// Write the content to the file inside the zip
			_, err = fileInZip.Write(content)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
		wg.Done()
	}()

	for {
		buffer := make([]byte, 1024)
		n, err := os.Stdin.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if n == 0 {
			break
		}

		ch <- buffer[:n]
	}
	close(ch)
	wg.Wait()

}
