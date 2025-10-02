package util

import (
	"crypto/rsa"
	"errors"

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
