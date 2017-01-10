// https://gist.github.com/cannium/c167a19030f2a3c6adbb5a5174bea3ff
package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/poly1305"
	"golang.org/x/crypto/scrypt"
)

var (
	ErrSaltLen = errors.New("Salt must be 16 length")
)

func NewAEAD(password, salt []byte) (xx20 cipher.AEAD, salt0 []byte, err error) {
	if salt == nil {
		salt0 = make([]byte, 16)
		if _, err = io.ReadFull(rand.Reader, salt0); err != nil {
			return
		}
	} else if len(salt) != 16 {
		err = ErrSaltLen
		return
	} else {
		salt0 = salt
	}

	var key []byte
	key, err = scrypt.Key([]byte(password), salt0, 16384, 8, 1, 32)
	if err != nil {
		return
	}

	if err == nil {
		xx20, err = chacha20poly1305.New(key)
	}
	return
}

func EncryptXX20p1305(xx20 cipher.AEAD, plaintext []byte) (ciphertext []byte, err error) {
	ciphertext = make([]byte, len(plaintext)+poly1305.TagSize+chacha20poly1305.NonceSize)
	nonce := ciphertext[:chacha20poly1305.NonceSize]

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	xx20.Seal(ciphertext[chacha20poly1305.NonceSize:chacha20poly1305.NonceSize], nonce, plaintext, nil)
	return
}

func DecryptXX20p1305(xx20 cipher.AEAD, ciphertext []byte) ([]byte, error) {
	return xx20.Open(
		ciphertext[chacha20poly1305.NonceSize:chacha20poly1305.NonceSize],
		ciphertext[:chacha20poly1305.NonceSize],
		ciphertext[chacha20poly1305.NonceSize:],
		nil,
	)
}
