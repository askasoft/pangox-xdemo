package args

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox/xwa/xargs"
	"github.com/askasoft/pangox/xwa/xerrs"
)

type IDArg = xargs.IDArg
type PKArg = xargs.PKArg

type Integers = xargs.Integers
type Decimals = xargs.Decimals
type Keywords = xargs.Keywords
type Wildcards = xargs.Wildcards

func ParseIntegers(val string) (Integers, error) {
	return xargs.ParseIntegers(val)
}

func ParseUintegers(val string) (Integers, error) {
	return xargs.ParseUintegers(val)
}

func ParseDecimals(val string) (Decimals, error) {
	return xargs.ParseDecimals(val)
}

func ParseUdecimals(val string) (Decimals, error) {
	return xargs.ParseUdecimals(val)
}

func ParseKeywords(val string) Keywords {
	return xargs.ParseKeywords(val)
}

func NextKeyword(val string) (string, string, bool) {
	return xargs.NextKeyword(val)
}

func ParseWildcards(val string) Wildcards {
	return xargs.ParseWildcards(val)
}

type ParamError = xerrs.ParamError

var (
	ErrInvalidID      = xerrs.ErrInvalidID
	ErrInvalidRequest = xerrs.ErrInvalidRequest
)

// FormatBindErrors translate bind or validate errors and merge it to a new error
func FormatBindErrors(locale string, err error, ns string) error {
	return xerrs.FormatBindErrors(locale, err, ns)
}

// AddBindErrors translate bind or validate errors and add it to context
func AddBindErrors(c *xin.Context, err error, ns string) {
	xerrs.AddBindErrors(c, err, ns)
}

func InvalidIDError(c *xin.Context) error {
	return xerrs.InvalidIDError(c.Locale)
}

func InvalidRequestError(c *xin.Context) error {
	return xerrs.InvalidRequestError(c.Locale)
}

func InvalidFieldError(c *xin.Context, ns, field string) error {
	return xerrs.InvalidParamError(c.Locale, ns, field)
}
