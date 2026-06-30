package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

const (
	format    = "2006/01/02"
)

var (
	extensions = []string{"*.jpg", "*.JPG", "*.jpeg", "*.JPEG", "*.arw", "*.ARW"}
)

// makeNewDirectory makes a directory and all parent directories
func makeNewDirectory(directory string) error {
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		slog.Error("error", "failed to create directory", err.Error())
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}

	return nil
}

// getFiles lists all JPG files in the current directory
func getFiles(src *string) ([]string, error) {
	files := make([]string, 0)
	for _, ext := range extensions {
		matches, err := filepath.Glob(*src + "/" + ext)
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

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	return os.Rename(file, fmt.Sprintf("%s/%s", destination, file))
}

func main() {
	src := flag.String("src", ".", "source directory")
	dst := flag.String("dst", ".", "destination directory")
	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Fprintf(os.Stderr, "Usage: ./organize_photos -src /path/to/images -dst /path/to/destination")
		os.Exit(1)
	}

	fmt.Println("\n Organizing JPG Files...")
	slog.Info("info", "Organizing image files...", "")

	files, err := getFiles(src)
	if err != nil {
		slog.Error("error", "something has gone wrong", nil)
	}
	
	if len(files) == 0 {
		slog.Error("error", "no image files to process", nil)
		return
	}

	errorCount := 0
	successCount := 0
	duplicateCount := 0

	for _, img := range files {
		directory, err := readMetadata(img)
		if err != nil {
			slog.Error("error", "getting date for", err.Error())
			errorCount++
			continue
		}

		if err := makeNewDirectory(directory); err != nil {
			errorCount++
			slog.Error("error", "creating directory for", err.Error())
			continue
		}
		err = moveFile(img, directory)
		if err != nil {
			duplicateCount++
			slog.Warn("warn", "could not move file", err.Error())
		}

		successCount++
	}

	template := ` Organization complete.
 Successfully moved files: %d
 Duplicate files: %d
 Errors: %d
 Total files moved: %d.
 `

	fmt.Printf(template, successCount, duplicateCount, errorCount, len(files))
}

