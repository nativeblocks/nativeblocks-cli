package auth

import (
	"errors"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const authCacheFileName = "auth"

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
		return nil, errors.New("authentication failed: " + err.Error())
	}

	var authResponse AuthResponse
	err = graphqlutil.Parse(apiResponse, &authResponse)
	if err != nil {
		return nil, err
	}

	if authResponse.AuthLogin.AccessToken == "" {
		return nil, errors.New("no access token received in response")
	}

	authConfig := AuthModel{
		AccessToken: authResponse.AuthLogin.AccessToken,
		Email:       authResponse.AuthLogin.Email,
	}

	if err := fm.SaveToFile(authCacheFileName, authConfig); err != nil {
		return nil, errors.New("failed to save auth config: " + err.Error())
	}

	return &authConfig, nil
}

func AuthGet(fm fileutil.FileManager) (*AuthModel, error) {
	var model AuthModel
	if err := fm.LoadFromFile(authCacheFileName, &model); err != nil {
		return nil, errors.New("not authenticated. Please login first using 'nativeblocks auth'")
	}
	return &model, nil
}
