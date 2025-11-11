package super

import (
	"net/http"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
)

func ConfigsIndex(c *xin.Context) {
	h := middles.H(c)

	h["Sections"] = ini.Sections()

	c.HTML(http.StatusOK, "super/configs", h)
}
