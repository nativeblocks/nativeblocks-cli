package codeGenModule

import (
	"encoding/json"
	"errors"
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

	cmd.AddCommand(genTypescriptCmd())
	cmd.AddCommand(genPHPCmd())

	return cmd
}

func baseCodeGen(path string, integrationSchema string, kind string, language string) error {
	baseDir := fileutil.GetFileDir(path + "/")

	fm, err := fileutil.NewFileManager(&baseDir)
	if err != nil {
		return err
	}

	integrations := make(map[string]interface{})
	err = jsonutil.FetchJSONFromURL(integrationSchema, &integrations)
	if err != nil {
		return err
	}

	for key, value := range integrations {
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

		var name string
		if kind == "BLOCK" {
			name = key + "-Block"
		} else {
			name = key + "-Action"
		}

		if language == "TS" {
			integration := generateTypescriptClass(strcase.ToCamel(name), component, kind)
			if err := fm.SaveByteToFile(strcase.ToCamel(name)+".ts", []byte(integration)); err != nil {
				return err
			}
		} else if language == "PHP" {
			block := generatePHPClass(strcase.ToCamel(name), component, kind)
			if err := fm.SaveByteToFile(strcase.ToCamel(name)+".php", []byte(block)); err != nil {
				return err
			}
		} else {
			return errors.New("unsupported language: " + language)
		}
	}
	return nil
}

func genTypescriptCmd() *cobra.Command {
	var path string
	var blocksSchema string
	var actionsSchema string
	cmd := &cobra.Command{
		Use:   "typescript",
		Short: "Generate typescript",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := baseCodeGen(path, blocksSchema, "BLOCK", "TS")
			if err != nil {
				return err
			}
			err = baseCodeGen(path, actionsSchema, "ACTION", "TS")
			if err != nil {
				return err
			}
			fmt.Printf("Typescript classes generated: %v \n", path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "Output path")
	cmd.Flags().StringVarP(&blocksSchema, "blocksSchemaUrl", "b", "", "Blocks schema url")
	cmd.Flags().StringVarP(&actionsSchema, "actionsSchemaUrl", "a", "", "Actions schema url")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("blocksSchemaUrl")
	_ = cmd.MarkFlagRequired("actionsSchemaUrl")
	return cmd
}

func genPHPCmd() *cobra.Command {
	var path string
	var blocksSchema string
	var actionsSchema string
	cmd := &cobra.Command{
		Use:   "php",
		Short: "Generate php",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := baseCodeGen(path, blocksSchema, "BLOCK", "PHP")
			if err != nil {
				return err
			}
			err = baseCodeGen(path, actionsSchema, "ACTION", "PHP")
			if err != nil {
				return err
			}
			fmt.Printf("PHP classes generated: %v \n", path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "Output path")
	cmd.Flags().StringVarP(&blocksSchema, "blocksSchemaUrl", "b", "", "Blocks schema url")
	cmd.Flags().StringVarP(&actionsSchema, "actionsSchemaUrl", "a", "", "Actions schema url")
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
