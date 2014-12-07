// sso
package models

import (
	"time"
)

type SSOToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expiry       int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SSOVerification struct {
	CharacterId        int64     `json:"CharacterID"`
	CharacterName      string    `json:"CharacterName"`
	ExpiresOn          time.Time `json:"ExpiresOn"`
	Scopes             string    `json:"Scopes"`
	TokenType          string    `json:"TokenType"`
	CharacterOwnerHash string    `json:"CharacterOwnerHash"`
}
