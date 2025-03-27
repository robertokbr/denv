package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// createZipArchive creates a zip archive of the specified directory
func createZipArchive(sourceDir, zipPath string) error {
	// Create the zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the directory
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if filePath == sourceDir {
			return nil
		}

		// Create a relative path for the file in the zip
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// Create a zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		// If it's a directory, just create the header
		if info.IsDir() {
			header.Name += "/"
			_, err = zipWriter.CreateHeader(header)
			return err
		}

		// Create a writer for the file
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open the file to copy its contents
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file contents to the zip
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// extractZipArchive extracts a zip file to the specified directory
func extractZipArchive(zipPath, extractDir string) error {
	// Open the zip file
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipFile.Close()

	// Create the extraction directory
	err = os.MkdirAll(extractDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create extraction directory: %v", err)
	}

	// Extract each file
	for _, file := range zipFile.File {
		// Construct the full path for the file
		filePath := filepath.Join(extractDir, file.Name)

		// Check for path traversal
		if !strings.HasPrefix(filePath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		// Create directories as needed
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, file.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		// Create the file
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		// Create the file
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		// Open the source file
		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open zip file: %v", err)
		}

		// Copy the contents
		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()

		if err != nil {
			return fmt.Errorf("failed to copy file contents: %v", err)
		}
	}

	return nil
}
