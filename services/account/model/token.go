package model

type TokenPair struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken SRefreshToken `json:"refresh_token"`
}

type SRefreshToken struct {
	TokenString string `json:"token_string"`
	ID          string `json:"-"`
	UUID        string `json:"-"`
}
