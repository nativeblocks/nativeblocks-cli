package project

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
	"github.com/spf13/cobra"
)

const (
	ProjectFileName      = "project"
	RegionFileName       = "region"
	AuthFileName         = "auth"
	OrganizationFileName = "organization"
)

type Project struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Platform string   `json:"platform"`
	APIKeys  []APIKey `json:"apiKeys"`
}

type APIKey struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

type ProjectConfig struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	APIKey  string `json:"apiKey"`
	KeyName string `json:"keyName"`
}

type RegionConfig struct {
	Url string `json:"url"`
}

type AuthConfig struct {
	AccessToken string `json:"accessToken"`
}

type OrganizationConfig struct {
	Id string `json:"id"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

const projectsQuery = `
  query projects($organizationId: String!) {
    projects(organizationId: $organizationId) {
      id
      name
      platform
      apiKeys {
        name
        apiKey
      }
    }
  }
`

type IntegrationProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type IntegrationData struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type IntegrationEvent struct {
	Event string `json:"event"`
}

type IntegrationSlot struct {
	Slot string `json:"slot"`
}

type Integration struct {
	IntegrationKeyType         string                `json:"integrationKeyType"`
	IntegrationVersion         int8                  `json:"integrationVersion"`
	IntegrationID              string                `json:"integrationId"`
	IntegrationPlatformSupport string                `json:"integrationPlatformSupport"`
	IntegrationKind            string                `json:"integrationKind"`
	IntegrationProperties      []IntegrationProperty `json:"integrationProperties"`
	IntegrationData            []IntegrationData     `json:"integrationData"`
	IntegrationEvents          []IntegrationEvent    `json:"integrationEvents"`
	IntegrationSlots           []IntegrationSlot     `json:"integrationSlots"`
}

type InstalledIntegrationResponse struct {
	IntegrationsInstalled []Integration `json:"integrationsInstalled"`
}

const installedIntegrationsQuery = `
	query integrationsInstalled($organizationId: String!, $projectId: String!, $kind: String!) {
		integrationsInstalled(organizationId: $organizationId, projectId: $projectId, kind: $kind) {
			integrationKeyType
			integrationVersion
			integrationId
			integrationPlatformSupport
			integrationKind
			integrationProperties {
				key
				value
				type
			}
			integrationData {
				key
				type
			}
			integrationEvents {
				event
			}
			integrationSlots {
				slot
			}
		}
	}
`

func ProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	cmd.AddCommand(projectListCmd())
	cmd.AddCommand(projectGetCmd())
	cmd.AddCommand(projectSchemaGenCmd())
	return cmd
}

func projectListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List and select a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			var regionConfig RegionConfig
			if err := fm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
				return fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
			}

			var authConfig AuthConfig
			if err := fm.LoadFromFile(AuthFileName, &authConfig); err != nil {
				return fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
			}

			var orgConfig OrganizationConfig
			if err := fm.LoadFromFile(OrganizationFileName, &orgConfig); err != nil {
				return fmt.Errorf("organization not set. Please select an organization first using 'nativeblocks organization list'")
			}

			client := graphqlutil.NewClient()

			headers := map[string]string{
				"Authorization": "Bearer " + authConfig.AccessToken,
			}

			variables := map[string]interface{}{
				"organizationId": orgConfig.Id,
			}

			resp, err := client.Execute(
				regionConfig.Url,
				headers,
				projectsQuery,
				variables,
			)
			if err != nil {
				return fmt.Errorf("failed to fetch projects: %v", err)
			}

			responseData, err := json.Marshal(resp.Data)
			if err != nil {
				return fmt.Errorf("failed to process response: %v", err)
			}

			var projResp ProjectsResponse
			if err := json.Unmarshal(responseData, &projResp); err != nil {
				fmt.Printf("Debug - Raw response: %s\n", string(responseData))
				return fmt.Errorf("failed to parse projects response: %v", err)
			}

			if len(projResp.Projects) == 0 {
				return fmt.Errorf("no projects found")
			}

			var options []string
			optionMap := make(map[string]ProjectConfig)

			for _, proj := range projResp.Projects {
				apiKeyCount := len(proj.APIKeys)
				optionText := fmt.Sprintf("%s (%s) - %s", proj.Name, proj.Id, proj.Platform)
				options = append(options, optionText)

				projConfig := ProjectConfig{
					Id:   proj.Id,
					Name: proj.Name,
				}
				if apiKeyCount > 0 {
					projConfig.APIKey = proj.APIKeys[0].APIKey
					projConfig.KeyName = proj.APIKeys[0].Name
				}
				optionMap[optionText] = projConfig
			}

			var selection string
			prompt := &survey.Select{
				Message: "Choose a project:",
				Options: options,
			}

			if err := survey.AskOne(prompt, &selection); err != nil {
				return fmt.Errorf("selection cancelled: %v", err)
			}

			selectedProj := optionMap[selection]
			if err := fm.SaveToFile(ProjectFileName, selectedProj); err != nil {
				return fmt.Errorf("failed to save project config: %v", err)
			}

			fmt.Printf("Selected project: %s (%s)\n", selectedProj.Name, selectedProj.Id)
			if selectedProj.APIKey != "" {
				fmt.Printf("API Key '%s' is configured for use\n", selectedProj.KeyName)
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

			var config ProjectConfig
			if err := fm.LoadFromFile(ProjectFileName, &config); err != nil {
				return err
			}

			fmt.Printf("Current project: %s \n", config.Name)
			return nil
		},
	}
}

func projectSchemaGenCmd() *cobra.Command {
	var edition string
	var directory string
	cmd := &cobra.Command{
		Use:   "gen-schema",
		Short: "Generate a JSON schema file",
		RunE: func(cmd *cobra.Command, args []string) error {
			finalDir := directory + "./nativeblocks"
			inputFm, err := fileutil.NewFileManager(&finalDir)
			if err != nil {
				return err
			}
			blockKeyTypes := make([]string, 0)
			blockProperties := make([]MetaItem, 0)
			blockData := make([]MetaItem, 0)

			actionKeyTypes := make([]string, 0)
			actionProperties := make([]MetaItem, 0)
			actionData := make([]MetaItem, 0)

			if edition == "cloud" || edition == "Cloud" || edition == "CLOUD" {
				installedBlocks, err := getInstalledIntegration("BLOCK")
				if err != nil {
					return err
				}

				for _, installedIntegration := range installedBlocks.IntegrationsInstalled {
					blockKeyTypes = append(blockKeyTypes, installedIntegration.IntegrationKeyType)
					for _, property := range installedIntegration.IntegrationProperties {
						meta := MetaItem(property)
						blockProperties = append(blockProperties, meta)
					}
					for _, dataItem := range installedIntegration.IntegrationData {
						meta := MetaItem{Key: dataItem.Key, Value: "", Type: dataItem.Type}
						blockData = append(blockData, meta)
					}
				}

				installedActions, err := getInstalledIntegration("ACTION")
				if err != nil {
					return err
				}

				for _, installedIntegration := range installedActions.IntegrationsInstalled {
					actionKeyTypes = append(actionKeyTypes, installedIntegration.IntegrationKeyType)
					for _, property := range installedIntegration.IntegrationProperties {
						meta := MetaItem(property)
						actionProperties = append(actionProperties, meta)
					}
					for _, dataItem := range installedIntegration.IntegrationData {
						meta := MetaItem{Key: dataItem.Key, Value: "", Type: dataItem.Type}
						actionData = append(actionData, meta)
					}
				}
			} else {
				blockExist := inputFm.FileExists("integrations/block")
				if blockExist {
					blockKeyTypes = FindKeyTypes(inputFm.BaseDir + "/integrations/block")
					blockProperties = FindProperties(inputFm.BaseDir + "/integrations/block")
					blockData = FindData(inputFm.BaseDir + "/integrations/block")
				}

				actionExist := inputFm.FileExists("integrations/action")
				if actionExist {
					actionKeyTypes = FindKeyTypes(inputFm.BaseDir + "/integrations/action")
					actionProperties = FindProperties(inputFm.BaseDir + "/integrations/action")
					actionData = FindData(inputFm.BaseDir + "/integrations/action")
				}

				blockKeyTypes = append(blockKeyTypes, "ROOT")
			}

			schema, err := generateBaseSchema(blockKeyTypes, actionKeyTypes, blockProperties, blockData, actionProperties, actionData)
			if err != nil {
				return nil
			}

			if err := inputFm.SaveToFile("schema.json", schema); err != nil {
				return err
			}
			fmt.Printf("Schema file generated successfully at %s \n", directory)
			return nil
		},
	}
	cmd.Flags().StringVarP(&edition, "edition", "e", "", "Edition type (cloud or community)")
	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Output directory path")
	cmd.MarkFlagRequired("edition")
	cmd.MarkFlagRequired("directory")
	return cmd
}

func getInstalledIntegration(kind string) (*InstalledIntegrationResponse, error) {
	fm, err := fileutil.NewFileManager(nil)
	if err != nil {
		return nil, err
	}

	var regionConfig RegionConfig
	if err := fm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
		return nil, fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
	}

	var authConfig AuthConfig
	if err := fm.LoadFromFile(AuthFileName, &authConfig); err != nil {
		return nil, fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
	}

	var orgConfig OrganizationConfig
	if err := fm.LoadFromFile(OrganizationFileName, &orgConfig); err != nil {
		return nil, fmt.Errorf("organization not set. Please select an organization first using 'nativeblocks organization list'")
	}

	var projConfig ProjectConfig
	if err := fm.LoadFromFile(ProjectFileName, &projConfig); err != nil {
		return nil, fmt.Errorf("project not set. Please select a project first using 'nativeblocks project list'")
	}

	client := graphqlutil.NewClient()

	headers := map[string]string{
		"Authorization": "Bearer " + authConfig.AccessToken,
	}

	variables := map[string]interface{}{
		"organizationId": orgConfig.Id,
		"projectId":      projConfig.Id,
		"kind":           kind,
	}

	resp, err := client.Execute(
		regionConfig.Url,
		headers,
		installedIntegrationsQuery,
		variables,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch installed integrations: %v", err)
	}

	responseData, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to process response: %v", err)
	}

	var installedIntegrationResponse InstalledIntegrationResponse
	if err := json.Unmarshal(responseData, &installedIntegrationResponse); err != nil {
		fmt.Printf("Debug - Raw response: %s\n", string(responseData))
		return nil, fmt.Errorf("failed to parse projects response: %v", err)
	}
	return &installedIntegrationResponse, nil
}
