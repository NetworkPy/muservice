package models

import "context"

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

func (t *TokenPair) NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error) {
	panic("s")
}
