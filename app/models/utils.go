package models

import (
	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/sqx"
)

const (
	TagSetFile = "s"
	TagTmpFile = "t"
	TagJobFile = "j"
	TagPetFile = "p"
)

type Strings = sqx.JSONStringArray

func toString(o any) string {
	return jsonx.Prettify(o)
}
