package util

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTStandardClaims struct {
	*jwt.StandardClaims
	Role      string `json:"role,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
}

func ParseRSAPrivateKeyFromPEM(der []byte) (*rsa.PrivateKey, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(der)
	if err != nil {
		return nil, err
	}
	return privateKey, nil

}

func ParseWithClaims(tokenString string, claims jwt.Claims, signKey *rsa.PrivateKey) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return &signKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func EncodeJWTAccessToken(claims jwt.Claims, signKey *rsa.PrivateKey) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS512"))
	t.Claims = claims
	return t.SignedString(signKey)
}

func DecodeJWTAccessToken(tokenString string, claims jwt.Claims, signKey *rsa.PrivateKey) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return &signKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, errors.New("decode token error")
}

func GenerateNewAccessTokenRepo(sessionID, preRegisterUserID string, signKey *rsa.PrivateKey) (string, string, error) {
	issAt := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location())
	expiresAt := issAt.Add(24 * time.Hour)
	claims := JWTStandardClaims{
		StandardClaims: &jwt.StandardClaims{
			Audience:  preRegisterUserID,
			IssuedAt:  issAt.Unix(),
			Subject:   sessionID,
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token, err := EncodeJWTAccessToken(claims, signKey)
	if err != nil {
		return "", "", err
	}
	return token, sessionID, nil
}
