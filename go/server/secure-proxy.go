package server

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"golang.org/x/crypto/acme/autocert"
	"net/url"
)

// Server Ports, zu denen  Ports 80 und 443
// vom Internet Router (z.B. FritzBox) mit Port Forwarding weitergeleitet wird
const httpport = ":8070"
const tlsport = ":8443"

// redirect to hotPuma's NAT port
const target = "http://192.168.11.103:4567/"

// parsed target url
var targetUrl *url.URL
// create the reverse proxy
var proxy *httputil.ReverseProxy

// the relative location of project files
const Base = "src/github.com/geobe/https-proxy"

// setting up viper configuration lib
func Setup(cfgfile string) {
	if cfgfile == "" {
		cfgfile = "config"
	}
	viper.SetConfigName(cfgfile)
	viper.AddConfigPath(Base + "/config")    // for config in the working directory
	viper.AddConfigPath("../../config")    // for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}


func initVar() {
	targetUrl, _ = url.Parse(target)
	proxy = httputil.NewSingleHostReverseProxy(targetUrl)

}

// Serve a reverse proxy for the given url
func serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	// parse the url
	//url, _ := url.Parse(target)

	// Update the headers to allow for SSL redirection
	//req.URL.Host = targetUrl.Host
	//req.URL.Scheme = targetUrl.Scheme
	//req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	//req.Host = targetUrl.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func main() {
	initVar()

	// die zugelassenen host namen
	allowedHosts := []string{"geobe.spdns.org"}

	// der Verwalter der LetsEncrypt Zertifikate
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(allowedHosts...), //your domain here
		Email:            "geobe.whz@gmail.com",
		Cache:      autocert.DirCache("certs"), //folder for storing certificates
	}

	http.HandleFunc("/", serveReverseProxy)
	http.ListenAndServe("0.0.0.0:80", nil)
}
