package server

import (
	"fmt"
	"os"
	"time"

	"github.com/askasoft/pango/gwp"
	"github.com/askasoft/pango/imc"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/srv"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/jobs"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox/xwa"
	"github.com/askasoft/pangox/xwa/xhsvs"
	"github.com/askasoft/pangox/xwa/xxins"
)

// SRV service instance
var SRV = &service{}

// service srv.App, srv.Cmd implement
type service struct {
	outdir string
	debug  bool
}

// -----------------------------------
// srv.App implement

// Name app/service name
func (s *service) Name() string {
	return "xdemo"
}

// DisplayName app/service display name
func (s *service) DisplayName() string {
	return "Pangox Xdemo"
}

// Description app/service description
func (s *service) Description() string {
	return "Pangox Xdemo Service"
}

// Version app version
func (s *service) Version() string {
	return app.Versions()
}

// Init initialize the app
func (s *service) Init() {
	Init()
}

// Relead reload the app
func (s *service) Reload() {
	Reload()
}

// Run run the app
func (s *service) Run() {
	Run()
}

// Shutdown shutdown the app
func (s *service) Shutdown() {
	Shutdown()
}

// Wait wait signal for reload or shutdown the app
func (s *service) Wait() {
	srv.Wait(s)
}

// ------------------------------------------------------

// Init initialize the app
func Init() {
	initLogs()

	initConfigs()

	initCertificate()

	initCaches()

	initMessages()

	initTemplates()

	initDatabase()

	initRouter()

	initServers()

	initFileWatch()

	initStatsMonitor()

	initScheduler()
}

// Relead reload the app
func Reload() {
	log.Info("Reloading ...")

	if err := xwa.ReloadLogs("RELOAD"); err != nil {
		log.Error(err)
	}

	reloadConfigs()
}

// Run start the http server
func Run() {
	// start serve http servers in go-routines (non-blocking)
	xhsvs.Serves()

	// Start jobs (Resume interrupted jobs)
	if ini.GetBool("job", "startAtStartup") {
		go jobs.Starts()
	}
}

// Shutdown shutdown the app
func Shutdown() {
	log.Info("Shutting down ...")

	// stop scheduler
	stopScheduler()

	// close file watch
	closeFileWatch()

	// gracefully shutdown the http servers with timeout '[server] shutdownTimeout'.
	xhsvs.Shutdowns()

	log.Info("EXIT.")

	// close log
	log.Close()
}

// ------------------------------------------------------

func initLogs() {
	if err := xwa.InitLogs(); err != nil {
		fmt.Println(err)
		os.Exit(app.ExitErrLOG)
	}
}

func initConfigs() {
	if err := xwa.InitConfigs(); err != nil {
		log.Fatal(app.ExitErrCFG, err)
	}
}

func initCaches() {
	app.SCMAS = imc.New[string, bool](ini.GetDuration("cache", "schemaCacheExpires", time.Minute), time.Minute)
	app.TSETS = imc.New[string, map[string]string](ini.GetDuration("cache", "settingsCacheExpires", time.Minute), time.Minute)
	app.WORKS = imc.New[string, *gwp.WorkerPool](ini.GetDuration("cache", "workerCacheExpires", time.Minute), time.Minute)
	app.USERS = imc.New[string, *models.User](ini.GetDuration("cache", "userCacheExpires", time.Minute), time.Minute)
	app.AFIPS = imc.New[string, int](ini.GetDuration("cache", "afipCacheExpires", time.Minute*30), time.Minute)
}

func reloadCaches() {
	app.SCMAS.SetTTL(ini.GetDuration("cache", "schemaCacheExpires", time.Minute))
	app.TSETS.SetTTL(ini.GetDuration("cache", "settingsCacheExpires", time.Minute))
	app.WORKS.SetTTL(ini.GetDuration("cache", "workerCacheExpires", time.Minute))
	app.USERS.SetTTL(ini.GetDuration("cache", "userCacheExpires", time.Minute))
	app.AFIPS.SetTTL(ini.GetDuration("cache", "afipCacheExpires", time.Minute*30))
}

func initServers() {
	if err := xhsvs.InitServers(xxins.XIN, getCertificate); err != nil {
		log.Fatal(app.ExitErrSRV, err)
	}
}

func reloadServers() {
	if err := xhsvs.ReloadServers(); err != nil {
		log.Error(err)
	}
}

func reloadConfigs() {
	if err := xwa.InitConfigs(); err != nil {
		log.Error(err)
		return
	}

	reloadCertificate()

	reloadCaches()

	reloadServers()

	reloadDatabase()

	configRouter()

	configMiddleware()

	reloadMessages()

	reloadTemplates()

	runStatsMonitor()

	reloadFileWatch()

	reScheduler()
}
