package tasks

import (
	"time"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xfs"
)

func CleanTemporaryFiles() {
	before := time.Now().Add(-1 * ini.GetDuration("app", "tempfileExpires", time.Hour*2))

	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		tfs := tt.FS()
		logger := tt.Logger("XFS")

		sqb := app.SDB().Builder()
		sqb.Eq("tag", models.TagSetFile)
		sqb.Lt("time", before)
		sqb.NotIn("id", append(tt.SVs("sample_files"), tt.SV("sample_file")))
		sql, args := sqb.SQLWhere(), sqb.Params()

		cnt, err := tfs.DeleteWhere(sql, args...)
		if err != nil {
			logger.Debugf("CleanOutdatedSettingFiles('%s')", before.Format(time.RFC3339))
		} else {
			logger.Infof("CleanOutdatedSettingFiles('%s'): %d", before.Format(time.RFC3339), cnt)
		}

		xfs.CleanOutdatedTaggedFiles(tfs, models.TagTmpFile, before, logger)

		return nil
	})
}
