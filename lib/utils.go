package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CreateFile writes the given content to a file with the specified name.
func CreateFile(filePath string, content string) (string, error) {

	if filePath == "" {
		// Create a temporary file in the default temp directory
		tempFile, err := ioutil.TempFile("", "ansible_wrapper-*.sh")
		if err != nil {
			return "", fmt.Errorf("error creating file: %w", err)
		}
		// Schedule the file to be deleted when the program exits
		defer os.Remove(tempFile.Name())

		// Write some data to the file
		if _, err := tempFile.Write([]byte("Temporary file content")); err != nil {
			return "", fmt.Errorf("error writing to temporary file: %w", err)
		}

		// Close the file
		if err := tempFile.Close(); err != nil {
			return "", fmt.Errorf("error closing temporary file: %w", err)
		}

		// Do other tasks
		fmt.Println("Temporary file will be deleted upon program exit.")
		return tempFile.Name(), nil
	} else {
		var fullyQualifiedFilePath string
		if strings.HasPrefix(filePath, "~/") {
			dirname, _ := os.UserHomeDir()
			fullyQualifiedFilePath = filepath.Join(dirname, filePath[2:])
		}
		file, err := os.Create(fullyQualifiedFilePath)
		if err != nil {
			return fullyQualifiedFilePath, fmt.Errorf("error creating file: %w", err)
		}
		defer file.Close() // Ensure the file is closed when the function exits
		// Write the content to the file
		_, err = file.WriteString(content)
		if err != nil {
			return fullyQualifiedFilePath, fmt.Errorf("error writing to file: %w", err)
		}
		return fullyQualifiedFilePath, nil
	}

}
