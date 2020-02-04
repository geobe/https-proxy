package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	scc "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const resourcedir = Base + "/web/"

func InitRouter() *mux.Router {
	fmt.Printf("initializing router \n")

	router := mux.NewRouter()

	// finde Working directory = GOPATH
	//docbase, _ := os.Getwd()
	//docbase += "/"
	//var staticFiles string
	//info, errFile := os.Stat(docbase + "web")
	//if os.IsNotExist(errFile) || !info.IsDir() {
	//	staticFiles = docbase + resourcedir
	//} else {
	//	staticFiles = docbase + "web/"
	//}
	staticFiles := ResourceBase(Base)
	fmt.Printf("static files / is %v\n", staticFiles)
	resources := http.FileServer(http.Dir(staticFiles))
	// access to resource folders using regular expression
	router.PathPrefix("/{dir:(?:css|fonts|images)}/").Handler(resources).MatcherFunc(isNotRedirect)

	// path to login/logout/server selection page
	router.HandleFunc(GetLogInOutPath(), accesshandler)
	//
	router.HandleFunc("/{path:.*}", selectServer).Queries("target", "{target:.*}").MatcherFunc(isSession)
	//router.HandleFunc("/{target}", selectServer).MatcherFunc(isSession)
	router.HandleFunc("/err", err)
	router.HandleFunc("/{dispatch:.*}", baseDispatcher)

	router.Use(loggingMiddleware)
	fmt.Printf("initializing router done \n")

	return router
}

func isNotRedirect(request *http.Request, match *mux.RouteMatch) bool {
	session, err := SessionStore().Get(request, S_PROXY)
	if err != nil {
		return true
	}
	uri := request.RequestURI
	redirect := session.Values["redirect"]
	result := redirect == nil
	log.Printf("redirect is %v, uri %v isNotRedirect = %v\n", redirect, uri, result)
	return result
}

