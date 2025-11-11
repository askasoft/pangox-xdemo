package schema

import (
	"errors"
	"time"

	"github.com/askasoft/pango/doc/csvx"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/ran"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
)

func ReadSettingsFile() ([]*models.Setting, error) {
	file := app.SettingsCsvFile()

	log.Infof("Read settings file '%s'", file)

	settings := []*models.Setting{}
	if err := csvx.ScanFile(file, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (sm Schema) InitSchema() error {
	log.Infof("Initialize schema %q", sm)

	if err := sm.ExecSchemaSQL(); err != nil {
		return err
	}

	if err := sm.MigrateSuper(); err != nil {
		return err
	}

	settings, err := ReadSettingsFile()
	if err != nil {
		return err
	}

	if err := sm.MigrateSettings(settings); err != nil {
		return err
	}

	return nil
}

func (sm Schema) ExecSchemaSQL() error {
	file := app.SchemaSQLFile()

	log.Infof("Execute Schema SQL file '%s'", file)

	sqls, err := fsu.ReadString(file)
	if err != nil {
		return err
	}

	return sm.ExecSQL(sqls)
}

func (sm Schema) MigrateSettings(settings []*models.Setting) error {
	tb := sm.TableSettings()

	log.Infof("Migrate Settings %q", tb)

	db := app.SDB()

	sqb := db.Builder()
	sqb.Select().From(tb)
	sql, args := sqb.Build()

	osettings := make(map[string]*models.Setting)
	rows, err := db.Queryx(sql, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		var stg models.Setting
		if err := rows.StructScan(&stg); err != nil {
			rows.Close()
			return err
		}
		osettings[stg.Name] = &stg
	}
	rows.Close()

	sqbu := db.Builder()
	sqbu.Update(tb)
	sqbu.Names("style", "order", "required", "secret", "viewer", "editor", "validation")
	sqbu.Where("name = :name")
	sqlu := sqbu.SQL()

	stmu, err := db.PrepareNamed(sqlu)
	if err != nil {
		return err
	}
	defer stmu.Close()

	sqbc := db.Builder()
	sqbc.Insert(tb)
	sqbc.StructNames(&models.Setting{})
	sqlc := sqbc.SQL()
	stmc, err := db.PrepareNamed(sqlc)
	if err != nil {
		return err
	}
	defer stmc.Close()

	for _, stg := range settings {
		if ostg, ok := osettings[stg.Name]; ok {
			if ostg.IsSameMeta(stg) {
				continue
			}

			if _, err := stmu.Exec(stg); err != nil {
				return err
			}
			continue
		}

		stg.CreatedAt = time.Now()
		stg.UpdatedAt = stg.CreatedAt
		if _, err := stmc.Exec(stg); err != nil {
			return err
		}
	}

	return nil
}

func (sm Schema) MigrateSuper() error {
	suc := ini.GetSection("super")
	if suc == nil {
		return errors.New("missing [super] settings")
	}

	emails := str.Fields(suc.GetString("email"))
	if len(emails) == 0 {
		return errors.New("missing [super] email settings")
	}

	tb := sm.TableUsers()

	log.Infof("Migrate Super %q", tb)

	err := app.SDB().Transaction(func(tx *sqlx.Tx) (err error) {
		var uid int64

		sqb := tx.Builder()
		sqb.Select("COALESCE(MAX(id), 0)").From(tb).Lt("id", models.UserStartID)
		sql, args := sqb.Build()
		if err = tx.Get(&uid, sql, args...); err != nil {
			return
		}

		for _, email := range emails {
			sqb.Reset()
			sqb.Select().From(tb).Eq("email", email)
			sql, args = sqb.Build()

			user := &models.User{}
			err = tx.Get(user, sql, args...)
			if err == nil {
				if user.Role != models.RoleSuper || user.Status != models.UserActive {
					sqb.Reset()
					sqb.Update(tb)
					sqb.Setc("role", models.RoleSuper)
					sqb.Setc("status", models.UserActive)
					sqb.Setc("updated_at", time.Now())
					sql, args = sqb.Build()

					_, err = tx.Exec(sql, args...)
					if err != nil {
						return
					}
				}
				continue
			}

			if !errors.Is(err, sqlx.ErrNoRows) {
				return
			}

			uid++
			user.ID = uid
			user.Email = email
			user.Name = str.SubstrBefore(email, "@")
			user.SetPassword(suc.GetString("password", "changeme"))
			user.Role = models.RoleSuper
			user.Status = models.UserActive
			user.Secret = ran.RandInt63()
			user.LoginMFA = suc.GetString("loginmfa")
			user.CIDR = suc.GetString("cidr", "0.0.0.0/0\n::/0")
			user.CreatedAt = time.Now()
			user.UpdatedAt = user.CreatedAt

			sqb.Reset()
			sqb.Insert(tb)
			sqb.StructNames(user)
			sql = sqb.SQL()

			_, err = tx.NamedExec(sql, user)
			if err != nil {
				return err
			}
		}

		return sm.ResetUsersAutoIncrement(tx)
	})

	return err
}
