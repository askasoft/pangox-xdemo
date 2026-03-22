package middles

import (
	"net/http"
	"time"

	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pango/xin/middleware"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xwa/xargs"
	"github.com/askasoft/pangox/xwa/xerrs"
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
		"Debug":   app.IsDebug(),
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

func Elapsed(c *xin.Context) string {
	return tmu.HumanDuration(time.Since(c.GetTime(middleware.AccessLogStartKey)))
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

func InternalServerRecover(c *xin.Context, r any) {
	if xin.IsBrokenPipeError(r) {
		c.Logger.Warnf("Broken (//%s%s): %v", c.Request.Host, c.Request.URL, r)

		// connection is dead, we can't write a status to it.
		c.Abort()
		return
	}

	c.Logger.Errorf("Panic (//%s%s): %v", c.Request.Host, c.Request.URL, r)

	c.AddError(xerrs.PanicError(r))
	InternalServerError(c)
}

func RedirectToLogin(c *xin.Context) {
	if url := middleware.BuildRedirectURL(c, app.Base()+"/login", middleware.AuthOriginQuery); url != "" {
		c.Redirect(http.StatusTemporaryRedirect, url)
		c.Abort()
		return
	}

	Forbidden(c)
}
