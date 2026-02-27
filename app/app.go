package app

import (
	"crypto/tls"
	"time"

	"github.com/askasoft/pango/gwp"
	"github.com/askasoft/pango/ids/snowflake"
	"github.com/askasoft/pango/imc"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pango/vad"
	"github.com/askasoft/pango/xin/middleware"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox/xwa"
	"github.com/askasoft/pangox/xwa/xpwds"
	"github.com/askasoft/pangox/xwa/xsqls"
)

const (
	ExitErrARG int = iota + 10
	ExitErrCFG
	ExitErrDB
	ExitErrFSW
	ExitErrLOG
	ExitErrSCH
	ExitErrSRV
	ExitErrTPL
	ExitErrTXT
	ExitErrXIN
)

const (
	LOGIN_MFA_UNSET  = ""
	LOGIN_MFA_NONE   = "-"
	LOGIN_MFA_EMAIL  = "E"
	LOGIN_MFA_MOBILE = "M"
)

var (
	// VAD global validate
	VAD *vad.Validate

	// XBA global basic auth middleware
	XBA *middleware.BasicAuth

	// XCA global cookie auth middleware
	XCA *middleware.CookieAuth

	// XCN global cookie auth middleware (no failure)
	XCN *middleware.CookieAuth

	// Certificate X509 KeyPair
	SAMLCertificate *tls.Certificate

	// SCMAS schema cache
	SCMAS *imc.Cache[string, bool]

	// TSETS tenant setting cache
	TSETS *imc.Cache[string, map[string]string]

	// WORKS tenant worker pool cache
	WORKS *imc.Cache[string, *gwp.WorkerPool]

	// USERS tenant user cache
	USERS *imc.Cache[string, *models.User]

	// AFIPS authenticate failure ip cache
	AFIPS *imc.Cache[string, int]
)

func Version() string {
	return xwa.Version
}

func Revision() string {
	return xwa.Revision
}

func Versions() string {
	return xwa.Versions()
}

func BuildTime() time.Time {
	return xwa.BuildTime
}

func StartupTime() time.Time {
	return xwa.StartupTime
}

func InstanceID() int64 {
	return xwa.InstanceID
}

func Sequencer() *snowflake.Node {
	return xwa.Sequencer
}

func CFG() map[string]map[string]string {
	return xwa.CFG
}

func Base() string {
	return xwa.Base
}

func Domain() string {
	return xwa.Domain
}

func Secret() string {
	return xwa.Secret
}

func Locales() []string {
	return xwa.Locales
}

func SDB() *sqlx.DB {
	return xsqls.SDB
}

func DBType() string {
	d := xsqls.Driver()
	return str.If(d == "pgx", "postgres", d)
}

func FormatDate(a any) string {
	return tmu.LocalFormatDate(a)
}

func FormatTime(a any) string {
	return tmu.LocalFormatDateTime(a)
}

func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation(time.DateOnly, s, time.Local)
}

func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(time.DateTime, s, time.Local)
}

func RandomPassword() string {
	return xpwds.RandomPassword(64)
}

func MakeFileID(prefix, name string) string {
	return xwa.MakeFileID(prefix, name)
}
