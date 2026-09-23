package server

import (
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/utils/oldcpt"
)

func dbFixSettingCrypt(encrypt, exec bool, schemas ...string) error {
	if !exec {
		log.Info("DRY-RUN: no changes will be made (add -exec to apply)")
	}

	return dbIterateSchemas(func(sm schema.Schema) error {
		return app.SDB().Transaction(func(tx *sqlx.Tx) error {
			return fixSettingSecretValues(sm, tx, encrypt, exec)
		})
	}, schemas...)
}

func fixSettingSecretValues(sm schema.Schema, tx sqlx.Sqlx, encrypt, exec bool) error {
	stgs, err := sm.SelectSecretSettings(tx)
	if err != nil {
		return err
	}

	var done, skip int
	for _, stg := range stgs {
		plain := stg.DisplayValue()

		ok, err := fixSettingSecretValue(stg, encrypt)
		if err != nil {
			return err
		}
		if !ok {
			skip++
			continue
		}

		log.Infof("  [%s] %s: %s -> %s", sm, stg.Name, plain, stg.Value)

		if exec {
			if err := updateSettingSecretValue(sm, tx, stg); err != nil {
				return err
			}
		}
		done++
	}

	if done > 0 || skip > 0 {
		log.Infof("[%s]: %d %s, %d skipped", sm, done, map[bool]string{true: "encrypted", false: "decrypted"}[encrypt], skip)
	}

	return nil
}

func fixSettingSecretValue(stg *models.Setting, encrypt bool) (bool, error) {
	if stg.Value == "" {
		return false, nil
	}

	if encrypt {
		if stg.IsSecretEncrypted() {
			return false, nil
		}
		return true, stg.EncryptSecretValue()
	}

	if !stg.IsSecretEncrypted() {
		return false, nil
	}
	return true, stg.DecryptSecretValue()
}

func updateSettingSecretValue(sm schema.Schema, tx sqlx.Sqlx, stg *models.Setting) error {
	sqb := tx.Builder()

	sqb.Update(sm.TableSettings())
	sqb.Setc("value", stg.Value)
	sqb.Eq("name", stg.Name)

	sql, args := sqb.Build()

	_, err := tx.Update(sql, args...)
	return err
}

func dbFixUserPasswords(exec bool, schemas ...string) error {
	if !exec {
		log.Info("DRY-RUN: no changes will be made (add -exec to apply)")
	}

	return dbIterateSchemas(func(sm schema.Schema) error {
		return app.SDB().Transaction(func(tx *sqlx.Tx) error {
			return fixUserPasswords(tx, sm, exec)
		})
	}, schemas...)
}

func fixUserPasswords(tx sqlx.Sqlx, sm schema.Schema, exec bool) error {
	users, err := sm.FindUsers(tx, models.RoleSuper, &args.UserQueryArg{})
	if err != nil {
		return err
	}

	sqbu := tx.Builder()
	sqbu.Update(sm.TableUsers())
	sqbu.Names("password")
	sqbu.Where("id = :id")
	sqlu := sqbu.SQL()

	stmu, err := tx.NamedPrepare(sqlu)
	if err != nil {
		return err
	}
	defer stmu.Close()

	for _, user := range users {
		oldCryptor := oldcpt.NewAes128CBCCryptor(user.Email)
		pwd, err := oldCryptor.DecryptString(user.Password)
		if err != nil {
			return err
		}

		user.SetPassword(pwd)

		log.Warnf("%s [%s] #%d <%s> : (%s -> %s)", str.If(exec, "Fix", "Try"), sm, user.ID, user.Email, pwd, user.Password)

		if exec {
			if _, err := stmu.Exec(user); err != nil {
				return err
			}
		}
	}

	return nil
}
