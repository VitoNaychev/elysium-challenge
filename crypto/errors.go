package crypto

type CryptoError struct {
	msg string
}

func NewCryptoError(msg string) *CryptoError {
	return &CryptoError{
		msg: msg,
	}
}

func (c *CryptoError) Error() string {
	return c.msg
}

var (
	ErrMissingSubject    = NewCryptoError("missing subject in JWT")
	ErrNonintegerSubject = NewCryptoError("cannot convert subject ID to integer")
)
