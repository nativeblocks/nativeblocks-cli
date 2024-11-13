package integration

import (
	"fmt"
	"os"

	"github.com/nativeblocks/cli/cmd/auth"
	"github.com/nativeblocks/cli/cmd/organization"
	"github.com/nativeblocks/cli/cmd/region"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	OrganizationFileName = "organization"
	RegionFileName       = "region"
	AuthFileName         = "auth"
)

func IntegrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integration",
		Short: "Manage integrations",
	}

	cmd.AddCommand(integrationListCmd())
	cmd.AddCommand(integrationSyncCmd())

	return cmd
}

func integrationListCmd() *cobra.Command {
	var kind, platform string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get integration list",
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

			organization, err := organization.GetOrganization(*fm)
			if err != nil {
				return err
			}

			integrations, err := GetIntegrations(*fm, region.Url, auth.AccessToken, organization.Id, kind, platform)
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "KeyType", "Version", "Kind", "PlatformSupport"})

			for _, integration := range integrations {
				table.Append([]string{
					integration.Id,
					integration.Name,
					integration.KeyType,
					fmt.Sprintf("%v", integration.Version),
					integration.Kind,
					integration.PlatformSupport,
				})
			}
			table.Render()

			return nil
		},
	}
	cmd.Flags().StringVarP(&kind, "kind", "k", "", "Integration kind")
	cmd.Flags().StringVarP(&platform, "platform", "p", "", "Integration platform")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("platform")
	return cmd
}

func integrationSyncCmd() *cobra.Command {
	var directory string
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Get integration sync",
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

			organization, err := organization.GetOrganization(*fm)
			if err != nil {
				return err
			}

			baseDir := fileutil.GetFileDir(directory)
			fileName := fileutil.GetFileName(directory)

			inputFm, err := fileutil.NewFileManager(&baseDir)
			if err != nil {
				return err
			}

			fileExists := inputFm.FileExists(fileName)
			if !fileExists {
				return fmt.Errorf("could not find the file under: %v", directory)
			}

			var jsonInput IntegrationModel
			err = inputFm.LoadFromFile(fileName, &jsonInput)
			if err != nil {
				return err
			}

			if jsonInput.KeyType == "" {
				return fmt.Errorf("could not find integration keyType")
			}

			err = SyncIntegration(*fm, region.Url, auth.AccessToken, organization.Id, jsonInput)
			if err != nil {
				return err
			}

			fmt.Printf("Integration successfully synced \n")

			return nil
		},
	}
	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Integration working directory")
	cmd.MarkFlagRequired("keyType")
	cmd.MarkFlagRequired("directory")
	return cmd
}
