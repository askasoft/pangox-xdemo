package demos

import (
	"errors"
	"net/http"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/utils/tbsutil"
)

type tagsArg struct {
	Text       string   `form:"text"`
	Label      string   `form:"label"`
	Hchecks    []string `form:"hchecks"`
	Vchecks    []string `form:"vchecks"`
	Ochecks    []string `form:"ochecks"`
	Hradios    string   `form:"hradios"`
	Vradios    string   `form:"vradios"`
	Fselect    string   `form:"fselect"`
	Nselect    string   `form:"nselect"`
	Mselect    []string `form:"mselect"`
	Textarea   string   `form:"textarea"`
	Htmledit   string   `form:"htmledit"`
	Singlefile string   `form:"singlefile"`
	Multifiles []string `form:"multifiles"`
}

func TagsIndex(c *xin.Context) {
	h := middles.H(c)

	a := &tagsArg{
		Label:    "label-right",
		Ochecks:  []string{"c2"},
		Hradios:  "r1",
		Vradios:  "r2",
		Htmledit: "<pre>HTML本文</pre>",
	}
	_ = c.Bind(a)

	labels := tbsutil.GetLinkedHashMap(c.Locale, "demos.tags.labels")
	checks := tbsutil.GetLinkedHashMap(c.Locale, "demos.tags.checks")
	radios := tbsutil.GetLinkedHashMap(c.Locale, "demos.tags.radios")
	selects := tbsutil.GetLinkedHashMap(c.Locale, "demos.tags.selects")

	h["LabelsList"] = labels
	h["ChecksList"] = checks
	h["RadiosList"] = radios
	h["SelectList"] = selects
	h["Arg"] = a

	c.AddError(errors.New(str.Repeat("Error message. ", 20)))
	h["Warning"] = str.Repeat("Warning message. ", 20)
	h["Warnings"] = []string{
		"Warning message 1.",
		"Warning message 2.",
	}
	h["Message"] = str.Repeat("Information message. ", 20)
	h["Messages"] = []string{
		"Information message 1.",
		"Information message 2.",
	}
	h["Success"] = str.Repeat("Success message. ", 20)
	h["Successes"] = []string{
		"Success message 1.",
		"Success message 2.",
	}

	c.HTML(http.StatusOK, "demos/tags", h)
}
