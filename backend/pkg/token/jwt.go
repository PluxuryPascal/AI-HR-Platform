package token

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTtoken struct {
	Issuer   string
	ExpireAt time.Time
	PrvKey   jwk.Key
}

func NewJWTtoken(issuer string, expireAt time.Time, prvKey jwk.Key) *JWTtoken {
	return &JWTtoken{
		Issuer:   issuer,
		ExpireAt: expireAt,
		PrvKey:   prvKey,
	}
}

func (t *JWTtoken) GenerateToken(subject string) ([]byte, error) {
	tokenPrepare := jwt.NewBuilder().
		Issuer(t.Issuer).
		Subject(subject).
		IssuedAt(time.Now()).
		Expiration(t.ExpireAt)

	token, err := tokenPrepare.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.EdDSA, t.PrvKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func (t *JWTtoken) VerifyToken(token []byte, pubKey jwk.Key) (jwt.Token, error) {
	parsed, err := jwt.Parse(token, jwt.WithKey(jwa.EdDSA, pubKey), jwt.WithIssuer(t.Issuer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return parsed, nil
}
