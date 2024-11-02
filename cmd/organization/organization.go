package organization

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
	"github.com/spf13/cobra"
)

const (
	OrgFileName    = "organization"
	RegionFileName = "region"
	AuthFileName   = "auth"
)

type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationConfig struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RegionConfig struct {
	URL string `json:"url"`
}

type AuthConfig struct {
	AccessToken string `json:"accessToken"`
}

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
}

const organizationsQuery = `
  query organizations {
    organizations {
      id
      name
    }
  }
`

func NewOrganizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Manage organizations",
	}

	cmd.AddCommand(organizationListCmd())
	cmd.AddCommand(organizationGetCmd())
	return cmd
}

func organizationListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List and select an organization",
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

			client := graphqlutil.NewClient()

			headers := map[string]string{
				"Authorization": "Bearer " + authConfig.AccessToken,
			}

			resp, err := client.Execute(
				regionConfig.URL,
				headers,
				organizationsQuery,
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to fetch organizations: %v", err)
			}

			responseData, err := json.Marshal(resp.Data)
			if err != nil {
				return fmt.Errorf("failed to process response: %v", err)
			}

			var orgResp OrganizationsResponse
			if err := json.Unmarshal(responseData, &orgResp); err != nil {
				fmt.Printf("Debug - Raw response: %s\n", string(responseData))
				return fmt.Errorf("failed to parse organizations response: %v", err)
			}

			if len(orgResp.Organizations) == 0 {
				return fmt.Errorf("no organizations found")
			}

			var options []string
			optionMap := make(map[string]Organization)

			for _, org := range orgResp.Organizations {
				optionText := fmt.Sprintf("%s (%s)", org.Name, org.Id)
				options = append(options, optionText)
				optionMap[optionText] = org
			}

			var selection string
			prompt := &survey.Select{
				Message: "Choose an organization:",
				Options: options,
			}

			if err := survey.AskOne(prompt, &selection); err != nil {
				return fmt.Errorf("selection cancelled: %v", err)
			}

			selectedOrg := optionMap[selection]
			orgConfig := OrganizationConfig{
				Id:   selectedOrg.Id,
				Name: selectedOrg.Name,
			}

			if err := fm.SaveToFile(OrgFileName, orgConfig); err != nil {
				return fmt.Errorf("failed to save organization config: %v", err)
			}

			fmt.Printf("Selected organization: %s (%s)\n", selectedOrg.Name, selectedOrg.Id)
			return nil
		},
	}
}

func organizationGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager()
			if err != nil {
				return err
			}

			var config OrganizationConfig
			if err := fm.LoadFromFile(OrgFileName, &config); err != nil {
				return err
			}

			fmt.Printf("Current organization: %s\n", config)
			return nil
		},
	}
}
