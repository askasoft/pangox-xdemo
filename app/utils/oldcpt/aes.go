package oldcpt

import (
	"crypto/aes"

	"github.com/askasoft/pango/cpt"
)

var CutPadKey = cpt.CutPadKey

func NewAes128CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 128)
}

func NewAes192CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 192)
}

func NewAes256CBCCryptor(secret string) Cryptor {
	return NewAesCBCCryptor(secret, 256)
}

func NewAesCBCCryptor(secret string, bits int) Cryptor {
	k := []byte(CutPadKey(secret, bits/8))

	c, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}

	return &cryptor{
		cipher:  c,
		blocker: cbcBlocker{c, k[:c.BlockSize()]},
		padder:  cpt.NewPkcs7Padding(c.BlockSize()),
	}
}
