package auth

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
)

const AuthFileName = "auth"

type AuthModel struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
}

type AuthResponse struct {
	AuthLogin struct {
		AccessToken string `json:"accessToken"`
		Email       string `json:"email"`
	} `json:"authLogin"`
}

func (regionModel *AuthModel) Get(fm fileutil.FileManager) (*AuthModel, error) {
	var model AuthModel
	if err := fm.LoadFromFile(AuthFileName, &model); err != nil {
		return nil, fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
	}
	return &model, nil
}
