package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	t.Run("valid secret key", func(t *testing.T) {
		maker, err := NewJWTMaker("abcdefghijklmnopqrstuvwxyz123456")
		require.NoError(t, err)
		require.NotNil(t, maker)
	})

	t.Run("short secret key", func(t *testing.T) {
		maker, err := NewJWTMaker("short")
		require.Error(t, err)
		require.Nil(t, maker)
	})
}

func TestCreateAndVerifyToken(t *testing.T) {
	maker, err := NewJWTMaker("abcdefghijklmnopqrstuvwxyz123456")
	require.NoError(t, err)

	username := "testuser"
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker("abcdefghijklmnopqrstuvwxyz123456")
	require.NoError(t, err)

	token, err := maker.CreateToken("testuser", -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, ErrExpiredToken, err)
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewJWTMaker("abcdefghijklmnopqrstuvwxyz123456")
	require.NoError(t, err)

	payload, err := maker.VerifyToken("this.is.not.a.jwt")
	require.Error(t, err)
	require.Equal(t, ErrInvalidToken, err)
	require.Nil(t, payload)
}
