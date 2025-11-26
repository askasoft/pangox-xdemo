package tasks

import (
	"time"

	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func CleanOutdatedData() {
	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		return cleanOutdatedAuditLogs(tt)
	})
}

func cleanOutdatedAuditLogs(tt *tenant.Tenant) error {
	retention := tt.SI("secure_auditlog_retention", 10)
	before := tmu.AddYear(time.Now(), -retention)

	cnt, err := tt.DeleteAuditLogsBefore(app.SDB(), before)
	if err != nil {
		return err
	}
	if cnt > 0 {
		tt.Logger("SCH").Infof("[%s] cleanOutdatedAuditLogs(%q): %d", tt.Schema, before.Format(time.RFC3339), cnt)
	}
	return nil
}
