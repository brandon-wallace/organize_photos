package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

const (
	format    = "2006/01/02"
	extension = "*.jpg"
)

// makeNewDirectory makes a directory and all parent directories
func makeNewDirectory(directory string) {
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// getFiles lists all JPG files in the current directory
func getFiles() ([]string, error) {
	files, err := filepath.Glob(extension)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// readMetadata extracts the datetime information from the metadata
func readMetadata(image string) string {
	img, err := os.Open(image)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	metadata, err := exif.Decode(img)
	if err != nil {
		log.Fatal(err)
	}
	datetime, _ := metadata.DateTime()

	return datetime.Format(format)
}

// moveFile moves a file from a source to a destination
func moveFile(file string, destination string) error {
	return os.Rename(file, fmt.Sprintf("%s/%s", destination, file))
}

func main() {
	fmt.Println("Organizing JPG Files")
	files, err := getFiles()
	if err != nil {
		log.Fatal("error: something has gone wrong")
	}

	count := 0
	for _, img := range files {
		directory := readMetadata(img)
		makeNewDirectory(directory)
		err := moveFile(img, directory)
		if err != nil {
			log.Fatal("error: ", err)
		}
		count++
	}
	fmt.Printf("%d files moved.\n", count)
}
