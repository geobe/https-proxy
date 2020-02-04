package main

import (
	"fmt"
	"github.com/geobe/https-proxy/go/controller"
	"net/http"
)

func main() {

	controller.SetupConfig("testconfig")
	router := controller.InitRouter()

	// die zugelassenen host namen
	//allowedHosts := []string{"geobe.spdns.org"}

	// der Verwalter der LetsEncrypt Zertifikate
	//certManager := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist(allowedHosts...), //your domain here
	//	Email:            "geobe.whz@gmail.com",
	//	Cache:      autocert.DirCache("certs"), //folder for storing certificates
	//}

	//http.HandleFunc("/puma", serveReverseProxy)
	fmt.Printf("Reverse proxy started\n")
	http.ListenAndServe("0.0.0.0:8070", router)
	fmt.Printf("Reverse proxy done\n")
}
