package varutil

import (
	"regexp"
	"strings"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/vad"
	"github.com/askasoft/pangox/xwa/xerrs"
)

var reVarName = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,49}$")

func IsValidVarName(name string) bool {
	return reVarName.Match(str.UnsafeBytes(name))
}

func ValidateVars(fl vad.FieldLevel) error {
	vad.MustStringField("vars", fl)

	ini := ini.NewIni()
	err := ini.LoadData(str.NewReader(fl.Field().String()))
	if err != nil {
		return err
	}

	for _, k := range ini.Section("").Keys() {
		if !IsValidVarName(k) {
			return xerrs.NewLocaleError("setting.variable.error.name", k)
		}
	}

	return nil
}

type Vars map[string]string

func BuildVariables(v string) (Vars, error) {
	i := ini.NewIni()

	err := i.LoadData(str.NewReader(v))
	if err != nil {
		return Vars{}, err
	}

	return i.Section("").StringMap(), nil
}

func BuildVarReplacer(vars Vars) *strings.Replacer {
	var kvs []string
	for k, v := range vars {
		kvs = append(kvs, "{{"+k+"}}", v)
	}
	return str.NewReplacer(kvs...)
}
