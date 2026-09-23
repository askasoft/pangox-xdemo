package cptutil

import (
	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/str"
)

const prefix = "enc:"

// CPT golbal cryptor
var CPT cpt.Cryptor

func Init(secret string) error {
	cpt, err := cpt.NewAes256GCMCryptor(secret)
	if err != nil {
		return err
	}

	CPT = cpt
	return nil
}

// IsEncrypted は s が暗号化済みの値かどうかを返す。
func IsEncrypted(s string) bool {
	return str.StartsWith(s, prefix)
}

func MustEncrypt(s string) string {
	return gog.Must(Encrypt(s))
}

func MustDecrypt(s string) string {
	return gog.Must(Decrypt(s))
}

func Encrypt(plain string) (string, error) {
	if plain == "" || IsEncrypted(plain) {
		return plain, nil
	}

	enc, err := CPT.EncryptString(plain)
	if err != nil {
		return plain, err
	}

	return prefix + enc, nil
}

func Decrypt(stored string) (string, error) {
	if !IsEncrypted(stored) {
		return stored, nil
	}

	dec, err := CPT.DecryptString(stored[len(prefix):])
	if err != nil {
		return stored, err
	}

	return dec, nil
}
