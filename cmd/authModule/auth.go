package authModule

import (
	"fmt"
	"github.com/nativeblocks/cli/cmd/regionModule"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func AuthCmd() *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with username and password",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := regionModule.GetRegion(*fm)
			if err != nil {
				return err
			}

			authModel, err := Authenticate(*fm, region.Url, username, password)
			if err != nil {
				return err
			}

			if authModel.Email == "" {
				fmt.Println("Successfully authenticated")
			} else {
				fmt.Printf("Successfully authenticated as %s\n", authModel.Email)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username/Email")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
