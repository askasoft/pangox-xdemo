package settings

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/handlers/files"
	"github.com/askasoft/pangox-xdemo/app/models"
)

func SettingFileUpload(c *xin.Context) {
	files.UploadFile(c, models.TagSetFile)
}

func SettingFilePreview(c *xin.Context) {
	fid := c.Param("id")
	dnloadURL := app.Base() + "/a/settings/files/dnload" + fid

	files.PreviewFile(c, fid, dnloadURL)
}
