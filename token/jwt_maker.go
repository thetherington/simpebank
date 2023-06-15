package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken creates a new token for specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	// Create claims with multiple fields populated
	claims := JWTCustomClaims{
		payload.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ID:        payload.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", payload, fmt.Errorf("failed to sign token in createJWT: %w", err)
	}

	return signedToken, payload, nil

}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("problem with parsing token: %w", err)
	}

	if claims, ok := jwtToken.Claims.(*JWTCustomClaims); ok && jwtToken.Valid {
		id, err := uuid.Parse(claims.ID)
		if err != nil {
			return nil, err
		}

		return &Payload{
			ID:        id,
			Username:  claims.Username,
			IssuedAt:  claims.IssuedAt.Local(),
			ExpiredAt: claims.ExpiresAt.Local(),
		}, nil
	}

	return nil, fmt.Errorf("problem with token claims: %w", err)
}
