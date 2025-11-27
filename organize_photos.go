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
)

var (
	extensions = []string{"*.jpg", "*.JPG", "*.jpeg", "*.JPEG"}
)

// makeNewDirectory makes a directory and all parent directories
func makeNewDirectory(directory string) error {
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}

	return nil
}

// getFiles lists all JPG files in the current directory
func getFiles() ([]string, error) {
	files := make([]string, 0)
	for _, ext := range extensions {
		matches, err := filepath.Glob(ext)
		if err != nil {
			return nil, fmt.Errorf("failed to glob pattern %s: %w", ext, err)
		}
		files = append(files, matches...)
	}

	return files, nil
}

// readMetadata extracts the datetime information from the metadata
func readMetadata(image string) (string, error) {
	img, err := os.Open(image)
	if err != nil {
		return "", fmt.Errorf("failed to open image %s %w", image, err)
	}
	defer img.Close()

	metadata, err := exif.Decode(img)
	if err != nil {
		return "", fmt.Errorf("failed to decode EXIF data %s %w", image, err)
	}
	datetime, err := metadata.DateTime()
	if err != nil {
		return "", fmt.Errorf("failed to get datetime from %s %w", image, err)
	}

	return datetime.Format(format), nil
}

// moveFile moves a file from a source to a destination
func moveFile(file string, destination string) error {
	filename := filepath.Base(file)
	path := filepath.Join(destination, filename)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	return os.Rename(file, fmt.Sprintf("%s/%s", destination, file))
}

func main() {
	fmt.Println("\n Organizing JPG Files...")

	files, err := getFiles()
	if err != nil {
		log.Fatal("error: something has gone wrong")
	}
	
	if len(files) == 0 {
		fmt.Println("No image files to process.")
		return
	}

	errorCount := 0
	successCount := 0

	for _, img := range files {
		directory, err := readMetadata(img)
		if err != nil {
			log.Printf("error getting date for %s: %v", img, err)
			errorCount++
			continue
		}

		if err := makeNewDirectory(directory); err != nil {
			errorCount++
			fmt.Printf("error creating directory for %s: %v", img, err)
			continue
		}
		err = moveFile(img, directory)
		if err != nil {
			log.Fatal("error: ", err)
		}

		successCount++
	}

	fmt.Printf("\n Organization complete.\n")
	fmt.Printf(" Successfully moved files: %d\n", successCount)
	fmt.Printf(" Errors: %d\n", errorCount)
	fmt.Printf(" Total files moved: %d.\n", len(files))
}

