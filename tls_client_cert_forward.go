package tlsclientcertforward

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//nolint:gochecknoglobals // TODO exchange for traefik log when available
var (
	LoggerWARN  = log.New(ioutil.Discard, "WARN:  tlsclientcertforward: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerINFO  = log.New(ioutil.Discard, "INFO:  tlsclientcertforward: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerDEBUG = log.New(ioutil.Discard, "DEBUG: tlsclientcertforward: ", log.Ldate|log.Ltime|log.Lshortfile)
)

type Config struct {
	LogLevel     string `yaml:"logLevel"`
	SnHeaderName string `yaml:"snHeaderName"`
}

func CreateConfig() *Config {
	return &Config{
		SnHeaderName: "Forwarded-Tls-Client-Cert-Dec-Sn",
		LogLevel:     "info",
	}
}

type TLSClientCertForward struct {
	next   http.Handler
	config *Config
	name   string
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	switch config.LogLevel {
	case "warn":
		LoggerWARN.SetOutput(os.Stdout)
	case "info":
		LoggerWARN.SetOutput(os.Stdout)
		LoggerINFO.SetOutput(os.Stdout)
	case "debug":
		LoggerWARN.SetOutput(os.Stdout)
		LoggerINFO.SetOutput(os.Stdout)
		LoggerDEBUG.SetOutput(os.Stdout)
	default:
		return nil, fmt.Errorf("ERROR: tlsclientcertforward: %s", config.LogLevel)
	}

	return &TLSClientCertForward{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (t *TLSClientCertForward) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	LoggerDEBUG.Printf("ServeHTTP started")

	snHeaderName := t.config.SnHeaderName

	LoggerDEBUG.Printf("ServeHTTP config read")

	if req.TLS != nil && len(req.TLS.PeerCertificates) > 0 {
		peerCert := req.TLS.PeerCertificates[0]
		headerContent := extractSn(peerCert)
		headerEscapedContent := url.QueryEscape(headerContent)
		req.Header.Set(snHeaderName, headerEscapedContent)
		LoggerINFO.Printf("set %s header to %s", snHeaderName, headerEscapedContent)
	} else {
		LoggerWARN.Printf("Tried to extract a certificate on a request without mutual TLS")
	}

	t.next.ServeHTTP(rw, req)
}

func extractSn(peerCert *x509.Certificate) string {
	if peerCert.SerialNumber != nil {
		sn := peerCert.SerialNumber.String()
		LoggerDEBUG.Printf("peer cert sn = %s", sn)

		return sn
	}

	LoggerDEBUG.Printf("peer cert sn is empty")

	return ""
}
