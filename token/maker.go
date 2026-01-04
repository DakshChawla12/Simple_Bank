package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a valid for a specific duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks the validity of the token
	VerifyToken(token string) (*Payload, error)
}
