package zebedee

import (
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"
)

const defaultRandomAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandomString generates a cryptographically random string with the specified length.
//
// The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
func RandomString(length int) string {
	return RandomStringWithAlphabet(length, defaultRandomAlphabet)
}

// RandomStringWithAlphabet generates a cryptographically random string
// with the specified length and characters set.
//
// It panics if for some reason rand.Int returns a non-nil error.
func RandomStringWithAlphabet(length int, alphabet string) string {
	b := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))

	for i := range b {
		n, err := cryptoRand.Int(cryptoRand.Reader, max)
		if err != nil {
			panic(err)
		}
		b[i] = alphabet[n.Int64()]
	}

	return string(b)
}

func S256Challenge(code string) (challenge string, verifier string) {
	verifier = base64.URLEncoding.EncodeToString([]byte(code))
	h := sha256.New()
	h.Write([]byte(code))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "="), verifier
}
