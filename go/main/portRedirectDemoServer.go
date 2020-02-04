package main

import (
	"crypto/tls"
	"fmt"
	"github.com/geobe/https-proxy/go/controller"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

// Server ports that are redirection targets of our internet router (e.g. a FritzBox)
const demohttpport = ":18080"
const demotlsport = ":18443"

func main() {

	mux := mux.NewRouter()
	mux.HandleFunc("/{target:.*}", demoShowPath)

	// host names we want to allow
	allowedHosts := []string{"iot.geobe.de", "geobe.spdns.org"}

	// manage LetsEncrypt certificates
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(allowedHosts...), //your domain here
		Email:      "georg.beier@fh-zwickau.de",
		Cache:      autocert.DirCache("democerts"), //folder for storing certificates
	}

	// configure secure server
	server := &http.Server{
		Addr: "0.0.0.0" + demotlsport,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: mux,
	}

	// switching redirect handler
	handlerSwitch := &controller.HandlerSwitch{
		Mux:          mux,
		Redirect:     http.HandlerFunc(controller.RedirectHTTP),
		AllowedHosts: allowedHosts,
	}

	// configure local/redirect server
	redirectserver := &http.Server{
		Addr:    "0.0.0.0" + demohttpport,
		Handler: handlerSwitch, //http.HandlerFunc(RedirectHTTP),
	}
	// start redirect server asynchronously on HTTP
	go redirectserver.ListenAndServe()

	// and start primary server on HTTPS
	log.Printf("server starting\n")
	server.ListenAndServeTLS("", "")

}

func demoShowPath(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	target := vars["target"]
	fmt.Fprintf(writer, "path was %v", target)
}
