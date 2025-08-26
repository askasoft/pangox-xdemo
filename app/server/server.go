package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/askasoft/pango/gwp"
	"github.com/askasoft/pango/imc"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sch"
	"github.com/askasoft/pango/srv"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/jobs"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox/xwa"
	"github.com/askasoft/pangox/xwa/xfsws"
	"github.com/askasoft/pangox/xwa/xhsvs"
)

// SRV service instance
var SRV = &service{}

// service srv.App, srv.Cmd implement
type service struct {
	debug bool
}

// -----------------------------------
// srv.App implement

// Name app/service name
func (s *service) Name() string {
	return "xdemo"
}

// DisplayName app/service display name
func (s *service) DisplayName() string {
	return "Pango Xdemo"
}

// Description app/service description
func (s *service) Description() string {
	return "Pango Xdemo Service"
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
	// Starting the server in a goroutine so that
	// it won't block the graceful shutdown handling below
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
	sch.Stop()

	// close fs watch
	_ = xfsws.CloseFileWatch()

	// gracefully shutdown the http servers with timeout '[server] shutdownTimeout' (defautl 5 seconds).
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

func initCertificate() {
	xcert, err := loadCertificate()
	if err != nil {
		log.Fatal(app.ExitErrCFG, err)
	}

	app.Certificate = xcert
}

func initCaches() {
	app.SCMAS = imc.New[string, bool](ini.GetDuration("cache", "schemaCacheExpires", time.Minute), time.Minute)
	app.CONFS = imc.New[string, map[string]string](ini.GetDuration("cache", "configCacheExpires", time.Minute), time.Minute)
	app.WORKS = imc.New[string, *gwp.WorkerPool](ini.GetDuration("cache", "workerCacheExpires", time.Minute), time.Minute)
	app.USERS = imc.New[string, *models.User](ini.GetDuration("cache", "userCacheExpires", time.Minute), time.Minute)
	app.AFIPS = imc.New[string, int](ini.GetDuration("cache", "afipCacheExpires", time.Minute*30), time.Minute)
}

func reloadCaches() {
	app.SCMAS.SetTTL(ini.GetDuration("cache", "schemaCacheExpires", time.Minute))
	app.CONFS.SetTTL(ini.GetDuration("cache", "configCacheExpires", time.Minute))
	app.WORKS.SetTTL(ini.GetDuration("cache", "workerCacheExpires", time.Minute))
	app.USERS.SetTTL(ini.GetDuration("cache", "userCacheExpires", time.Minute))
	app.AFIPS.SetTTL(ini.GetDuration("cache", "afipCacheExpires", time.Minute*30))
}

func initServers() {
	if err := xhsvs.InitServers(app.XIN, getCertificate); err != nil {
		log.Fatal(app.ExitErrSRV, err)
	}
}

func reloadServers() {
	xhsvs.ConfigServers()
}

func getCertificate(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return app.Certificate, nil
}

// ------------------------------------------------------

func loadCertificate() (*tls.Certificate, error) {
	certificate := ini.GetString("server", "certificate")
	certkeyfile := ini.GetString("server", "certkeyfile")

	xcert, err := tls.LoadX509KeyPair(certificate, certkeyfile)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate (%q, %q): %w", certificate, certkeyfile, err)
	}

	xcert.Leaf, err = x509.ParseCertificate(xcert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("invalid certificate (%q, %q): %w", certificate, certkeyfile, err)
	}

	return &xcert, nil
}

func reloadCertificate() {
	xcert, err := loadCertificate()
	if err != nil {
		log.Error(err)
		return
	}

	app.Certificate = xcert
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

	configMiddleware()

	reloadMessages()

	reloadTemplates()

	runStatsMonitor()

	reloadFileWatch()

	reScheduler()
}
