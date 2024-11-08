package auth

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const authLoginMutation = `
  mutation authLogin($email: String!, $password: String!) {
    authLogin(email: $email, password: $password) {
      accessToken
      email
    }
  }
`

func Authenticate(fm fileutil.FileManager, regionUrl, username, password string) (*AuthModel, error) {
	client := graphqlutil.NewClient()

	variables := map[string]interface{}{
		"email":    username,
		"password": password,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		nil,
		authLoginMutation,
		variables,
	)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	var authResponse AuthResponse
	err = graphqlutil.Parse(apiResponse, &authResponse)
	if err != nil {
		return nil, err
	}

	if authResponse.AuthLogin.AccessToken == "" {
		return nil, fmt.Errorf("no access token received in response")
	}

	authConfig := AuthModel{
		AccessToken: authResponse.AuthLogin.AccessToken,
		Email:       authResponse.AuthLogin.Email,
	}

	if err := fm.SaveToFile(AuthFileName, authConfig); err != nil {
		return nil, fmt.Errorf("failed to save auth config: %v", err)
	}

	return &authConfig, nil
}

func (authModel *AuthModel) AuthGet(fm fileutil.FileManager) (*AuthModel, error) {
	var model AuthModel
	if err := fm.LoadFromFile(AuthFileName, &model); err != nil {
		return nil, fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
	}
	return &model, nil
}
