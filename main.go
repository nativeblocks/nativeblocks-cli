package main

import (
	"os"

	"github.com/nativeblocks/cli/cmd/auth"
	"github.com/nativeblocks/cli/cmd/frame"
	"github.com/nativeblocks/cli/cmd/integration"
	"github.com/nativeblocks/cli/cmd/organization"
	"github.com/nativeblocks/cli/cmd/project"
	"github.com/nativeblocks/cli/cmd/region"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nativeblocks",
		Short: "Nativeblocks cli",
	}

	rootCmd.AddCommand(
		region.RegionCmd(),
		auth.AuthCmd(),
		organization.OrganizationCmd(),
		project.ProjectCmd(),
		frame.FrameCmd(),
		integration.IntegrationCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
