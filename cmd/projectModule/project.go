package projectModule

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nativeblocks/cli/cmd/authModule"
	"github.com/nativeblocks/cli/cmd/organizationModule"
	"github.com/nativeblocks/cli/cmd/regionModule"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func ProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	cmd.AddCommand(projectSetCmd())
	cmd.AddCommand(projectGetCmd())
	cmd.AddCommand(projectSchemaGenCmd())
	return cmd
}

func projectSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Select a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := regionModule.GetRegion(*fm)
			if err != nil {
				return err
			}

			auth, err := authModule.AuthGet(*fm)
			if err != nil {
				return err
			}

			organization, err := organizationModule.GetOrganization(*fm)
			if err != nil {
				return err
			}

			projects, err := GetProjects(region.Url, auth.AccessToken, organization.Id)
			if err != nil {
				return err
			}

			var options []string
			optionMap := make(map[string]ProjectModel)

			for _, proj := range projects {
				optionText := fmt.Sprintf("%s (%s) - %s", proj.Name, proj.Id, proj.Platform)
				options = append(options, optionText)
				optionMap[optionText] = proj
			}

			var selection string
			prompt := &survey.Select{
				Message: "Choose a project:",
				Options: options,
			}

			if err := survey.AskOne(prompt, &selection); err != nil {
				return errors.New("selection cancelled: " + err.Error())
			}

			selectedProj := optionMap[selection]
			err = SelectProject(*fm, &selectedProj)
			if err != nil {
				return err
			}

			fmt.Printf("Selected project: %s (%s)\n", selectedProj.Name, selectedProj.Id)
			if selectedProj.Id != "" {
				fmt.Printf("API Key '%s' is configured for use\n", selectedProj.APIKeys[0].Name)
			} else {
				fmt.Printf("Warning: No API keys available for this project\n")
			}
			return nil
		},
	}
}

func projectGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}
			project, err := GetProject(*fm)
			if err != nil {
				return err
			}
			fmt.Printf("Current project: %s \n", project.Name)
			return nil
		},
	}
}

