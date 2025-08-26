package server

import (
	"io/fs"

	"github.com/askasoft/pango/fsw"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/txts"
	"github.com/askasoft/pangox/xwa/xfsws"
	"github.com/askasoft/pangox/xwa/xtxts"
)

func init() {
	xtxts.FSs = []fs.FS{txts.FS}
	xfsws.ReloadMsgsOnChange = reloadMessagesOnChange
}

func initMessages() {
	if err := xtxts.InitMessages(); err != nil {
		log.Fatal(app.ExitErrTXT, err)
	}
}

func reloadMessages() {
	if err := xtxts.ReloadMessages(); err != nil {
		log.Error(err)
	}
}

func reloadMessagesOnChange(path string, op fsw.Op) {
	if err := xtxts.ReloadMessagesOnChange(path, op.String()); err != nil {
		log.Error(err)
	}
}
