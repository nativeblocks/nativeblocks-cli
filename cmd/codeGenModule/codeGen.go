package codeGenModule

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/jsonutil"
	"github.com/spf13/cobra"
	"strings"
)

func CodeGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code-gen",
		Short: "Generate blocks and actions codes",
	}

	cmd.AddCommand(genJSCmd())

	return cmd
}

func genJSCmd() *cobra.Command {
	var path string
	var blocksSchema string
	var actionsSchema string
	cmd := &cobra.Command{
		Use:   "js",
		Short: "Generate JS",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := fileutil.GetFileDir(path + "/")

			fm, err := fileutil.NewFileManager(&baseDir)
			if err != nil {
				return err
			}

			blocks := make(map[string]interface{})
			err = jsonutil.FetchJSONFromURL(blocksSchema, &blocks)
			if err != nil {
				return err
			}
			actions := make(map[string]interface{})
			err = jsonutil.FetchJSONFromURL(actionsSchema, &actions)
			if err != nil {
				return err
			}

			for key, value := range blocks {
				if key == "schema-version" {
					continue
				}
				componentBytes, err := json.Marshal(value)
				if err != nil {
					fmt.Println("Error marshaling component:", err)
					continue
				}

				var component Integration
				err = json.Unmarshal(componentBytes, &component)
				if err != nil {
					fmt.Println("Error unmarshalling component:", err)
					continue
				}
				name := key + "-Block"
				block := generateJSClass(strcase.ToLowerCamel(name), component, "BLOCK")
				if err := fm.SaveByteToFile(strcase.ToLowerCamel(name)+".js", []byte(block)); err != nil {
					return err
				}
			}
			for key, value := range actions {
				if key == "schema-version" {
					continue
				}
				componentBytes, err := json.Marshal(value)
				if err != nil {
					fmt.Println("Error marshaling component:", err)
					continue
				}

				var component Integration
				err = json.Unmarshal(componentBytes, &component)
				if err != nil {
					fmt.Println("Error unmarshalling component:", err)
					continue
				}
				name := key + "-Action"
				action := generateJSClass(strcase.ToLowerCamel(name), component, "ACTION")
				if err := fm.SaveByteToFile(strcase.ToLowerCamel(name)+".js", []byte(action)); err != nil {
					return err
				}
			}
			fmt.Printf("JS classes generated: %v \n", baseDir)
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "Output path")
	cmd.Flags().StringVarP(&blocksSchema, "blocksSchemaUrl", "b", "", "Blocks schema url")
	cmd.Flags().StringVarP(&actionsSchema, "actionsSchemaUrl", "a", "", "Blocks schema url")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("blocksSchemaUrl")
	_ = cmd.MarkFlagRequired("actionsSchemaUrl")
	return cmd
}

