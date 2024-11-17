package project

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Schema struct {
	Schema      string      `json:"$schema"`
	Type        string      `json:"type"`
	Required    interface{} `json:"required"`
	Properties  interface{} `json:"properties"`
	Definitions interface{} `json:"definitions"`
}

type Block struct {
	KeyType            string           `json:"keyType"`
	Key                string           `json:"key"`
	VisibilityKey      string           `json:"visibilityKey"`
	Slot               string           `json:"slot"`
	Slots              []map[string]any `json:"slots"`
	IntegrationVersion int              `json:"integrationVersion"`
	Data               []map[string]any `json:"data"`
	Properties         []map[string]any `json:"properties"`
	Actions            []map[string]any `json:"actions"`
	Blocks             []Block          `json:"blocks"`
}

type Trigger struct {
	KeyType            string           `json:"keyType"`
	Then               string           `json:"then"`
	IntegrationVersion int              `json:"integrationVersion"`
	Name               string           `json:"name"`
	Properties         []map[string]any `json:"properties"`
	Data               []map[string]any `json:"data"`
	Triggers           []Trigger        `json:"triggers"`
}

func generateBaseSchema(blockKeyTypes, actionKeyTypes, blockProperties, blockData, blockSlots, blockEvents, actionProperties, actionData []string) (Schema, error) {
	baseSchema := Schema{
		Schema:   "http://json-schema.org/draft-07/schema#",
		Type:     "object",
		Required: []string{"name", "route", "isStarter", "type", "variables", "blocks"},
		Properties: map[string]interface{}{
			"name": map[string]string{
				"type": "string",
			},
			"route": map[string]string{
				"type": "string",
			},
			"type": map[string]interface{}{
				"type": "string",
				"enum": []string{"FRAME", "BOTTOM_SHEET", "DIALOG"},
			},
			"isStarter": map[string]string{
				"type": "boolean",
			},
			"variables": map[string]interface{}{
				"type":        "array",
				"uniqueItems": true,
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"key", "value", "type"},
					"properties": map[string]interface{}{
						"key": map[string]string{
							"type": "string",
						},
						"value": map[string]string{
							"type": "string",
						},
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"},
						},
					},
				},
			},
			"blocks": map[string]interface{}{
				"type":     "array",
				"maxItems": 1,
				"items": map[string]interface{}{
					"$ref": "#/definitions/block",
				},
			},
		},
		Definitions: map[string]interface{}{
			"block": map[string]interface{}{
				"type": "object",
				"required": []string{
					"keyType", "key", "visibilityKey", "slot", "slots", "integrationVersion",
					"data", "properties", "actions", "blocks",
				},
				"properties": map[string]interface{}{
					"keyType": map[string]interface{}{
						"type": "string",
						"enum": blockKeyTypes,
					},
					"key": map[string]string{
						"type": "string",
					},
					"visibilityKey": map[string]string{
						"type": "string",
					},
					"slot": map[string]string{
						"type": "string",
					},
					"slots": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"slot": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(blockSlots),
								},
							},
						},
					},
					"integrationVersion": map[string]string{
						"type": "integer",
					},
					"data": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"key", "value", "type"},
							"properties": map[string]interface{}{
								"key": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(blockData),
								},
								"value": map[string]string{
									"type": "string",
								},
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"},
								},
							},
						},
					},
					"properties": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"key", "valueMobile", "valueTablet", "valueDesktop", "type"},
							"properties": map[string]interface{}{
								"key": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(blockProperties),
								},
								"valueMobile": map[string]string{
									"type": "string",
								},
								"valueTablet": map[string]string{
									"type": "string",
								},
								"valueDesktop": map[string]string{
									"type": "string",
								},
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"},
								},
							},
						},
					},
					"actions": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"event", "triggers"},
							"properties": map[string]interface{}{
								"event": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(blockEvents),
								},
								"triggers": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"$ref": "#/definitions/trigger",
									},
								},
							},
						},
					},
					"blocks": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"$ref": "#/definitions/block",
						},
					},
				},
			},
			"trigger": map[string]interface{}{
				"type":     "object",
				"required": []string{"keyType", "then", "name", "integrationVersion", "properties", "data", "triggers"},
				"properties": map[string]interface{}{
					"keyType": map[string]interface{}{
						"type": "string",
						"enum": actionKeyTypes,
					},
					"then": map[string]interface{}{
						"type": "string",
						"enum": []string{"NEXT", "END", "SUCCESS", "FAILURE"},
					},
					"integrationVersion": map[string]string{
						"type": "integer",
					},
					"name": map[string]string{
						"type": "string",
					},
					"properties": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"key", "value", "type"},
							"properties": map[string]interface{}{
								"key": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(actionProperties),
								},
								"value": map[string]string{
									"type": "string",
								},
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"},
								},
							},
						},
					},
					"data": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":     "object",
							"required": []string{"key", "value", "type"},
							"properties": map[string]interface{}{
								"key": map[string]interface{}{
									"type": "string",
									"enum": getUniqueKeys(actionData),
								},
								"value": map[string]string{
									"type": "string",
								},
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"},
								},
							},
						},
					},
					"triggers": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"$ref": "#/definitions/trigger",
						},
					},
				},
			},
		},
	}

	return baseSchema, nil
}

func getUniqueKeys[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func findKeyTypes(dirPath string) []string {
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
		fmt.Printf("Error walking path: %v\n", err)
		return nil
	}

	return keyTypes
}

func findData(dirPath string) []string {
	var data []string
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
				data = append(data, key)
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
		return nil
	}

	return data
}

func findProperties(dirPath string) []string {
	var properties []string
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
				properties = append(properties, key)
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
		return nil
	}

	return properties
}

func findSlots(dirPath string) []string {
	var slots []string
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == "slots.json" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData []map[string]interface{}
			if err := json.Unmarshal(fileContent, &jsonData); err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}

			for _, item := range jsonData {
				key, _ := item["slot"].(string)
				slots = append(slots, key)
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
		return nil
	}

	return slots
}

func findEvents(dirPath string) []string {
	var events []string
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == "events.json" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData []map[string]interface{}
			if err := json.Unmarshal(fileContent, &jsonData); err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}

			for _, item := range jsonData {
				key, _ := item["event"].(string)
				events = append(events, key)
			}
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
		return nil
	}

	return events
}
