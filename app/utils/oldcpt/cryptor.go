package oldcpt

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/str"
)

var (
	ErrCipherBytesTooShort = errors.New("cipher bytes too short")
)

type (
	Cryptor = cpt.Cryptor
	Blocker = cpt.Blocker
	Padding = cpt.Padding
)

type cryptor struct {
	cipher  cipher.Block
	blocker Blocker
	padder  Padding
}

func (c *cryptor) EncryptString(src string) (string, error) {
	bs, err := c.EncryptBytes(str.UnsafeBytes(src))
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bs), nil
}

func (c *cryptor) EncryptBytes(src []byte) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cpt: encrypt panic: %v", r)
		}
	}()

	if c.padder != nil {
		src = c.padder.Pad(src)
	}

	return c.blocker.EncryptBlocks(src)
}

func (c *cryptor) DecryptString(src string) (string, error) {
	bs, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	dst, err := c.DecryptBytes(bs)
	if err != nil {
		return "", err
	}

	return str.UnsafeString(dst), nil
}

func (c *cryptor) DecryptBytes(src []byte) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			dst = nil
			err = fmt.Errorf("cpt: decrypt panic: %v", r)
		}
	}()

	dst, err = c.blocker.DecryptBlocks(src)
	if err != nil {
		return
	}

	if c.padder != nil {
		dst, err = c.padder.Unpad(dst)
	}
	return
}

type aeadBlocker struct {
	aeader cipher.AEAD
}

func (ab aeadBlocker) EncryptBlocks(src []byte) ([]byte, error) {
	nonce := make([]byte, ab.aeader.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	dst := ab.aeader.Seal(nil, nonce, src, nil)
	dst = append(nonce, dst...)
	return dst, nil
}

func (ab aeadBlocker) DecryptBlocks(src []byte) ([]byte, error) {
	nonceSize := ab.aeader.NonceSize()

	if len(src) < nonceSize {
		return nil, ErrCipherBytesTooShort
	}

	nonce, data := src[:nonceSize], src[nonceSize:]

	return ab.aeader.Open(nil, nonce, data, nil)
}

type cbcBlocker struct {
	b  cipher.Block
	iv []byte
}

func (cb cbcBlocker) EncryptBlocks(src []byte) ([]byte, error) {
	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(cb.b, cb.iv).CryptBlocks(dst, src)
	return dst, nil
}

func (cb cbcBlocker) DecryptBlocks(src []byte) ([]byte, error) {
	dst := make([]byte, len(src))
	cipher.NewCBCDecrypter(cb.b, cb.iv).CryptBlocks(dst, src)
	return dst, nil
}
