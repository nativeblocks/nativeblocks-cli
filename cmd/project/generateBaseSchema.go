package project

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

func generateBaseSchema(
	blockKeyTypes, actionKeyTypes []string,
	blockProperties, blockData, actionProperties, actionData []MetaItem,
) (Schema, error) {
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
								"slot": map[string]string{
									"type": "string",
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
								"event": map[string]string{
									"type": "string",
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
