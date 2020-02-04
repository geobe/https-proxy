// Package controller holds all handlers and handler functions
// as well as necessary infrastructure for session management
// and security
package controller

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"github.com/geobe/https-proxy/go/model"
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// the relative location of project files
const Base = "src/github.com/geobe/https-proxy"

// keys for the session store
const S_PROXY = "Proxy-App-Session"

// config parameter values
/** map username -> user struct */
var users map[string]model.User

/** map target key -> ["target link name", "target url"] */
var targets map[string][]string

/** map target key -> reverse proxy to that target host */
var proxies map[string]*httputil.ReverseProxy

/** Server Ports, zu denen  Ports 80 und 443
  vom Internet Router (z.B. FritzBox) mit Port Forwarding weitergeleitet wird */
var httpport string
var tlsport string

/** no unauthorised response for any access to other URL than logInOutPath */
var isStealth bool

/** an arbitrary path to the login/logout page */
var logInOutPath string

/**
set up viper configuration lib and initialize data from configuration file
*/
func SetupConfig(cfgfile string) {
	if cfgfile == "" {
		cfgfile = "config"
	}
	viper.SetConfigName(cfgfile)
	viper.AddConfigPath(Base + "/config") // for config in the working directory
	viper.AddConfigPath("../../config")   // for config in the working directory
	viper.AddConfigPath(".")              // for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	users = getUsers()
	targets = getTargets()
	proxies = getProxies()
	isStealth = IsStealthMode()
	logInOutPath = getLogInOutPath()
	getPorts(&httpport, &tlsport)
	// prepare random number generator
	rand.Seed(time.Now().UnixNano())
	// setup template path
	Templates(Base)
}

/** accessor for the gorilla session store */
func SessionStore() sessions.Store {
	return sessionStore
}

func GetHttpPort() string {
	return httpport
}

func GetTlsPort() string {
	return tlsport
}

func GetUsers() map[string]model.User {
	return users
}

func GetTargets() map[string][]string {
	return targets
}

func GetLogInOutPath() string {
	return logInOutPath
}

func GetProxies() map[string]*httputil.ReverseProxy {
	return proxies
}

func getProxies() map[string]*httputil.ReverseProxy {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	prox := make(map[string]*httputil.ReverseProxy)
	for targetName, linkinfo := range targets {
		target := linkinfo[1]
		targetUrl, err := url.Parse(target)
		if err != nil {
			fmt.Printf("URL parse error %v\n", err)
		} else {
			prox[targetName] = httputil.NewSingleHostReverseProxy(targetUrl)
		}
	}
	return prox
}

func getTargets() map[string][]string {
	targets := viper.GetStringMapStringSlice("targets")
	return targets
}

func IsStealthMode() bool {
	stealth := viper.GetBool("stealthmode")
	return stealth
}

// read users from config file
func getUsers() map[string]model.User {
	umap := make(map[string]model.User)
	uValues := viper.Get("users").([]interface{})
	for _, uValue := range uValues {
		switch userEntry := uValue.(type) {
		case map[string]interface{}:
			accraw := userEntry["access"].([]interface{})
			acclist := make([]string, len(accraw))
			for index, value := range accraw {
				acclist[index] = value.(string)
			}
			nu := model.NewUser(
				userEntry["login"].(string),
				userEntry["password"].(string),
				acclist...)
			umap[nu.Login] = *nu
		default:
			for k1, v1 := range uValue.(map[string]interface{}) {
				fmt.Errorf("Error in configuration file %s: %s\n", k1, v1)
			}
		}
	}
	return umap
}

func getPorts(httpport, tlsport *string) {
	*httpport = viper.GetString("httpport")
	*tlsport = viper.GetString("httpsport")
}

func getLogInOutPath() string {
	lio := viper.GetString("loginout")
	if strings.HasPrefix(lio, "/") {
		return lio
	} else {
		return "/" + lio
	}
}

// keep session store variable private
var sessionStore = makeStore()

// helper function to create a gorilla session store with
// a strong set of keys
func makeStore() sessions.Store {
	store := sessions.NewCookieStore(
		scc.GenerateRandomKey(32),
		scc.GenerateRandomKey(32))
	registerTypes()
	return store
}

// register application types for serialization/deserialization
// necessary for session store
func registerTypes() {
	gob.Register(model.User{})
	gob.Register(Refmap{})
}
