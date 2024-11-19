package jwt

import (
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWT_Success(t *testing.T) {
	key := entities.Secret("test-secret-key")
	dur := 5 * time.Minute
	login := "test-login"

	tokenString, err := CreateJWT(key, login, dur)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestCreateJWT_TokenContent(t *testing.T) {
	key := entities.Secret("test-secret-key")
	login := "test-login"
	dur := 5 * time.Minute

	tokenString, err := CreateJWT(key, login, dur)
	assert.NoError(t, err)

	parsedToken, _, err := ParseJWT(key, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, login, parsedToken)
}

func TestParseJWT_Success(t *testing.T) {
	key := entities.Secret("test-secret-key")
	login := "test-login"
	dur := 5 * time.Minute

	tokenString, err := CreateJWT(key, login, dur)
	assert.NoError(t, err)

	parsedLogin, expiresAt, err := ParseJWT(key, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, login, parsedLogin)
	assert.WithinDuration(t, time.Now().Add(dur), expiresAt, time.Second*2)
}

func TestParseJWT_InvalidSignature(t *testing.T) {
	key := entities.Secret("test-secret-key")
	wrongKey := entities.Secret("wrong-secret-key")
	login := "test-login"
	dur := 5 * time.Minute

	tokenString, err := CreateJWT(key, login, dur)
	assert.NoError(t, err)

	_, _, err = ParseJWT(wrongKey, tokenString)
	assert.Error(t, err)
}

func TestParseJWT_InvalidToken(t *testing.T) {
	key := entities.Secret("test-secret-key")

	_, _, err := ParseJWT(key, "invalid-token")
	assert.Error(t, err)
}

func TestParseJWT_ExpiredToken(t *testing.T) {
	key := entities.Secret("test-secret-key")
	login := "test-login"

	tokenString, err := CreateJWT(key, login, time.Millisecond*100)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 200)

	_, _, err = ParseJWT(key, tokenString)
	assert.Error(t, err)
}
