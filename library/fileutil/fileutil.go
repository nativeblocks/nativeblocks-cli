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
	BaseDir string
}

func NewFileManager(customDir *string) (*FileManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home path: %v", err)
	}

	var baseDir string
	if customDir != nil {
		baseDir = filepath.Join(*customDir)
	} else {
		baseDir = filepath.Join(homeDir, ConfigDirName, CliDirName)
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config path: %v", err)
	}

	return &FileManager{BaseDir: baseDir}, nil
}

func (fm *FileManager) SaveToFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	filePath := filepath.Join(fm.BaseDir, filename)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

func (fm *FileManager) LoadFromFile(filename string, target interface{}) error {
	filePath := filepath.Join(fm.BaseDir, filename)
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
	filePath := filepath.Join(fm.BaseDir, filename)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

func (fm *FileManager) FileExists(filename string) bool {
	filePath := filepath.Join(fm.BaseDir, filename)
	_, err := os.Stat(filePath)
	return err == nil
}

func (fm *FileManager) GetFilePath(filename string) string {
	return filepath.Join(fm.BaseDir, filename)
}

func GetFileDir(filePath string) string {
	return filepath.Dir(filePath)
}

func GetFileName(filePath string) string {
	return filepath.Base(filePath)
}
