package utils

import (
	"fmt"
	"io"
	"os"
)

// Utility function to close an io.Closer and log errors without returning them
func CloseAndLogError(closer io.Closer) {
	if closer == nil {
		return
	}

	// Attempt to close the resource (e.g., an HTTP response or a file).
	// If an error occurs during the close operation, the error is captured.
	if err := closer.Close(); err != nil {
		fmt.Printf("Error closing resource: %v", err)
	}
}

// utility function to simplified function to remove a directory and log any errors
func RemoveDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		fmt.Printf("Error removing directory %s: %v", path, err)
	}
}
