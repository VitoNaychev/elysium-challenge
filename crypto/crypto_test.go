package crypto_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/crypto"
)

var jwtConfig = crypto.JWTConfig{
	Secret:    []byte("very-secret-key"),
	ExpiresAt: time.Second,
}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestJWTVerification(t *testing.T) {

	t.Run("returns Token on valid JWT ", func(t *testing.T) {
		jwtString, _ := crypto.GenerateJWT(jwtConfig, 0)

		_, err := crypto.VerifyJWT(jwtConfig, jwtString)
		assert.RequireNoError(t, err)
	})

	t.Run("returns error on invalid JWT", func(t *testing.T) {
		jwtString, _ := crypto.GenerateJWT(jwtConfig, 0)

		jwtByteArr := []byte(jwtString)
		if jwtByteArr[10] == 'A' {
			jwtByteArr[10] = 'B'
		} else {
			jwtByteArr[10] = 'A'
		}
		jwtString = string(jwtByteArr)

		_, err := crypto.VerifyJWT(jwtConfig, jwtString)
		if err == nil {
			t.Errorf("did not get error but expected one")
		}
	})
}
