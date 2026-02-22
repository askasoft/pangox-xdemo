package demos

import (
	"net/http"

	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
)

func UploadsIndex(c *xin.Context) {
	h := middles.H(c)

	c.HTML(http.StatusOK, "demos/uploads", h)
}
