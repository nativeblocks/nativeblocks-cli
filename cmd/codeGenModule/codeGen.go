package codeGenModule

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/jsonutil"
	"github.com/spf13/cobra"
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

type Integration struct {
	Data       []DataItem     `json:"data"`
	Events     []EventItem    `json:"events"`
	KeyType    string         `json:"keyType"`
	Properties []PropertyItem `json:"properties"`
	Slots      []SlotItem     `json:"slots"`
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

type EventItem struct {
	Event string `json:"event"`
}

type SlotItem struct {
	Slot string `json:"slot"`
}
