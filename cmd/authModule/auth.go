package authModule

import (
	"fmt"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func AuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
	}

	cmd.AddCommand(authTokenCmd())

	return cmd
}

func authTokenCmd() *cobra.Command {
	var accessToken string

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Authenticate with token",
		RunE: func(cmd *cobra.Command, args []string) error {
			fm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			authModel, err := AuthenticateWithToken(*fm, accessToken)
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

	cmd.Flags().StringVarP(&accessToken, "accessToken", "a", "", "access token")
	_ = cmd.MarkFlagRequired("accessToken")

	return cmd
}
