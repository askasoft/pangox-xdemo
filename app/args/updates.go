package args

import (
	"time"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox/xwa/xargs"
)

type UpdatedAtArg struct {
	UpdatedAt *time.Time `json:"updated_at,omitempty" form:"-"`
}

func (uaa *UpdatedAtArg) SetUpdatedAt(t time.Time) {
	uaa.UpdatedAt = &t
}

func (uaa *UpdatedAtArg) ToString(a any) string {
	u := uaa.UpdatedAt
	uaa.UpdatedAt = nil
	s := jsonx.Stringify(a)
	uaa.UpdatedAt = u
	return s
}

type UserUpdatesArg struct {
	IDArg
	UpdatedAtArg

	Role     *string `json:"role,omitempty" form:"role,strip"`
	Status   *string `json:"status,omitempty" form:"status,strip"`
	LoginMFA *string `json:"login_mfa,omitempty" form:"login_mfa,strip"`
	CIDR     *string `json:"cidr,omitempty" form:"cidr,strip" validate:"omitempty,cidrs"`
}

func (uua *UserUpdatesArg) String() string {
	return uua.ToString(uua)
}

func (uua *UserUpdatesArg) isEmpty() bool {
	return uua.Role == nil && uua.Status == nil && uua.LoginMFA == nil && uua.CIDR == nil
}

func (uua *UserUpdatesArg) Bind(c *xin.Context) error {
	if err := c.Bind(uua); err != nil {
		return err
	}
	if err := uua.ParseID(); err != nil {
		return err
	}
	if uua.isEmpty() {
		return xargs.ErrInvalidUpdates
	}
	return nil
}

func (uua *UserUpdatesArg) AddUpdates(sqb *sqlx.Builder) {
	if uua.Role != nil {
		sqb.Setc("role", *uua.Role)
	}
	if uua.Status != nil {
		sqb.Setc("status", *uua.Status)
	}
	if uua.LoginMFA != nil {
		sqb.Setc("login_mfa", *uua.LoginMFA)
	}
	if uua.CIDR != nil {
		sqb.Setc("cidr", *uua.CIDR)
	}

	uua.SetUpdatedAt(time.Now())
	sqb.Setc("updated_at", *uua.UpdatedAt)
}

type FileUpdatesArg struct {
	PKArg
	UpdatedAtArg

	Tag *string `json:"tag,omitempty" form:"tag,strip"`
}

func (fua *FileUpdatesArg) String() string {
	return fua.ToString(fua)
}

func (fua *FileUpdatesArg) isEmpty() bool {
	return fua.Tag == nil
}

func (fua *FileUpdatesArg) Bind(c *xin.Context) error {
	if err := c.Bind(fua); err != nil {
		return err
	}
	if err := fua.ParseID(); err != nil {
		return err
	}
	if fua.isEmpty() {
		return xargs.ErrInvalidUpdates
	}
	return nil
}

func (fua *FileUpdatesArg) AddUpdates(sqb *sqlx.Builder) {
	if fua.Tag != nil {
		sqb.Setc("tag", *fua.Tag)
	}
}

type PetUpdatesArg struct {
	IDArg
	UpdatedAtArg

	Gender *string    `json:"gender,omitempty" form:"gender,strip"`
	BornAt *time.Time `json:"born_at,omitempty" form:"born_at"`
	Origin *string    `json:"origin,omitempty" form:"origin,strip"`
	Temper *string    `json:"temper,omitempty" form:"temper,strip"`
	Habits *[]string  `json:"habits,omitempty" form:"habits,strip"`
}

func (pua *PetUpdatesArg) String() string {
	return pua.ToString(pua)
}

func (pua *PetUpdatesArg) isEmpty() bool {
	return pua.Gender == nil && pua.BornAt == nil && pua.Origin == nil && pua.Temper == nil && pua.Habits == nil
}

func (pua *PetUpdatesArg) Bind(c *xin.Context) error {
	if err := c.Bind(pua); err != nil {
		return err
	}
	if err := pua.ParseID(); err != nil {
		return err
	}
	if pua.isEmpty() {
		return xargs.ErrInvalidUpdates
	}
	return nil
}

func (pua *PetUpdatesArg) AddUpdates(sqb *sqlx.Builder) {
	if pua.Gender != nil {
		sqb.Setc("gender", *pua.Gender)
	}
	if pua.BornAt != nil {
		sqb.Setc("born_at", *pua.BornAt)
	}
	if pua.Origin != nil {
		sqb.Setc("origin", *pua.Origin)
	}
	if pua.Temper != nil {
		sqb.Setc("temper", *pua.Temper)
	}
	if pua.Habits != nil {
		sqb.Setc("habits", sqx.JSONStringArray(*pua.Habits))
	}

	pua.SetUpdatedAt(time.Now())
	sqb.Setc("updated_at", *pua.UpdatedAt)
}
