package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/NetworkPy/muserv/muservice/account/models"
	"github.com/NetworkPy/muserv/muservice/account/models/apperrors"
	"github.com/NetworkPy/muserv/muservice/account/security"
	"github.com/google/uuid"
)

// tokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository       models.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository       models.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) models.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenService) NewPairFromUser(ctx context.Context, u *models.User, prevTokenID string) (*models.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := security.GenerateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := security.GenerateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevTokenID)
		}
	}

	return &models.TokenPair{
		IDToken: models.IDToken{
			SS: idToken,
		},
		RefreshToken: models.RefreshToken{
			ID:  refreshToken.ID,
			UID: u.UID,
			SS:  refreshToken.SS,
		},
	}, nil
}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the IDTokenCustomClaims
func (s *tokenService) ValidateIDToken(tokenString string) (*models.User, error) {
	claims, err := security.ValidateIDToken(tokenString, s.PubKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken validates the id token jwt string
// It returns the refreshToken
func (s *tokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	claims, err := security.ValidateRefreshToken(tokenString, s.RefreshSecret)

	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string %s\n%v", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID cloud not be parsed as UUID %s\n%v", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &models.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
