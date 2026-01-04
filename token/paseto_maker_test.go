package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPasetoMaker(t *testing.T) {
	t.Run("valid symmetric key", func(t *testing.T) {
		key := "abcdefghijklmnopqrstuvwxyz123456" // 32 bytes
		maker, err := NewPasetoMaker(key)

		require.NoError(t, err)
		require.NotNil(t, maker)
	})

	t.Run("invalid symmetric key length", func(t *testing.T) {
		key := "short-key"
		maker, err := NewPasetoMaker(key)

		require.Error(t, err)
		require.Nil(t, maker)
	})
}

func TestCreateAndVerifyPasetoToken(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz123456"
	maker, err := NewPasetoMaker(key)
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

func TestExpiredPasetoToken(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz123456"
	maker, err := NewPasetoMaker(key)
	require.NoError(t, err)

	token, err := maker.CreateToken("testuser", -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	key := "abcdefghijklmnopqrstuvwxyz123456"
	maker, err := NewPasetoMaker(key)
	require.NoError(t, err)

	payload, err := maker.VerifyToken("this.is.not.a.paseto")
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestPasetoTokenWithWrongKey(t *testing.T) {
	key1 := "abcdefghijklmnopqrstuvwxyz123456"
	key2 := "123456abcdefghijklmnopqrstuvwxyz"

	maker1, err := NewPasetoMaker(key1)
	require.NoError(t, err)

	maker2, err := NewPasetoMaker(key2)
	require.NoError(t, err)

	token, err := maker1.CreateToken("testuser", time.Minute)
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
