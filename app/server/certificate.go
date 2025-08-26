package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app"
)

func getCertificate(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return app.Certificate, nil
}

func initCertificate() {
	xcert, err := loadCertificate()
	if err != nil {
		log.Fatal(app.ExitErrCFG, err)
	}

	app.Certificate = xcert
}

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
