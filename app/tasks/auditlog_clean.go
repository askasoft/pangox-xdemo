package tasks

import (
	"time"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func CleanOutdatedAuditLogs() {
	before := time.Now().Add(-1 * ini.GetDuration("auditlog", "outdatedBefore", time.Hour*8760))

	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		cnt, err := tt.DeleteAuditLogsBefore(app.SDB(), before)
		if err != nil {
			return err
		}

		tt.Logger("SCH").Infof("[%s] CleanOutdatedAuditLogs(%q): %d", tt.Schema, before.Format(time.RFC3339), cnt)
		return nil
	})
}
