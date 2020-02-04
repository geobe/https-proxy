package main

import (
	"fmt"
	"github.com/geobe/https-proxy/go/controller"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"net/http/httputil"
	"os"
)

const resourcedir = controller.Base + "/web/"
const puma = "http://192.168.111.81:4567/"

var proxy *httputil.ReverseProxy

func main() {

	controller.SetupConfig("testconfig")
	tplrouter := mux.NewRouter()
	// finde Working directory = GOPATH
	docbase, _ := os.Getwd()
	docbase += "/"
	resources := http.FileServer(http.Dir(docbase + resourcedir))
	// Zugriff auf die Resourcen-Verzeichnisse mit regular expression
	tplrouter.PathPrefix("/{dir:(?:css|fonts|images)}/").Handler(resources)
	//fmt.Printf("log in at %s\n", controller.GetLogInOutPath())
	tplrouter.HandleFunc(controller.GetLogInOutPath(), accesshandler)
	//tplrouter.HandleFunc("/", reroute)
	http.ListenAndServe("0.0.0.0:8050", tplrouter)
}

func accesshandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		values := controller.Viewmodel{
			"submitto": controller.GetLogInOutPath(),
		}
		controller.Views().ExecuteTemplate(writer, "access", values)
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		login := html.EscapeString(request.PostFormValue("login"))
		passwd := html.EscapeString(request.PostFormValue("password"))
		user, ok := controller.GetUsers()[login]
		if ok && user.Password == passwd {
			_, targets := controller.MakeLinks(login)
			values := controller.Viewmodel{
				"submitto": controller.GetLogInOutPath(),
				"targets":  targets,
			}
			controller.Views().ExecuteTemplate(writer, "links", values)
		} else {
			values := controller.Viewmodel{
				"submitto": controller.GetLogInOutPath(),
				"failure":  "Failure, wrong login or password",
			}
			controller.Views().ExecuteTemplate(writer, "access", values)
		}

	}
}

//func reroute(writer http.ResponseWriter, request *http.Request) {
//	res := mux.Vars(request)["path"]
//	fmt.Printf("reroute res, m: %v, %v\n%v\n\n", res, request.Method, request.Header)
//	if request.Method == http.MethodOptions {
//		writer.WriteHeader(http.StatusNoContent)
//		writer.Header().Set("Access-Control-Allow-Origin", "null")
//		writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
//		writer.Header().Set("Access-Control-Allow-Headers", "*")
//		//writer.Header().Set("Origin", "HTTP://192.168.101.100:8050")
//		//writer.Header().Set("Allow", "OPTIONS, GET, POST")
//	} else {
//		proxy.ServeHTTP(writer, request)
//		writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.101.100:8050")
//		writer.Header().Set("Access-Control-Allow-Methods", "*")
//		fmt.Printf("reroute Header: %v\n", writer.Header())
//	}
//	fmt.Printf("reroute Header: %v\n", writer.Header())
//}

func showTemplate(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("request m, h, u: %v, %v, %v\n", request.Method, request.Host, request.URL)
	vars := mux.Vars(request)
	values := controller.Viewmodel{
		"page": vars["get"],
	}
	controller.Views().ExecuteTemplate(writer, "access", values)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Printf("template response Header: %v\n", writer.Header())
}
