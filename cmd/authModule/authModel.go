package authModule

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
