package cptutil

import (
	"github.com/askasoft/pangox/xwa/xcpts"
)

var (
	LoginCryptor   *xcpts.XCryptor
	TokenCryptor   *xcpts.XCryptor
	SettingCryptor *xcpts.XCryptor
)

func Init(secret string) error {
	loginCryptor, err := xcpts.NewCryptor(secret, "login")
	if err != nil {
		return err
	}

	tokenCryptor, err := xcpts.NewCryptor(secret, "token")
	if err != nil {
		return err
	}

	settingCryptor, err := xcpts.NewCryptor(secret, "setting")
	if err != nil {
		return err
	}

	LoginCryptor, TokenCryptor, SettingCryptor = loginCryptor, tokenCryptor, settingCryptor
	return nil
}
