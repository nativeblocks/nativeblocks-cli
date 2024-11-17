package main

import (
	"os"

	"github.com/nativeblocks/cli/cmd/authModule"
	"github.com/nativeblocks/cli/cmd/frameModule"
	"github.com/nativeblocks/cli/cmd/integrationModule"
	"github.com/nativeblocks/cli/cmd/organizationModule"
	"github.com/nativeblocks/cli/cmd/projectModule"
	"github.com/nativeblocks/cli/cmd/regionModule"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nativeblocks",
		Short: "Nativeblocks cli",
	}

	rootCmd.AddCommand(
		regionModule.RegionCmd(),
		authModule.AuthCmd(),
		organizationModule.OrganizationCmd(),
		projectModule.ProjectCmd(),
		frameModule.FrameCmd(),
		integrationModule.IntegrationCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
