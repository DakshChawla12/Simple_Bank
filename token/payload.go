package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

var (
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
)

// --- jwt.Claims interface implementation ---

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	// token is valid immediately after issuance
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p Payload) GetIssuer() (string, error) {
	// optional — return empty if not used
	return "", nil
}

func (p Payload) GetSubject() (string, error) {
	// usually user identifier
	return p.Username, nil
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {
	// optional — empty audience
	return jwt.ClaimStrings{}, nil
}

// --- constructor ---

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}

	return payload, nil
}
