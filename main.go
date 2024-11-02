package main

import (
	"fmt"
	"os"

	region "github.com/nativeblocks/cli/cmd/region"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nativeblocks",
		Short: "NativeBlocks CLI",
	}

	rootCmd.AddCommand(
		region.RegionCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