func generateJSClass(blockName string, component Integration, kind string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// PLEASE DO NOT EDIT THIS FILE, THIS IS GENERATED BY NATIVEBLOCKS \n"))
	sb.WriteString(fmt.Sprintf("class %s {\n", blockName))

	sb.WriteString("  keyType;\n")
	sb.WriteString("  key;\n")
	sb.WriteString("  visibilityKey;\n")
	sb.WriteString("  slot;\n")
	sb.WriteString("  integrationVersion;\n")
	sb.WriteString("  properties;\n")
	sb.WriteString("  data;\n")
	sb.WriteString("  events;\n")

	if kind == "BLOCK" {
		sb.WriteString("  actions;\n")
		sb.WriteString("  blocks;\n")
		sb.WriteString("  slots;\n")
	} else {
		sb.WriteString("  triggers;\n")
	}
	sb.WriteString("\n")

	sb.WriteString("  constructor(config = {}) {\n")
	sb.WriteString(fmt.Sprintf("    this.keyType = \"%s\";\n", component.KeyType))
	sb.WriteString("    this.key = config.key;\n")
	sb.WriteString("    this.visibilityKey = config.visibilityKey;\n")

	if kind == "BLOCK" {
		sb.WriteString("    this.slot = config.slot || \"content\";\n")
		sb.WriteString("    this.actions = [];\n")
		sb.WriteString("    this.blocks = [];\n")
		sb.WriteString("    this.slots = [{slot: this.slot}];\n")
	} else {
		sb.WriteString("    this.triggers = [];\n")
	}

	sb.WriteString(fmt.Sprintf("    this.integrationVersion = %v;\n", component.Version))
	sb.WriteString("    this.properties = [];\n")
	sb.WriteString("    this.data = [];\n")
	sb.WriteString("    this.events = [];\n")

	sb.WriteString("    const initialProperties = [\n")
	for _, prop := range component.Properties {
		if kind == "BLOCK" {
			sb.WriteString(fmt.Sprintf("      { key: \"%s\", valueMobile: \"%s\", valueTablet: \"%s\", valueDesktop: \"%s\", type: \"%s\" },\n",
				prop.Key,
				prop.Value,
				prop.Value,
				prop.Value,
				prop.Type))
		} else {
			sb.WriteString(fmt.Sprintf("      { key: \"%s\", value: \"%s\", type: \"%s\" },\n",
				prop.Key,
				prop.Value,
				prop.Type))
		}
	}
	sb.WriteString("    ];\n")
	sb.WriteString("    this.properties.push(...initialProperties);\n")

	sb.WriteString("    const initialData = [\n")
	for _, data := range component.Data {
		sb.WriteString(fmt.Sprintf("      { key: \"%s\", value: null, type: \"%s\" },\n",
			data.Key,
			data.Type))
	}
	sb.WriteString("    ];\n")
	sb.WriteString("    this.data.push(...initialData);\n")

	sb.WriteString("    const initialEvents = [\n")
	for _, event := range component.Events {
		eventName := event.(map[string]interface{})["event"].(string)
		sb.WriteString(fmt.Sprintf("      { event: \"%s\", triggers: [] },\n",
			eventName))
	}
	sb.WriteString("    ];\n")
	sb.WriteString("    this.events.push(...initialEvents);\n")
	sb.WriteString("  }\n\n")

	sb.WriteString("  isValidEvent(eventName) {\n")
	sb.WriteString("    return this.events.some(e => e.event === eventName);\n")
	sb.WriteString("  }\n\n")

	if kind == "BLOCK" {
		sb.WriteString("  isValidBlock(block) {\n")
		sb.WriteString("    return block && \n")
		sb.WriteString("           typeof block === \"object\" && \n")
		sb.WriteString("           block.key && \n")
		sb.WriteString("           block.keyType;\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  addAction(event) {\n")
		sb.WriteString("    if (!this.isValidEvent(event)) {\n")
		sb.WriteString("      throw new Error(`Invalid event: ${event}. Must be one of: ${this.events.map(e => e.event).join(\", \")}`);\n")
		sb.WriteString("    }\n")
		sb.WriteString("    const action = { event, triggers: [] };\n")
		sb.WriteString("    this.actions.push(action);\n")
		sb.WriteString("    return this;\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  getAction(event) {\n")
		sb.WriteString("    return this.actions.find(a => a.event === event);\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  addBlock(block) {\n")
		sb.WriteString("    if (!this.isValidBlock(block)) {\n")
		sb.WriteString("      throw new Error(\"Invalid block: Block must be an object with at least key and keyType\");\n")
		sb.WriteString("    }\n")
		sb.WriteString("    this.blocks.push(block);\n")
		sb.WriteString("    return this;\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  getBlock(key) {\n")
		sb.WriteString("    return this.blocks.find(b => b.key === key);\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  getBlocks() {\n")
		sb.WriteString("    return [...this.blocks];\n")
		sb.WriteString("  }\n\n")
	} else {
		sb.WriteString("  isValidTrigger(trigger) {\n")
		sb.WriteString("    return trigger && \n")
		sb.WriteString("           typeof trigger === \"object\" && \n")
		sb.WriteString("           trigger.key && \n")
		sb.WriteString("           trigger.keyType && \n")
		sb.WriteString("           trigger.then;\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  addTrigger(trigger) {\n")
		sb.WriteString("    if (!this.isValidTrigger(trigger)) {\n")
		sb.WriteString("      throw new Error(\"Invalid trigger: Trigger must be an object with key, keyType, and then properties\");\n")
		sb.WriteString("    }\n")
		sb.WriteString("    this.triggers.push({\n")
		sb.WriteString("      key: trigger.key,\n")
		sb.WriteString("      keyType: trigger.keyType,\n")
		sb.WriteString("      then: trigger.then,\n")
		sb.WriteString("      properties: trigger.properties || [],\n")
		sb.WriteString("      data: trigger.data || []\n")
		sb.WriteString("    });\n")
		sb.WriteString("    return this;\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  getTrigger(key) {\n")
		sb.WriteString("    return this.triggers.find(t => t.key === key);\n")
		sb.WriteString("  }\n\n")

		sb.WriteString("  getTriggers() {\n")
		sb.WriteString("    return [...this.triggers];\n")
		sb.WriteString("  }\n\n")
	}

	sb.WriteString("  getAvailableEvents() {\n")
	sb.WriteString("    return this.events.map(e => e.event);\n")
	sb.WriteString("  }\n\n")

	sb.WriteString("  getProperty(key) {\n")
	sb.WriteString("    return this.properties.find(prop => prop.key === key);\n")
	sb.WriteString("  }\n\n")

	if kind == "BLOCK" {
		sb.WriteString("  modifyProperty(key, valueMobile, valueTablet, valueDesktop) {\n")
		sb.WriteString("    const propIndex = this.properties.findIndex(p => p.key === key);\n")
		sb.WriteString("    if (propIndex !== -1) {\n")
		sb.WriteString("      this.properties[propIndex] = { ...this.properties[propIndex], valueMobile, valueTablet, valueDesktop };\n")
		sb.WriteString("    }\n")
		sb.WriteString("    return this;\n")
		sb.WriteString("  }\n\n")
	} else {
		sb.WriteString("  modifyProperty(key, value) {\n")
		sb.WriteString("    const propIndex = this.properties.findIndex(p => p.key === key);\n")
		sb.WriteString("    if (propIndex !== -1) {\n")
		sb.WriteString("      this.properties[propIndex] = { ...this.properties[propIndex], value };\n")
		sb.WriteString("    }\n")
		sb.WriteString("    return this;\n")
		sb.WriteString("  }\n\n")
	}

	sb.WriteString("  getData(key) {\n")
	sb.WriteString("    return this.data.find(d => d.key === key);\n")
	sb.WriteString("  }\n\n")

	sb.WriteString("  assignData(key, value) {\n")
	sb.WriteString("    const dataIndex = this.data.findIndex(d => d.key === key);\n")
	sb.WriteString("    if (dataIndex !== -1) {\n")
	sb.WriteString("      this.data[dataIndex] = { ...this.data[dataIndex], value };\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return this;\n")
	sb.WriteString("  }\n\n")

	sb.WriteString("  build() {\n")
	sb.WriteString("    return {\n")
	sb.WriteString("      keyType: this.keyType,\n")
	sb.WriteString("      key: this.key,\n")
	sb.WriteString("      visibilityKey: this.visibilityKey,\n")
	if kind == "BLOCK" {
		sb.WriteString("      slot: this.slot,\n")
		sb.WriteString("      slots: this.slots,\n")
	}
	sb.WriteString("      data: this.data,\n")
	sb.WriteString("      properties: this.properties,\n")
	sb.WriteString("      integrationVersion: this.integrationVersion,\n")
	if kind == "BLOCK" {
		sb.WriteString("      actions: this.actions,\n")
		sb.WriteString("      blocks: this.blocks\n")
	} else {
		sb.WriteString("      triggers: this.triggers\n")
	}
	sb.WriteString("    };\n")
	sb.WriteString("  }\n")

	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf("module.exports = %s;\n", blockName))

	return sb.String()
}

type Integration struct {
	Data       []DataItem     `json:"data"`
	Events     []interface{}  `json:"events"`
	KeyType    string         `json:"keyType"`
	Properties []PropertyItem `json:"properties"`
	Slots      []interface{}  `json:"slots"`
	Version    int            `json:"version"`
}

type DataItem struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type PropertyItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