func projectSchemaGenCmd() *cobra.Command {
	var edition string
	var path string
	cmd := &cobra.Command{
		Use:   "gen-schema",
		Short: "Generate a JSON schema file",
		RunE: func(cmd *cobra.Command, args []string) error {
			finalDir := path + "/.nativeblocks"
			inputFm, err := fileutil.NewFileManager(&finalDir)
			if err != nil {
				return err
			}
			blockKeyTypes := make([]string, 0)
			blockProperties := make([]string, 0)
			blockData := make([]string, 0)
			blockSlots := make([]string, 0)
			blockEvents := make([]string, 0)

			actionKeyTypes := make([]string, 0)
			actionProperties := make([]string, 0)
			actionData := make([]string, 0)

			blocks := make(map[string]interface{})
			actions := make(map[string]interface{})

			version := time.Now().UTC().Format("2006-01-02 15:04")

			if edition == "cloud" || edition == "Cloud" || edition == "CLOUD" {
				fm, err := fileutil.NewFileManager(nil)
				if err != nil {
					return err
				}
				region, err := regionModule.GetRegion(*fm)
				if err != nil {
					return err
				}

				auth, err := authModule.AuthGet(*fm)
				if err != nil {
					return err
				}

				organization, err := organizationModule.GetOrganization(*fm)
				if err != nil {
					return err
				}

				project, err := GetProject(*fm)
				if err != nil {
					return err
				}

				installedBlocks, err := GetInstalledIntegration(region.Url, auth.AccessToken, organization.Id, project.Id, "BLOCK")
				if err != nil {
					return err
				}

				for _, installedIntegration := range installedBlocks {
					blockKeyTypes = append(blockKeyTypes, installedIntegration.IntegrationKeyType)
					for _, property := range installedIntegration.IntegrationProperties {
						blockProperties = append(blockProperties, property.Key)
					}
					for _, dataItem := range installedIntegration.IntegrationData {
						blockData = append(blockData, dataItem.Key)
					}
					for _, slot := range installedIntegration.IntegrationSlots {
						blockSlots = append(blockSlots, slot.Slot)
					}
					for _, event := range installedIntegration.IntegrationEvents {
						blockEvents = append(blockEvents, event.Event)
					}
					block := map[string]interface{}{
						"keyType":    installedIntegration.IntegrationKeyType,
						"version":    installedIntegration.IntegrationVersion,
						"data":       installedIntegration.IntegrationData,
						"properties": installedIntegration.IntegrationProperties,
						"slots":      installedIntegration.IntegrationSlots,
						"events":     installedIntegration.IntegrationEvents,
					}

					if installedIntegration.IntegrationData == nil {
						block["data"] = []interface{}{}
					}
					if installedIntegration.IntegrationProperties == nil {
						block["properties"] = []interface{}{}
					}
					if installedIntegration.IntegrationSlots == nil {
						block["slots"] = []interface{}{}
					}
					if installedIntegration.IntegrationEvents == nil {
						block["events"] = []interface{}{}
					}

					blocks[installedIntegration.IntegrationKeyType] = block
				}

				installedActions, err := GetInstalledIntegration(region.Url, auth.AccessToken, organization.Id, project.Id, "ACTION")
				if err != nil {
					return err
				}

				for _, installedIntegration := range installedActions {
					actionKeyTypes = append(actionKeyTypes, installedIntegration.IntegrationKeyType)
					for _, property := range installedIntegration.IntegrationProperties {
						actionProperties = append(actionProperties, property.Key)
					}
					for _, dataItem := range installedIntegration.IntegrationData {
						actionData = append(actionData, dataItem.Key)
					}
					action := map[string]interface{}{
						"keyType":    installedIntegration.IntegrationKeyType,
						"version":    installedIntegration.IntegrationVersion,
						"data":       installedIntegration.IntegrationData,
						"properties": installedIntegration.IntegrationProperties,
						"events":     installedIntegration.IntegrationEvents,
					}

					if installedIntegration.IntegrationData == nil {
						action["data"] = []interface{}{}
					}
					if installedIntegration.IntegrationProperties == nil {
						action["properties"] = []interface{}{}
					}
					if installedIntegration.IntegrationEvents == nil {
						action["events"] = []interface{}{}
					}

					actions[installedIntegration.IntegrationKeyType] = action
				}
			} else {
				blockExist := inputFm.FileExists("integrations/block")
				if blockExist {
					blockKeyTypes = findKeyTypes(inputFm.BaseDir + "/integrations/block")
					blockProperties = findProperties(inputFm.BaseDir + "/integrations/block")
					blockData = findData(inputFm.BaseDir + "/integrations/block")
					blockSlots = findSlots(inputFm.BaseDir + "/integrations/block")
					blockEvents = findEvents(inputFm.BaseDir + "/integrations/block")

					blocks = findIntegrations(inputFm.BaseDir + "/integrations/block")
				}
				actionExist := inputFm.FileExists("integrations/action")
				if actionExist {
					actionKeyTypes = findKeyTypes(inputFm.BaseDir + "/integrations/action")
					actionProperties = findProperties(inputFm.BaseDir + "/integrations/action")
					actionData = findData(inputFm.BaseDir + "/integrations/action")
					actions = findIntegrations(inputFm.BaseDir + "/integrations/action")
				}
				blockKeyTypes = append(blockKeyTypes, "ROOT")
			}

			blocks["schema-version"] = version
			actions["schema-version"] = version

			schema, err := generateBaseSchema(version, blockKeyTypes, actionKeyTypes, blockProperties, blockData, blockSlots, blockEvents, actionProperties, actionData)
			if err != nil {
				return nil
			}

			if err := inputFm.SaveToFile("schema.json", schema); err != nil {
				return err
			}

			if err := inputFm.SaveToFile("blocks.json", blocks); err != nil {
				return err
			}
			if err := inputFm.SaveToFile("actions.json", actions); err != nil {
				return err
			}

			fmt.Printf("Schema file generated successfully at %s \n", inputFm.BaseDir)
			return nil
		},
	}
	cmd.Flags().StringVarP(&edition, "edition", "e", "", "Edition type (cloud or community)")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Output path")
	_ = cmd.MarkFlagRequired("edition")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}
