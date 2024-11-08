package region

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

const (
	ProjectFileName      = "project"
	OrganizationFileName = "organization"
	RegionFileName       = "region"
	AuthFileName         = "auth"
)

func RegionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "Manage API region",
	}

	cmd.AddCommand(regionSetCmd())
	cmd.AddCommand(regionGetCmd())

	return cmd
}

func regionSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Set API region URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			SetRegion(*fm, args[0])

			fmt.Printf("Region URL set to: %s\n", args[0])
			return nil
		},
	}
}

func regionGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current API region URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := GetRegion(*fm)
			if err != nil {
				return err
			}
			fmt.Printf("Current region URL: %s\n", region.Url)

			return nil
		},
	}
}
