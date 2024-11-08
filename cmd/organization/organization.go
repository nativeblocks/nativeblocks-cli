package organization

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	auth "github.com/nativeblocks/cli/cmd/auth"
	region "github.com/nativeblocks/cli/cmd/region"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func OrganizationCmd() *cobra.Command {
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
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			var regionModel region.RegionModel
			region, err := regionModel.GetRegion(*fm)
			if err != nil {
				return nil
			}

			var authModel auth.AuthModel
			auth, err := authModel.AuthGet(*fm)
			if err != nil {
				return nil
			}

			orgs, err := GetOrganizations(*fm, region.Url, auth.AccessToken)
			if err != nil {
				return nil
			}

			var options []string
			optionMap := make(map[string]OrganizationModel)

			for _, org := range orgs {
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
			orgModel := OrganizationModel{
				Id:   selectedOrg.Id,
				Name: selectedOrg.Name,
			}
			SelectOrganization(fm, &orgModel)

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
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			var organizationModel OrganizationModel
			organization, err := organizationModel.GetOrganization(*fm)
			if err != nil {
				return nil
			}

			fmt.Printf("Current organization: %s \n", organization.Name)
			return nil
		},
	}
}
