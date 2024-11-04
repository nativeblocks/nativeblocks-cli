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

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	cmd.AddCommand(projectListCmd())
	cmd.AddCommand(projectGetCmd())
	return cmd
}

func projectListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List and select a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager()
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
			fm, err := fileutil.NewFileManager()
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