func accesshandler(writer http.ResponseWriter, request *http.Request) {
	session, err := SessionStore().Get(request, S_PROXY)
	if err != nil {
		fmt.Errorf("Session error %v\n", err.Error())
		//http.Error(writer, err.Error(), http.StatusInternalServerError)
		//return
	}
	if request.Method == http.MethodGet {
		// logout if already logged in
		if session.Values["user"] != nil {
			// delete the session
			session.Options.MaxAge = -1
			session.Values["user"] = nil
			session.Values["redirect"] = nil
			session.Values["refs"] = nil
			if err := session.Save(request, writer); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		qparam := request.URL.Query()
		loginq, lok := qparam["login"]
		pwq, pok := qparam["pw"]
		targetq, tok := qparam["target"]
		if lok && pok {
			if checkLogin(loginq[0], pwq[0], session, request, writer) && tok {
				findTargetAndDispatch(targetq[0], session, writer, request)
			} else {
				showSelectView(loginq[0], session, request, writer)
			}
		} else {
			values := Viewmodel{
				"submitto": GetLogInOutPath(),
			}
			showTemplate(writer, request, "access", values)
		}
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		login := html.EscapeString(request.PostFormValue("login"))
		passwd := html.EscapeString(request.PostFormValue("password"))
		if checkLogin(login, passwd, session, request, writer) {
			showSelectView(login, session, request, writer)
			return
		}
	}
}

func showSelectView(login string, session *sessions.Session, request *http.Request, writer http.ResponseWriter) {
	refs := NewRefmap(MakeLinks(login))
	session.Values["refs"] = refs
	if err := session.Save(request, writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	values := Viewmodel{
		"submitto": GetLogInOutPath(),
		"targets":  refs.Refs,
	}
	showTemplate(writer, request, "links", values)
	return
}

func checkLogin(login string, passwd string, session *sessions.Session, request *http.Request, writer http.ResponseWriter) bool {
	user, ok := GetUsers()[login]
	if ok && user.Password == passwd {
		session.Values["user"] = user
		session.Options.MaxAge = 600
		return true
	} else {
		values := Viewmodel{
			"submitto": GetLogInOutPath(),
			"failure":  "Failure, wrong login or password",
		}
		showTemplate(writer, request, "access", values)
		//Views().ExecuteTemplate(writer, "access", values)
	}
	return false
}

// set session values and redirect to server
func selectServer(res http.ResponseWriter, req *http.Request) {
	session, err := SessionStore().Get(req, S_PROXY)
	if err != nil && !err.(scc.Error).IsDecode() {
		http.Error(res, "selectServer#1"+err.Error(), http.StatusInternalServerError)
		return
	}
	redirect := ""
	// get query parameter
	vars := mux.Vars(req)
	targetRef := vars["target"]
	// extract reference map from session
	refs, ok := session.Values["refs"].(Refmap)
	if ok {
		redirect = refs.Targets[targetRef]
	}
	// get key for proxy redirection
	if !ok || redirect == "" {
		backToLogin(res, req)
		return
	}
	findTargetAndDispatch(redirect, session, res, req)
}

func findTargetAndDispatch(redirect string, session *sessions.Session, res http.ResponseWriter, req *http.Request) {
	// adjust target url
	target, ok := GetTargets()[redirect]
	if !ok {
		backToLogin(res, req)
		return
	}
	session.Values["redirect"] = redirect
	err := session.Save(req, res)
	if err != nil {
		http.Error(res, "SessionStore error "+err.Error(), http.StatusInternalServerError)
		return
	}
	targetUri := target[1]
	targetUrl, err := url.Parse(targetUri)
	if err != nil {
		http.Error(res, "find target "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		req.URL.Path = targetUrl.Path
		dispatch(redirect, res, req)
	}
}

/**
evaluate session values to dispatch to the right proxy
if no user is logged in, just be quiet in stealth mode
*/
func baseDispatcher(res http.ResponseWriter, req *http.Request) {
	session, err := SessionStore().Get(req, S_PROXY)
	if err != nil {
		http.Error(res, "baseDispatcher#1"+err.Error(), http.StatusInternalServerError)
		return
	}
	user := session.Values["user"]
	redirect := session.Values["redirect"]
	age := session.Options.MaxAge
	if age > -1 && user != nil && redirect != nil {
		dispatch(redirect.(string), res, req)
	} else {
		if isStealth {
			return
		}
		http.Error(res, "", http.StatusMovedPermanently)
		return
	}
}

/**
dispatch request to the appropriate proxy
*/
func dispatch(redirect string, res http.ResponseWriter, req *http.Request) {
	target := GetTargets()[redirect]
	targetUri := target[1]
	targetUrl, err := url.Parse(targetUri)
	if err != nil {
		fmt.Printf("URL parse error %v\n", err)
	} else {
		req.URL.Host = targetUrl.Host
		req.URL.Scheme = targetUrl.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = targetUrl.Host
		proxy := proxies[redirect]
		proxy.ServeHTTP(res, req)
		//fmt.Printf("redirect is %v, Request is M: %v H: '%v' P: %v URI: %v\n", redirect, req.Method, targetUrl.Host, req.Proto, req.RequestURI)
	}
}

func backToLogin(res http.ResponseWriter, req *http.Request) {
	// error, back to login page
	http.Redirect(res, req, GetLogInOutPath(), http.StatusMovedPermanently)
	return
}

func isSession(r *http.Request, rm *mux.RouteMatch) bool {
	session, err := SessionStore().Get(r, S_PROXY)
	if err != nil {
		fmt.Printf("isSession#1: %v\n", err)
		return false
	}
	t := session.Values["user"]
	mxa := session.Options.MaxAge
	//fmt.Printf("isSession#2: %v, age: %v\n", t, mxa)
	return t != nil && mxa > -1
}

func noSession(r *http.Request, rm *mux.RouteMatch) bool {
	ret := !isSession(r, rm)
	return ret
}

func showTemplate(writer http.ResponseWriter, request *http.Request, view string, model Viewmodel) {
	Views().ExecuteTemplate(writer, view, model)
	//writer.Header().Set("Access-Control-Allow-Origin", "*")
	//writer.Header().Set("X-Content-Type-Options", "")
	//fmt.Printf("template response Header: %v\n", writer.Header())
}

func err(writer http.ResponseWriter, request *http.Request) {
	//error := request.Header.Get("Status")
	fmt.Fprintf(writer, "Error with path %s",
		request.URL.Path[1:])
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		if strings.HasSuffix(r.RequestURI, "css") || strings.HasSuffix(r.RequestURI, "ico") {
			log.Println(r.RequestURI)
		} else {
			log.Println(r.RequestURI)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		//if strings.HasSuffix(r.RequestURI, "css") {
		//	w.Header().Set("Content-Type", "text/css, charset=utf-8")
		//	w.Header().Set("X-Content-Type-Options", "")
		//	log.Println(w.Header())
		//	log.Println(w)
		//}
	})
}
