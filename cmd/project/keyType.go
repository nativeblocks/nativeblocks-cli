package project

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type MetaItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func FindKeyTypes(dirPath string) []string {
	var keyTypes []string
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == "integration.json" {
			fileContent, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			var jsonData map[string]interface{}
			if err := json.Unmarshal(fileContent, &jsonData); err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}

			if keyType, ok := jsonData["keyType"].(string); ok {
				keyTypes = append(keyTypes, keyType)
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return nil
	}

	return keyTypes
}

func FindData(dirPath string) []MetaItem {
	var data []MetaItem
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == "data.json" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData []map[string]interface{}
			if err := json.Unmarshal(fileContent, &jsonData); err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}

			for _, item := range jsonData {
				key, _ := item["key"].(string)
				dataType, _ := item["type"].(string)
				data = append(data, MetaItem{
					Key:   key,
					Value: "",
					Type:  dataType,
				})
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return nil
	}

	return data
}

func FindProperties(dirPath string) []MetaItem {
	var properties []MetaItem
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == "properties.json" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData []map[string]interface{}
			if err := json.Unmarshal(fileContent, &jsonData); err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}

			for _, item := range jsonData {
				key, _ := item["key"].(string)
				value, _ := item["value"].(string)
				dataType, _ := item["type"].(string)
				properties = append(properties, MetaItem{
					Key:   key,
					Value: value,
					Type:  dataType,
				})
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return nil
	}

	return properties
}
