package middles

import (
	"net/http"
	"time"

	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xwa/xargs"
)

func H(c *xin.Context) xin.H {
	tt := tenant.FromCtx(c)
	au := tenant.GetAuthUser(c)

	h := xin.H{
		"TT":      tt,
		"AU":      au,
		"CFG":     app.CFG(),
		"VER":     app.Version(),
		"REV":     app.Revision(),
		"Base":    app.Base(),
		"Domain":  app.Domain(),
		"Locales": app.Locales(),
		"Now":     time.Now(),
		"Ctx":     c,
		"Loc":     c.Locale,
		"Host":    c.Request.Host,
		"Token":   RefreshToken(c),
	}
	return h
}

func E(c *xin.Context) xin.H {
	return xargs.E(c)
}

func NotFound(c *xin.Context) {
	if xin.IsAjax(c) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusNotFound, "404", H(c))
	c.Abort()
}

func Forbidden(c *xin.Context) {
	if xin.IsAjax(c) {
		c.JSON(http.StatusForbidden, E(c))
	} else {
		c.HTML(http.StatusForbidden, "403", H(c))
	}
	c.Abort()
}

func InternalServerError(c *xin.Context) {
	if xin.IsAjax(c) {
		c.JSON(http.StatusInternalServerError, E(c))
	} else {
		c.HTML(http.StatusInternalServerError, "500", H(c))
	}
	c.Abort()
}
