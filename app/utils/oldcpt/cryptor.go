package oldcpt

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/str"
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
