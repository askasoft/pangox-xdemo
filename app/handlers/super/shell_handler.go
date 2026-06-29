package super

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/utils/tbsutil"
)

func ShellIndex(c *xin.Context) {
	h := middles.H(c)

	h["OS"] = str.Capitalize(runtime.GOOS)
	h["Timeouts"] = tbsutil.GetStrings(c.Locale, "super.shell.timeouts")

	labels := linkedhashmap.NewLinkedHashMap(
		cog.KV("code", tbs.GetText(c.Locale, "super.shell.label.code")),
		cog.KV("time", tbs.GetText(c.Locale, "super.shell.label.time")),
		cog.KV("output", tbs.GetText(c.Locale, "super.shell.label.output")),
	)
	h["Labels"] = labels

	c.HTML(http.StatusOK, "super/shell", h)
}

type ShellArg struct {
	Command string        `form:"command,strip"`
	Timeout time.Duration `form:"timeout"`
}

type ShellResult struct {
	Code   int    `json:"code,omitempty"`
	Time   string `json:"time,omitempty"`
	Output string `json:"output,omitempty"`
}

func ShellExec(c *xin.Context) {
	arg := &ShellArg{}
	_ = c.Bind(arg)

	sr := shellExec(c, arg.Command, arg.Timeout)

	c.JSON(http.StatusOK, sr)
}

func shellExec(c context.Context, command string, timeout time.Duration) (sr ShellResult) {
	if command == "" {
		return
	}

	timeout = min(300*time.Second, max(time.Second, timeout))

	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	start := time.Now()
	sr.Code, sr.Output = shellExecCommand(ctx, command)
	sr.Time = time.Since(start).String()
	return
}
