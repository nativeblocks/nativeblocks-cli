package integrationModule

import (
	"fmt"
	"os"

	"github.com/nativeblocks/cli/cmd/authModule"
	"github.com/nativeblocks/cli/cmd/organizationModule"
	"github.com/nativeblocks/cli/cmd/regionModule"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func IntegrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integration",
		Short: "Manage integrations",
	}

	cmd.AddCommand(integrationListCmd())
	cmd.AddCommand(integrationSyncCmd())
	cmd.AddCommand(integrationGetCmd())

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

			integrations, err := GetIntegrations(region.Url, auth.AccessToken, organization.Id, kind, platform)
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.Header([]string{"Id", "Name", "KeyType", "Version", "Kind", "PlatformSupport"})

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
	_ = cmd.MarkFlagRequired("kind")
	_ = cmd.MarkFlagRequired("platform")
	return cmd
}

func integrationSyncCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync integration",
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

			err = SyncIntegration(region.Url, auth.AccessToken, organization.Id, path)
			if err != nil {
				return err
			}

			fmt.Printf("Integration successfully synced \n")

			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "Integration working path")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func integrationGetCmd() *cobra.Command {
	var id, path string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get integration",
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

			err = GetIntegration(region.Url, auth.AccessToken, organization.Id, path, id)
			if err != nil {
				return err
			}

			fmt.Printf("Integration successfully synced \n")

			return nil
		},
	}
	cmd.Flags().StringVarP(&id, "integrationId", "i", "", "Integration id")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Integration working path")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}
