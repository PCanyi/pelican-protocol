package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"math/big"
	"os"
)

type KeyPayload struct {
	FirstName      string
	MiddleInitial  string
	LastName       string
	ValidEmail     string
	ValidSMSNumber string
	PublicKey      string

	// AcctUsername is the
	// sha1 hmac of PublicKey, prefixed with "p" (for pelican) and
	// encoded in base 36; A regex for it would be: "p[0-9a-z]{31}"
	// This is to conform to the requirements of linux usernames.
	// See man page for useradd; there is a 32 character limit,
	// and usernames must start with a letter and then contain
	// be only lowercase letters and digits. Underscores and
	// dashes are allowed too but we don't use them.
	AcctUsername string
}

func Sha256HMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// CheckMAC returns true if messageMAC is a valid HMAC tag for message.
func CheckSha256HMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func Sha1HMAC(message, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// CheckMAC returns true if messageMAC is a valid HMAC tag for message.
func CheckSha1HMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func SendKeyPayload(p *KeyPayload) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write to buf

	// Encode (send) some values.
	err := enc.Encode(p)
	if err != nil {
		panic(fmt.Sprintf("encode error:", err))
	}

	toWrite := len(buf.Bytes())

	n, err := io.Copy(os.Stdout, &buf)
	if err != nil {
		panic(err)
	}
	if n != int64(toWrite) {
		panic(fmt.Sprintf("did not write all of buf; n=%v\n, len(buf)=%v\n", n, toWrite))
	}
}

func main() {
	gob.Register(KeyPayload{})

	key := &KeyPayload{PublicKey: "0123456789abcdef-hello-public-key"}

	hmac := Sha1HMAC([]byte(key.PublicKey), []byte{})

	key.AcctUsername = encodeSha1AsUsername(hmac)

	SendKeyPayload(key)
	fmt.Printf("\n done sending: '%#v'.\n", key)
}

func encodeSha1AsUsername(sha1 []byte) string {
	i := new(big.Int)
	i.SetBytes(sha1)
	return "p" + bigIntToBase36string(i)
}

var enc36 string = "0123456789abcdefghijklmnopqrstuvwxyz"
var e36 []rune = []rune(enc36)

// i must be between 0 and 35 inclusive.
func encode36(i int64) rune {
	return e36[i]
}

func bigIntToBase36string(val *big.Int) string {
	const N = 31 // ceiling(log(2^160,36))
	res := make([]rune, N)
	left := new(big.Int)
	quo := new(big.Int)
	rem := new(big.Int)
	*left = *val

	div := big.NewInt(36)

	for i := 0; i < N; i++ {
		quo.QuoRem(left, div, rem)
		*left = *quo
		r := rem.Int64()
		e := encode36(r)
		res[N-1-i] = e
	}

	return string(res)
}
