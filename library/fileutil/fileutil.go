package fileutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName = ".nativeblocks"
	CliDirName    = "cli"
)

type FileManager struct {
	baseDir string
}

func NewFileManager() (*FileManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}
	fmt.Printf("HHHH", homeDir)

	baseDir := filepath.Join(homeDir, ConfigDirName, CliDirName)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	return &FileManager{baseDir: baseDir}, nil
}

func (fm *FileManager) SaveToFile(filename string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	filePath := filepath.Join(fm.baseDir, filename)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

func (fm *FileManager) LoadFromFile(filename string, target interface{}) error {
	filePath := filepath.Join(fm.baseDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return nil
}

func (fm *FileManager) DeleteFile(filename string) error {
	filePath := filepath.Join(fm.baseDir, filename)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

func (fm *FileManager) FileExists(filename string) bool {
	filePath := filepath.Join(fm.baseDir, filename)
	_, err := os.Stat(filePath)
	return err == nil
}

func (fm *FileManager) GetFilePath(filename string) string {
	return filepath.Join(fm.baseDir, filename)
}
