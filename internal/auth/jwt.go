package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Scope string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}

type Verifier struct {
	secret []byte
	issuer string
}

func NewJwtVerifier(secret, issuer string) (*Verifier, error) {

	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 bytes")
	}

	if issuer == "" {
		return nil, errors.New("JWT issuer is required")
	}

	return &Verifier{
		secret: []byte(secret),
		issuer: issuer,
	}, nil
}

func (vrf *Verifier) Verify(raw string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(tkn *jwt.Token) (any, error) {
			if tkn.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid token")
			}
			return vrf.secret, nil
		},

		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(vrf.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("invalid token")

	}

	return claims, nil
}
