package model

import "github.com/gofrs/uuid"

type TokenPair struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken SRefreshToken `json:"refresh_token"`
}

type SRefreshToken struct {
	TokenString string    `json:"token_string"`
	ID          uuid.UUID `json:"-"`
	UUID        uuid.UUID `json:"-"`
}
