package service

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type authService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func newAuthService(privateKey, publicKey string) *authService {
	privateKeyData, err := os.ReadFile(privateKey)
	if err != nil {
		panic(err)
	}
	publicKeyData, err := os.ReadFile(publicKey)
	if err != nil {
		panic(err)
	}
	private, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		panic(err)
	}

	return &authService{
		privateKey: private,
		publicKey:  public,
	}
}

func (s *authService) ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("incorrect sign method")
		}
		return s.publicKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

func (s *authService) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
