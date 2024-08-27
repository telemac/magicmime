package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/zRedShift/mimemagic"
)

func main() {
	// Check if a file path is provided as a command-line argument
	if len(os.Args) < 2 {
		slog.Error("No file path provided. Usage: go run main.go <file-path>")
		return
	}

	// Get the file path from the command-line arguments
	filePath := os.Args[1]

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("Error opening file", slog.String("file", filePath), slog.Any("error", err))
		return
	}
	defer file.Close()

	// Match the file contents against MIME types using MatchReader
	mimeType, err := mimemagic.MatchReader(file, "")
	if err != nil {
		slog.Error("Error detecting MIME type", slog.String("file", filePath), slog.Any("error", err))
		return
	}

	// Log MIME type information
	slog.Info("MIME type detected", slog.String("file", filePath), slog.Any("Extensions", mimeType.Extensions))

	// Check if any extensions were detected
	if len(mimeType.Extensions) == 0 {
		slog.Error("No valid MIME type detected", slog.String("file", filePath))
		return
	}

	// Get the first extension from the detected MIME type
	newExtension := mimeType.Extensions[0]

	// Ensure the new extension has a leading dot
	if !strings.HasPrefix(newExtension, ".") {
		newExtension = "." + newExtension
	}

	// Get the directory and base name of the original file path
	dir := filepath.Dir(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// Create a new file path with the detected extension
	newFilePath := filepath.Join(dir, baseName+newExtension)

	// Rename the file with the new extension
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		slog.Error("Error renaming file", slog.String("oldPath", filePath), slog.String("newPath", newFilePath), slog.Any("error", err))
		return
	}

	// Log the successful file renaming
	slog.Info("File renamed successfully", slog.String("oldPath", filePath), slog.String("newPath", newFilePath))
}
