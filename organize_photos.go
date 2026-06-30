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
	dateFormat = "2006/01/02"
)

var (
	extensions = []string{"*.jpg", "*.JPG", "*.jpeg", "*.JPEG", "*.arw", "*.ARW"}
)

// stats tracks the results of the organization process
type Stats struct {
	success    int
	duplicates int
	errors     int
}

// makeNewDirectory makes a directory and all parent directories
func makeNewDirectory(directory string) error {
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}

	return nil
}

// getFiles lists all JPG files in the current directory
func getFiles(src string) ([]string, error) {
	files := make([]string, 0)
	for _, ext := range extensions {
		pattern := filepath.Join(src, ext)
		matches, err := filepath.Glob(pattern)
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

	return datetime.Format(dateFormat), nil
}

// moveFile moves a file from a source to a destination without overwriting files
func moveFile(source string, destination string) error {
	if _, err := os.Stat(destination); err == nil {
		return fmt.Errorf("file already exists: %s", destination)
	}

	return os.Rename(source, destination)
}

func main() {
	src := flag.String("src", ".", "source directory")
	dst := flag.String("dst", ".", "destination directory")
	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Fprintf(os.Stderr, "Usage: ./organize_photos -src /path/to/images -dst /path/to/destination")
		os.Exit(1)
	}

	slog.Info("Organizing images...", "source", *src, "destination", *dst)

	files, err := getFiles(*src)
	if err != nil {
		slog.Error("Failed to read source directory", "error", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		slog.Warn("No image files to process")
		return
	}

	stats := &Stats{}

	for _, img := range files {
		dateString, err := readMetadata(img)
		if err != nil {
			slog.Debug("Missing EXIF data", "error", err.Error())
			dateString = "unknown"
		}

		targetDir := filepath.Join(*dst, dateString)

		if err := makeNewDirectory(dateString); err != nil {
			stats.errors++
			slog.Error("Failed to create directory", "error", err.Error())
			continue
		}

		filename := filepath.Base(img)
		targetPath := filepath.Join(targetDir, filename)

		if err := moveFile(img, targetPath); err != nil {
			stats.duplicates++
			slog.Warn("Could not move file", "warn", err.Error())
			continue
		}

		stats.success++
	}

	fmt.Printf("\nOrganization complete.\n")
	fmt.Printf(" Successfully moved: %d\n", stats.success)
	fmt.Printf(" Duplicate files:    %d\n", stats.duplicates)
	fmt.Printf(" Errors:             %d\n", stats.errors)
	fmt.Printf(" Total processed:    %d\n\n", len(files))
}
