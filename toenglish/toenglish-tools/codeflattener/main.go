package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Change this to your root directory path
	rootDir := "/Users/vishalvaibhav/Code/github-projects/kahoot-clone-nodejs/server"
	// Change this to your desired output file path
	outputFile := "combined_code.txt"

	// Check if file exists and remove it
	if _, err := os.Stat(outputFile); err == nil {
		err = os.Remove(outputFile)
		if err != nil {
			fmt.Printf("Error removing existing file: %v\n", err)
			return
		}
	}

	// Create/truncate output file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	// Walk through the directory
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files and specific file types you want to exclude
		if strings.HasPrefix(info.Name(), ".") || isExcludedFile(info.Name()) {
			return nil
		}

		// Read file contents
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", path, err)
		}

		// Write file path as comment
		relativePath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %v", path, err)
		}

		header := fmt.Sprintf("\n//%s\n\n", relativePath)
		if _, err := output.WriteString(header); err != nil {
			return fmt.Errorf("error writing header: %v", err)
		}

		// Write file contents
		if _, err := output.Write(content); err != nil {
			return fmt.Errorf("error writing content: %v", err)
		}

		// Add a newline after each file
		if _, err := output.WriteString("\n"); err != nil {
			return fmt.Errorf("error writing newline: %v", err)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Println("Successfully combined all code files!")
}

func isExcludedFile(filename string) bool {
	// Add file extensions or patterns you want to exclude
	excludedExtensions := []string{
		".exe", ".dll", ".so", ".dylib",
		".zip", ".tar", ".gz",
		".jpg", ".png", ".gif",
		".pdf", ".doc", ".docx",
	}

	for _, ext := range excludedExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
