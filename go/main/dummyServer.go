package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{target:.*}", showPath)

	http.ListenAndServe("0.0.0.0:8080", router)

}

func showPath(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	target := vars["target"]
	fmt.Fprintf(writer, "path was %v \nfrom request url host %v, remote address %v, url path %v", target, request.URL.Host, request.RemoteAddr, request.URL.Path)
}
