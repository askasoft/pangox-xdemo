package tasks

import (
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func VacuumSchemas() {
	if app.DBType() == "postgres" {
		_ = tenant.Iterate(func(tt *tenant.Tenant) error {
			return tt.VacuumSchema(app.SDB())
		})
	}
}
