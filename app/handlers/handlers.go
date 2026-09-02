package handlers

import (
	"net/http"
	"runtime"

	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/jobs"
	"github.com/askasoft/pangox-xdemo/app/middles"
)

func Index(c *xin.Context) {
	c.HTML(http.StatusOK, "index", middles.H(c))
}

func HealthCheck(c *xin.Context) {
	if err := app.SDB().Ping(); err != nil {
		c.Logger.Errorf("Healthcheck: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "OK\n")
}

func ServerStats(c *xin.Context) {
	h := xin.H{}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	h["memory"] = ms.Alloc

	h["dbstats"] = app.SDB().Stats()
	h["jobs"] = jobs.Running()

	c.JSON(http.StatusOK, h)
}
