package auth

import (
	"encoding/json"
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
	"github.com/spf13/cobra"
)

const (
	AuthFileName   = "auth"
	RegionFileName = "region"
)

type AuthConfig struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
}

type RegionConfig struct {
	URL string `json:"url"`
}

type AuthResponse struct {
	AuthLogin struct {
		AccessToken string `json:"accessToken"`
		Email       string `json:"email"`
	} `json:"authLogin"`
}

const authLoginMutation = `
  mutation authLogin($email: String!, $password: String!) {
    authLogin(email: $email, password: $password) {
      accessToken
      email
    }
  }
`

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

			var regionConfig RegionConfig
			if err := fm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
				return fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
			}

			client := graphqlutil.NewClient()

			variables := map[string]interface{}{
				"email":    username,
				"password": password,
			}

			resp, err := client.Execute(
				regionConfig.URL,
				nil,
				authLoginMutation,
				variables,
			)
			if err != nil {
				return fmt.Errorf("authentication failed: %v", err)
			}

			responseData, err := json.Marshal(resp.Data)
			if err != nil {
				return fmt.Errorf("failed to process response: %v", err)
			}

			var authResp AuthResponse
			if err := json.Unmarshal(responseData, &authResp); err != nil {
				fmt.Printf("Debug - Raw response: %s\n", string(responseData))
				return fmt.Errorf("failed to parse auth response: %v", err)
			}

			if authResp.AuthLogin.AccessToken == "" {
				return fmt.Errorf("no access token received in response")
			}

			authConfig := AuthConfig{
				AccessToken: authResp.AuthLogin.AccessToken,
				Email:       authResp.AuthLogin.Email,
			}

			if err := fm.SaveToFile(AuthFileName, authConfig); err != nil {
				return fmt.Errorf("failed to save auth config: %v", err)
			}

			if authConfig.Email == "" {
				fmt.Println("Successfully authenticated")
			} else {
				fmt.Printf("Successfully authenticated as %s\n", authConfig.Email)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username/Email")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}
