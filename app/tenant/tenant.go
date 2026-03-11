package tenant

import (
	"sync"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/utils/varutil"
	"github.com/askasoft/pangox/xwa/xerrs"
)

type Tenant struct {
	schema.Schema
	settings  map[string]string
	variables map[string]string
}

func (tt *Tenant) Logger(name string) log.Logger {
	logger := log.GetLogger(name)
	logger.SetProp("TENANT", string(tt.Schema))
	return logger
}

func NewTenant(name string) *Tenant {
	tt := &Tenant{Schema: schema.Schema(name)}
	tt.settings = tt.getSettings()

	vars, err := varutil.BuildVariables(tt.SV("environ_variables"))
	if err != nil {
		tt.Logger("SET").Errorf("invalid setting environ_variables: %v", err)
	}
	tt.variables = vars
	return tt
}

func GetSubdomain(c *xin.Context) (string, bool) {
	if !app.IsMultiTenant() {
		return "", true
	}

	domain, hostname := app.Domain(), c.RequestHostname()

	if hostname == domain {
		return "", true
	}

	if str.EndsWith(hostname, domain) {
		sub := hostname[:len(hostname)-len(domain)]
		if str.EndsWithByte(sub, '.') {
			return sub[:len(sub)-1], true
		}
	}

	return "", false
}

const TENANT_CTXKEY = "TENANT"

func FromCtx(c *xin.Context) *Tenant {
	tt, ok := c.Get(TENANT_CTXKEY)
	if !ok {
		panic("Invalid Tenant!")
	}
	return tt.(*Tenant)
}

func FindAndSetTenant(c *xin.Context) (*Tenant, error) {
	if tt, ok := c.Get(TENANT_CTXKEY); ok {
		return tt.(*Tenant), nil
	}

	s, ok := GetSubdomain(c)
	if !ok {
		return nil, xerrs.NewHostnameError(c.Request.Host)
	}

	if s == "" {
		s = app.DefaultSchema()
	}

	if app.IsMultiTenant() {
		ok, err := CheckSchema(s)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, xerrs.NewHostnameError(c.Request.Host)
		}
	}

	tt := NewTenant(s)
	c.Set(TENANT_CTXKEY, tt)
	return tt, nil
}

func Iterate(itf func(tt *Tenant) error) error {
	return schema.Iterate(func(sm schema.Schema) error {
		tt := NewTenant(string(sm))
		return itf(tt)
	})
}

func Create(name string, comment string) error {
	if err := schema.CreateSchema(name, comment); err != nil {
		return err
	}

	if err := schema.Schema(name).InitSchema(); err != nil {
		_ = schema.DeleteSchema(name)
		return err
	}

	return nil
}

// ---------------------------
var muSCMAS sync.Mutex

func CheckSchema(name string) (bool, error) {
	if v, ok := app.SCMAS.Get(name); ok {
		return v, nil
	}

	muSCMAS.Lock()
	defer muSCMAS.Unlock()

	// get again to prevent duplicated load
	if v, ok := app.SCMAS.Get(name); ok {
		return v, nil
	}

	exists, err := schema.ExistsSchema(name)
	if err != nil {
		return false, err
	}

	app.SCMAS.Set(name, exists)
	return exists, nil
}
