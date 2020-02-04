package controller

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// RedirectHTTP is an HTTP handler (suitable for use with http.HandleFunc)
// that responds to all requests by redirecting to the same URL served over HTTPS.
// It should only be invoked for requests received over HTTP.
func RedirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil || r.Host == "" {
		http.Error(w, "", http.StatusNotFound)
	}

	var u *url.URL
	u = r.URL
	host := r.Host
	u.Host = strings.Split(host, ":")[0]
	u.Scheme = "https"
	log.Printf("redirect to u.host  %s -> %s\n", r.Host, u.String())
	http.Redirect(w, r, u.String(), 302)
}

//struct to hold different route multiplexer. Internal calls from the same local net will be directly routed by Mux,
//else redirect to Redirect that shall redirect to HTTPS
type HandlerSwitch struct {
	Mux          http.Handler
	Redirect     http.Handler
	AllowedHosts []string
}

// not correct for 172.16.0.0/12
var matcher = regexp.MustCompile("(192\\.168.*)|(localhost)|(10\\..*)|(172\\..*)")

// Handler function that redirects access from external internet to redirect handler.
// Local is directly routed to  MUX.
func (h *HandlerSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	local := matcher.MatchString(host)
	if local {
		h.Mux.ServeHTTP(w, r)
	} else {
		var serve bool
		for _, v := range h.AllowedHosts {
			if strings.Contains(host, v) {
				serve = true
				break
			}
		}
		if serve {
			h.Redirect.ServeHTTP(w, r)
		} else {
			http.Error(w, "", http.StatusGone)
		}
	}
}
