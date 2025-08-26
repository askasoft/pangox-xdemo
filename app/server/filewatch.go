package server

import (
	"github.com/askasoft/pango/fsw"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox/xwa"
	"github.com/askasoft/pangox/xwa/xfsws"
)

func init() {
	xfsws.ReloadLogsOnChange = reloadLogsOnChange
	xfsws.ReloadCfgsOnChange = reloadConfigsOnChange
}

func reloadLogsOnChange(path string, op fsw.Op) {
	if err := xwa.ReloadLogs(op.String()); err != nil {
		log.Error(err)
	}
}

func reloadConfigsOnChange(path string, op fsw.Op) {
	log.Infof("Reloading configurations for '%s' [%v]", path, op)
	reloadConfigs()
}

// initFileWatch initialize file watch
func initFileWatch() {
	if err := xfsws.InitFileWatch(); err != nil {
		log.Fatal(app.ExitErrFSW, err)
	}
}

func reloadFileWatch() {
	if err := xfsws.ReloadFileWatch(); err != nil {
		log.Error(err)
	}
}
