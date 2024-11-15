package organization

import (
	"errors"
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

	cmd.AddCommand(organizationSetCmd())
	cmd.AddCommand(organizationGetCmd())
	return cmd
}

func organizationSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Select an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := region.GetRegion(*fm)
			if err != nil {
				return err
			}

			auth, err := auth.AuthGet(*fm)
			if err != nil {
				return err
			}

			orgs, err := GetOrganizations(*fm, region.Url, auth.AccessToken)
			if err != nil {
				return err
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
				return errors.New("selection cancelled: " + err.Error())
			}

			selectedOrg := optionMap[selection]
			SelectOrganization(fm, &selectedOrg)

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

			organization, err := GetOrganization(*fm)
			if err != nil {
				return err
			}

			fmt.Printf("Current organization: %s \n", organization.Name)
			return nil
		},
	}
}
