package server

import (
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/utils/oldcpt"
)

func dbFixUserPasswords(exec bool, schemas ...string) error {
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

		log.Warnf("Fix [%s] #%d: %s (%s -> %s)", sm, user.ID, user.Email, user.Password, pwd)

		user.SetPassword(pwd)

		if _, err := stmu.Exec(user); err != nil {
			return err
		}
	}

	return nil
}
