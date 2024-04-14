package crypto_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/crypto"
	"github.com/golang-jwt/jwt/v5"
)

var jwtConfig = crypto.JWTConfig{
	Secret:    []byte("very-secret-key"),
	ExpiresAt: time.Second,
}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestJWTVerification(t *testing.T) {

	t.Run("returns subject ID on valid JWT ", func(t *testing.T) {
		wantSubjectID := 10
		jwtString, _ := crypto.GenerateJWT(jwtConfig, wantSubjectID)

		gotSubjectID, err := crypto.VerifyJWT(jwtConfig, jwtString)
		assert.RequireNoError(t, err)

		assert.Equal(t, gotSubjectID, wantSubjectID)
	})

	t.Run("returns ErrMissingSubject on missing subject ", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.ExpiresAt)),
		})

		jwt, err := token.SignedString(jwtConfig.Secret)
		assert.RequireNoError(t, err)

		_, err = crypto.VerifyJWT(jwtConfig, jwt)
		assert.Equal(t, err, (error)(crypto.ErrMissingSubject))
	})

	t.Run("returns ErrNonintegerSubject on noninteger subject ", func(t *testing.T) {
		subject := "nonintegerSubject"

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.ExpiresAt)),
		})

		jwt, err := token.SignedString(jwtConfig.Secret)
		assert.RequireNoError(t, err)

		_, err = crypto.VerifyJWT(jwtConfig, jwt)
		assert.Equal(t, err, (error)(crypto.ErrNonintegerSubject))
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
