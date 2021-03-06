package security

import (
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	"github.com/NetworkPy/muserv/muservice/account/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// IDTokenCustomClaims holds structure of jwt claims of idToken
type IDTokenCustomClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

// RefreshToken holds the actual signed jwt string along with the ID
// We return the id so it can be used without re-parsing the JWT from signed string
type RefreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

// RefreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// GenerateIDToken generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateIDTokenString, but the signature makes this fairly clear
func GenerateIDToken(u *models.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := IDTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

// GenerateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func GenerateRefreshToken(uid uuid.UUID, key string, exp int64) (*RefreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &RefreshTokenData{
		SS:        ss,
		ID:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

// ValidateIDToken returns the token's claims if the token is valid
func ValidateIDToken(tokenString string, key *rsa.PublicKey) (*IDTokenCustomClaims, error) {
	claims := &IDTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*IDTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

// ValidateRefresh returns the token's claims if the token is valid
func ValidateRefreshToken(tokenString string, key string) (*RefreshTokenCustomClaims, error) {
	claims := &RefreshTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("refresh token is invalid")
	}

	return claims, nil
}
