package server

import (
	"crypto/tls"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox/xwa/xcert"
)

func getCertificate(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return xcert.GetCertificate(chi)
}

func initCertificate() {
	xcert.InitCertificate()

	initSAMLCertificate()
}

func reloadCertificate() {
	xcert.ReloadCertificate()

	reloadSAMLCertificate()
}

func initSAMLCertificate() {
	certificate := ini.GetString("saml", "certificate")
	certkeyfile := ini.GetString("saml", "certkeyfile")

	cert, err := xcert.LoadCertificate(certificate, certkeyfile)
	if err != nil {
		log.Fatal(app.ExitErrCFG, err)
	}
	app.SAMLCertificate = cert
}

func reloadSAMLCertificate() {
	certificate := ini.GetString("saml", "certificate")
	certkeyfile := ini.GetString("saml", "certkeyfile")

	cert, err := xcert.LoadCertificate(certificate, certkeyfile)
	if err != nil {
		log.Error(err)
		return
	}
	app.SAMLCertificate = cert
}
