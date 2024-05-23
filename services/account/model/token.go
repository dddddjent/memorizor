package model

type TokenPair struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}
