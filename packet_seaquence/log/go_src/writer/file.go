// Package writer provides functionality to write files
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteFile provides functionality to write files
type WriteFile struct {
	OutputFile string
}

// NewWriteFile creates a new WriteFile instance
func NewWriteFile(outputFile string) *WriteFile {
	return &WriteFile{
		OutputFile: outputFile,
	}
}

// Write writes content to a file
func (w *WriteFile) Write(content string) error {
	// Ensure the directory exists
	dir := filepath.Dir(w.OutputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ディレクトリを作成できませんでした: %w", err)
	}
	
	// Write the content to the file
	if err := os.WriteFile(w.OutputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイルに書き込めませんでした: %w", err)
	}
	
	return nil
}

// WriteSplit writes multiple contents to multiple files
// The output files will be named with a suffix like "_1", "_2", etc.
func (w *WriteFile) WriteSplit(contents []string) ([]string, error) {
	if len(contents) == 0 {
		return nil, fmt.Errorf("コンテンツがありません")
	}
	
	// Ensure the directory exists
	dir := filepath.Dir(w.OutputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("ディレクトリを作成できませんでした: %w", err)
	}
	
	// Prepare the base name and extension
	ext := filepath.Ext(w.OutputFile)
	baseName := strings.TrimSuffix(w.OutputFile, ext)
	
	// Write each content to a separate file
	outputFiles := make([]string, 0, len(contents))
	for i, content := range contents {
		// Create the output file name
		outputFile := fmt.Sprintf("%s_%d%s", baseName, i+1, ext)
		
		// Write the content to the file
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return outputFiles, fmt.Errorf("ファイル %s に書き込めませんでした: %w", outputFile, err)
		}
		
		outputFiles = append(outputFiles, outputFile)
	}
	
	return outputFiles, nil
}