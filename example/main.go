package main

import (
	"bytes"
	"fmt"
	"os"

	ar "github.com/MStoykov/go-libarchive"
)

func printContents(filename string) {
	fmt.Println("file ", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening file:\n %s\n", err)
		return
	}
	archive, err := ar.NewArchive(file)
	if err != nil {
		fmt.Printf("Error on NewArchive\n %s\n", err)
	}
	defer archive.Free()
	defer archive.Close()
	for {
		entry, err := archive.Next()
		if err != nil {
			fmt.Printf("Error on reader.Next():\n%s\n", err)
			return
		}
		fmt.Printf("Name %s\n", entry.PathName())
		var buf bytes.Buffer
		size, err := buf.ReadFrom(archive)

		if err != nil {
			fmt.Printf("Error on reading entry from archive:\n%s\n", err)
		}
		if size > 0 {
			fmt.Println("Contents:\n***************", buf.String(), "*********************")
		}
	}
}

func main() {
	for _, filename := range os.Args[1:] {
		printContents(filename)
	}
}
