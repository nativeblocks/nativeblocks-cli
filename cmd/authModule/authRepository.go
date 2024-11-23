package authModule

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"

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

const (
	ProjectFileName      = "project"
	OrganizationFileName = "organization"
)

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

	_ = fm.DeleteFile(OrganizationFileName)
	_ = fm.DeleteFile(ProjectFileName)

	if err := fm.SaveToFile(authCacheFileName, authConfig); err != nil {
		return nil, errors.New("failed to save auth config: " + err.Error())
	}

	return &authConfig, nil
}

func AuthenticateWithToken(fm fileutil.FileManager, accessToken string) (*AuthModel, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		log.Fatal("Invalid JWT token")
	}
	payload := parts[1]
	padding := len(payload) % 4
	if padding != 0 {
		padding = 4 - padding
	}
	payload = payload + strings.Repeat("=", padding)
	decodedPayload, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal("Error decoding payload: ", err)
	}
	var claims map[string]interface{}
	err = json.Unmarshal(decodedPayload, &claims)
	if err != nil {
		log.Fatal("Error unmarshaling payload: ", err)
	}
	eml, _ := claims["eml"].(string)
	authConfig := AuthModel{
		AccessToken: accessToken,
		Email:       eml,
	}

	_ = fm.DeleteFile(OrganizationFileName)
	_ = fm.DeleteFile(ProjectFileName)

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
